package postgres

import (
	"GoNews/pkg/config"
	"GoNews/pkg/storage"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

// Конструктор
func New() (*Storage, error) {

	var ctx context.Context = context.Background()

	// Подключение к БД. Функция возвращает объект БД.
	// Строка для подключения к БД записана
	// В файле config.go
	// Подробности читай в readme
	postgresString := config.PostgresString
	db, err := pgxpool.Connect(ctx,
		postgresString)
	fmt.Println("все гуд")
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)

	}

	//читаем файл схемы бд из корня проекта
	//Путь к файлу схемы указан в файле Config.go
	//Подробности читать в readme
	schema, err := ioutil.ReadFile(config.SchemaSQLPath)
	if err != nil {
		log.Fatal("Ошибка при чтении файла schema.sql:", err)
	}

	s := Storage{
		db: db,
	}
	//выполняем запросы из файла sсhema
	_, err = s.db.Exec(context.Background(), string(schema))
	if err != nil {
		log.Fatal("Ошибка при выполнении SQL-запросов:", err)
	}

	return &s, nil

}

//Получаем все посты из бд
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

//Добавляем пост в бд
func (s *Storage) AddPost(p storage.Post) error {

	var author_name string
	//Забираем имя автора из связаной таблицы authors
	rows, err := s.db.Query(context.Background(), `
	SELECT name
	FROM authors 
	WHERE id = $1`,
		p.AuthorID)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&author_name)
		if err != nil {
			return err
		}
	}

	//Фиксируем текущее время
	time := time.Now()

	//Записываем информацию о посте в бд
	_, err = s.db.Exec(context.Background(), `
	INSERT INTO posts(author_id, author_name, title, content, created_at, published_at)
	VALUES ($1,$2, $3, $4, $5, $6)`,
		p.AuthorID, author_name, p.Title, p.Content, time, time)

	if err != nil {
		return err
	}

	return nil
}

//Обновляем данные в посте. Изменяем название и текст, время публикации
//изменяется автоматически
func (s *Storage) UpdatePost(post storage.Post) error {
	//Начало транзакции
	tx, err := s.db.Begin(context.Background())
	if err != nil {
		return nil
	}

	time := time.Now()
	fmt.Print(1)
	_, err = s.db.Exec(context.Background(), `
	UPDATE posts
	SET title = $2, 
	content = $3, 
	published_at = $4
	WHERE id = $1
	`,
		post.ID, post.Title, post.Content, time)
	fmt.Print(2)
	if err != nil {
		//В случае ошибки откатываемся назад
		tx.Rollback(context.Background())
		return err
	}
	//Фиксация изменений
	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil
}
func (s *Storage) DeletePost(storage.Post) error {
	return nil
}

var Post1 = storage.Post{
	ID:       1,
	AuthorID: 1,
	Title:    "E",
	Content:  "2Go is a new language. Although it borrows ideas from existing languages, it has unusual properties that make effective Go programs different in character from programs written in its relatives. A straightforward translation of a C++ or Java program into Go is unlikely to produce a satisfactory result—Java programs are written in Java, not Go. On the other hand, thinking about the problem from a Go perspective could produce a successful but quite different program. In other words, to write Go well, it's important to understand its properties and idioms. It's also important to know the established conventions for programming in Go, such as naming, formatting, program construction, and so on, so that programs you write will be easy for other Go programmers to understand.",
}

var UpdPost1 = storage.Post{
	ID:       1,
	AuthorID: 1,
	Title:    "123",
	Content:  "321",
}
