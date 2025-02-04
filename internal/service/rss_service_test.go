package service

import (
	"news_aggregator/internal/models"
	"testing"
	"time"
)

func TestRSSService_ProcessFeed(t *testing.T) {
	rssService := NewRSSService([]string{"https://test.com"})
	if rssService == nil {
		t.Error("Сервис не создан")
	}
}

func TestRSSService_StartFeedProcessing(t *testing.T) {
	service := NewRSSService([]string{"https://test.com"})
	postsChan := make(chan *models.Post, 1)

	service.StartFeedProcessing(time.Second, postsChan)
	time.Sleep(time.Millisecond * 100)

	close(postsChan)
}

func TestRSSService_FetchFeed(t *testing.T) {
	service := NewRSSService([]string{"https://habr.com/ru/rss/all/all/"})
	_, err := service.FetchFeed("https://habr.com/ru/rss/all/all/")
	if err != nil {
		t.Skipf("Пропуск теста: %v", err)
	}
}
