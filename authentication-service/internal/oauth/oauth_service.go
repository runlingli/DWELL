package oauth

import (
	"authentication/internal/store"
	"authentication/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/oauth2"
)

// const accessTokenTTL = 15 * time.Minute
// const refreshTokenTTL = 7 * 24 * time.Hour

type OauthService struct {
	refreshStore      *store.RefreshStore
	googleOauthConfig *oauth2.Config
}

type UserInfo struct {
	Email     string `json:"email"`
	ID        string `json:"sub"`
	FirstName string `json:"given_name"`
	LastName  string `json:"family_name"`
}

func NewOauthService(refreshStore *store.RefreshStore, googleOauthConfig *oauth2.Config) *OauthService {
	return &OauthService{
		refreshStore:      refreshStore,
		googleOauthConfig: googleOauthConfig,
	}
}

func (s *OauthService) GoogleLoginURL() (string, error) {
	state, err := utils.GenerateJTI()
	if err != nil {
		return "", err
	}

	key := "oauth:google:state:" + state
	ctx, cancel := context.WithTimeout(context.Background(), utils.RedisTimeout)
	defer cancel()
	err = s.refreshStore.Save(ctx, key, 6*time.Minute)
	if err != nil {
		return "", err
	}

	return s.googleOauthConfig.AuthCodeURL(state), nil
}

func (s *OauthService) ValidateGoogleState(state string) (bool, error) {
	key := "oauth:google:state:" + state
	ctx, cancel := context.WithTimeout(context.Background(), utils.RedisTimeout)
	defer cancel()
	exists, err := s.refreshStore.Exists(ctx, key)
	if err != nil || !exists {
		return false, err
	}
	_ = s.refreshStore.Delete(ctx, key)
	return true, nil
}

func (s *OauthService) GetUserInfo(code string) (UserInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	token, err := s.googleOauthConfig.Exchange(ctx, code)
	if err != nil {
		return UserInfo{}, err
	}
	client := s.googleOauthConfig.Client(ctx, token)
	resp, err := client.Get("https://openidconnect.googleapis.com/v1/userinfo")
	if err != nil {
		return UserInfo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return UserInfo{}, fmt.Errorf("failed to get user info: %s", resp.Status)
	}

	var info UserInfo
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return UserInfo{}, err
	}
	return info, nil
}
