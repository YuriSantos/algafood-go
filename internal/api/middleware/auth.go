package middleware

import (
	"net/http"
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
)

// AuthMiddleware valida o token JWT e carrega o usuário completo no contexto.
func AuthMiddleware(cfg *config.JWTConfig, usuarioSvc *service.UsuarioService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			exceptionhandler.HandleUnauthorized(c)
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			exceptionhandler.HandleUnauthorized(c)
			c.Abort()
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, exception.NewAuthenticationException("Metodo de assinatura de token invalido")
			}
			return []byte(cfg.SecretKey), nil
		})

		if err != nil || !token.Valid {
			exceptionhandler.HandleUnauthorized(c)
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			exceptionhandler.HandleUnauthorized(c)
			c.Abort()
			return
		}

		// Extrai o ID do usuário e carrega o objeto de usuário completo
		if userIDFloat, ok := claims["usuario_id"].(float64); ok {
			userID := uint64(userIDFloat)
			usuario, err := usuarioSvc.FindByID(userID)
			if err != nil {
				// Se o usuário não for encontrado no DB (ex: foi deletado), a autenticação falha.
				exceptionhandler.HandleUnauthorized(c)
				c.Abort()
				return
			}
			// Armazena o objeto de usuário completo no contexto
			c.Set(string(currentUserKey), usuario)
		} else {
			exceptionhandler.HandleUnauthorized(c)
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetCurrentUser extrai o usuário autenticado do contexto do Gin.
// Isso é análogo ao SecurityContextHolder.getContext().getAuthentication().getPrincipal() do Spring.
func GetCurrentUser(c *gin.Context) (*model.Usuario, bool) {
	user, exists := c.Get(string(currentUserKey))
	if !exists {
		return nil, false
	}
	usuario, ok := user.(*model.Usuario)
	return usuario, ok
}

// HasAuthority verifica se o usuário autenticado tem uma permissão específica.
func HasAuthority(c *gin.Context, authority string) bool {
	usuario, exists := GetCurrentUser(c)
	if !exists {
		return false
	}

	for _, grupo := range usuario.Grupos {
		for _, permissao := range grupo.Permissoes {
			if permissao.Nome == authority {
				return true
			}
		}
	}

	return false
}

// RequireAuthority é um middleware que exige uma permissão específica.
func RequireAuthority(authority string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !HasAuthority(c, authority) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Acesso negado"})
			c.Abort()
			return
		}
		c.Next()
	}
}
