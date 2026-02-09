package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/yurisasc/algafood-go/internal/api"
	"github.com/yurisasc/algafood-go/internal/api/handler"
	"github.com/yurisasc/algafood-go/internal/config"
	"github.com/yurisasc/algafood-go/internal/domain/service"
	infraRepo "github.com/yurisasc/algafood-go/internal/infrastructure/repository"
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
	fluxoPedidoSvc := service.NewFluxoPedidoService(pedidoRepo, pedidoSvc)

	// Initialize handlers
	estadoHandler := handler.NewEstadoHandler(estadoSvc)
	cidadeHandler := handler.NewCidadeHandler(cidadeSvc)
	cozinhaHandler := handler.NewCozinhaHandler(cozinhaSvc)
	formaPagamentoHandler := handler.NewFormaPagamentoHandler(formaPagamentoSvc)
	permissaoHandler := handler.NewPermissaoHandler(permissaoSvc)
	grupoHandler := handler.NewGrupoHandler(grupoSvc)
	usuarioHandler := handler.NewUsuarioHandler(usuarioSvc, authSvc)
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
		usuarioSvc, // Modificado
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
