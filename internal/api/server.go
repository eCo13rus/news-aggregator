package api

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Server struct {
	router  *mux.Router
	handler *Handler
}

func NewServer(handler *Handler) *Server {
	return &Server{
		router:  mux.NewRouter(),
		handler: handler,
	}
}

func (s *Server) SetupRoutes() {
	s.router.Use(RequestIDMiddleware)
	s.router.Use(LoggingMiddleware)

	s.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	s.router.HandleFunc("/api/news/{n}", s.handler.GetNews).Methods(http.MethodGet, http.MethodOptions)
	s.router.HandleFunc("/api/news/detail/{id:[0-9]+}", s.handler.GetNewsDetail).Methods(http.MethodGet)
}

func (s *Server) Start(addr string) error {
	log.Printf("Сервер запущен на %s\n", addr)
	return http.ListenAndServe(addr, s.router)
}
