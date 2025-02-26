package models

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type ForumModel struct {
	DB *sql.DB
}

func InitializeDB() (*sql.DB, error) {
	dataBase, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return nil, err
	}

	// defer dataBase.Close()

	// create the database queries tabCategory: category,les
	query := `
	CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_uuid TEXT UNIQUE NOT NULL,
        email TEXT UNIQUE NOT NULL,
        username TEXT UNIQUE NOT NULL,
        password TEXT NOT NULL,
        image_path TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		
	CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		user_uuid TEXT NOT NULL,
		expires_at DATETIME NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_uuid) REFERENCES users(id)
    );
	
	CREATE TABLE IF NOT EXISTS posts(
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id TEXT UNIQUE NOT NULL,
		user_uuid TEXT NOT NULL,
		title TEXT NOT NULL,
		content TEXT NOT NULL, 
		media BLOB, --binary data - video, image and GIFs
		content_type TEXT, --content type tracking(text, image, gif)
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_uuid) REFERENCES users(user_uuid)
		);
	
	CREATE TABLE IF NOT EXISTS comments(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTERGER,
		user_id INTEGER NOT NULL,
		parent_id TEXT,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(user_uuid) ON DELETE CASCADE,
		FOREIGN KEY (parent_id) REFERENCES comments(id) ON DELETE CASCADE
	);

		
	CREATE TABLE IF NOT EXISTS post_likes(
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		post_id TEXT,
		type TEXT CHECK(type IN ('like', 'dislike')),
		FOREIGN KEY (user_id) REFERENCES users(user_uuid),
		FOREIGN KEY (post_id) REFERENCES posts(id)
	);

	CREATE TABLE IF NOT EXISTS comment_likes(
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		comment_id TEXT NOT NULL,
		type TEXT CHECK(type IN('like', 'dislike')),
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (comment_id) REFERENCES comments(id)
	);
		
	CREATE TABLE IF NOT EXISTS categories (
		id TEXT PRIMARY KEY,
		name TEXT UNIQUE NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
		
		CREATE TABLE IF NOT EXISTS post_categories(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id TEXT NOT NULL,
		category_id TEXT,
		FOREIGN KEY (post_id) REFERENCES posts(id),
		FOREIGN KEY (category_id) REFERENCES categories(id)
		);
		
		INSERT OR IGNORE INTO categories (id, name) VALUES
		('tech','tech'),
		('lifestyle', 'lifestyle'),
		('gaming', 'gaming'),
		('food', 'food'),
		('travel', 'travel'),
		('other', 'other');

	
	`

	if _, err := dataBase.Exec(query); err != nil {
		dataBase.Close()
		return nil, err
	}

	return dataBase, nil
}

/*
	INSERT into categories(category_value) VALUES
		('education'),
		('politics'),
		('sports'),
		('lifestyle'),
		('religion'),
		('relationship and family'),
		('Health'),
		('Real-estate'),
		('Governance'),
		('technology'),
    	('gaming'),
    	('food'),
    	('travel'),
    	('other');
*/
