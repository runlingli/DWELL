package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
)

// ForgotPassword handles sending a password reset code to the user's email
func (app *Config) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	log.Printf("Forgot password request received")

	var requestPayload struct {
		Email string `json:"email"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		log.Printf("Error reading forgot password request: %v", err)
		app.errorJSON(w, errors.New("invalid request payload"), http.StatusBadRequest)
		return
	}

	if requestPayload.Email == "" {
		app.errorJSON(w, errors.New("email is required"), http.StatusBadRequest)
		return
	}

	// Check if the user exists
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			// Don't reveal if user exists or not for security
			// Still return success to prevent email enumeration
			log.Printf("Forgot password: user not found for email %s", requestPayload.Email)
			payload := jsonResponse{
				Error:   false,
				Message: "If an account exists with this email, a reset code has been sent",
			}
			app.writeJSON(w, http.StatusOK, payload)
			return
		}
		log.Printf("Error checking user for forgot password: %v", err)
		app.errorJSON(w, errors.New("server error"), http.StatusInternalServerError)
		return
	}

	// Check if user has a password (not OAuth-only user)
	if user.Password == "" {
		log.Printf("User %s is OAuth-only, cannot reset password", requestPayload.Email)
		app.errorJSON(w, errors.New("this account uses Google sign-in. Please login with Google"), http.StatusBadRequest)
		return
	}

	// Send password reset code
	err = app.MailService.SendPasswordResetCode(requestPayload.Email)
	if err != nil {
		log.Printf("Error sending password reset code: %v", err)
		// Don't expose internal errors, but log them
		if err.Error() == "Too many requests. Please try again later." {
			app.errorJSON(w, err, http.StatusTooManyRequests)
			return
		}
		app.errorJSON(w, errors.New("failed to send reset code"), http.StatusInternalServerError)
		return
	}

	// Log the password reset request
	go func() {
		err := app.logRequest("password-reset-request", fmt.Sprintf("Password reset requested for %s", requestPayload.Email))
		if err != nil {
			log.Printf("Failed to log password reset request: %v", err)
		}
	}()

	payload := jsonResponse{
		Error:   false,
		Message: "Password reset code sent successfully",
	}
	app.writeJSON(w, http.StatusOK, payload)
}

// ResetPassword handles resetting the user's password with a verification code
func (app *Config) ResetPassword(w http.ResponseWriter, r *http.Request) {
	log.Printf("Reset password request received")

	var requestPayload struct {
		Email            string `json:"email"`
		VerificationCode string `json:"verification_code"`
		NewPassword      string `json:"new_password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		log.Printf("Error reading reset password request: %v", err)
		app.errorJSON(w, errors.New("invalid request payload"), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if requestPayload.Email == "" {
		app.errorJSON(w, errors.New("email is required"), http.StatusBadRequest)
		return
	}
	if requestPayload.VerificationCode == "" {
		app.errorJSON(w, errors.New("verification code is required"), http.StatusBadRequest)
		return
	}
	if requestPayload.NewPassword == "" {
		app.errorJSON(w, errors.New("new password is required"), http.StatusBadRequest)
		return
	}
	if len(requestPayload.NewPassword) < 6 {
		app.errorJSON(w, errors.New("password must be at least 6 characters"), http.StatusBadRequest)
		return
	}

	// Verify the reset code
	isValid, err := app.MailService.VerifyPasswordResetCode(requestPayload.Email, requestPayload.VerificationCode)
	if err != nil {
		log.Printf("Error verifying password reset code: %v", err)
		app.errorJSON(w, errors.New("invalid or expired reset code"), http.StatusBadRequest)
		return
	}
	if !isValid {
		log.Printf("Invalid password reset code for %s", requestPayload.Email)
		app.errorJSON(w, errors.New("invalid verification code"), http.StatusBadRequest)
		return
	}

	// Get the user
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		log.Printf("Error getting user for password reset: %v", err)
		app.errorJSON(w, errors.New("user not found"), http.StatusNotFound)
		return
	}

	// Update the password
	err = user.ResetPassword(requestPayload.NewPassword)
	if err != nil {
		log.Printf("Error resetting password: %v", err)
		app.errorJSON(w, errors.New("failed to reset password"), http.StatusInternalServerError)
		return
	}

	// Log the password reset
	go func() {
		err := app.logRequest("password-reset-complete", fmt.Sprintf("Password reset completed for %s", requestPayload.Email))
		if err != nil {
			log.Printf("Failed to log password reset completion: %v", err)
		}
	}()

	payload := jsonResponse{
		Error:   false,
		Message: "Password reset successfully",
	}
	app.writeJSON(w, http.StatusOK, payload)
}
