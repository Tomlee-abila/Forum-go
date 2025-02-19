package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"learn.zone01kisumu.ke/git/clomollo/forum/internal/database"
)

type Session struct {
	ID        string
	UserID    int
	ExpiresAt time.Time
	Forum *database.ForumModel
}

func (s *Session) CreateSession(userID int) (string, error) {
	querry := `DELETE FROM Sessions WHERE user_id=?`
	_, err := s.Forum.DB.Exec(querry, userID)
	if err != nil {
		return "", fmt.Errorf("failed to delete existing sessions: %w", err)
	}
	SessionID := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)

	query := "INSERT INTO Sessions(id, user_id, expires_at) VALUES(?, ?, ?)"
	_, err = s.Forum.DB.Exec(query, SessionID, userID, expiresAt)
	if err != nil {
		return "", fmt.Errorf("failed to insert session: %w", err)
	}

	return SessionID, nil
}

func (s *Session)GetSession(sessionID string) (*Session, error) {
	query := `SELECT id, user_id, expires_at FROM Sessions WHERE id=?`
	row := s.Forum.DB.QueryRow(query, sessionID)
	var session Session

	err := row.Scan(&session.ID, &session.UserID, &session.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	return &session, nil
}

func (s *Session)DeleteSession(sessionID string) error {
	query := "DELETE FROM Sessions WHERE id=?"
	_, err := s.Forum.DB.Exec(query, sessionID)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}
