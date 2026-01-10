package main

import (
	"authentication/data"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type UserInfo struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (app *Config) GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Starting Google OAuth Login Handler")
	url, err := app.OAuthService.GoogleLoginURL()
	if err != nil {
		log.Printf("Error generating Google login URL: %v", err)
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	log.Printf("Redirecting to Google OAuth URL: %s", url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (app *Config) GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	code := r.FormValue("code")

	validState, err := app.OAuthService.ValidateGoogleState(state)
	if err != nil || !validState {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}

	userInfo, err := app.OAuthService.GetUserInfo(code)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get user from google: %s", err), http.StatusInternalServerError)
		return
	}

	var userID int
	user, err := app.Models.User.GetByGoogleID(userInfo.ID)
	if err != nil {
		log.Printf("User with Google ID %s not found, creating new user", userInfo.ID)
		if user, err = app.Models.User.GetByEmail(userInfo.Email); err == nil && user != nil {
			log.Printf("Existing user with email %s found, updating Google ID", userInfo.Email)
			user.GoogleID = userInfo.ID
			user.Update()

		}
		user := data.User{
			Email:     userInfo.Email,
			FirstName: userInfo.FirstName,
			LastName:  userInfo.LastName,
			Password:  "",
			GoogleID:  userInfo.ID,
			Active:    1,
		}
		userID, err = app.Models.User.Insert(user)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create user: %s", err), http.StatusInternalServerError)
			return
		}
	} else {
		userID = user.ID
	}

	atExp := time.Now().Add(accessTokenTime)
	rtExp := time.Now().Add(refreshTokenTime)

	tokenPair, err :=
		app.TokenService.GenerateTokenPair(
			int64(userID),
			atExp,
			rtExp,
		)
	if err != nil {
		app.errorJSON(w, errors.New("failed to generate token"), http.StatusInternalServerError)
		return
	}
	log.Printf("Generated token pair for user ID %d", userID)
	app.setAccessTokenCookie(w, tokenPair.AccessToken, atExp)
	app.setRefreshTokenCookie(w, tokenPair.RefreshToken, rtExp)

	redirectURL := os.Getenv("REDIRECT_URL")

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
