package token

import (
	"authentication/internal/store"
	"authentication/internal/utils"
	"context"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// const accessTokenTTL = 15 * time.Minute
// const refreshTokenTTL = 7 * 24 * time.Hour

type Service struct {
	refreshStore  *store.RefreshStore
	accessSecret  []byte
	refreshSecret []byte
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

func NewService(refreshStore *store.RefreshStore, accessSecret string, refreshSecret string) *Service {
	return &Service{
		refreshStore:  refreshStore,
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
	}
}

func (s *Service) GenerateTokenPair(userID int64, atExp, rtExp time.Time) (*TokenPair, error) {
	// 生成 jti
	var jti string
	if s.refreshStore != nil {
		var err error
		jti, err = utils.GenerateJTI()
		if err != nil {
			return nil, err
		}

		ctx, cancel := context.WithTimeout(context.Background(), utils.RedisTimeout)
		defer cancel()
		err = s.refreshStore.Save(ctx, jti, time.Until(rtExp))
		if err != nil {
			return nil, err
		}
	}

	// 生成 Access Token

	atClaims := jwt.MapClaims{
		"sub": strconv.FormatInt(userID, 10),
		"exp": atExp.Unix(),
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	atStr, err := at.SignedString(s.accessSecret)
	if err != nil {
		return nil, err
	}

	// 生成 Refresh Token（如果支持）
	// type RegisteredClaims struct {
	//	Issuer    string `json:"iss,omitempty"`
	//	Subject   string `json:"sub,omitempty"`
	//	Audience  ClaimStrings `json:"aud,omitempty"`
	//	ExpiresAt *NumericDate `json:"exp,omitempty"`
	//	NotBefore *NumericDate `json:"nbf,omitempty"`
	//	IssuedAt  *NumericDate `json:"iat,omitempty"`
	//	ID        string `json:"jti,omitempty"`
	//}

	var rtStr string

	if jti != "" {
		log.Printf("generating refresh token with exp: %v", rtExp)
		rtClaims := jwt.MapClaims{
			// convert int64 to string
			"sub": strconv.FormatInt(userID, 10),
			"exp": rtExp.Unix(),
			"jti": jti,
		}
		rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
		rtStr, err = rt.SignedString(s.refreshSecret)
		if err != nil {
			return nil, err
		}
	}

	return &TokenPair{
		AccessToken:  atStr,
		RefreshToken: rtStr,
	}, nil
}

// return new access token
func (s *Service) Refresh(refreshToken string, atExp time.Time) (string, error) {
	// 1. 如果系统不支持 refresh
	if s.refreshStore == nil {
		return "", errors.New("refresh token not supported")
	}

	// 2. 解析并校验 JWT
	//自动识别sub, exp, jti等标准字段到RegisteredClaims
	claims := &jwt.RegisteredClaims{}

	// 自动用token.Valid判断签名和过期时间
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("unexpected signing method")
			return nil, errors.New("unexpected signing method")
		}
		return s.refreshSecret, nil
	})
	log.Printf("refresh token exp: %v, now: %v", claims.ExpiresAt, time.Now())
	if err != nil || !token.Valid {
		log.Printf("invalid refresh token")
		return "", errors.New("invalid refresh token")
	}

	// 从claims.subject取出
	userID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		log.Printf("invalid subject")
		return "", errors.New("invalid subject")
	}

	jti := claims.ID
	if jti == "" {
		log.Printf("invalid jti")
		return "", errors.New("invalid jti")
	}

	// 6. 查 Redis
	ctx, cancel := context.WithTimeout(context.Background(), utils.RedisTimeout)
	defer cancel()

	exists, err := s.refreshStore.Exists(ctx, jti)
	if err != nil {
		return "", err
	}
	if !exists {
		log.Printf("jti does not exist")
		return "", errors.New("refresh token already used or revoked")
	}
	// generate new at
	atClaims := jwt.MapClaims{
		"sub": strconv.FormatInt(userID, 10),
		"exp": atExp.Unix(),
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	atStr, err := at.SignedString(s.accessSecret)

	if err != nil {
		return "", err
	}

	return atStr, nil
}

func (s *Service) ValidateAccessToken(accessToken string) (int64, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(accessToken, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.accessSecret, nil
	})

	if err != nil || !token.Valid {
		return 0, errors.New("invalid access token")
	}

	// 检查是否过期
	if claims.ExpiresAt == nil || time.Now().After(claims.ExpiresAt.Time) {
		log.Printf("access token expired")
		return 0, errors.New("access token expired")
	}

	// 解析 subject
	userID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return 0, errors.New("invalid subject in token")
	}
	return userID, nil
}
