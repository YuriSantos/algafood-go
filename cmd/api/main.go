package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/yurisasc/algafood-go/internal/api"
	"github.com/yurisasc/algafood-go/internal/api/handler"
	"github.com/yurisasc/algafood-go/internal/config"
	"github.com/yurisasc/algafood-go/internal/domain/service"
	"github.com/yurisasc/algafood-go/internal/infrastructure/email"
	"github.com/yurisasc/algafood-go/internal/infrastructure/eventbridge"
	"github.com/yurisasc/algafood-go/internal/infrastructure/notification"
	infraRepo "github.com/yurisasc/algafood-go/internal/infrastructure/repository"
	"github.com/yurisasc/algafood-go/internal/infrastructure/sqs"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := config.NewDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Connected to database successfully")

	// Initialize repositories
	estadoRepo := infraRepo.NewEstadoRepository(db)
	cidadeRepo := infraRepo.NewCidadeRepository(db)
	cozinhaRepo := infraRepo.NewCozinhaRepository(db)
	formaPagamentoRepo := infraRepo.NewFormaPagamentoRepository(db)
	permissaoRepo := infraRepo.NewPermissaoRepository(db)
	grupoRepo := infraRepo.NewGrupoRepository(db)
	usuarioRepo := infraRepo.NewUsuarioRepository(db)
	restauranteRepo := infraRepo.NewRestauranteRepository(db)
	produtoRepo := infraRepo.NewProdutoRepository(db)
	pedidoRepo := infraRepo.NewPedidoRepository(db)
	vendaQueryRepo := infraRepo.NewVendaQueryRepository(db)

	// Initialize services
	authSvc := service.NewAuthService(&cfg.JWT)
	tokenBlacklistSvc := service.NewTokenBlacklistService(&cfg.Redis, &cfg.JWT)

	// Verifica conexão com Redis
	if err := tokenBlacklistSvc.Ping(); err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v. Token blacklist will not work correctly.", err)
	} else {
		log.Println("Connected to Redis successfully")
	}

	estadoSvc := service.NewEstadoService(estadoRepo)
	cidadeSvc := service.NewCidadeService(cidadeRepo, estadoSvc)
	cozinhaSvc := service.NewCozinhaService(cozinhaRepo)
	formaPagamentoSvc := service.NewFormaPagamentoService(formaPagamentoRepo)
	permissaoSvc := service.NewPermissaoService(permissaoRepo)
	grupoSvc := service.NewGrupoService(grupoRepo, permissaoSvc)
	usuarioSvc := service.NewUsuarioService(usuarioRepo, grupoSvc)
	restauranteSvc := service.NewRestauranteService(restauranteRepo, cozinhaSvc, cidadeSvc, formaPagamentoSvc, usuarioSvc)
	produtoSvc := service.NewProdutoService(produtoRepo, restauranteSvc)
	pedidoSvc := service.NewPedidoService(pedidoRepo, restauranteSvc, cidadeSvc, usuarioSvc, produtoSvc, formaPagamentoSvc)

	// Initialize event publisher
	eventPublisher, err := eventbridge.NewEventPublisher(&cfg.EventBridge, &cfg.SQS, &cfg.AWS)
	if err != nil {
		log.Printf("Warning: Failed to initialize EventBridge publisher: %v. Using fake publisher.", err)
		eventPublisher = eventbridge.NewFakeEventPublisher()
	} else {
		log.Println("EventBridge publisher initialized successfully")
	}

	fluxoPedidoSvc := service.NewFluxoPedidoService(pedidoRepo, pedidoSvc, eventPublisher)

	// Initialize email service
	emailSvc, err := email.NewEmailService(&cfg.Email, &cfg.AWS)
	if err != nil {
		log.Printf("Warning: Failed to initialize email service: %v. Using fake email service.", err)
		emailSvc = email.NewFakeEmailService()
	} else {
		log.Println("Email service initialized successfully")
	}

	// Initialize SQS listener for notifications
	notificationHandler := notification.NewNotificationHandler(emailSvc)
	sqsListener, err := sqs.NewSQSListenerFromConfig(&cfg.SQS, &cfg.AWS, notificationHandler)
	if err != nil {
		log.Printf("Warning: Failed to initialize SQS listener: %v", err)
	} else {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		sqsListener.Start(ctx)
		log.Println("SQS listener started successfully")

		// Graceful shutdown
		go func() {
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
			<-sigChan
			log.Println("Shutting down SQS listener...")
			sqsListener.Stop()
			cancel()
		}()
	}

	// Initialize handlers
	estadoHandler := handler.NewEstadoHandler(estadoSvc)
	cidadeHandler := handler.NewCidadeHandler(cidadeSvc)
	cozinhaHandler := handler.NewCozinhaHandler(cozinhaSvc)
	formaPagamentoHandler := handler.NewFormaPagamentoHandler(formaPagamentoSvc)
	permissaoHandler := handler.NewPermissaoHandler(permissaoSvc)
	grupoHandler := handler.NewGrupoHandler(grupoSvc)
	usuarioHandler := handler.NewUsuarioHandler(usuarioSvc, authSvc, tokenBlacklistSvc)
	restauranteHandler := handler.NewRestauranteHandler(restauranteSvc)
	produtoHandler := handler.NewProdutoHandler(produtoSvc)
	pedidoHandler := handler.NewPedidoHandler(pedidoSvc, fluxoPedidoSvc)
	estatisticaHandler := handler.NewEstatisticaHandler(vendaQueryRepo)

	// Setup Gin
	gin.SetMode(cfg.Server.Mode)
	engine := gin.New()

	// Setup router
	router := api.NewRouter(
		estadoHandler,
		cidadeHandler,
		cozinhaHandler,
		formaPagamentoHandler,
		permissaoHandler,
		grupoHandler,
		usuarioHandler,
		restauranteHandler,
		produtoHandler,
		pedidoHandler,
		estatisticaHandler,
		usuarioSvc,
		tokenBlacklistSvc,
		cfg,
	)
	router.Setup(engine)

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Starting server on %s", addr)
	if err := engine.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
