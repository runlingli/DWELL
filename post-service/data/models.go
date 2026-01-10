package data

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/lib/pq"
)

const dbTimeout = time.Second * 3

var db *sql.DB

func New(dbPool *sql.DB) Models {
	db = dbPool
	return Models{
		Post: Post{},
	}
}

type Models struct {
	Post Post
}

// Post represents a rental listing
type Post struct {
	ID               int       `json:"id"`
	Title            string    `json:"title"`
	Price            float64   `json:"price"`
	Location         string    `json:"location,omitempty"`
	Neighborhood     string    `json:"neighborhood"`
	Lat              float64   `json:"lat"`
	Lng              float64   `json:"lng"`
	Radius           int       `json:"radius"`
	Type             string    `json:"type"`
	ImageURL         string    `json:"imageUrl"`
	AdditionalImages []string  `json:"additionalImages,omitempty"`
	Description      string    `json:"description"`
	Bedrooms         int       `json:"bedrooms"`
	Bathrooms        int       `json:"bathrooms"`
	AvailableFrom    time.Time `json:"availableFrom"`
	AvailableTo      time.Time `json:"availableTo"`
	AuthorID         int       `json:"authorId"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

// PostWithAuthor includes author information for API responses
type PostWithAuthor struct {
	Post
	Author Author `json:"author"`
}

// Author represents the post author's public info
type Author struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar,omitempty"`
}

// GetAll returns all posts with author info, ordered by creation date
func (p *Post) GetAll() ([]*PostWithAuthor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT
			p.id, p.title, p.price, COALESCE(p.location, ''), p.neighborhood,
			p.lat, p.lng, p.radius, p.type, COALESCE(p.image_url, ''),
			COALESCE(p.additional_images, '{}'), COALESCE(p.description, ''),
			p.bedrooms, p.bathrooms, p.available_from, p.available_to,
			p.author_id, p.created_at, p.updated_at,
			u.first_name, u.last_name
		FROM posts p
		JOIN users u ON p.author_id = u.id
		ORDER BY p.created_at DESC
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*PostWithAuthor

	for rows.Next() {
		var post PostWithAuthor
		var firstName, lastName string

		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Price,
			&post.Location,
			&post.Neighborhood,
			&post.Lat,
			&post.Lng,
			&post.Radius,
			&post.Type,
			&post.ImageURL,
			pq.Array(&post.AdditionalImages),
			&post.Description,
			&post.Bedrooms,
			&post.Bathrooms,
			&post.AvailableFrom,
			&post.AvailableTo,
			&post.AuthorID,
			&post.CreatedAt,
			&post.UpdatedAt,
			&firstName,
			&lastName,
		)
		if err != nil {
			log.Println("Error scanning post:", err)
			return nil, err
		}

		post.Author = Author{
			Name: firstName + " " + lastName,
		}

		posts = append(posts, &post)
	}

	return posts, nil
}

