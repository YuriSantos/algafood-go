package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yurisasc/algafood-go/internal/config"
	"github.com/yurisasc/algafood-go/internal/domain/model"
)

const (
	// Prefixo para chaves no Redis
	userCachePrefix = "user:cache:"

	// TTL padrão para cache de usuário
	userCacheTTL = 5 * time.Minute

	// Timeout para operações de cache
	userCacheTimeout = 100 * time.Millisecond
)

// CachedUser representa os dados do usuário armazenados em cache
type CachedUser struct {
	ID          uint64        `json:"id"`
	Nome        string        `json:"nome"`
	Email       string        `json:"email"`
	Grupos      []CachedGrupo `json:"grupos"`
	Authorities []string      `json:"authorities"`
}

// CachedGrupo representa um grupo em cache
type CachedGrupo struct {
	ID         uint64            `json:"id"`
	Nome       string            `json:"nome"`
	Permissoes []CachedPermissao `json:"permissoes"`
}

// CachedPermissao representa uma permissão em cache
type CachedPermissao struct {
	ID        uint64 `json:"id"`
	Nome      string `json:"nome"`
	Descricao string `json:"descricao"`
}

// UserCacheService gerencia o cache de usuários no Redis
type UserCacheService struct {
	redisClient *redis.Client
}

// NewUserCacheService cria um novo serviço de cache de usuários
func NewUserCacheService(redisCfg *config.RedisConfig) *UserCacheService {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", redisCfg.Host, redisCfg.Port),
		Password:     redisCfg.Password,
		DB:           redisCfg.DB,
		PoolSize:     10,
		MinIdleConns: 2,
		DialTimeout:  2 * time.Second,
		ReadTimeout:  userCacheTimeout,
		WriteTimeout: userCacheTimeout,
	})

	return &UserCacheService{redisClient: client}
}

// GetUser obtém um usuário do cache
func (s *UserCacheService) GetUser(userID uint64) (*CachedUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), userCacheTimeout)
	defer cancel()

	key := fmt.Sprintf("%s%d", userCachePrefix, userID)
	data, err := s.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var cachedUser CachedUser
	if err := json.Unmarshal(data, &cachedUser); err != nil {
		return nil, err
	}

	return &cachedUser, nil
}

// SetUser armazena um usuário no cache
func (s *UserCacheService) SetUser(user *model.Usuario) error {
	ctx, cancel := context.WithTimeout(context.Background(), userCacheTimeout)
	defer cancel()

	cachedUser := s.toCachedUser(user)
	data, err := json.Marshal(cachedUser)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s%d", userCachePrefix, user.ID)
	return s.redisClient.Set(ctx, key, data, userCacheTTL).Err()
}

// InvalidateUser remove um usuário do cache
func (s *UserCacheService) InvalidateUser(userID uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), userCacheTimeout)
	defer cancel()

	key := fmt.Sprintf("%s%d", userCachePrefix, userID)
	return s.redisClient.Del(ctx, key).Err()
}

// GetAuthorities obtém as authorities de um usuário do cache
func (s *UserCacheService) GetAuthorities(userID uint64) ([]string, error) {
	cachedUser, err := s.GetUser(userID)
	if err != nil || cachedUser == nil {
		return nil, err
	}
	return cachedUser.Authorities, nil
}

// toCachedUser converte um modelo de usuário para a versão em cache
func (s *UserCacheService) toCachedUser(user *model.Usuario) *CachedUser {
	cachedUser := &CachedUser{
		ID:    user.ID,
		Nome:  user.Nome,
		Email: user.Email,
	}

	authoritiesMap := make(map[string]bool)
	for _, grupo := range user.Grupos {
		cachedGrupo := CachedGrupo{
			ID:   grupo.ID,
			Nome: grupo.Nome,
		}

		for _, permissao := range grupo.Permissoes {
			cachedGrupo.Permissoes = append(cachedGrupo.Permissoes, CachedPermissao{
				ID:        permissao.ID,
				Nome:      permissao.Nome,
				Descricao: permissao.Descricao,
			})
			authoritiesMap[permissao.Nome] = true
		}

		cachedUser.Grupos = append(cachedUser.Grupos, cachedGrupo)
	}

	for auth := range authoritiesMap {
		cachedUser.Authorities = append(cachedUser.Authorities, auth)
	}

	return cachedUser
}

// ToModel converte um CachedUser para model.Usuario
func (c *CachedUser) ToModel() *model.Usuario {
	usuario := &model.Usuario{
		ID:    c.ID,
		Nome:  c.Nome,
		Email: c.Email,
	}

	for _, cg := range c.Grupos {
		grupo := model.Grupo{
			ID:   cg.ID,
			Nome: cg.Nome,
		}

		for _, cp := range cg.Permissoes {
			grupo.Permissoes = append(grupo.Permissoes, model.Permissao{
				ID:        cp.ID,
				Nome:      cp.Nome,
				Descricao: cp.Descricao,
			})
		}

		usuario.Grupos = append(usuario.Grupos, grupo)
	}

	return usuario
}

// Ping verifica a conexão com o Redis
func (s *UserCacheService) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), userCacheTimeout)
	defer cancel()
	return s.redisClient.Ping(ctx).Err()
}

// Close fecha a conexão com o Redis
func (s *UserCacheService) Close() error {
	return s.redisClient.Close()
}

// WarmUpCache pré-carrega usuários no cache
func (s *UserCacheService) WarmUpCache(users []*model.Usuario) {
	for _, user := range users {
		if err := s.SetUser(user); err != nil {
			log.Printf("Aviso: Falha ao pré-carregar usuário %d no cache: %v", user.ID, err)
		}
	}
	log.Printf("Cache de usuários aquecido com %d usuários", len(users))
}
