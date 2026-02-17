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
	estadoCachePrefix    = "estado:cache:"
	estadoAllCacheKey    = "estado:all"
	cidadeCachePrefix    = "cidade:cache:"
	cidadeAllCacheKey    = "cidade:all"
	cidadeByEstadoPrefix = "cidade:estado:"

	// TTL padrão para cache de localização (cidades e estados mudam raramente)
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
	ttl         time.Duration
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

	return &LocationCacheService{
		redisClient: client,
		ttl:         locationCacheTTL,
	}
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
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var cached CachedEstado
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil, err
	}

	return cached.ToModel(), nil
}

// SetEstado armazena um estado no cache
func (s *LocationCacheService) SetEstado(estado *model.Estado) error {
	ctx, cancel := context.WithTimeout(context.Background(), locationCacheTimeout)
	defer cancel()

	cached := &CachedEstado{
		ID:   estado.ID,
		Nome: estado.Nome,
	}

	data, err := json.Marshal(cached)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s%d", estadoCachePrefix, estado.ID)
	return s.redisClient.Set(ctx, key, data, s.ttl).Err()
}

// GetAllEstados obtém todos os estados do cache
func (s *LocationCacheService) GetAllEstados() ([]model.Estado, error) {
	ctx, cancel := context.WithTimeout(context.Background(), locationCacheTimeout)
	defer cancel()

	data, err := s.redisClient.Get(ctx, estadoAllCacheKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var cached []CachedEstado
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil, err
	}

	estados := make([]model.Estado, len(cached))
	for i, c := range cached {
		estados[i] = *c.ToModel()
	}

	return estados, nil
}

// SetAllEstados armazena todos os estados no cache
func (s *LocationCacheService) SetAllEstados(estados []model.Estado) error {
	ctx, cancel := context.WithTimeout(context.Background(), locationCacheTimeout)
	defer cancel()

	cached := make([]CachedEstado, len(estados))
	for i, e := range estados {
		cached[i] = CachedEstado{
			ID:   e.ID,
			Nome: e.Nome,
		}
	}

	data, err := json.Marshal(cached)
	if err != nil {
		return err
	}

	return s.redisClient.Set(ctx, estadoAllCacheKey, data, s.ttl).Err()
}

// InvalidateEstado remove um estado do cache
func (s *LocationCacheService) InvalidateEstado(id uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), locationCacheTimeout)
	defer cancel()

	key := fmt.Sprintf("%s%d", estadoCachePrefix, id)
	pipe := s.redisClient.Pipeline()
	pipe.Del(ctx, key)
	pipe.Del(ctx, estadoAllCacheKey)
	// Também invalida cache de cidades do estado
	pipe.Del(ctx, fmt.Sprintf("%s%d", cidadeByEstadoPrefix, id))
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
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var cached CachedCidade
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil, err
	}

	return cached.ToModel(), nil
}

// SetCidade armazena uma cidade no cache
func (s *LocationCacheService) SetCidade(cidade *model.Cidade) error {
	ctx, cancel := context.WithTimeout(context.Background(), locationCacheTimeout)
	defer cancel()

	cached := s.toCachedCidade(cidade)

	data, err := json.Marshal(cached)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s%d", cidadeCachePrefix, cidade.ID)
	return s.redisClient.Set(ctx, key, data, s.ttl).Err()
}

// GetAllCidades obtém todas as cidades do cache
func (s *LocationCacheService) GetAllCidades() ([]model.Cidade, error) {
	ctx, cancel := context.WithTimeout(context.Background(), locationCacheTimeout)
	defer cancel()

	data, err := s.redisClient.Get(ctx, cidadeAllCacheKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var cached []CachedCidade
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil, err
	}

	cidades := make([]model.Cidade, len(cached))
	for i, c := range cached {
		cidades[i] = *c.ToModel()
	}

	return cidades, nil
}

