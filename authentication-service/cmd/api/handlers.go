package main

import (
	"authentication/data"
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

type TokenData struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

type userData struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

const accessTokenTime = 15 * time.Minute    // 15 minutes
const refreshTokenTime = 7 * 24 * time.Hour // 7 days
//const accessTokenTime = 10 * time.Second
//const refreshTokenTime = 30 * time.Second

// Authenticate 是用户登录认证处理器
func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start authentication")
	var requestPayload struct {
		Email    string `json:"email"`    // 用户邮箱
		Password string `json:"password"` // 用户密码（明文）
	}

	// =============================
	// 读取 JSON 请求体
	// =============================
	// 调用之前写好的 readJSON 函数
	// 将请求体解析到 requestPayload 结构体
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		// 如果解析失败，返回 400 错误
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// =============================
	// 验证用户是否存在
	// =============================
	// 从数据库根据邮箱查找用户
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		// 用户不存在或数据库查询失败
		// 不暴露具体原因，统一返回 "invalid credentials"
		log.Printf("User not found by email: %v", err)
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	// =============================
	// 校验密码
	// =============================
	// 调用 User 结构体的方法 PasswordMatches
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		// 密码不匹配或校验出错
		log.Printf("Invalid password attempt for user %s: %v", requestPayload.Email, err)
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}
	go func() {
		err := app.logRequest("authentication", fmt.Sprintf("%s logged in", user.Email))
		if err != nil {
			log.Printf("Failed to send log to logger-service during authentication: %v", err)
		}
	}()

	atExp := time.Now().Add(accessTokenTime)
	rtExp := time.Now().Add(refreshTokenTime)

	tokenPair, err :=
		app.TokenService.GenerateTokenPair(
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
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data: userData{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		},
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) Register(w http.ResponseWriter, r *http.Request) {

	log.Printf("Register service begin\n")
	var requestPayload struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		FirstName        string `json:"first_name"`
		LastName         string `json:"last_name"`
		VerificationCode string `json:"verification_code"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, errors.New("Error reading user information"), http.StatusBadRequest)
		return
	}

	existingUser, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err == nil && existingUser != nil {
		app.errorJSON(w, errors.New("Email already exists"), http.StatusConflict)
		return
	} else if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// 验证验证码
	isValid, err := app.MailService.VerifyCode(requestPayload.Email, requestPayload.VerificationCode)
	if err != nil || !isValid {
		app.errorJSON(w, errors.New("Invalid verification code"), http.StatusBadRequest)
		return
	}
	// 创建新用户

	user := data.User{
		Email:     requestPayload.Email,
		FirstName: requestPayload.FirstName,
		LastName:  requestPayload.LastName,
		Password:  requestPayload.Password,
		Active:    1,
	}

	log.Printf("Inserting user(inactive): %+v\n", user)

	_, err = app.Models.User.Insert(user)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	go func() {
		err := app.logRequest("registration attempt", fmt.Sprintf("%s registered", user.Email))
		if err != nil {
			log.Printf("Failed to send log to logger-service during registration: %v", err)
		}
	}()

	log.Printf("Verify email service begin\n")

	atExp := time.Now().Add(accessTokenTime)
	rtExp := time.Now().Add(refreshTokenTime)

	tokenPair, err :=
		app.TokenService.GenerateTokenPair(
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
	}

	app.writeJSON(w, http.StatusCreated, payload)

}

func (app *Config) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email string `json:"email"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, errors.New("Error reading email information"), http.StatusBadRequest)
		return
	}

	existingUser, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err == nil && existingUser != nil {
		app.errorJSON(w, errors.New("Email already exists"), http.StatusConflict)
		return
	} else if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	app.MailService.SendCode(requestPayload.Email, "DWELL Verification Code")
}

func (app *Config) setAccessTokenCookie(
	w http.ResponseWriter,
	accessToken string,
	accessTokenExp time.Time,
) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: false, // 调试阶段
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  accessTokenExp,
	})
}

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

func (app *Config) Profile(w http.ResponseWriter, r *http.Request) {
	var userID int64
	// 尝试取 access token
	atCookie, err := r.Cookie("access_token")
	if atCookie != nil {
		log.Printf("access token exists")
		userID, err = app.TokenService.ValidateAccessToken(atCookie.Value)
	}

	if err != nil {
		log.Printf("access token invalid, trying refresh token")
		// 尝试取 refresh token
		rtCookie, rtErr := r.Cookie("refresh_token")
		if rtErr != nil {
			log.Printf("refresh token also missing")
			app.errorJSON(w, errors.New("access token and refresh token missing"), http.StatusUnauthorized)
			return
		}

		// 调用 TokenService 刷新 token
		newAT, refreshErr := app.TokenService.Refresh(rtCookie.Value, time.Now().Add(accessTokenTime))
		if refreshErr != nil {
			log.Printf("refresh token invalid")
			app.errorJSON(w, errors.New("refresh token invalid"), http.StatusUnauthorized)
			return
		}

		// 刷新 cookie
		log.Printf("set auth cookies")
		app.setAccessTokenCookie(w, newAT, time.Now().Add(accessTokenTime))

		// 用新的 AT 继续校验
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
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
		},
	}

	app.writeJSON(w, http.StatusOK, payload)
}

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
