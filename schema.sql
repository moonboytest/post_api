DROP TABLE IF EXISTS posts, authors;

CREATE TABLE authors (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    author_id INTEGER REFERENCES authors(id) NOT NULL,
    author_name TEXT not NULL,
    title TEXT  NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    published_at TIMESTAMP NOT NULL DEFAULT NOW()

);


INSERT INTO authors (id, name) VALUES (1, 'Дмитрий');
INSERT INTO authors (id, name) VALUES (2, 'михаил');
INSERT INTO authors (id, name) VALUES (3, 'roman');
/*
INSERT INTO posts (id, author_id, title, content, author_name) VALUES (0, 0, 'Статья', 'Содержание статьи', 'IDI NAHUY'); */