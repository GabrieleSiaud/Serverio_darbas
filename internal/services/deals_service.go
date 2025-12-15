package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"serverio_darbas/internal/generated/repository"
)

type DealsService struct {
	db     *repository.Queries
	client *http.Client
	ttl    time.Duration
}

func NewDealsService(db *repository.Queries) *DealsService {
	return &DealsService{
		db: db,
		client: &http.Client{
			Timeout: 8 * time.Second,
		},
		ttl: 10 * time.Minute, // cache 10 min
	}
}

func (s *DealsService) GetDeals(ctx context.Context, title string) (any, error) {
	cacheKey := "cheapshark:deals:title=" + title

	// 1) cache hit
	cached, err := s.db.GetExternalCache(ctx, cacheKey)
	if err == nil {
		var out any
		if err := json.Unmarshal(cached.ResponseJson, &out); err == nil {
			return out, nil
		}
	}

	// 2) call external API
	u := "https://www.cheapshark.com/api/1.0/deals?title=" + url.QueryEscape(title) + "&pageSize=10"
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("external api error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("external api status %d: %s", resp.StatusCode, string(b))
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var parsed any
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, err
	}

	// 3) save cache
	expires := time.Now().Add(s.ttl)
	_ = s.db.UpsertExternalCache(ctx, repository.UpsertExternalCacheParams{
		CacheKey:     cacheKey,
		ResponseJson: raw,
		ExpiresAt:    expires,
	})

	return parsed, nil
}
