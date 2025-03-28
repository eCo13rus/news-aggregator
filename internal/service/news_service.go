package service

import (
	"fmt"
	"log"
	"news_aggregator/internal/models"
	"news_aggregator/internal/storage"
	"sync"
	"time"
)

// NewsService представляет сервис для работы с новостями
type NewsService struct {
	storage    storage.Storage
	rssService *RSSService
	postsChan  chan *models.Post
	wg         sync.WaitGroup
}

// NewNewsService создает новый экземпляр сервиса новостей
func NewNewsService(storage storage.Storage, rssUrls []string) *NewsService {
	return &NewsService{
		storage:    storage,
		rssService: NewRSSService(rssUrls),
		postsChan:  make(chan *models.Post, 100),
	}
}

// Start запускает процесс агрегации новостей
func (s *NewsService) Start(period time.Duration) error {
	s.rssService.StartFeedProcessing(period, s.postsChan)

	s.wg.Add(1)
	go s.handlePosts()

	log.Println("Сервис новостей запущен")
	return nil
}

// handlePosts обрабатывает полученные новости
func (s *NewsService) handlePosts() {
	defer s.wg.Done()

	for post := range s.postsChan {
		if err := s.storage.AddPost(post); err != nil {
			log.Printf("Ошибка сохранения новости: %v\n", err)
			continue
		}
		log.Printf("Сохранена новость: %s\n", post.Title)
	}
}

func (s *NewsService) GetLatestPostsWithPagination(limit, page int, searchQuery string) ([]*models.Post, int, error) {
	posts, totalCount, err := s.storage.GetPosts(limit, page, searchQuery)
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка получения новостей: %v", err)
	}
	return posts, totalCount, nil
}

func (s *NewsService) Stop() {
	close(s.postsChan)
	s.wg.Wait()
	log.Println("Сервис новостей остановлен")
}

func (s *NewsService) GetLatestPosts(limit int) ([]*models.Post, error) {
	posts, _, err := s.GetLatestPostsWithPagination(limit, 1, "")
	if err != nil {
		return nil, fmt.Errorf("ошибка получения новостей: %v", err)
	}
	return posts, nil
}

func (s *NewsService) GetPostByID(id int) (*models.Post, error) {
	post, err := s.storage.GetPostByID(id)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения новости: %v", err)
	}
	return post, nil
}
