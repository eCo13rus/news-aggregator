package models

import "testing"

func TestRSSItem_ParsePubDate(t *testing.T) {
	item := &RSSItem{PubDate: "Mon, 02 Jan 2006 15:04:05 GMT"}
	_, err := item.ParsePubDate()
	if err != nil {
		t.Errorf("Ошибка парсинга корректной даты: %v", err)
	}
}
