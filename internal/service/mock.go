package service

import (
	"news_aggregator/internal/models"
	"time"
)

type MockNewsService struct{}

func (m *MockNewsService) GetLatestPosts(limit int) ([]*models.Post, error) {
	return []*models.Post{{Title: "Test"}}, nil
}

func (m *MockNewsService) Start(period time.Duration) error {
	return nil
}

func (m *MockNewsService) Stop() {}
