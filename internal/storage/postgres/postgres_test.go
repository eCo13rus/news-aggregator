package postgres

import (
	"news_aggregator/internal/models"
	"testing"
	"time"
)

func TestAddPost(t *testing.T) {
	cfg := &models.DatabaseConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "eco",
		Password: "Horror4199.",
		Name:     "news_aggregator",
		SSLMode:  "disable",
	}

	storage, err := New(cfg)
	if err != nil {
		t.Skip("Пропуск теста: нет подключения к БД")
	}
	defer storage.Close()

	post := &models.Post{
		Title:   "Тест",
		Content: "Контент",
		PubTime: 123456789,
		Link:    "http://test.com",
	}

	err = storage.AddPost(post)
	if err != nil {
		t.Errorf("Ошибка добавления поста: %v", err)
	}
}

func TestStorage_GetPosts(t *testing.T) {
	cfg := &models.DatabaseConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "eco",
		Password: "Horror4199.",
		Name:     "news_aggregator",
		SSLMode:  "disable",
	}

	storage, err := New(cfg)
	if err != nil {
		t.Skip("Пропуск теста: нет подключения к БД")
	}
	defer storage.Close()

	testPost := &models.Post{
		Title:   "Тест GetPosts",
		Content: "Контент теста GetPosts",
		PubTime: time.Now().Unix(),
		Link:    "http://test-get.com",
	}

	err = storage.AddPost(testPost)
	if err != nil {
		t.Errorf("Ошибка добавления поста: %v", err)
	}

	posts, err := storage.GetPosts(1)
	if err != nil {
		t.Errorf("Ошибка получения постов: %v", err)
	}

	if len(posts) == 0 {
		t.Error("Ожидался хотя бы один пост")
	}
}
