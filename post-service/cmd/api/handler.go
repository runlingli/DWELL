package main

import (
	"errors"
	"log"
	"net/http"
	"post-service/data"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

// PostPayload represents the incoming post data from clients
type PostPayload struct {
	ID               int      `json:"id,omitempty"`
	Title            string   `json:"title"`
	Price            float64  `json:"price"`
	Location         string   `json:"location,omitempty"`
	Neighborhood     string   `json:"neighborhood"`
	Lat              float64  `json:"lat"`
	Lng              float64  `json:"lng"`
	Radius           int      `json:"radius"`
	Type             string   `json:"type"`
	ImageURL         string   `json:"imageUrl"`
	AdditionalImages []string `json:"additionalImages,omitempty"`
	Description      string   `json:"description"`
	Bedrooms         int      `json:"bedrooms"`
	Bathrooms        int      `json:"bathrooms"`
	AvailableFrom    int64    `json:"availableFrom"` // Unix timestamp from frontend
	AvailableTo      int64    `json:"availableTo"`   // Unix timestamp from frontend
	AuthorID         int      `json:"authorId"`
}

// GetAllPosts returns all posts
func (app *Config) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	log.Println("========== GetAllPosts START ==========")

	posts, err := app.Models.Post.GetAll()
	if err != nil {
		log.Printf("ERROR getting posts from database: %v", err)
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	log.Printf("Found %d posts in database", len(posts))

	// Convert to frontend format
	response := make([]map[string]any, 0, len(posts))
	for _, post := range posts {
		response = append(response, convertPostToFrontend(post))
	}

	payload := jsonResponse{
		Error: false,
		Data:  response,
	}

	log.Printf("Returning %d posts to client", len(response))
	app.writeJSON(w, http.StatusOK, payload)
}

// GetPostByID returns a single post
func (app *Config) GetPostByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.errorJSON(w, errors.New("invalid post ID"), http.StatusBadRequest)
		return
	}

	log.Printf("Getting post by ID: %d", id)

	post, err := app.Models.Post.GetByID(id)
	if err != nil {
		log.Printf("Error getting post: %v", err)
		app.errorJSON(w, errors.New("post not found"), http.StatusNotFound)
		return
	}

	payload := jsonResponse{
		Error: false,
		Data:  convertPostToFrontend(post),
	}

	app.writeJSON(w, http.StatusOK, payload)
}

