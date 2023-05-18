package postgres

import (
	"GoNews/pkg/config"
	"GoNews/pkg/storage"
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

// Конструктор
func New() (*Storage, error) {

	var ctx context.Context = context.Background()

	// Подключение к БД. Функция возвращает объект БД.
	pwd := config.Pwd
	db, err := pgxpool.Connect(ctx,
		"host=localhost port=5432 dbname=catalog user=postgres password="+
			pwd+
			"sslmode=prefer connect_timeout=10")
	fmt.Println("все гуд")
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)

	}

	//читаем файл схемы бд из корня проекта
	schema, err := ioutil.ReadFile("schema.sql")
	if err != nil {
		log.Fatal("Ошибка при чтении файла schema.sql:", err)
	}

	s := Storage{
		db: db,
	}

	_, err = s.db.Exec(context.Background(), string(schema))
	if err != nil {
		log.Fatal("Ошибка при выполнении SQL-запросов:", err)
	}

	return &s, nil

}

//выполняем запросы из файла sсhema
func (s *Storage) Posts() ([]storage.Post, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT 
			id,
			author_id,
			author_name,
			title,
			content,
			created_at,
			published_at
		FROM posts
		ORDER BY id;
		`,
	)
	if err != nil {
		return nil, err
	}

	var posts []storage.Post

	for rows.Next() {
		var p storage.Post
		err := rows.Scan(
			&p.ID,
			&p.AuthorID,
			&p.AuthorName,
			&p.Title,
			&p.Content,
			&p.CreatedAt,
			&p.PublishedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()

}

func (s *Storage) AddPost(storage.Post) error {
	return nil
}
func (s *Storage) UpdatePost(storage.Post) error {
	return nil
}
func (s *Storage) DeletePost(storage.Post) error {
	return nil
}

var posts = []storage.Post{
	{
		ID:      1,
		Title:   "Effective Go123123",
		Content: "Go is a new language. Although it borrows ideas from existing languages, it has unusual properties that make effective Go programs different in character from programs written in its relatives. A straightforward translation of a C++ or Java program into Go is unlikely to produce a satisfactory result—Java programs are written in Java, not Go. On the other hand, thinking about the problem from a Go perspective could produce a successful but quite different program. In other words, to write Go well, it's important to understand its properties and idioms. It's also important to know the established conventions for programming in Go, such as naming, formatting, program construction, and so on, so that programs you write will be easy for other Go programmers to understand.",
	},
	{
		ID:      2,
		Title:   "The Go Memory Model",
		Content: "The Go memory model specifies the conditions under which reads of a variable in one goroutine can be guaranteed to observe values produced by writes to the same variable in a different goroutine.",
	},
}
