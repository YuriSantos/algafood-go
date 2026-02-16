package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yurisasc/algafood-go/internal/api/handler"
	"github.com/yurisasc/algafood-go/internal/api/middleware"
	"github.com/yurisasc/algafood-go/internal/config"
	"github.com/yurisasc/algafood-go/internal/domain/service"
)

type Router struct {
	estadoHandler         *handler.EstadoHandler
	cidadeHandler         *handler.CidadeHandler
	cozinhaHandler        *handler.CozinhaHandler
	formaPagamentoHandler *handler.FormaPagamentoHandler
	permissaoHandler      *handler.PermissaoHandler
	grupoHandler          *handler.GrupoHandler
	usuarioHandler        *handler.UsuarioHandler
	restauranteHandler    *handler.RestauranteHandler
	produtoHandler        *handler.ProdutoHandler
	pedidoHandler         *handler.PedidoHandler
	estatisticaHandler    *handler.EstatisticaHandler
	usuarioSvc            *service.UsuarioService
	tokenBlacklistSvc     *service.TokenBlacklistService
	cfg                   *config.Config
}

func NewRouter(
	estadoHandler *handler.EstadoHandler,
	cidadeHandler *handler.CidadeHandler,
	cozinhaHandler *handler.CozinhaHandler,
	formaPagamentoHandler *handler.FormaPagamentoHandler,
	permissaoHandler *handler.PermissaoHandler,
	grupoHandler *handler.GrupoHandler,
	usuarioHandler *handler.UsuarioHandler,
	restauranteHandler *handler.RestauranteHandler,
	produtoHandler *handler.ProdutoHandler,
	pedidoHandler *handler.PedidoHandler,
	estatisticaHandler *handler.EstatisticaHandler,
	usuarioSvc *service.UsuarioService,
	tokenBlacklistSvc *service.TokenBlacklistService,
	cfg *config.Config,
) *Router {
	return &Router{
		estadoHandler:         estadoHandler,
		cidadeHandler:         cidadeHandler,
		cozinhaHandler:        cozinhaHandler,
		formaPagamentoHandler: formaPagamentoHandler,
		permissaoHandler:      permissaoHandler,
		grupoHandler:          grupoHandler,
		usuarioHandler:        usuarioHandler,
		restauranteHandler:    restauranteHandler,
		produtoHandler:        produtoHandler,
		pedidoHandler:         pedidoHandler,
		estatisticaHandler:    estatisticaHandler,
		usuarioSvc:            usuarioSvc,
		tokenBlacklistSvc:     tokenBlacklistSvc,
		cfg:                   cfg,
	}
}

func (r *Router) Setup(engine *gin.Engine) {
	// Global middleware
	engine.Use(middleware.CorsMiddleware())
	engine.Use(middleware.LoggerMiddleware())
	engine.Use(middleware.RecoveryMiddleware())

	// API v1 routes
	v1 := engine.Group("/v1")
	{
		// Public routes (no auth required)
		r.setupPublicRoutes(v1)

		// Protected routes (auth required)
		v1.Use(middleware.AuthMiddleware(&r.cfg.JWT, r.usuarioSvc, r.tokenBlacklistSvc))
		r.setupProtectedRoutes(v1)
	}
}

func (r *Router) setupPublicRoutes(rg *gin.RouterGroup) {
	rg.POST("/login", r.usuarioHandler.Login)
	rg.POST("/usuarios", r.usuarioHandler.Adicionar)
}

