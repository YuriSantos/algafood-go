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
	// Prefixos para chaves no Redis
	estadoCachePrefix = "estado:cache:"
	estadoAllCacheKey = "estado:all"
	cidadeCachePrefix = "cidade:cache:"
	cidadeAllCacheKey = "cidade:all"

	// TTL para cache de localização
	locationCacheTTL = 30 * time.Minute

	// Timeout para operações de cache
	locationCacheTimeout = 100 * time.Millisecond
)

// CachedEstado representa um estado em cache
type CachedEstado struct {
	ID   uint64 `json:"id"`
	Nome string `json:"nome"`
}

// CachedCidade representa uma cidade em cache
type CachedCidade struct {
	ID       uint64        `json:"id"`
	Nome     string        `json:"nome"`
	EstadoID uint64        `json:"estadoId"`
	Estado   *CachedEstado `json:"estado,omitempty"`
}

// LocationCacheService gerencia o cache de cidades e estados no Redis
type LocationCacheService struct {
	redisClient *redis.Client
}

// NewLocationCacheService cria um novo serviço de cache de localização
func NewLocationCacheService(redisCfg *config.RedisConfig) *LocationCacheService {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", redisCfg.Host, redisCfg.Port),
		Password:     redisCfg.Password,
		DB:           redisCfg.DB,
		PoolSize:     10,
		MinIdleConns: 2,
		DialTimeout:  2 * time.Second,
		ReadTimeout:  locationCacheTimeout,
		WriteTimeout: locationCacheTimeout,
	})

	return &LocationCacheService{redisClient: client}
}

// ==================== ESTADOS ====================

// GetEstado obtém um estado do cache
func (s *LocationCacheService) GetEstado(id uint64) (*model.Estado, error) {
	ctx, cancel := context.WithTimeout(context.Background(), locationCacheTimeout)
	defer cancel()

	key := fmt.Sprintf("%s%d", estadoCachePrefix, id)
	data, err := s.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var cached CachedEstado
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil, err
	}

	return &model.Estado{ID: cached.ID, Nome: cached.Nome}, nil
}

// SetEstado armazena um estado no cache
func (s *LocationCacheService) SetEstado(estado *model.Estado) error {
	ctx, cancel := context.WithTimeout(context.Background(), locationCacheTimeout)
	defer cancel()

	cached := CachedEstado{ID: estado.ID, Nome: estado.Nome}
	data, err := json.Marshal(cached)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s%d", estadoCachePrefix, estado.ID)
	return s.redisClient.Set(ctx, key, data, locationCacheTTL).Err()
}

// GetAllEstados obtém todos os estados do cache
func (s *LocationCacheService) GetAllEstados() ([]model.Estado, error) {
	ctx, cancel := context.WithTimeout(context.Background(), locationCacheTimeout)
	defer cancel()

	data, err := s.redisClient.Get(ctx, estadoAllCacheKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var cached []CachedEstado
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil, err
	}

	estados := make([]model.Estado, len(cached))
	for i, c := range cached {
		estados[i] = model.Estado{ID: c.ID, Nome: c.Nome}
	}

	return estados, nil
}

// SetAllEstados armazena todos os estados no cache
func (s *LocationCacheService) SetAllEstados(estados []model.Estado) error {
	ctx, cancel := context.WithTimeout(context.Background(), locationCacheTimeout)
	defer cancel()

	cached := make([]CachedEstado, len(estados))
	for i, e := range estados {
		cached[i] = CachedEstado{ID: e.ID, Nome: e.Nome}
	}

	data, err := json.Marshal(cached)
	if err != nil {
		return err
	}

	return s.redisClient.Set(ctx, estadoAllCacheKey, data, locationCacheTTL).Err()
}

// InvalidateEstado remove um estado do cache
func (s *LocationCacheService) InvalidateEstado(id uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), locationCacheTimeout)
	defer cancel()

	pipe := s.redisClient.Pipeline()
	pipe.Del(ctx, fmt.Sprintf("%s%d", estadoCachePrefix, id))
	pipe.Del(ctx, estadoAllCacheKey)
	_, err := pipe.Exec(ctx)
	return err
}

// ==================== CIDADES ====================

