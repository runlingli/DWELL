package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

// Shared types
type TokenData struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

type userData struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

// Token expiration constants
const accessTokenTime = 15 * time.Minute    // 15 minutes
const refreshTokenTime = 7 * 24 * time.Hour // 7 days

// setAccessTokenCookie sets the access token as a cookie
func (app *Config) setAccessTokenCookie(
	w http.ResponseWriter,
	accessToken string,
	accessTokenExp time.Time,
) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  accessTokenExp,
	})
}

// setRefreshTokenCookie sets the refresh token as a cookie
func (app *Config) setRefreshTokenCookie(
	w http.ResponseWriter,
	refreshToken string,
	refreshTokenExp time.Time,
) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Expires:  refreshTokenExp,
	})
}

// logRequest sends a log entry to the logger service
func (app *Config) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