func (r *Router) setupProtectedRoutes(rg *gin.RouterGroup) {
	// Logout
	rg.POST("/logout", r.usuarioHandler.Logout)

	// Usuarios
	usuarios := rg.Group("/usuarios")
	{
		usuarios.GET("/eu", r.usuarioHandler.Eu) // Nova rota
		usuarios.GET("", r.usuarioHandler.Listar)
		usuarios.GET("/:usuarioId", r.usuarioHandler.Buscar)
		usuarios.PUT("/:usuarioId", r.usuarioHandler.Atualizar)
		usuarios.PUT("/:usuarioId/senha", r.usuarioHandler.AlterarSenha)

		// Usuario Grupos
		usuarios.GET("/:usuarioId/grupos", r.usuarioHandler.ListarGrupos)
		usuarios.PUT("/:usuarioId/grupos/:grupoId", r.usuarioHandler.AssociarGrupo)
		usuarios.DELETE("/:usuarioId/grupos/:grupoId", r.usuarioHandler.DesassociarGrupo)
	}

	// Estados
	estados := rg.Group("/estados")
	{
		estados.GET("", r.estadoHandler.Listar)
		estados.GET("/:estadoId", r.estadoHandler.Buscar)
		estados.POST("", r.estadoHandler.Adicionar)
		estados.PUT("/:estadoId", r.estadoHandler.Atualizar)
		estados.DELETE("/:estadoId", r.estadoHandler.Remover)
	}

	// Cidades
	cidades := rg.Group("/cidades")
	{
		cidades.GET("", r.cidadeHandler.Listar)
		cidades.GET("/:cidadeId", r.cidadeHandler.Buscar)
		cidades.POST("", r.cidadeHandler.Adicionar)
		cidades.PUT("/:cidadeId", r.cidadeHandler.Atualizar)
		cidades.DELETE("/:cidadeId", r.cidadeHandler.Remover)
	}

	// Cozinhas
	cozinhas := rg.Group("/cozinhas")
	{
		cozinhas.GET("", r.cozinhaHandler.Listar)
		cozinhas.GET("/:cozinhaId", r.cozinhaHandler.Buscar)
		cozinhas.POST("", r.cozinhaHandler.Adicionar)
		cozinhas.PUT("/:cozinhaId", r.cozinhaHandler.Atualizar)
		cozinhas.DELETE("/:cozinhaId", r.cozinhaHandler.Remover)
	}

	// Formas de Pagamento
	formasPagamento := rg.Group("/formas-pagamento")
	{
		formasPagamento.GET("", r.formaPagamentoHandler.Listar)
		formasPagamento.GET("/:formaPagamentoId", r.formaPagamentoHandler.Buscar)
		formasPagamento.POST("", r.formaPagamentoHandler.Adicionar)
		formasPagamento.PUT("/:formaPagamentoId", r.formaPagamentoHandler.Atualizar)
		formasPagamento.DELETE("/:formaPagamentoId", r.formaPagamentoHandler.Remover)
	}

	// Permissoes
	permissoes := rg.Group("/permissoes")
	{
		permissoes.GET("", r.permissaoHandler.Listar)
		permissoes.GET("/:permissaoId", r.permissaoHandler.Buscar)
	}

	// Grupos
	grupos := rg.Group("/grupos")
	{
		grupos.GET("", r.grupoHandler.Listar)
		grupos.GET("/:grupoId", r.grupoHandler.Buscar)
		grupos.POST("", r.grupoHandler.Adicionar)
		grupos.PUT("/:grupoId", r.grupoHandler.Atualizar)
		grupos.DELETE("/:grupoId", r.grupoHandler.Remover)

		// Grupo Permissoes
		grupos.GET("/:grupoId/permissoes", r.grupoHandler.ListarPermissoes)
		grupos.PUT("/:grupoId/permissoes/:permissaoId", r.grupoHandler.AssociarPermissao)
		grupos.DELETE("/:grupoId/permissoes/:permissaoId", r.grupoHandler.DesassociarPermissao)
	}

	// Restaurantes
	restaurantes := rg.Group("/restaurantes")
	{
		restaurantes.GET("", r.restauranteHandler.Listar)
		restaurantes.GET("/:restauranteId", r.restauranteHandler.Buscar)
		restaurantes.POST("", r.restauranteHandler.Adicionar)
		restaurantes.PUT("/:restauranteId", r.restauranteHandler.Atualizar)
		restaurantes.PUT("/:restauranteId/ativo", r.restauranteHandler.Ativar)
		restaurantes.DELETE("/:restauranteId/ativo", r.restauranteHandler.Inativar)
		restaurantes.PUT("/ativacoes", r.restauranteHandler.AtivarEmMassa)
		restaurantes.DELETE("/ativacoes", r.restauranteHandler.InativarEmMassa)
		restaurantes.PUT("/:restauranteId/abertura", r.restauranteHandler.Abrir)
		restaurantes.PUT("/:restauranteId/fechamento", r.restauranteHandler.Fechar)

		// Restaurante Formas Pagamento
		restaurantes.GET("/:restauranteId/formas-pagamento", r.restauranteHandler.ListarFormasPagamento)
		restaurantes.PUT("/:restauranteId/formas-pagamento/:formaPagamentoId", r.restauranteHandler.AssociarFormaPagamento)
		restaurantes.DELETE("/:restauranteId/formas-pagamento/:formaPagamentoId", r.restauranteHandler.DesassociarFormaPagamento)

		// Restaurante Responsaveis
		restaurantes.GET("/:restauranteId/responsaveis", r.restauranteHandler.ListarResponsaveis)
		restaurantes.PUT("/:restauranteId/responsaveis/:usuarioId", r.restauranteHandler.AssociarResponsavel)
		restaurantes.DELETE("/:restauranteId/responsaveis/:usuarioId", r.restauranteHandler.DesassociarResponsavel)

		// Restaurante Produtos
		restaurantes.GET("/:restauranteId/produtos", r.produtoHandler.Listar)
		restaurantes.GET("/:restauranteId/produtos/:produtoId", r.produtoHandler.Buscar)
		restaurantes.POST("/:restauranteId/produtos", r.produtoHandler.Adicionar)
		restaurantes.PUT("/:restauranteId/produtos/:produtoId", r.produtoHandler.Atualizar)
	}

	// Pedidos
	pedidos := rg.Group("/pedidos")
	{
		pedidos.GET("", r.pedidoHandler.Pesquisar)
		pedidos.GET("/:codigoPedido", r.pedidoHandler.Buscar)
		pedidos.POST("", r.pedidoHandler.Adicionar)
		pedidos.PUT("/:codigoPedido/confirmacao", r.pedidoHandler.Confirmar)
		pedidos.PUT("/:codigoPedido/cancelamento", r.pedidoHandler.Cancelar)
		pedidos.PUT("/:codigoPedido/entrega", r.pedidoHandler.Entregar)
	}

	// Estatisticas
	estatisticas := rg.Group("/estatisticas")
	{
		estatisticas.GET("/vendas-diarias", r.estatisticaHandler.ConsultarVendasDiarias)
	}
}
