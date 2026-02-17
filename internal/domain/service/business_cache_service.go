package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	"github.com/yurisasc/algafood-go/internal/config"
	"github.com/yurisasc/algafood-go/internal/domain/model"
)

const (
	// Prefixos para chaves no Redis
	restauranteCachePrefix    = "restaurante:cache:"
	restauranteAllCacheKey    = "restaurante:all"
	cozinhaCachePrefix        = "cozinha:cache:"
	cozinhaAllCacheKey        = "cozinha:all"
	formaPagamentoCachePrefix = "forma_pagamento:cache:"
	formaPagamentoAllCacheKey = "forma_pagamento:all"

	// TTL para cache de restaurantes (mais curto pois mudam mais frequentemente)
	restauranteCacheTTL    = 10 * time.Minute
	cozinhaCacheTTL        = 30 * time.Minute
	formaPagamentoCacheTTL = 30 * time.Minute

	// Timeout para operações de cache
	businessCacheTimeout = 100 * time.Millisecond
)

// CachedCozinha representa uma cozinha em cache
type CachedCozinha struct {
	ID   uint64 `json:"id"`
	Nome string `json:"nome"`
}

// CachedFormaPagamento representa uma forma de pagamento em cache
type CachedFormaPagamento struct {
	ID        uint64 `json:"id"`
	Descricao string `json:"descricao"`
}

// CachedRestaurante representa um restaurante em cache (versão simplificada)
type CachedRestaurante struct {
	ID                 uint64          `json:"id"`
	Nome               string          `json:"nome"`
	TaxaFrete          decimal.Decimal `json:"taxaFrete"`
	Ativo              bool            `json:"ativo"`
	Aberto             bool            `json:"aberto"`
	CozinhaID          uint64          `json:"cozinhaId"`
	Cozinha            *CachedCozinha  `json:"cozinha,omitempty"`
	EnderecoCidadeID   uint64          `json:"enderecoCidadeId,omitempty"`
	FormasPagamentoIDs []uint64        `json:"formasPagamentoIds,omitempty"`
	ResponsaveisIDs    []uint64        `json:"responsaveisIds,omitempty"`
}

// BusinessCacheService gerencia o cache de entidades de negócio (restaurantes, cozinhas, etc.)
type BusinessCacheService struct {
	redisClient *redis.Client
}

// NewBusinessCacheService cria um novo serviço de cache de negócio
func NewBusinessCacheService(redisCfg *config.RedisConfig) *BusinessCacheService {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", redisCfg.Host, redisCfg.Port),
		Password:     redisCfg.Password,
		DB:           redisCfg.DB,
		PoolSize:     10,
		MinIdleConns: 2,
		DialTimeout:  2 * time.Second,
		ReadTimeout:  businessCacheTimeout,
		WriteTimeout: businessCacheTimeout,
	})

	return &BusinessCacheService{
		redisClient: client,
	}
}

// ==================== COZINHAS ====================

// GetCozinha obtém uma cozinha do cache
func (s *BusinessCacheService) GetCozinha(id uint64) (*model.Cozinha, error) {
	ctx, cancel := context.WithTimeout(context.Background(), businessCacheTimeout)
	defer cancel()

	key := fmt.Sprintf("%s%d", cozinhaCachePrefix, id)
	data, err := s.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var cached CachedCozinha
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil, err
	}

	return &model.Cozinha{ID: cached.ID, Nome: cached.Nome}, nil
}

// SetCozinha armazena uma cozinha no cache
func (s *BusinessCacheService) SetCozinha(cozinha *model.Cozinha) error {
	ctx, cancel := context.WithTimeout(context.Background(), businessCacheTimeout)
	defer cancel()

	cached := CachedCozinha{ID: cozinha.ID, Nome: cozinha.Nome}
	data, err := json.Marshal(cached)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s%d", cozinhaCachePrefix, cozinha.ID)
	return s.redisClient.Set(ctx, key, data, cozinhaCacheTTL).Err()
}

