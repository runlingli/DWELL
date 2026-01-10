package main

import (
	"authentication/data"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Register handles new user registration
func (app *Config) Register(w http.ResponseWriter, r *http.Request) {
	log.Printf("Register service begin")

	var requestPayload struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		FirstName        string `json:"first_name"`
		LastName         string `json:"last_name"`
		VerificationCode string `json:"verification_code"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, errors.New("error reading user information"), http.StatusBadRequest)
		return
	}

	// Check if email already exists
	existingUser, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err == nil && existingUser != nil {
		app.errorJSON(w, errors.New("email already exists"), http.StatusConflict)
		return
	} else if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Verify the verification code
	isValid, err := app.MailService.VerifyCode(requestPayload.Email, requestPayload.VerificationCode)
	if err != nil || !isValid {
		app.errorJSON(w, errors.New("invalid verification code"), http.StatusBadRequest)
		return
	}

	// Create new user
	user := data.User{
		Email:     requestPayload.Email,
		FirstName: requestPayload.FirstName,
		LastName:  requestPayload.LastName,
		Password:  requestPayload.Password,
		Active:    1,
	}

	log.Printf("Inserting user: %+v", user)

	userID, err := app.Models.User.Insert(user)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	user.ID = userID

	// Log the registration
	go func() {
		err := app.logRequest("registration attempt", fmt.Sprintf("%s registered", user.Email))
		if err != nil {
			log.Printf("Failed to send log to logger-service during registration: %v", err)
		}
	}()

	// Generate token pair for auto-login after registration
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
		Message: fmt.Sprintf("%s registered successfully!", user.Email),
		Data: map[string]any{
			"id":         user.ID,
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
		},
	}

	app.writeJSON(w, http.StatusCreated, payload)
}

// VerifyEmail sends a verification code to the user's email
func (app *Config) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	log.Printf("Verify email service begin")

	var requestPayload struct {
		Email string `json:"email"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		log.Printf("Error reading email information: %v", err)
		app.errorJSON(w, errors.New("error reading email information"), http.StatusBadRequest)
		return
	}

	// Check if email already exists
	existingUser, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err == nil && existingUser != nil {
		log.Printf("Email already exists: %s", requestPayload.Email)
		app.errorJSON(w, errors.New("email already exists"), http.StatusConflict)
		return
	} else if err != nil && err != sql.ErrNoRows {
		log.Printf("Server error when checking existing email: %v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Send verification code
	log.Printf("Sending verification code to email: %s", requestPayload.Email)
	err = app.MailService.SendCode(requestPayload.Email, "DWELL Verification Code")
	if err != nil {
		log.Printf("Error sending verification code: %v", err)
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("%s sent successfully!", requestPayload.Email),
	}
	app.writeJSON(w, http.StatusCreated, payload)
}
