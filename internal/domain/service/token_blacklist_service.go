package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/yurisasc/algafood-go/internal/config"
)

const (
	// Prefixo usado para as chaves de token invalidados no Redis
	tokenBlacklistPrefix = "token:blacklist:"
)

// TokenBlacklistService gerencia tokens JWT invalidados (logout) usando Redis.
type TokenBlacklistService struct {
	redisClient *redis.Client
	jwtCfg      *config.JWTConfig
}

func NewTokenBlacklistService(redisCfg *config.RedisConfig, jwtCfg *config.JWTConfig) *TokenBlacklistService {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisCfg.Host, redisCfg.Port),
		Password: redisCfg.Password,
		DB:       redisCfg.DB,
	})

	return &TokenBlacklistService{
		redisClient: client,
		jwtCfg:      jwtCfg,
	}
}

// InvalidateToken adiciona um token à blacklist no Redis.
// O token é armazenado com TTL baseado na sua data de expiração.
func (s *TokenBlacklistService) InvalidateToken(tokenString string) error {
	ctx := context.Background()

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtCfg.SecretKey), nil
	})
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil
	}

	// Calcula o TTL baseado na expiração do token
	var ttl time.Duration
	if exp, ok := claims["exp"].(float64); ok {
		expTime := time.Unix(int64(exp), 0)
		ttl = time.Until(expTime)
		if ttl <= 0 {
			// Token já expirado, não precisa adicionar à blacklist
			return nil
		}
	} else {
		// Se não houver expiração, usa 24 horas como padrão
		ttl = 24 * time.Hour
	}

	// Armazena o token na blacklist com TTL
	key := tokenBlacklistPrefix + tokenString
	return s.redisClient.Set(ctx, key, "1", ttl).Err()
}

// IsBlacklisted verifica se um token está na blacklist no Redis.
func (s *TokenBlacklistService) IsBlacklisted(tokenString string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	key := tokenBlacklistPrefix + tokenString

	exists, err := s.redisClient.Exists(ctx, key).Result()
	if err != nil {
		// Em caso de erro ou timeout, considera que o token não está na blacklist
		// para não bloquear usuários em caso de falha do Redis
		log.Printf("[BLACKLIST] Erro ao verificar blacklist: %v", err)
		return false
	}

	if exists > 0 {
		log.Printf("[BLACKLIST] Token está na blacklist.")
	}

	return exists > 0
}

// Close fecha a conexão com o Redis.
func (s *TokenBlacklistService) Close() error {
	return s.redisClient.Close()
}

// Ping verifica a conexão com o Redis.
func (s *TokenBlacklistService) Ping() error {
	ctx := context.Background()
	return s.redisClient.Ping(ctx).Err()
}
