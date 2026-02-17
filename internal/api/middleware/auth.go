package middleware

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yurisasc/algafood-go/internal/api/exceptionhandler"
	"github.com/yurisasc/algafood-go/internal/config"
	"github.com/yurisasc/algafood-go/internal/domain/exception"
	"github.com/yurisasc/algafood-go/internal/domain/model"
	"github.com/yurisasc/algafood-go/internal/domain/service"
)

type contextKey string

const (
	// currentUserKey é a chave usada para armazenar o objeto de usuário no contexto
	currentUserKey contextKey = "currentUser"
	// userIDKey é a chave usada para armazenar o ID do usuário no contexto
	userIDKey contextKey = "userId"
	// authoritiesKey é a chave usada para armazenar as authorities no contexto
	authoritiesKey contextKey = "authorities"
	// usuarioSvcKey é a chave para armazenar o serviço de usuário no contexto
	usuarioSvcKey contextKey = "usuarioSvc"
)

// AuthMiddleware valida o token JWT e extrai as informações do token.
// O usuário é obtido do cache Redis quando disponível, ou do banco quando necessário.
func AuthMiddleware(cfg *config.JWTConfig, usuarioSvc *service.UsuarioService, tokenBlacklistSvc *service.TokenBlacklistService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Printf("[AUTH] Falha: Header Authorization ausente")
			exceptionhandler.HandleUnauthorized(c)
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			log.Printf("[AUTH] Falha: Formato de Authorization inválido: %s", authHeader)
			exceptionhandler.HandleUnauthorized(c)
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Verifica se o token está na blacklist (logout) - operação rápida com timeout
		if tokenBlacklistSvc != nil && tokenBlacklistSvc.IsBlacklisted(tokenString) {
			log.Printf("[AUTH] Falha: Token está na blacklist")
			exceptionhandler.HandleUnauthorized(c)
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, exception.NewAuthenticationException("Método de assinatura de token inválido")
			}
			return []byte(cfg.SecretKey), nil
		})

		if err != nil || !token.Valid {
			log.Printf("[AUTH] Falha: Erro ao parsear token: %v, válido: %v", err, token != nil && token.Valid)
			exceptionhandler.HandleUnauthorized(c)
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Printf("[AUTH] Falha: Não foi possível extrair claims do token")
			exceptionhandler.HandleUnauthorized(c)
			c.Abort()
			return
		}

		// Extrai o ID do usuário do token
		userIDFloat, ok := claims["usuario_id"].(float64)
		if !ok {
			log.Printf("[AUTH] Falha: usuario_id não encontrado ou inválido nas claims: %v", claims)
			exceptionhandler.HandleUnauthorized(c)
			c.Abort()
			return
		}
		userID := uint64(userIDFloat)
		log.Printf("[AUTH] Sucesso: Usuário autenticado ID=%d", userID)

		// Extrai authorities do token JWT (já estão no token após login)
		var authorities []string
		if authList, ok := claims["authorities"].([]interface{}); ok {
			for _, auth := range authList {
				if authStr, ok := auth.(string); ok {
					authorities = append(authorities, authStr)
				}
			}
		}

		// Se não tem authorities no token, tenta obter do cache Redis
		if len(authorities) == 0 {
			if cachedAuthorities, err := usuarioSvc.GetAuthoritiesFromCache(userID); err == nil && len(cachedAuthorities) > 0 {
				authorities = cachedAuthorities
			}
		}

		// Armazena no contexto - NÃO carrega o usuário completo do banco ainda
		c.Set(string(userIDKey), userID)
		c.Set(string(authoritiesKey), authorities)
		c.Set(string(usuarioSvcKey), usuarioSvc)

		c.Next()
	}
}

// GetCurrentUserID retorna o ID do usuário autenticado do contexto.
func GetCurrentUserID(c *gin.Context) (uint64, bool) {
	userID, exists := c.Get(string(userIDKey))
	if !exists {
		return 0, false
	}
	id, ok := userID.(uint64)
	return id, ok
}

// GetCurrentUser carrega o usuário do cache Redis ou do banco de dados (lazy loading).
// Use apenas quando realmente precisar dos dados completos do usuário.
func GetCurrentUser(c *gin.Context) (*model.Usuario, bool) {
	// Verifica se já foi carregado anteriormente nesta requisição
	if user, exists := c.Get(string(currentUserKey)); exists {
		usuario, ok := user.(*model.Usuario)
		return usuario, ok
	}

	// Carrega o usuário (primeiro do cache, depois do banco)
	userID, ok := GetCurrentUserID(c)
	if !ok {
		return nil, false
	}

	svc, exists := c.Get(string(usuarioSvcKey))
	if !exists {
		return nil, false
	}
	usuarioSvc, ok := svc.(*service.UsuarioService)
	if !ok {
		return nil, false
	}

	// FindByID já usa cache internamente
	usuario, err := usuarioSvc.FindByID(userID)
	if err != nil {
		return nil, false
	}

	// Cache na requisição para evitar múltiplas consultas
	c.Set(string(currentUserKey), usuario)
	return usuario, true
}
