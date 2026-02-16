package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yurisasc/algafood-go/internal/api/assembler"
	"github.com/yurisasc/algafood-go/internal/api/dto"
	"github.com/yurisasc/algafood-go/internal/api/exceptionhandler"
	"github.com/yurisasc/algafood-go/internal/api/middleware" // Importado
	"github.com/yurisasc/algafood-go/internal/domain/service"
)

type UsuarioHandler struct {
	service               *service.UsuarioService
	authService           *service.AuthService
	tokenBlacklistService *service.TokenBlacklistService
}

func NewUsuarioHandler(service *service.UsuarioService, authService *service.AuthService, tokenBlacklistService *service.TokenBlacklistService) *UsuarioHandler {
	return &UsuarioHandler{
		service:               service,
		authService:           authService,
		tokenBlacklistService: tokenBlacklistService,
	}
}

// Eu retorna os detalhes do usuário autenticado.
func (h *UsuarioHandler) Eu(c *gin.Context) {
	// Extrai o usuário do contexto, que foi colocado lá pelo AuthMiddleware.
	usuario, ok := middleware.GetCurrentUser(c)
	if !ok {
		// Isso não deve acontecer em uma rota protegida, mas é uma boa prática verificar.
		exceptionhandler.HandleUnauthorized(c)
		return
	}

	// Usa o assembler para converter a entidade do usuário para o DTO de resposta.
	c.JSON(http.StatusOK, assembler.ToUsuarioModel(usuario))
}

func (h *UsuarioHandler) Login(c *gin.Context) {
	var input dto.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		exceptionhandler.HandleValidationError(c, err)
		return
	}

	user, err := h.service.Authenticate(input.Email, input.Senha)
	if err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}

	token, err := h.authService.GenerateToken(user)
	if err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Logout invalida o token JWT atual.
func (h *UsuarioHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.Status(http.StatusNoContent)
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		c.Status(http.StatusNoContent)
		return
	}

	tokenString := parts[1]
	if err := h.tokenBlacklistService.InvalidateToken(tokenString); err != nil {
		// Mesmo se houver erro ao invalidar, retornamos sucesso para o cliente
		c.Status(http.StatusNoContent)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *UsuarioHandler) Listar(c *gin.Context) {
	usuarios, err := h.service.FindAll()
	if err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, assembler.ToUsuarioModels(usuarios))
}

func (h *UsuarioHandler) Buscar(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("usuarioId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	usuario, err := h.service.FindByID(id)
	if err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, assembler.ToUsuarioModel(usuario))
}

func (h *UsuarioHandler) Adicionar(c *gin.Context) {
	var input dto.UsuarioComSenhaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		exceptionhandler.HandleValidationError(c, err)
		return
	}

	usuario := assembler.ToUsuarioEntity(&input)
	if err := h.service.Save(usuario); err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, assembler.ToUsuarioModel(usuario))
}

func (h *UsuarioHandler) Atualizar(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("usuarioId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	var input dto.UsuarioInput
	if err := c.ShouldBindJSON(&input); err != nil {
		exceptionhandler.HandleValidationError(c, err)
		return
	}

	usuario, err := h.service.FindByID(id)
	if err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}

	usuario.Nome = input.Nome
	usuario.Email = input.Email
	if err := h.service.Save(usuario); err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, assembler.ToUsuarioModel(usuario))
}

func (h *UsuarioHandler) AlterarSenha(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("usuarioId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	var input dto.SenhaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		exceptionhandler.HandleValidationError(c, err)
		return
	}

	if err := h.service.AlterarSenha(id, input.SenhaAtual, input.NovaSenha); err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ListarGrupos lists groups of a user
func (h *UsuarioHandler) ListarGrupos(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("usuarioId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	usuario, err := h.service.FindByID(id)
	if err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, assembler.ToGrupoModels(usuario.Grupos))
}

// AssociarGrupo adds a group to a user
func (h *UsuarioHandler) AssociarGrupo(c *gin.Context) {
	usuarioID, _ := strconv.ParseUint(c.Param("usuarioId"), 10, 64)
	grupoID, _ := strconv.ParseUint(c.Param("grupoId"), 10, 64)

	if err := h.service.AssociarGrupo(usuarioID, grupoID); err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// DesassociarGrupo removes a group from a user
func (h *UsuarioHandler) DesassociarGrupo(c *gin.Context) {
	usuarioID, _ := strconv.ParseUint(c.Param("usuarioId"), 10, 64)
	grupoID, _ := strconv.ParseUint(c.Param("grupoId"), 10, 64)

	if err := h.service.DesassociarGrupo(usuarioID, grupoID); err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
