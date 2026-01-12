package data

import (
	"context"
	"database/sql"
	"strconv"
	"time"
)

const dbTimeout = time.Second * 3

var db *sql.DB

func New(dbPool *sql.DB) Models {
	db = dbPool
	return Models{
		Favorite: Favorite{},
	}
}

type Models struct {
	Favorite Favorite
}

// Favorite represents a user's favorite post
type Favorite struct {
	UserID    int       `json:"userId"`
	PostID    int       `json:"postId"`
	CreatedAt time.Time `json:"createdAt"`
}

// GetPostIDsByUser returns all favorite post IDs for a user as strings
func (f *Favorite) GetPostIDsByUser(userID int) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT post_id FROM favorites WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var postIDs []string
	for rows.Next() {
		var postID int
		if err := rows.Scan(&postID); err != nil {
			return nil, err
		}
		postIDs = append(postIDs, strconv.Itoa(postID))
	}

	return postIDs, nil
}

// Insert adds a favorite
func (f *Favorite) Insert(userID, postID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `INSERT INTO favorites (user_id, post_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`

	_, err := db.ExecContext(ctx, stmt, userID, postID)
	return err
}

// Delete removes a favorite
func (f *Favorite) Delete(userID, postID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `DELETE FROM favorites WHERE user_id = $1 AND post_id = $2`

	_, err := db.ExecContext(ctx, stmt, userID, postID)
	return err
}

// Exists checks if a favorite exists
func (f *Favorite) Exists(userID, postID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT EXISTS(SELECT 1 FROM favorites WHERE user_id = $1 AND post_id = $2)`

	var exists bool
	err := db.QueryRowContext(ctx, query, userID, postID).Scan(&exists)
	return exists, err
}

// BulkInsert inserts multiple favorites (for sync)
func (f *Favorite) BulkInsert(userID int, postIDs []int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `INSERT INTO favorites (user_id, post_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`

	for _, postID := range postIDs {
		_, err := db.ExecContext(ctx, stmt, userID, postID)
		if err != nil {
			return err
		}
	}

	return nil
}
