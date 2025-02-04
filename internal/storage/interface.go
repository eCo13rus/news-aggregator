package storage

import (
	"news_aggregator/internal/models"
)

// Storage определяет интерфейс для работы с хранилищем данных
type Storage interface {
	AddPost(post *models.Post) error
	GetPosts(limit int) ([]*models.Post, error)
}
