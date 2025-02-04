package service

import (
	"news_aggregator/internal/models"
	"testing"
	"time"
)

type mockStorage struct {
	posts []*models.Post
}

func (m *mockStorage) AddPost(post *models.Post) error {
	m.posts = append(m.posts, post)
	return nil
}

func (m *mockStorage) GetPosts(limit int) ([]*models.Post, error) {
	if limit > len(m.posts) {
		limit = len(m.posts)
	}
	return m.posts[:limit], nil
}

func TestNewsService_GetLatestPosts(t *testing.T) {
	storage := &mockStorage{
		posts: []*models.Post{
			{
				ID:      1,
				Title:   "Test Post 1",
				Content: "Content 1",
				PubTime: time.Now().Unix(),
			},
			{
				ID:      2,
				Title:   "Test Post 2",
				Content: "Content 2",
				PubTime: time.Now().Unix() - 3600,
			},
		},
	}

	service := NewNewsService(storage, []string{"https://test.com"})
	posts, err := service.GetLatestPosts(1)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(posts) != 1 {
		t.Errorf("Expected 1 post, got %d", len(posts))
	}

	if posts[0].Title != "Test Post 1" {
		t.Errorf("Expected title 'Test Post 1', got %s", posts[0].Title)
	}
}
