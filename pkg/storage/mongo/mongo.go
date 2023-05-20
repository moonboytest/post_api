package mongo

import (
	"GoNews/pkg/config"
	"GoNews/pkg/storage"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Хранилище данных.
type MongoDB struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// Конструктор объекта хранилища.
func New() (*MongoDB, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(config.MongoDBString))
	if err != nil {
		return nil, err
	}
	db := client.Database(config.MongoDBName)
	collection := db.Collection(config.MongoCollectionName)

	return &MongoDB{
		client:     client,
		collection: collection,
	}, nil
}

//Получаем все посты из бд
func (db *MongoDB) Posts() ([]storage.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()
	cur, err := db.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	var posts []storage.Post
	if err := cur.All(ctx, &posts); err != nil {
		return nil, err
	}

	for i := range posts {
		posts[i].CreatedAtFormatted = posts[i].CreatedAt.Format("06-01-02 15-04")
		posts[i].PublishedAtFormatted = posts[i].PublishedAt.Format("06-01-02 15-04")
	}

	return posts, nil
}

//Добавляем пост в бд
func (db *MongoDB) AddPost(post storage.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	post.CreatedAt = time.Now()
	post.PublishedAt = time.Now()

	_, err := db.collection.InsertOne(ctx, post)
	if err != nil {
		return err
	}

	return nil
}

//Обновлеяем пост - title, content, author id и author name. Время создания поста
//остается неизменным. Время публикации изменяется автоматически
func (db *MongoDB) UpdatePost(post storage.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": post.ID}
	update := bson.M{
		"$set": bson.M{
			"Title":       post.Title,
			"Content":     post.Content,
			"AuthorID":    post.AuthorID,
			"AuthorName":  post.AuthorName,
			"PublishedAt": time.Now(),
		},
	}

	_, err := db.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

//удаления поста из бд по его id
func (db *MongoDB) DeletePost(post storage.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": post.ID}

	_, err := db.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
