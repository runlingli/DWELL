package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

// GetUserFavoriteIDs returns all favorite post IDs for a user
func (app *Config) GetUserFavoriteIDs(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userId")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		app.errorJSON(w, errors.New("invalid user ID"), http.StatusBadRequest)
		return
	}

	log.Printf("Getting favorite IDs for user: %d", userID)

	postIDs, err := app.Models.Favorite.GetPostIDsByUser(userID)
	if err != nil {
		log.Printf("Error getting favorites: %v", err)
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if postIDs == nil {
		postIDs = []string{}
	}

	payload := jsonResponse{
		Error: false,
		Data:  postIDs,
	}

	app.writeJSON(w, http.StatusOK, payload)
}

// AddFavorite adds a post to user's favorites
func (app *Config) AddFavorite(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		UserID int `json:"userId"`
		PostID int `json:"postId"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	log.Printf("Adding favorite - User: %d, Post: %d", requestPayload.UserID, requestPayload.PostID)

	err = app.Models.Favorite.Insert(requestPayload.UserID, requestPayload.PostID)
	if err != nil {
		log.Printf("Error adding favorite: %v", err)
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Favorite added successfully",
	}

	app.writeJSON(w, http.StatusCreated, payload)
}

// RemoveFavorite removes a post from user's favorites
func (app *Config) RemoveFavorite(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userId")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		app.errorJSON(w, errors.New("invalid user ID"), http.StatusBadRequest)
		return
	}

	postIDStr := chi.URLParam(r, "postId")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		app.errorJSON(w, errors.New("invalid post ID"), http.StatusBadRequest)
		return
	}

	log.Printf("Removing favorite - User: %d, Post: %d", userID, postID)

	err = app.Models.Favorite.Delete(userID, postID)
	if err != nil {
		log.Printf("Error removing favorite: %v", err)
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Favorite removed successfully",
	}

	app.writeJSON(w, http.StatusOK, payload)
}

// SyncFavorites syncs localStorage favorites to database
func (app *Config) SyncFavorites(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		UserID  int   `json:"userId"`
		PostIDs []int `json:"postIds"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	log.Printf("Syncing favorites - User: %d, Posts: %v", requestPayload.UserID, requestPayload.PostIDs)

	err = app.Models.Favorite.BulkInsert(requestPayload.UserID, requestPayload.PostIDs)
	if err != nil {
		log.Printf("Error syncing favorites: %v", err)
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Return updated list
	postIDs, err := app.Models.Favorite.GetPostIDsByUser(requestPayload.UserID)
	if err != nil {
		log.Printf("Error getting favorites after sync: %v", err)
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if postIDs == nil {
		postIDs = []string{}
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Favorites synced successfully",
		Data:    postIDs,
	}

	app.writeJSON(w, http.StatusOK, payload)
}
