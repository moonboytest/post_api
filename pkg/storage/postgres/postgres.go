package postgres

import (
	"GoNews/pkg/config"
	"GoNews/pkg/storage"
	"context"
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
		// Форматируем время в нужный формат
		p.CreatedAtFormatted = p.CreatedAt.Format("06-01-02 15-04")
		p.PublishedAtFormatted = p.PublishedAt.Format("06-01-02 15-04")
		posts = append(posts, p)
	}
	return posts, rows.Err()

}

//Добавляем пост в бд
func (s *Storage) AddPost(p storage.Post) error {
	//Начало транзакции
	tx, err := s.db.Begin(context.Background())
	if err != nil {
		return nil
	}

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

//Обновляем данные в посте. Изменяем название и текст, время публикации
//изменяется автоматически
func (s *Storage) UpdatePost(p storage.Post) error {
	//Начало транзакции
	tx, err := s.db.Begin(context.Background())
	if err != nil {
		return nil
	}

	time := time.Now()
	_, err = s.db.Exec(context.Background(), `
	UPDATE posts
	SET title = $2, 
	content = $3, 
	published_at = $4
	WHERE id = $1
	`,
		p.ID, p.Title, p.Content, time)
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

//Удаление поста
func (s *Storage) DeletePost(p storage.Post) error {
	_, err := s.db.Exec(context.Background(), `
	DELETE FROM posts
	WHERE id = $1 `,
		p.ID)
	if err != nil {
		return err
	}
	return nil
}


