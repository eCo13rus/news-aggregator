package models

// Config представляет структуру конфигурации приложения
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	RSS      []string       `json:"rss"`
	Period   int            `json:"request_period"`
}

// ServerConfig конфигурация сервера
type ServerConfig struct {
	Port string `json:"port"`
	Host string `json:"host"`
}

// DatabaseConfig конфигурация базы данных
type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
	SSLMode  string `json:"ssl_mode"`
}
