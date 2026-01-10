package mail

import (
	"authentication/internal/store"
	"authentication/internal/utils"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

// const accessTokenTTL = 15 * time.Minute
// const refreshTokenTTL = 7 * 24 * time.Hour

type MailService struct {
	refreshStore *store.RefreshStore
}

func NewMailService(refreshStore *store.RefreshStore) *MailService {
	return &MailService{
		refreshStore: refreshStore,
	}
}

func (s *MailService) SendCode(to string, subject string) error {
	var mail struct {
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}
	key := "mail:verification:" + to

	ctxGet, cancelGet := context.WithTimeout(context.Background(), utils.RedisTimeout)
	defer cancelGet()

	code, err := s.refreshStore.GetValue(ctxGet, key)

	if err != nil {
		return err
	}

	if code != "" {
		log.Printf("Too many requests for email: %s", to)
		return errors.New("Too many requests. Please try again later.")
	}

	mail.To = to
	mail.Subject = subject

	verificationVode, err := utils.GenerateCode()
	if err != nil {
		log.Printf("Error generating verification code: %v", err)
		return err
	}

	ctxSave, cancelSave := context.WithTimeout(context.Background(), utils.RedisTimeout)
	defer cancelSave()
	err = s.refreshStore.SavePair(ctxSave, key, verificationVode, 5*time.Minute)
	if err != nil {
		log.Printf("Error saving verification code to store: %v", err)
		return err
	}

	mail.Message = fmt.Sprintf("Your verification code is: %s", verificationVode)

	jsonData, _ := json.MarshalIndent(mail, "", "\t")
	logServiceURL := "http://mailer-service/send"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating mail service request: %v", err)
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: utils.MailTimeout}
	resp, err := client.Do(request)
	if err != nil {
		log.Printf("Error calling mail service: %v", err)
		return err
	}

	defer resp.Body.Close()

	// make sure we get back the right status code
	if resp.StatusCode != http.StatusAccepted {
		return errors.New("error calling mail service")
	}
	log.Printf("Verification code sent to email: %s", to)

	return nil
}

func (s *MailService) VerifyCode(email string, code string) (bool, error) {
	key := "mail:verification:" + email
	ctx, cancel := context.WithTimeout(context.Background(), utils.RedisTimeout)
	defer cancel()
	storedCode, err := s.refreshStore.GetValue(ctx, key)
	if err != nil {
		return false, err
	}
	if storedCode != code {
		return false, nil
	}
	_ = s.refreshStore.Delete(ctx, key)
	return true, nil
}

// SendPasswordResetCode sends a password reset code to the email
// Uses a different Redis key prefix than verification codes
func (s *MailService) SendPasswordResetCode(to string) error {
	key := "mail:password-reset:" + to

	ctxGet, cancelGet := context.WithTimeout(context.Background(), utils.RedisTimeout)
	defer cancelGet()

	code, err := s.refreshStore.GetValue(ctxGet, key)
	if err != nil {
		return err
	}

	if code != "" {
		log.Printf("Too many password reset requests for email: %s", to)
		return errors.New("Too many requests. Please try again later.")
	}

	resetCode, err := utils.GenerateCode()
	if err != nil {
		log.Printf("Error generating password reset code: %v", err)
		return err
	}

	ctxSave, cancelSave := context.WithTimeout(context.Background(), utils.RedisTimeout)
	defer cancelSave()
	// Password reset codes valid for 10 minutes
	err = s.refreshStore.SavePair(ctxSave, key, resetCode, 10*time.Minute)
	if err != nil {
		log.Printf("Error saving password reset code to store: %v", err)
		return err
	}

	mail := struct {
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}{
		To:      to,
		Subject: "DWELL Password Reset",
		Message: fmt.Sprintf("Your password reset code is: %s\n\nThis code will expire in 10 minutes.", resetCode),
	}

	jsonData, _ := json.MarshalIndent(mail, "", "\t")
	mailServiceURL := "http://mailer-service/send"

	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating mail service request: %v", err)
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: utils.MailTimeout}
	resp, err := client.Do(request)
	if err != nil {
		log.Printf("Error calling mail service: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return errors.New("error calling mail service")
	}
	log.Printf("Password reset code sent to email: %s", to)

	return nil
}

// VerifyPasswordResetCode verifies the password reset code
func (s *MailService) VerifyPasswordResetCode(email string, code string) (bool, error) {
	key := "mail:password-reset:" + email
	ctx, cancel := context.WithTimeout(context.Background(), utils.RedisTimeout)
	defer cancel()

	storedCode, err := s.refreshStore.GetValue(ctx, key)
	if err != nil {
		return false, err
	}
	if storedCode == "" {
		return false, errors.New("no reset code found or code expired")
	}
	if storedCode != code {
		return false, nil
	}
	// Delete the code after successful verification
	_ = s.refreshStore.Delete(ctx, key)
	return true, nil
}
