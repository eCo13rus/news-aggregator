package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"news_aggregator/internal/models"
	"news_aggregator/internal/service"
	"strconv"
	"strings"
)

type Handler struct {
	newsService service.NewsProvider
}

func NewHandler(newsService *service.NewsService) *Handler {
	return &Handler{
		newsService: newsService,
	}
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	requestID := "system"
	if id := r.Context().Value(RequestIDKey); id != nil {
		requestID = id.(string)
	}
	response := map[string]string{
		"status": "OK",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("[%s] Ошибка отправки ответа: %v\n", requestID, err)
	}
}

func (h *Handler) GetNews(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	n, err := strconv.Atoi(vars["n"])
	if err != nil {
		log.Printf("Ошибка парсинга количества новостей: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	searchQuery := r.URL.Query().Get("s")
	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	requestID := r.Context().Value(RequestIDKey).(string)

	log.Printf("[%s] Запрос новостей: лимит=%d, страница=%d, поиск='%s'",
		requestID, n, page, searchQuery)

	posts, totalCount, err := h.newsService.GetLatestPostsWithPagination(n, page, searchQuery)
	if err != nil {
		log.Printf("[%s] Ошибка получения новостей: %v", requestID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := (totalCount + n - 1) / n

	response := models.NewsResponse{
		News: posts,
		Pagination: models.Pagination{
			CurrentPage:  page,
			TotalPages:   totalPages,
			ItemsPerPage: n,
			TotalItems:   totalCount,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("[%s] Ошибка сериализации ответа: %v", requestID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetNewsDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	newsID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Некорректный ID новости", http.StatusBadRequest)
		return
	}

	requestID := r.Context().Value(RequestIDKey).(string)

	post, err := h.newsService.GetPostByID(newsID)
	if err != nil {
		log.Printf("[%s] Ошибка получения новости: %v", requestID, err)

		if strings.Contains(err.Error(), "не найдена") {
			http.Error(w, fmt.Sprintf("Новость с ID %d не найдена", newsID), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if post == nil {
		http.Error(w, fmt.Sprintf("Новость с ID %d не найдена", newsID), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}
