package service

import (
	"news_aggregator/internal/models"
	"time"
)

type NewsProvider interface {
	GetLatestPosts(limit int) ([]*models.Post, error)
	GetLatestPostsWithPagination(limit, page int, searchQuery string) ([]*models.Post, int, error)
	GetPostByID(id int) (*models.Post, error)
	Start(period time.Duration) error
	Stop()
}
