package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/yurisasc/algafood-go/internal/api/dto"
)

// RecoveryMiddleware recovers from panics and returns a proper error response
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v\n%s", err, debug.Stack())

				problem := dto.NewProblem(
					http.StatusInternalServerError,
					dto.ProblemTypeSystemError,
					"Ocorreu um erro interno inesperado",
					"Ocorreu um erro interno inesperado no sistema. Tente novamente e se o problema persistir, entre em contato com o administrador do sistema.",
				)

				c.AbortWithStatusJSON(http.StatusInternalServerError, problem)
			}
		}()
		c.Next()
	}
}
