package models

import (
	"fmt"

	"github.com/google/uuid"
)

// processLike handles the database operations for likes
func (f *ForumModel) ProcessLike(itemType, itemID, userID, likeType string) error {
	var tableName, idColumn string
	// Set table and column names based on item type
	if itemType == "post" {
		tableName = "post_likes"
		idColumn = "post_id"
	} else {
		tableName = "comment_likes"
		idColumn = "comment_id"
	}
	// First verify the item exists
	exists, err := f.checkItemExists(itemType, itemID)
	if err != nil {
		return fmt.Errorf("error checking item existence: %v", err)
	}
	if !exists {
		return fmt.Errorf("item not found")
	}
	// Check for existing like
	var existingType string
	query := fmt.Sprintf("SELECT type FROM %s WHERE user_id = ? AND %s = ?", tableName, idColumn)
	err = f.DB.QueryRow(query, userID, itemID).Scan(&existingType)
	if err == nil {
		// Update existing like
		updateQuery := fmt.Sprintf("UPDATE %s SET type = ? WHERE user_id = ? AND %s = ?", tableName, idColumn)
		_, err = f.DB.Exec(updateQuery, likeType, userID, itemID)
	} else {
		// Insert new like
		insertQuery := fmt.Sprintf("INSERT INTO %s (id, user_id, %s, type) VALUES (?, ?, ?, ?)", tableName, idColumn)
		_, err = f.DB.Exec(insertQuery, uuid.New().String(), userID, itemID, likeType)
	}
	return err
}

// checkItemExists verifies if the post or comment exists
func (f *ForumModel) checkItemExists(itemType, itemID string) (bool, error) {
	var exists bool
	var query string
	if itemType == "post" {
		query = "SELECT EXISTS(SELECT 1 FROM posts WHERE post_id = ?)"
	} else {
		query = "SELECT EXISTS(SELECT 1 FROM comments WHERE id = ?)"
	}
	err := f.DB.QueryRow(query, itemID).Scan(&exists)
	return exists, err
}
func (f *ForumModel) GetReactionCounts(itemType, itemID string) (map[string]int, error) {
    counts := make(map[string]int)
    var tableName, idColumn string

    switch itemType {
    case "post":
        tableName = "post_likes"
        idColumn = "post_id"
    case "comment":
        tableName = "comment_likes"
        idColumn = "comment_id"
    default:
        return nil, fmt.Errorf("invalid item type")
    }

    // Declare separate variables for scanning
    var likes, dislikes int

    query := fmt.Sprintf(`
        SELECT 
            SUM(CASE WHEN type = 'like' THEN 1 ELSE 0 END) as likes,
            SUM(CASE WHEN type = 'dislike' THEN 1 ELSE 0 END) as dislikes
        FROM %s 
        WHERE %s = ?`, tableName, idColumn)

    err := f.DB.QueryRow(query, itemID).Scan(&likes, &dislikes)
    if err != nil {
        return nil, fmt.Errorf("error getting reaction counts: %v", err)
    }

    // Assign the scanned values to the map
    counts["likes"] = likes
    counts["dislikes"] = dislikes

    return counts, nil
}