// GetAllCozinhas obtém todas as cozinhas do cache
func (s *BusinessCacheService) GetAllCozinhas() ([]model.Cozinha, error) {
	ctx, cancel := context.WithTimeout(context.Background(), businessCacheTimeout)
	defer cancel()

	data, err := s.redisClient.Get(ctx, cozinhaAllCacheKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var cached []CachedCozinha
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil, err
	}

	cozinhas := make([]model.Cozinha, len(cached))
	for i, c := range cached {
		cozinhas[i] = model.Cozinha{ID: c.ID, Nome: c.Nome}
	}

	return cozinhas, nil
}

// SetAllCozinhas armazena todas as cozinhas no cache
func (s *BusinessCacheService) SetAllCozinhas(cozinhas []model.Cozinha) error {
	ctx, cancel := context.WithTimeout(context.Background(), businessCacheTimeout)
	defer cancel()

	cached := make([]CachedCozinha, len(cozinhas))
	for i, c := range cozinhas {
		cached[i] = CachedCozinha{ID: c.ID, Nome: c.Nome}
	}

	data, err := json.Marshal(cached)
	if err != nil {
		return err
	}

	return s.redisClient.Set(ctx, cozinhaAllCacheKey, data, cozinhaCacheTTL).Err()
}

// InvalidateCozinha invalida o cache de uma cozinha
func (s *BusinessCacheService) InvalidateCozinha(id uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), businessCacheTimeout)
	defer cancel()

	pipe := s.redisClient.Pipeline()
	pipe.Del(ctx, fmt.Sprintf("%s%d", cozinhaCachePrefix, id))
	pipe.Del(ctx, cozinhaAllCacheKey)
	_, err := pipe.Exec(ctx)
	return err
}

// ==================== RESTAURANTES ====================

// GetRestaurante obtém um restaurante do cache
func (s *BusinessCacheService) GetRestaurante(id uint64) (*CachedRestaurante, error) {
	ctx, cancel := context.WithTimeout(context.Background(), businessCacheTimeout)
	defer cancel()

	key := fmt.Sprintf("%s%d", restauranteCachePrefix, id)
	data, err := s.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var cached CachedRestaurante
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil, err
	}

	return &cached, nil
}

// SetRestaurante armazena um restaurante no cache
func (s *BusinessCacheService) SetRestaurante(restaurante *model.Restaurante) error {
	ctx, cancel := context.WithTimeout(context.Background(), businessCacheTimeout)
	defer cancel()

	cached := s.toCachedRestaurante(restaurante)
	data, err := json.Marshal(cached)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s%d", restauranteCachePrefix, restaurante.ID)
	return s.redisClient.Set(ctx, key, data, restauranteCacheTTL).Err()
}

// InvalidateRestaurante invalida o cache de um restaurante
func (s *BusinessCacheService) InvalidateRestaurante(id uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), businessCacheTimeout)
	defer cancel()

	pipe := s.redisClient.Pipeline()
	pipe.Del(ctx, fmt.Sprintf("%s%d", restauranteCachePrefix, id))
	pipe.Del(ctx, restauranteAllCacheKey)
	_, err := pipe.Exec(ctx)
	return err
}

// toCachedRestaurante converte um restaurante para versão em cache
func (s *BusinessCacheService) toCachedRestaurante(r *model.Restaurante) *CachedRestaurante {
	cached := &CachedRestaurante{
		ID:        r.ID,
		Nome:      r.Nome,
		TaxaFrete: r.TaxaFrete,
		Ativo:     r.Ativo,
		Aberto:    r.Aberto,
		CozinhaID: r.CozinhaID,
	}

	if r.Cozinha.ID > 0 {
		cached.Cozinha = &CachedCozinha{
			ID:   r.Cozinha.ID,
			Nome: r.Cozinha.Nome,
		}
	}

	if r.Endereco.CidadeID > 0 {
		cached.EnderecoCidadeID = r.Endereco.CidadeID
	}

	for _, fp := range r.FormasPagamento {
		cached.FormasPagamentoIDs = append(cached.FormasPagamentoIDs, fp.ID)
	}

	for _, resp := range r.Responsaveis {
		cached.ResponsaveisIDs = append(cached.ResponsaveisIDs, resp.ID)
	}

	return cached
}

