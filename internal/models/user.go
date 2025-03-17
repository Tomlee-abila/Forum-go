package models

import (
	"database/sql"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
	// "learn.zone01kisumu.ke/git/clomollo/forum/internal/database"
)

type User struct {
	ID        int
	UserID    string
	Email     string
	Username  string
	Password  string
	// ImagePath string // This correctly stores 'image_path' from the database
}


func (f *ForumModel) CreateUser(userUuid, email, username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	querystatement := "INSERT INTO users(user_uuid,email, username, password) VALUES(?,?,?,?)"
	_, err = f.DB.Exec(querystatement, userUuid, email, username, string(hashedPassword))
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	return nil
}

func (f *ForumModel) GetUserByEmail(email string) (*User, error) {
	query := "SELECT id, user_uuid,email, username, password FROM users WHERE email = ?"
	row := f.DB.QueryRow(query, email)
	user := User{}
	err := row.Scan(&user.ID, &user.UserID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (f *ForumModel) GetUserByID(uuid string) (*User, error) {
	query :="SELECT id, user_uuid, email, username, password FROM users WHERE user_uuid=?"


	row := f.DB.QueryRow(query, uuid)
	user := User{}
	// var imagePath sql.NullString
	err := row.Scan(&user.ID, &user.UserID, &user.Email, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("Failed to get user by ID: %v", err)
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	   // Assign imagePath only if it's not NULL
	   
    // Convert NULL to an empty string
    // if imagePath.Valid {
    //     user.ImagePath = imagePath.String
    // } else {
    //     user.ImagePath = "" // If NULL, set empty string
    // }
	return &user, nil
}

func (f *ForumModel) GetUserByUsername(username string) (*User, error) {
	query := "SELECT id, email, username, password FROM users WHERE username = ?"
	row := f.DB.QueryRow(query, username)
	user := User{}
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("Failed to get user by username: %v", err)
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	return &user, nil
}
