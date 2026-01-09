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