// GetByID returns a single post by ID with author info
func (p *Post) GetByID(id int) (*PostWithAuthor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT
			p.id, p.title, p.price, COALESCE(p.location, ''), p.neighborhood,
			p.lat, p.lng, p.radius, p.type, COALESCE(p.image_url, ''),
			COALESCE(p.additional_images, '{}'), COALESCE(p.description, ''),
			p.bedrooms, p.bathrooms, p.available_from, p.available_to,
			p.author_id, p.created_at, p.updated_at,
			u.first_name, u.last_name
		FROM posts p
		JOIN users u ON p.author_id = u.id
		WHERE p.id = $1
	`

	var post PostWithAuthor
	var firstName, lastName string

	row := db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&post.ID,
		&post.Title,
		&post.Price,
		&post.Location,
		&post.Neighborhood,
		&post.Lat,
		&post.Lng,
		&post.Radius,
		&post.Type,
		&post.ImageURL,
		pq.Array(&post.AdditionalImages),
		&post.Description,
		&post.Bedrooms,
		&post.Bathrooms,
		&post.AvailableFrom,
		&post.AvailableTo,
		&post.AuthorID,
		&post.CreatedAt,
		&post.UpdatedAt,
		&firstName,
		&lastName,
	)

	if err != nil {
		return nil, err
	}

	post.Author = Author{
		Name: firstName + " " + lastName,
	}

	return &post, nil
}

// GetByAuthorID returns all posts by a specific author
func (p *Post) GetByAuthorID(authorID int) ([]*PostWithAuthor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT
			p.id, p.title, p.price, COALESCE(p.location, ''), p.neighborhood,
			p.lat, p.lng, p.radius, p.type, COALESCE(p.image_url, ''),
			COALESCE(p.additional_images, '{}'), COALESCE(p.description, ''),
			p.bedrooms, p.bathrooms, p.available_from, p.available_to,
			p.author_id, p.created_at, p.updated_at,
			u.first_name, u.last_name
		FROM posts p
		JOIN users u ON p.author_id = u.id
		WHERE p.author_id = $1
		ORDER BY p.created_at DESC
	`

	rows, err := db.QueryContext(ctx, query, authorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*PostWithAuthor

	for rows.Next() {
		var post PostWithAuthor
		var firstName, lastName string

		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Price,
			&post.Location,
			&post.Neighborhood,
			&post.Lat,
			&post.Lng,
			&post.Radius,
			&post.Type,
			&post.ImageURL,
			pq.Array(&post.AdditionalImages),
			&post.Description,
			&post.Bedrooms,
			&post.Bathrooms,
			&post.AvailableFrom,
			&post.AvailableTo,
			&post.AuthorID,
			&post.CreatedAt,
			&post.UpdatedAt,
			&firstName,
			&lastName,
		)
		if err != nil {
			log.Println("Error scanning post:", err)
			return nil, err
		}

		post.Author = Author{
			Name: firstName + " " + lastName,
		}

		posts = append(posts, &post)
	}

	return posts, nil
}

// Insert creates a new post and returns its ID
func (p *Post) Insert(post Post) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `
		INSERT INTO posts (
			title, price, location, neighborhood, lat, lng, radius, type,
			image_url, additional_images, description, bedrooms, bathrooms,
			available_from, available_to, author_id, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
		RETURNING id
	`

	var newID int
	err := db.QueryRowContext(ctx, stmt,
		post.Title,
		post.Price,
		post.Location,
		post.Neighborhood,
		post.Lat,
		post.Lng,
		post.Radius,
		post.Type,
		post.ImageURL,
		pq.Array(post.AdditionalImages),
		post.Description,
		post.Bedrooms,
		post.Bathrooms,
		post.AvailableFrom,
		post.AvailableTo,
		post.AuthorID,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

// Update updates an existing post
func (p *Post) Update(post Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `
		UPDATE posts SET
			title = $1,
			price = $2,
			location = $3,
			neighborhood = $4,
			lat = $5,
			lng = $6,
			radius = $7,
			type = $8,
			image_url = $9,
			additional_images = $10,
			description = $11,
			bedrooms = $12,
			bathrooms = $13,
			available_from = $14,
			available_to = $15,
			updated_at = $16
		WHERE id = $17 AND author_id = $18
	`

	result, err := db.ExecContext(ctx, stmt,
		post.Title,
		post.Price,
		post.Location,
		post.Neighborhood,
		post.Lat,
		post.Lng,
		post.Radius,
		post.Type,
		post.ImageURL,
		pq.Array(post.AdditionalImages),
		post.Description,
		post.Bedrooms,
		post.Bathrooms,
		post.AvailableFrom,
		post.AvailableTo,
		time.Now(),
		post.ID,
		post.AuthorID,
	)

	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Delete removes a post by ID (only if author matches)
func (p *Post) Delete(id, authorID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `DELETE FROM posts WHERE id = $1 AND author_id = $2`

	result, err := db.ExecContext(ctx, stmt, id, authorID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
