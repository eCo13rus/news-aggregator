package models

import (
	"encoding/xml"
	"time"
)

// RSS представляет структуру RSS-фида
type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

// Channel представляет канал в RSS-фиде
type Channel struct {
	Title       string    `xml:"title"`
	Description string    `xml:"description"`
	Link        string    `xml:"link"`
	Items       []RSSItem `xml:"item"`
}

// RSSItem представляет отдельную публикацию в RSS-фиде
type RSSItem struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Link        string `xml:"link"`
	GUID        string `xml:"guid"`
}

// ParsePubDate преобразует строку даты публикации в Unix timestamp
func (item *RSSItem) ParsePubDate() (int64, error) {
	// Поддерживаемые форматы дат в RSS
	layouts := []string{
		time.RFC1123,
		time.RFC1123Z,
		time.RFC822,
		time.RFC822Z,
		"Mon, 02 Jan 2006 15:04:05 -0700",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, item.PubDate); err == nil {
			return t.Unix(), nil
		}
	}
	return time.Now().Unix(), nil
}
