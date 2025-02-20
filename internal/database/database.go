package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitializeDB() error {
	dataBase, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return err
	}

	// defer dataBase.Close()

	// create the database queries tables
	query := `
	CREATE TABLE IF NOT EXISTS users(
		id TEXT PRIMARY KEY,
		useremail TEXT UNIQUE NOT NULL,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS posts(
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		category TEXT NOT NULL, 
		title TEXT NOT NULL,
		content TEXT NOT NULL, --discription of the post
		media BLOB, --binary data - video, image and GIFs
		content_type TEXT, --content type tracking(text, image, gif)
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (category) REFERENCES post_categories(category)
	);

	CREATE TABLE IF NOT EXISTS post_categories(
		post_id TEXT NOT NULL,
		category_id TEXT NOT NULL,
		FOREIGN KEY (post_id) REFERENCES posts(id),
		PRIMARY KEY (post_id, category_id)
	);

	CREATE TABLE IF NOT EXISTS comments(
		id TEXT PRIMARY KEY,
		post_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (post_id) REFERENCES posts(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	);

	CREATE TABLE IF NOT EXISTS post_likes(
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		post_id TEXT,
		type TEXT CHECK(type IN ('like', 'dislike')),
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (post_id) REFERENCES posts(id)
	);
	CREATE TABLE IF NOT EXISTS comment_likes(
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		comment_id TEXT NOT NULL,
		type TEXT CHECK(type IN('like', 'dislike')),
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (comment_id) REFERENCES comments(id)

	);`

	if _, err := dataBase.Exec(query); err != nil {
		dataBase.Close()
		return err
	}

	DB = dataBase
	return nil
}
