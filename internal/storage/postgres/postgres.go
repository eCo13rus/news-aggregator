package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"news_aggregator/internal/models"
	"strconv"
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

func (s *Storage) GetPosts(limit int, page int, searchQuery string) ([]*models.Post, int, error) {
	offset := (page - 1) * limit

	countQuery := `
		SELECT COUNT(*) FROM posts
		WHERE 1=1
	`

	dataQuery := `
		SELECT id, title, content, pub_time, link, created_at, updated_at
		FROM posts
		WHERE 1=1
	`

	var params []interface{}
	var searchCondition string

	if searchQuery != "" {
		searchCondition = " AND title ILIKE $1"
		params = append(params, "%"+searchQuery+"%")
	}

	var totalCount int
	err := s.db.QueryRow(countQuery+searchCondition, params...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка подсчета количества публикаций: %v", err)
	}

	dataQuery += searchCondition + " ORDER BY pub_time DESC LIMIT $" + strconv.Itoa(len(params)+1) + " OFFSET $" + strconv.Itoa(len(params)+2)
	params = append(params, limit, offset)

	rows, err := s.db.Query(dataQuery, params...)
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка получения публикаций: %v", err)
	}
	defer rows.Close()

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
			return nil, 0, fmt.Errorf("ошибка сканирования строки: %v", err)
		}
		posts = append(posts, post)
	}

	return posts, totalCount, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) GetPostByID(id int) (*models.Post, error) {
	query := `
        SELECT id, title, content, pub_time, link, created_at, updated_at
        FROM posts
        WHERE id = $1
    `

	post := &models.Post{}
	err := s.db.QueryRow(query, id).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.PubTime,
		&post.Link,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("новость с ID %d не найдена", id)
		}
		return nil, fmt.Errorf("ошибка получения новости: %v", err)
	}

	return post, nil
}
