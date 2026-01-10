package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Authenticate handles user login authentication
func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start authentication")
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Verify user exists
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		log.Printf("User not found by email: %v", err)
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	// Validate password
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		log.Printf("Invalid password attempt for user %s: %v", requestPayload.Email, err)
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	// Log the authentication
	go func() {
		err := app.logRequest("authentication", fmt.Sprintf("%s logged in", user.Email))
		if err != nil {
			log.Printf("Failed to send log to logger-service during authentication: %v", err)
		}
	}()

	// Generate token pair
	atExp := time.Now().Add(accessTokenTime)
	rtExp := time.Now().Add(refreshTokenTime)

	tokenPair, err := app.TokenService.GenerateTokenPair(
		int64(user.ID),
		atExp,
		rtExp,
	)
	if err != nil {
		app.errorJSON(w, errors.New("failed to generate token"), http.StatusInternalServerError)
		return
	}

	app.setAccessTokenCookie(w, tokenPair.AccessToken, atExp)
	app.setRefreshTokenCookie(w, tokenPair.RefreshToken, rtExp)

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data: userData{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		},
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}
