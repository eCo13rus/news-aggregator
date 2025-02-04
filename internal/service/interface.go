package service

import (
	"news_aggregator/internal/models"
	"time"
)

type NewsProvider interface {
	GetLatestPosts(limit int) ([]*models.Post, error)
	Start(period time.Duration) error
	Stop()
}
