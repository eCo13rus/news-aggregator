package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"news_aggregator/internal/models"
)

// Storage реализует интерфейс хранилища для PostgreSQL
type Storage struct {
	db *sql.DB
}

// New создает новый экземпляр хранилища PostgreSQL
func New(cfg *models.DatabaseConfig) (*Storage, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
	)

	// Открываем соединение с БД
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка проверки соединения с БД: %v", err)
	}

	log.Println("Успешное подключение к БД")

	return &Storage{db: db}, nil
}

// AddPost реализует добавление публикации
func (s *Storage) AddPost(post *models.Post) error {
	query := `
        INSERT INTO posts (title, content, pub_time, link)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (link) DO NOTHING
    `

	_, err := s.db.Exec(query, post.Title, post.Content, post.PubTime, post.Link)
	if err != nil {
		return fmt.Errorf("ошибка добавления публикации: %v", err)
	}

	return nil
}

// GetPosts реализует получение последних публикаций
func (s *Storage) GetPosts(limit int) ([]*models.Post, error) {
	query := `
        SELECT id, title, content, pub_time, link, created_at, updated_at
        FROM posts
        ORDER BY pub_time DESC
        LIMIT $1
    `

	rows, err := s.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения публикаций: %v", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.PubTime,
			&post.Link,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %v", err)
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// Close закрывает соединение с БД
func (s *Storage) Close() error {
	return s.db.Close()
}
