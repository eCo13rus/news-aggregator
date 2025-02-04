package models

// Post представляет структуру новостной публикации
type Post struct {
	ID        int    `json:"id" db:"id"`
	Title     string `json:"title" db:"title"`
	Content   string `json:"content" db:"content"`
	PubTime   int64  `json:"pub_time" db:"pub_time"`
	Link      string `json:"link" db:"link"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
}
