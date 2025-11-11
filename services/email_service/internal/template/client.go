package template

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/brainox/hng-group55-distributed-notification-system/services/email_service/internal/models"
	"github.com/brainox/hng-group55-distributed-notification-system/services/email_service/pkg/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	redis      *redis.Client
	cacheTTL   time.Duration
}

func NewClient(baseURL string, redis *redis.Client) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		redis:    redis,
		cacheTTL: 10 * time.Minute,
	}
}

// FetchTemplate fetches a template from the Template Service
func (c *Client) FetchTemplate(ctx context.Context, templateKey string) (*models.EmailTemplate, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("template:%s", templateKey)
	cached, err := c.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var template models.EmailTemplate
		if err := json.Unmarshal([]byte(cached), &template); err == nil {
			logger.Log.Info("template cache hit", zap.String("template_key", templateKey))
			return &template, nil
		}
	}

	// Fetch from Template Service
	url := fmt.Sprintf("%s/api/v1/templates/key/%s?language=en&version=latest", c.baseURL, templateKey)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch template: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("template service returned status %d: %s", resp.StatusCode, string(body))
	}

	var response models.TemplateResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("template not found: %s", templateKey)
	}

	template := &models.EmailTemplate{
		Subject:   response.Data.Version.Subject,
		Body:      response.Data.Version.Body,
		Variables: response.Data.Version.Variables,
	}

	// Cache the template
	templateJSON, _ := json.Marshal(template)
	_ = c.redis.Set(ctx, cacheKey, templateJSON, c.cacheTTL).Err()

	logger.Log.Info("template fetched from service", zap.String("template_key", templateKey))
	return template, nil
}
