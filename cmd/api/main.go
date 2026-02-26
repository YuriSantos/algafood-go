package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	appCtx, stopSignal := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stopSignal()

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

	// Initialize cache services
	userCacheSvc := service.NewUserCacheService(&cfg.Redis)
	locationCacheSvc := service.NewLocationCacheService(&cfg.Redis)
	businessCacheSvc := service.NewBusinessCacheService(&cfg.Redis)

	// Verifica conexão com Redis
	if err := tokenBlacklistSvc.Ping(); err != nil {
		log.Printf("Aviso: Falha ao conectar ao Redis: %v. Cache e blacklist não funcionarão.", err)
	} else {
		log.Println("Conectado ao Redis com sucesso")
	}

	estadoSvc := service.NewEstadoService(estadoRepo, locationCacheSvc)
	cidadeSvc := service.NewCidadeService(cidadeRepo, estadoSvc, locationCacheSvc)
	cozinhaSvc := service.NewCozinhaService(cozinhaRepo, businessCacheSvc)
	formaPagamentoSvc := service.NewFormaPagamentoService(formaPagamentoRepo, businessCacheSvc)
	permissaoSvc := service.NewPermissaoService(permissaoRepo)
	grupoSvc := service.NewGrupoService(grupoRepo, permissaoSvc)
	usuarioSvc := service.NewUsuarioService(usuarioRepo, grupoSvc, userCacheSvc)
	restauranteSvc := service.NewRestauranteService(restauranteRepo, cozinhaSvc, cidadeSvc, formaPagamentoSvc, usuarioSvc, businessCacheSvc)
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
	var sqsListener sqs.SQSListenerInterface
	sqsEnabled := false

	sqsListener, err = sqs.NewSQSListenerFromConfig(&cfg.SQS, &cfg.AWS, notificationHandler)
	if err != nil {
		log.Printf("Warning: Failed to initialize SQS listener: %v", err)
	} else {
		sqsListener.Start(appCtx)
		sqsEnabled = true
		log.Println("SQS listener started successfully")
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
	httpServer := &http.Server{
		Addr:    addr,
		Handler: engine,
	}

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("Iniciando servidor em %s", addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		if err != nil {
			log.Fatalf("Falha ao iniciar servidor: %v", err)
		}
	case <-appCtx.Done():
		log.Println("Sinal de encerramento recebido. Finalizando aplicação...")
	}

	if sqsEnabled {
		log.Println("Parando SQS listener...")
		sqsListener.Stop()
	}

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("Erro ao encerrar servidor HTTP: %v", err)
	} else {
		log.Println("Servidor HTTP encerrado com sucesso")
	}
}