// ToModel converte CachedRestaurante para model.Restaurante (parcial)
func (c *CachedRestaurante) ToModel() *model.Restaurante {
	r := &model.Restaurante{
		ID:        c.ID,
		Nome:      c.Nome,
		TaxaFrete: c.TaxaFrete,
		Ativo:     c.Ativo,
		Aberto:    c.Aberto,
		CozinhaID: c.CozinhaID,
	}

	if c.Cozinha != nil {
		r.Cozinha = model.Cozinha{
			ID:   c.Cozinha.ID,
			Nome: c.Cozinha.Nome,
		}
	}

	return r
}

// ==================== FORMAS DE PAGAMENTO ====================

// GetFormaPagamento obtém uma forma de pagamento do cache
func (s *BusinessCacheService) GetFormaPagamento(id uint64) (*model.FormaPagamento, error) {
	ctx, cancel := context.WithTimeout(context.Background(), businessCacheTimeout)
	defer cancel()

	key := fmt.Sprintf("%s%d", formaPagamentoCachePrefix, id)
	data, err := s.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var cached CachedFormaPagamento
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil, err
	}

	return &model.FormaPagamento{ID: cached.ID, Descricao: cached.Descricao}, nil
}

// SetFormaPagamento armazena uma forma de pagamento no cache
func (s *BusinessCacheService) SetFormaPagamento(fp *model.FormaPagamento) error {
	ctx, cancel := context.WithTimeout(context.Background(), businessCacheTimeout)
	defer cancel()

	cached := CachedFormaPagamento{ID: fp.ID, Descricao: fp.Descricao}
	data, err := json.Marshal(cached)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s%d", formaPagamentoCachePrefix, fp.ID)
	return s.redisClient.Set(ctx, key, data, formaPagamentoCacheTTL).Err()
}

// GetAllFormasPagamento obtém todas as formas de pagamento do cache
func (s *BusinessCacheService) GetAllFormasPagamento() ([]model.FormaPagamento, error) {
	ctx, cancel := context.WithTimeout(context.Background(), businessCacheTimeout)
	defer cancel()

	data, err := s.redisClient.Get(ctx, formaPagamentoAllCacheKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var cached []CachedFormaPagamento
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil, err
	}

	fps := make([]model.FormaPagamento, len(cached))
	for i, c := range cached {
		fps[i] = model.FormaPagamento{ID: c.ID, Descricao: c.Descricao}
	}

	return fps, nil
}

// SetAllFormasPagamento armazena todas as formas de pagamento no cache
func (s *BusinessCacheService) SetAllFormasPagamento(fps []model.FormaPagamento) error {
	ctx, cancel := context.WithTimeout(context.Background(), businessCacheTimeout)
	defer cancel()

	cached := make([]CachedFormaPagamento, len(fps))
	for i, fp := range fps {
		cached[i] = CachedFormaPagamento{ID: fp.ID, Descricao: fp.Descricao}
	}

	data, err := json.Marshal(cached)
	if err != nil {
		return err
	}

	return s.redisClient.Set(ctx, formaPagamentoAllCacheKey, data, formaPagamentoCacheTTL).Err()
}

// InvalidateFormaPagamento invalida o cache de uma forma de pagamento
func (s *BusinessCacheService) InvalidateFormaPagamento(id uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), businessCacheTimeout)
	defer cancel()

	pipe := s.redisClient.Pipeline()
	pipe.Del(ctx, fmt.Sprintf("%s%d", formaPagamentoCachePrefix, id))
	pipe.Del(ctx, formaPagamentoAllCacheKey)
	_, err := pipe.Exec(ctx)
	return err
}

// Ping verifica a conexão com o Redis
func (s *BusinessCacheService) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), businessCacheTimeout)
	defer cancel()
	return s.redisClient.Ping(ctx).Err()
}

// Close fecha a conexão com o Redis
func (s *BusinessCacheService) Close() error {
	return s.redisClient.Close()
}

// WarmUpCache pré-carrega cozinhas no cache
func (s *BusinessCacheService) WarmUpCozinhas(cozinhas []model.Cozinha) {
	if err := s.SetAllCozinhas(cozinhas); err != nil {
		log.Printf("Aviso: Falha ao pré-carregar cozinhas no cache: %v", err)
	} else {
		log.Printf("Cache aquecido com %d cozinhas", len(cozinhas))
	}
}
