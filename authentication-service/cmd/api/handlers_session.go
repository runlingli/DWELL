package main

import (
	"log"
	"net/http"
	"time"
)

// Logout handles user logout by clearing tokens
func (app *Config) Logout(w http.ResponseWriter, r *http.Request) {
	log.Printf("Logout request received")

	// Clear the access token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0), // Expire immediately
		MaxAge:   -1,
	})

	// Clear the refresh token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Unix(0, 0), // Expire immediately
		MaxAge:   -1,
	})

	// TODO: Optionally invalidate the refresh token in Redis
	// This would require extracting the JTI from the refresh token
	// and calling app.TokenService.RefreshStore.Delete(ctx, jti)

	payload := jsonResponse{
		Error:   false,
		Message: "Logged out successfully",
	}
	app.writeJSON(w, http.StatusOK, payload)
}
