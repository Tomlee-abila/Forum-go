package database

import "database/sql"

// function to initialize the database for the forum project
func initializeDatabase(databaseName string) (*sql.DB, error) {
	dataBase, err := sql.Open("sqlite3", databaseName)
	if err != nil {
		return nil, err
	}

	// create the database queries tables
	query := `
	CREATE TABLE IF NOT EXIST users(
		id TEXT PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		useremail TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		image_path TEXT
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS posts(
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		category TEXT NOT NULL, 
		title TEXT NOT NULL,
		content TEXT NOT NULL, //discription of the post
		media BLOB, //binary data - video, image and GIFs
		content_type TEXT, //content type tracking(text, image, gif)
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

	CREATE TABLE IF NOT EXISTS likes(
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		post_id TEXT,
		type TEXT CHECK(type IN ('like', 'dislike')),
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (post_id) REFERENCES posts(id)
	);

	`

	if _, err := dataBase.Exec(query); err != nil {
		dataBase.Close()
		return nil, err
	}
	return dataBase, nil
}
