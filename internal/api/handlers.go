package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"news_aggregator/internal/service"
	"strconv"
)

// Handler представляет обработчик HTTP-запросов
type Handler struct {
	newsService service.NewsProvider
}

// NewHandler создает новый экземпляр обработчика
func NewHandler(newsService *service.NewsService) *Handler {
	return &Handler{
		newsService: newsService,
	}
}

// HealthCheck проверяет работоспособность сервера
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "OK",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Ошибка отправки ответа: %v\n", err)
	}
}

// GetNews возвращает последние новости
func (h *Handler) GetNews(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	n, err := strconv.Atoi(vars["n"])
	if err != nil {
		log.Printf("Ошибка парсинга количества новостей: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	posts, err := h.newsService.GetLatestPosts(n)
	if err != nil {
		log.Printf("Ошибка получения новостей: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		log.Printf("Ошибка сериализации ответа: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
