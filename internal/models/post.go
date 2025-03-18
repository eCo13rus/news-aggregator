package models

type Post struct {
	ID        int    `json:"id" db:"id"`
	Title     string `json:"title" db:"title"`
	Content   string `json:"content" db:"content"`
	PubTime   int64  `json:"pub_time" db:"pub_time"`
	Link      string `json:"link" db:"link"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
}

type Pagination struct {
	CurrentPage  int `json:"current_page"`
	TotalPages   int `json:"total_pages"`
	ItemsPerPage int `json:"items_per_page"`
	TotalItems   int `json:"total_items"`
}

type NewsResponse struct {
	News       []*Post    `json:"news"`
	Pagination Pagination `json:"pagination"`
}
