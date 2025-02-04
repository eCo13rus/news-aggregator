package service

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"news_aggregator/internal/models"
	"time"
)

// RSSService представляет сервис для работы с RSS
type RSSService struct {
	feedURLs []string
	client   *http.Client
}

// NewRSSService создает новый экземпляр RSS сервиса
func NewRSSService(urls []string) *RSSService {
	return &RSSService{
		feedURLs: urls,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// FetchFeed получает и парсит RSS-фид по указанному URL
func (s *RSSService) FetchFeed(url string) (*models.RSS, error) {
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения RSS фида: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("некорректный статус ответа: %d", resp.StatusCode)
	}

	var feed models.RSS
	if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return nil, fmt.Errorf("ошибка парсинга RSS: %v", err)
	}

	return &feed, nil
}

// StartFeedProcessing запускает параллельную обработку RSS-фидов
func (s *RSSService) StartFeedProcessing(period time.Duration, postsChan chan<- *models.Post) {
	for _, url := range s.feedURLs {
		go s.processFeed(url, period, postsChan)
	}
}

// processFeed обрабатывает отдельный RSS-фид
func (s *RSSService) processFeed(url string, period time.Duration, postsChan chan<- *models.Post) {
	ticker := time.NewTicker(period)
	defer ticker.Stop()

	for ; true; <-ticker.C {
		feed, err := s.FetchFeed(url)
		if err != nil {
			fmt.Printf("Ошибка получения фида %s: %v\n", url, err)
			continue
		}

		for _, item := range feed.Channel.Items {
			pubTime, err := item.ParsePubDate()
			if err != nil {
				fmt.Printf("Ошибка парсинга даты публикации: %v\n", err)
				continue
			}

			post := &models.Post{
				Title:   item.Title,
				Content: item.Description,
				PubTime: pubTime,
				Link:    item.Link,
			}

			postsChan <- post
		}
	}
}
