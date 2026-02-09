package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yurisasc/algafood-go/internal/config"
	"github.com/yurisasc/algafood-go/internal/domain/model"
)

type AuthService struct {
	cfg *config.JWTConfig
}

func NewAuthService(cfg *config.JWTConfig) *AuthService {
	return &AuthService{cfg: cfg}
}

func (s *AuthService) GenerateToken(user *model.Usuario) (string, error) {
	claims := jwt.MapClaims{
		"usuario_id": user.ID,
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
		"iat":        time.Now().Unix(),
		"iss":        s.cfg.Issuer,
	}

	// Add authorities to claims
	var authorities []string
	for _, grupo := range user.Grupos {
		for _, permissao := range grupo.Permissoes {
			authorities = append(authorities, permissao.Nome)
		}
	}
	claims["authorities"] = authorities

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.SecretKey))
}