// GetPostsByAuthor returns all posts by a specific author
func (app *Config) GetPostsByAuthor(w http.ResponseWriter, r *http.Request) {
	authorIDStr := chi.URLParam(r, "authorId")
	authorID, err := strconv.Atoi(authorIDStr)
	if err != nil {
		app.errorJSON(w, errors.New("invalid author ID"), http.StatusBadRequest)
		return
	}

	log.Printf("Getting posts by author ID: %d", authorID)

	posts, err := app.Models.Post.GetByAuthorID(authorID)
	if err != nil {
		log.Printf("Error getting posts: %v", err)
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	response := make([]map[string]any, 0, len(posts))
	for _, post := range posts {
		response = append(response, convertPostToFrontend(post))
	}

	payload := jsonResponse{
		Error: false,
		Data:  response,
	}

	app.writeJSON(w, http.StatusOK, payload)
}

// CreatePost creates a new post
func (app *Config) CreatePost(w http.ResponseWriter, r *http.Request) {
	log.Println("========== CreatePost START ==========")

	var requestPayload PostPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		log.Printf("ERROR reading post payload: %v", err)
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	log.Printf("Received payload - Title: %s, AuthorID: %d, Price: %.2f", requestPayload.Title, requestPayload.AuthorID, requestPayload.Price)
	log.Printf("Full payload: %+v", requestPayload)

	// Validate authorId
	if requestPayload.AuthorID == 0 {
		log.Println("ERROR: AuthorID is 0 - this will fail foreign key constraint!")
		app.errorJSON(w, errors.New("authorId is required and must be a valid user ID"), http.StatusBadRequest)
		return
	}

	// Convert timestamps
	availableFrom := time.Unix(requestPayload.AvailableFrom/1000, 0)
	availableTo := time.Unix(requestPayload.AvailableTo/1000, 0)
	log.Printf("Converted timestamps - From: %v, To: %v", availableFrom, availableTo)

	post := data.Post{
		Title:            requestPayload.Title,
		Price:            requestPayload.Price,
		Location:         requestPayload.Location,
		Neighborhood:     requestPayload.Neighborhood,
		Lat:              requestPayload.Lat,
		Lng:              requestPayload.Lng,
		Radius:           requestPayload.Radius,
		Type:             requestPayload.Type,
		ImageURL:         requestPayload.ImageURL,
		AdditionalImages: requestPayload.AdditionalImages,
		Description:      requestPayload.Description,
		Bedrooms:         requestPayload.Bedrooms,
		Bathrooms:        requestPayload.Bathrooms,
		AvailableFrom:    availableFrom,
		AvailableTo:      availableTo,
		AuthorID:         requestPayload.AuthorID,
	}

	log.Println("Inserting post into database...")
	newID, err := app.Models.Post.Insert(post)
	if err != nil {
		log.Printf("ERROR inserting post into database: %v", err)
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	log.Printf("SUCCESS: Post created with ID: %d", newID)

	// Fetch the created post with author info
	createdPost, err := app.Models.Post.GetByID(newID)
	if err != nil {
		log.Printf("Error fetching created post: %v", err)
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Post created successfully",
		Data:    convertPostToFrontend(createdPost),
	}

	app.writeJSON(w, http.StatusCreated, payload)
}

// UpdatePost updates an existing post
func (app *Config) UpdatePost(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.errorJSON(w, errors.New("invalid post ID"), http.StatusBadRequest)
		return
	}

	var requestPayload PostPayload

	err = app.readJSON(w, r, &requestPayload)
	if err != nil {
		log.Printf("Error reading post payload: %v", err)
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	log.Printf("Updating post %d: %+v", id, requestPayload)

	// Convert timestamps
	availableFrom := time.Unix(requestPayload.AvailableFrom/1000, 0)
	availableTo := time.Unix(requestPayload.AvailableTo/1000, 0)

	post := data.Post{
		ID:               id,
		Title:            requestPayload.Title,
		Price:            requestPayload.Price,
		Location:         requestPayload.Location,
		Neighborhood:     requestPayload.Neighborhood,
		Lat:              requestPayload.Lat,
		Lng:              requestPayload.Lng,
		Radius:           requestPayload.Radius,
		Type:             requestPayload.Type,
		ImageURL:         requestPayload.ImageURL,
		AdditionalImages: requestPayload.AdditionalImages,
		Description:      requestPayload.Description,
		Bedrooms:         requestPayload.Bedrooms,
		Bathrooms:        requestPayload.Bathrooms,
		AvailableFrom:    availableFrom,
		AvailableTo:      availableTo,
		AuthorID:         requestPayload.AuthorID,
	}

	err = app.Models.Post.Update(post)
	if err != nil {
		log.Printf("Error updating post: %v", err)
		app.errorJSON(w, errors.New("post not found or unauthorized"), http.StatusNotFound)
		return
	}

	// Fetch updated post
	updatedPost, err := app.Models.Post.GetByID(id)
	if err != nil {
		log.Printf("Error fetching updated post: %v", err)
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Post updated successfully",
		Data:    convertPostToFrontend(updatedPost),
	}

	app.writeJSON(w, http.StatusOK, payload)
}

// DeletePost deletes a post
func (app *Config) DeletePost(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.errorJSON(w, errors.New("invalid post ID"), http.StatusBadRequest)
		return
	}

	// Get author ID from request body
	var requestPayload struct {
		AuthorID int `json:"authorId"`
	}

	err = app.readJSON(w, r, &requestPayload)
	if err != nil {
		log.Printf("Error reading delete payload: %v", err)
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	log.Printf("Deleting post %d by author %d", id, requestPayload.AuthorID)

	err = app.Models.Post.Delete(id, requestPayload.AuthorID)
	if err != nil {
		log.Printf("Error deleting post: %v", err)
		app.errorJSON(w, errors.New("post not found or unauthorized"), http.StatusNotFound)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Post deleted successfully",
	}

	app.writeJSON(w, http.StatusOK, payload)
}

// convertPostToFrontend converts a PostWithAuthor to frontend format
func convertPostToFrontend(post *data.PostWithAuthor) map[string]any {
	return map[string]any{
		"id":               strconv.Itoa(post.ID),
		"title":            post.Title,
		"price":            post.Price,
		"location":         post.Location,
		"neighborhood":     post.Neighborhood,
		"coordinates":      map[string]float64{"lat": post.Lat, "lng": post.Lng},
		"radius":           post.Radius,
		"type":             post.Type,
		"imageUrl":         post.ImageURL,
		"additionalImages": post.AdditionalImages,
		"description":      post.Description,
		"bedrooms":         post.Bedrooms,
		"bathrooms":        post.Bathrooms,
		"createdAt":        post.CreatedAt.UnixMilli(),
		"availableFrom":    post.AvailableFrom.UnixMilli(),
		"availableTo":      post.AvailableTo.UnixMilli(),
		"author": map[string]any{
			"name":   post.Author.Name,
			"avatar": post.Author.Avatar,
		},
	}
}
