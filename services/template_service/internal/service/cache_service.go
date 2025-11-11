package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/brainox/hng-group55-distributed-notification-system/services/template_service/internal/models"
	"github.com/redis/go-redis/v9"
)

type CacheService interface {
	GetTemplate(ctx context.Context, key string, language string, version string) (*models.TemplateResponse, error)
	SetTemplate(ctx context.Context, key string, language string, version string, template *models.TemplateResponse, ttl time.Duration) error
	InvalidateTemplate(ctx context.Context, key string) error
}

type cacheService struct {
	redis *redis.Client
}

func NewCacheService(redisClient *redis.Client) CacheService {
	return &cacheService{redis: redisClient}
}

func (s *cacheService) GetTemplate(ctx context.Context, key string, language string, version string) (*models.TemplateResponse, error) {
	cacheKey := fmt.Sprintf("template:%s:%s:%s", key, language, version)

	data, err := s.redis.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	}
	if err != nil {
		return nil, err
	}

	var template models.TemplateResponse
	if err := json.Unmarshal([]byte(data), &template); err != nil {
		return nil, err
	}

	return &template, nil
}

func (s *cacheService) SetTemplate(ctx context.Context, key string, language string, version string, template *models.TemplateResponse, ttl time.Duration) error {
	cacheKey := fmt.Sprintf("template:%s:%s:%s", key, language, version)

	data, err := json.Marshal(template)
	if err != nil {
		return err
	}

	return s.redis.Set(ctx, cacheKey, data, ttl).Err()
}

func (s *cacheService) InvalidateTemplate(ctx context.Context, key string) error {
	// Use pattern matching to delete all versions of this template
	pattern := fmt.Sprintf("template:%s:*", key)

	iter := s.redis.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := s.redis.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}

	return iter.Err()
}