// SetAllCidades armazena todas as cidades no cache
func (s *LocationCacheService) SetAllCidades(cidades []model.Cidade) error {
	ctx, cancel := context.WithTimeout(context.Background(), locationCacheTimeout)
	defer cancel()

	cached := make([]CachedCidade, len(cidades))
	for i, c := range cidades {
		cached[i] = *s.toCachedCidade(&c)
	}

	data, err := json.Marshal(cached)
	if err != nil {
		return err
	}

	return s.redisClient.Set(ctx, cidadeAllCacheKey, data, s.ttl).Err()
}

// GetCidadesByEstado obtém cidades de um estado do cache
func (s *LocationCacheService) GetCidadesByEstado(estadoID uint64) ([]model.Cidade, error) {
	ctx, cancel := context.WithTimeout(context.Background(), locationCacheTimeout)
	defer cancel()

	key := fmt.Sprintf("%s%d", cidadeByEstadoPrefix, estadoID)
	data, err := s.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var cached []CachedCidade
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil, err
	}

	cidades := make([]model.Cidade, len(cached))
	for i, c := range cached {
		cidades[i] = *c.ToModel()
	}

	return cidades, nil
}

// SetCidadesByEstado armazena cidades de um estado no cache
func (s *LocationCacheService) SetCidadesByEstado(estadoID uint64, cidades []model.Cidade) error {
	ctx, cancel := context.WithTimeout(context.Background(), locationCacheTimeout)
	defer cancel()

	cached := make([]CachedCidade, len(cidades))
	for i, c := range cidades {
		cached[i] = *s.toCachedCidade(&c)
	}

	data, err := json.Marshal(cached)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s%d", cidadeByEstadoPrefix, estadoID)
	return s.redisClient.Set(ctx, key, data, s.ttl).Err()
}

// InvalidateCidade remove uma cidade do cache
func (s *LocationCacheService) InvalidateCidade(cidade *model.Cidade) error {
	ctx, cancel := context.WithTimeout(context.Background(), locationCacheTimeout)
	defer cancel()

	pipe := s.redisClient.Pipeline()
	pipe.Del(ctx, fmt.Sprintf("%s%d", cidadeCachePrefix, cidade.ID))
	pipe.Del(ctx, cidadeAllCacheKey)
	if cidade.EstadoID > 0 {
		pipe.Del(ctx, fmt.Sprintf("%s%d", cidadeByEstadoPrefix, cidade.EstadoID))
	}
	_, err := pipe.Exec(ctx)
	return err
}

// InvalidateAllCidades remove todas as cidades do cache
func (s *LocationCacheService) InvalidateAllCidades() error {
	ctx, cancel := context.WithTimeout(context.Background(), locationCacheTimeout*10)
	defer cancel()

	// Busca todas as chaves de cidade
	patterns := []string{cidadeCachePrefix + "*", cidadeByEstadoPrefix + "*"}
	for _, pattern := range patterns {
		keys, err := s.redisClient.Keys(ctx, pattern).Result()
		if err != nil {
			continue
		}
		if len(keys) > 0 {
			s.redisClient.Del(ctx, keys...)
		}
	}
	s.redisClient.Del(ctx, cidadeAllCacheKey)
	return nil
}

// ==================== HELPERS ====================

func (s *LocationCacheService) toCachedCidade(cidade *model.Cidade) *CachedCidade {
	cached := &CachedCidade{
		ID:       cidade.ID,
		Nome:     cidade.Nome,
		EstadoID: cidade.EstadoID,
	}

	if cidade.Estado.ID > 0 {
		cached.Estado = &CachedEstado{
			ID:   cidade.Estado.ID,
			Nome: cidade.Estado.Nome,
		}
	}

	return cached
}

func (c *CachedEstado) ToModel() *model.Estado {
	return &model.Estado{
		ID:   c.ID,
		Nome: c.Nome,
	}
}

func (c *CachedCidade) ToModel() *model.Cidade {
	cidade := &model.Cidade{
		ID:       c.ID,
		Nome:     c.Nome,
		EstadoID: c.EstadoID,
	}

	if c.Estado != nil {
		cidade.Estado = *c.Estado.ToModel()
	}

	return cidade
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