// GetCidade obtém uma cidade do cache
func (s *LocationCacheService) GetCidade(id uint64) (*model.Cidade, error) {
	ctx, cancel := context.WithTimeout(context.Background(), locationCacheTimeout)
	defer cancel()

	key := fmt.Sprintf("%s%d", cidadeCachePrefix, id)
	data, err := s.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var cached CachedCidade
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil, err
	}

	cidade := &model.Cidade{
		ID:       cached.ID,
		Nome:     cached.Nome,
		EstadoID: cached.EstadoID,
	}
	if cached.Estado != nil {
		cidade.Estado = model.Estado{ID: cached.Estado.ID, Nome: cached.Estado.Nome}
	}

	return cidade, nil
}

// SetCidade armazena uma cidade no cache
func (s *LocationCacheService) SetCidade(cidade *model.Cidade) error {
	ctx, cancel := context.WithTimeout(context.Background(), locationCacheTimeout)
	defer cancel()

	cached := CachedCidade{
		ID:       cidade.ID,
		Nome:     cidade.Nome,
		EstadoID: cidade.EstadoID,
	}
	if cidade.Estado.ID > 0 {
		cached.Estado = &CachedEstado{ID: cidade.Estado.ID, Nome: cidade.Estado.Nome}
	}

	data, err := json.Marshal(cached)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s%d", cidadeCachePrefix, cidade.ID)
	return s.redisClient.Set(ctx, key, data, locationCacheTTL).Err()
}

// GetAllCidades obtém todas as cidades do cache
func (s *LocationCacheService) GetAllCidades() ([]model.Cidade, error) {
	ctx, cancel := context.WithTimeout(context.Background(), locationCacheTimeout)
	defer cancel()

	data, err := s.redisClient.Get(ctx, cidadeAllCacheKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var cached []CachedCidade
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil, err
	}

	cidades := make([]model.Cidade, len(cached))
	for i, c := range cached {
		cidades[i] = model.Cidade{
			ID:       c.ID,
			Nome:     c.Nome,
			EstadoID: c.EstadoID,
		}
		if c.Estado != nil {
			cidades[i].Estado = model.Estado{ID: c.Estado.ID, Nome: c.Estado.Nome}
		}
	}

	return cidades, nil
}

// SetAllCidades armazena todas as cidades no cache
func (s *LocationCacheService) SetAllCidades(cidades []model.Cidade) error {
	ctx, cancel := context.WithTimeout(context.Background(), locationCacheTimeout)
	defer cancel()

	cached := make([]CachedCidade, len(cidades))
	for i, c := range cidades {
		cached[i] = CachedCidade{
			ID:       c.ID,
			Nome:     c.Nome,
			EstadoID: c.EstadoID,
		}
		if c.Estado.ID > 0 {
			cached[i].Estado = &CachedEstado{ID: c.Estado.ID, Nome: c.Estado.Nome}
		}
	}

	data, err := json.Marshal(cached)
	if err != nil {
		return err
	}

	return s.redisClient.Set(ctx, cidadeAllCacheKey, data, locationCacheTTL).Err()
}

// InvalidateCidade remove uma cidade do cache
func (s *LocationCacheService) InvalidateCidade(cidade *model.Cidade) error {
	ctx, cancel := context.WithTimeout(context.Background(), locationCacheTimeout)
	defer cancel()

	pipe := s.redisClient.Pipeline()
	pipe.Del(ctx, fmt.Sprintf("%s%d", cidadeCachePrefix, cidade.ID))
	pipe.Del(ctx, cidadeAllCacheKey)
	_, err := pipe.Exec(ctx)
	return err
}

// Ping verifica a conexão com o Redis
func (s *LocationCacheService) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), locationCacheTimeout)
	defer cancel()
	return s.redisClient.Ping(ctx).Err()
}

// Close fecha a conexão com o Redis
func (s *LocationCacheService) Close() error {
	return s.redisClient.Close()
}

// WarmUpCache pré-carrega estados e cidades no cache
func (s *LocationCacheService) WarmUpCache(estados []model.Estado, cidades []model.Cidade) {
	if err := s.SetAllEstados(estados); err != nil {
		log.Printf("Aviso: Falha ao pré-carregar estados no cache: %v", err)
	} else {
		log.Printf("Cache aquecido com %d estados", len(estados))
	}

	if err := s.SetAllCidades(cidades); err != nil {
		log.Printf("Aviso: Falha ao pré-carregar cidades no cache: %v", err)
	} else {
		log.Printf("Cache aquecido com %d cidades", len(cidades))
	}
}
