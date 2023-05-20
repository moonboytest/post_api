package storage

import "time"

// Post - публикация.
type Post struct {
	ID                   int       `json:"ID"`
	Title                string    `json:"Title"`
	Content              string    `json:"Content"`
	AuthorID             int       `json:"Author_ID"`
	AuthorName           string    `json:"Author_Name"`
	CreatedAt            time.Time `json:"-"`
	PublishedAt          time.Time `json:"-"`
	CreatedAtFormatted   string    `json:"created_at,omitempty"`
	PublishedAtFormatted string    `json:"published_at,omitempty"`
}

// Interface задаёт контракт на работу с БД.
type Interface interface {
	Posts() ([]Post, error) // получение всех публикаций
	AddPost(Post) error     // создание новой публикации
	UpdatePost(Post) error  // обновление публикации
	DeletePost(Post) error  // удаление публикации по ID
}
