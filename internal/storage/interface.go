package storage

import (
	"news_aggregator/internal/models"
)

type Storage interface {
	AddPost(post *models.Post) error
	GetPosts(limit int, page int, searchQuery string) ([]*models.Post, int, error)
	GetPostByID(id int) (*models.Post, error)
}
