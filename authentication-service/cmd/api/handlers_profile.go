package main

import (
	"errors"
	"log"
	"net/http"
	"time"
)

func (app *Config) Profile(w http.ResponseWriter, r *http.Request) {
	var userID int64
	var err error

	atCookie, err := r.Cookie("access_token")
	if atCookie != nil {
		log.Printf("access token exists")
		userID, err = app.TokenService.ValidateAccessToken(atCookie.Value)
	}

	if err != nil {
		log.Printf("access token invalid, trying refresh token")
		rtCookie, rtErr := r.Cookie("refresh_token")
		if rtErr != nil {
			log.Printf("refresh token also missing")
			app.errorJSON(w, errors.New("access token and refresh token missing"), http.StatusUnauthorized)
			return
		}

		newAT, refreshErr := app.TokenService.Refresh(rtCookie.Value, time.Now().Add(accessTokenTime))
		if refreshErr != nil {
			log.Printf("refresh token invalid")
			app.errorJSON(w, errors.New("refresh token invalid"), http.StatusUnauthorized)
			return
		}

		log.Printf("set auth cookies")
		app.setAccessTokenCookie(w, newAT, time.Now().Add(accessTokenTime))

		userID, err = app.TokenService.ValidateAccessToken(newAT)
		if err != nil {
			log.Printf("new access token invalid")
			app.errorJSON(w, err, http.StatusUnauthorized)
			return
		}
	}

	log.Printf("fetching user profile for userID %d", userID)
	user, err := app.Models.User.GetOne(int(userID))
	if err != nil {
		app.errorJSON(w, errors.New("user not found"), http.StatusNotFound)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Profile fetched successfully",
		Data: map[string]any{
			"id":         user.ID,
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
		},
	}

	app.writeJSON(w, http.StatusOK, payload)
}
