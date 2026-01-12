package main

import (
	"broker/event"
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

// =======================
// 前端 → Broker 的请求结构
// =======================

// RequestPayload 表示：
// 前端发送给 broker 的完整 JSON 请求结构
type RequestPayload struct {
	Action         string               `json:"action"`
	Register       regPayload           `json:"register,omitempty"`
	Auth           AuthPayload          `json:"auth,omitempty"`
	Log            LogPayload           `json:"log,omitempty"`
	Mail           MailPayload          `json:"mail,omitempty"`
	Verify         VerifyCodePayload    `json:"verify,omitempty"`
	Resource       string               `json:"resource,omitempty"`
	ForgotPassword ForgotPasswordPaylod `json:"forgot_password,omitempty"`
	ResetPassword  ResetPasswordPayload `json:"reset_password,omitempty"`
	Post           PostPayload          `json:"post,omitempty"`
	DeletePost     DeletePostPayload    `json:"delete_post,omitempty"`
}

type regPayload struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	VerificationCode string `json:"verification_code"`
}

type AuthPayload struct {
	Email    string `json:"email"`    // 用户邮箱
	Password string `json:"password"` // 用户密码（明文）
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type VerifyCodePayload struct {
	Email string `json:"email"`
}

type MailPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type ForgotPasswordPaylod struct {
	Email string `json:"email"`
}

type ResetPasswordPayload struct {
	Email            string `json:"email"`
	VerificationCode string `json:"verification_code"`
	NewPassword      string `json:"new_password"`
}

// Post-related payloads
type PostPayload struct {
	ID               int      `json:"id,omitempty"`
	Title            string   `json:"title"`
	Price            float64  `json:"price"`
	Location         string   `json:"location,omitempty"`
	Neighborhood     string   `json:"neighborhood"`
	Lat              float64  `json:"lat"`
	Lng              float64  `json:"lng"`
	Radius           int      `json:"radius"`
	Type             string   `json:"type"`
	ImageURL         string   `json:"imageUrl"`
	AdditionalImages []string `json:"additionalImages,omitempty"`
	Description      string   `json:"description"`
	Bedrooms         int      `json:"bedrooms"`
	Bathrooms        int      `json:"bathrooms"`
	AvailableFrom    int64    `json:"availableFrom"`
	AvailableTo      int64    `json:"availableTo"`
	AuthorID         int      `json:"authorId"`
}

type DeletePostPayload struct {
	ID       int `json:"id"`
	AuthorID int `json:"authorId"`
}

// =======================
// Broker 测试接口
// =======================

// Broker 是一个 HTTP handler
func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {

	// 构造一个标准 JSON 响应
	payload := jsonResponse{
		Error:   false,            // 没有错误
		Message: "Hit the broker", // 只是返回一条提示信息
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

// =======================
// Broker 核心入口
// =======================

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {

	// 用于接收前端 JSON 请求
	var requestPayload RequestPayload

	// 读取并解析请求体 JSON
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		// JSON 解析失败（格式错误、字段不对等）
		app.errorJSON(w, err)
		return
	}
	log.Printf("Received request for action: %s", requestPayload.Action)
	// 根据 action 决定调用哪个微服务
	switch requestPayload.Action {

	case "register":
		// 注册 → 转发给 authentication-service
		app.register(w, requestPayload.Register)

	case "auth":
		// 登录认证 → 转发给 authentication-service
		app.authenticate(w, requestPayload.Auth)

	case "logout":
		// 登出 → 转发给 authentication-service
		app.logout(w, r)

	case "log":
		app.logEventViaRabbit(w, requestPayload.Log)

	case "mail":
		app.sendMail(w, requestPayload.Mail)

	case "verify":
		// 验证码验证 → 转发给 authentication-service
		app.verifyCode(w, requestPayload.Verify)

	case "resource":
		app.getResource(w, r, requestPayload.Resource)

	case "forgot-password":
		// 忘记密码 → 转发给 authentication-service
		app.forgotPassword(w, requestPayload.ForgotPassword)

	case "reset-password":
		// 重置密码 → 转发给 authentication-service
		app.resetPassword(w, requestPayload.ResetPassword)

	case "get-posts":
		// 获取所有帖子 → 转发给 post-service
		app.getAllPosts(w)

	case "create-post":
		// 创建帖子 → 转发给 post-service
		app.createPost(w, requestPayload.Post)

	case "update-post":
		// 更新帖子 → 转发给 post-service
		app.updatePost(w, requestPayload.Post)

	case "delete-post":
		// 删除帖子 → 转发给 post-service
		app.deletePost(w, requestPayload.DeletePost)

	default:
		app.errorJSON(w, errors.New("unknown action"))
	}
}

func (app *Config) register(w http.ResponseWriter, a regPayload) {
	app.forwardToAuthService(w, nil, "POST", "http://authentication-service/register", a)
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	log.Printf("Authenticating user: %s", a.Email)
	app.forwardToAuthService(w, nil, "POST", "http://authentication-service/authenticate", a)
}

func (app *Config) verifyCode(w http.ResponseWriter, v VerifyCodePayload) {
	log.Printf("Verifying code for user: %s", v.Email)
	app.forwardToAuthService(w, nil, "POST", "http://authentication-service/verify-email", v)
}

func (app *Config) getResource(w http.ResponseWriter, r *http.Request, resource string) {
	var url string
	switch resource {
	case "profile":
		url = "http://authentication-service/resource/profile"
	default:
		app.errorJSON(w, errors.New("unknown resource"))
		return
	}

	// GET 请求不需要 body，传 r 让 Cookie 可以复制
	app.forwardToAuthService(w, r, "GET", url, nil)
}

func (app *Config) oauthGoogleLogin(w http.ResponseWriter, r *http.Request) {
	log.Printf("Forwarding browser to authentication-service Google OAuth login")

	// 直接用 http.Client 发请求获取 OAuth URL
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // 不自动跟随重定向
		},
	}

	resp, err := client.Get("http://authentication-service/authenticate/google")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer resp.Body.Close()

	// Broker 拿到 Location 后直接 302 给浏览器
	if resp.StatusCode == http.StatusTemporaryRedirect || resp.StatusCode == http.StatusFound {
		location := resp.Header.Get("Location")
		log.Printf("Redirecting browser to %s", location)
		http.Redirect(w, r, location, http.StatusFound) // 直接 302 浏览器跳转
		return
	}

	// 如果认证服务返回错误 JSON
	var payload struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, resp.StatusCode, payload)
}

// forwardToAuthService 通用转发函数
func (app *Config) forwardToAuthService(
	w http.ResponseWriter,
	r *http.Request,
	method string,
	url string,
	body any, // 可以传结构体，内部会 marshal
) {
	var reader *bytes.Reader
	if body != nil {
		jsonData, err := json.MarshalIndent(body, "", "\t")
		if err != nil {
			app.errorJSON(w, err)
			return
		}
		log.Printf("Forwarding request to %s, body: %s", url, string(jsonData))
		reader = bytes.NewReader(jsonData)
	} else {
		reader = bytes.NewReader(nil)
	}

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if r != nil { // 只有 r 不为 nil 时才复制 cookie
		for _, c := range r.Cookies() {
			req.AddCookie(c)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer resp.Body.Close()

	// 把认证服务返回的 Cookie 写回浏览器
	for _, c := range resp.Cookies() {
		http.SetCookie(w, c)
	}

	// 解析 JSON 响应
	var payload jsonResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		app.errorJSON(w, err)
		return
	}

	// 根据认证服务返回的 error 和状态码统一处理
	if payload.Error {
		log.Printf("Authentication service returned error: %s", payload.Message)
		statusCode := resp.StatusCode
		if statusCode == 0 {
			statusCode = http.StatusInternalServerError
		}
		app.errorJSON(w, errors.New(payload.Message), statusCode)
		return
	}

	// 转发认证服务的状态码和响应
	app.writeJSON(w, resp.StatusCode, payload)
}

func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {

}

func (app *Config) forgotPassword(w http.ResponseWriter, p ForgotPasswordPaylod) {
	log.Printf("Forwarding forgot password request for: %s", p.Email)
	app.forwardToAuthService(w, nil, "POST", "http://authentication-service/forgot-password", p)
}

func (app *Config) resetPassword(w http.ResponseWriter, p ResetPasswordPayload) {
	log.Printf("Forwarding reset password request for: %s", p.Email)
	app.forwardToAuthService(w, nil, "POST", "http://authentication-service/reset-password", p)
}

func (app *Config) logout(w http.ResponseWriter, r *http.Request) {
	log.Printf("Forwarding logout request")
	app.forwardToAuthService(w, r, "POST", "http://authentication-service/logout", nil)
}

// Post service handlers
func (app *Config) getAllPosts(w http.ResponseWriter) {
	log.Printf("Forwarding get all posts request")
	app.forwardToPostService(w, "GET", "http://post-service/posts", nil)
}

func (app *Config) createPost(w http.ResponseWriter, p PostPayload) {
	log.Printf("Forwarding create post request")
	app.forwardToPostService(w, "POST", "http://post-service/posts", p)
}

func (app *Config) updatePost(w http.ResponseWriter, p PostPayload) {
	log.Printf("Forwarding update post request for ID: %d", p.ID)
	url := "http://post-service/posts/" + strconv.Itoa(p.ID)
	app.forwardToPostService(w, "PUT", url, p)
}

func (app *Config) deletePost(w http.ResponseWriter, p DeletePostPayload) {
	log.Printf("Forwarding delete post request for ID: %d", p.ID)
	url := "http://post-service/posts/" + strconv.Itoa(p.ID)
	body := map[string]int{"authorId": p.AuthorID}
	app.forwardToPostService(w, "DELETE", url, body)
}


// RESTful API handlers for posts
func (app *Config) GetAllPostsREST(w http.ResponseWriter, r *http.Request) {
	log.Printf("RESTful: GET all posts")
	app.forwardToPostService(w, "GET", "http://post-service/posts", nil)
}

func (app *Config) GetPostByIDREST(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	log.Printf("RESTful: GET post by ID: %s", id)
	url := "http://post-service/posts/" + id
	app.forwardToPostService(w, "GET", url, nil)
}

func (app *Config) CreatePostREST(w http.ResponseWriter, r *http.Request) {
	log.Println("========== Broker: CreatePostREST START ==========")
	var post PostPayload
	if err := app.readJSON(w, r, &post); err != nil {
		log.Printf("ERROR reading post from request: %v", err)
		app.errorJSON(w, err)
		return
	}
	log.Printf("Received post - Title: %s, AuthorID: %d", post.Title, post.AuthorID)
	if post.AuthorID == 0 {
		log.Println("WARNING: AuthorID is 0!")
	}
	app.forwardToPostService(w, "POST", "http://post-service/posts", post)
}

func (app *Config) UpdatePostREST(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var post PostPayload
	if err := app.readJSON(w, r, &post); err != nil {
		app.errorJSON(w, err)
		return
	}
	log.Printf("RESTful: UPDATE post ID: %s", id)
	url := "http://post-service/posts/" + id
	app.forwardToPostService(w, "PUT", url, post)
}

func (app *Config) DeletePostREST(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var body struct {
		AuthorID int `json:"authorId"`
	}
	if err := app.readJSON(w, r, &body); err != nil {
		app.errorJSON(w, err)
		return
	}
	log.Printf("RESTful: DELETE post ID: %s", id)
	url := "http://post-service/posts/" + id
	app.forwardToPostService(w, "DELETE", url, body)
}

func (app *Config) forwardToPostService(w http.ResponseWriter, method, url string, body any) {
	var reader *bytes.Reader
	if body != nil {
		jsonData, err := json.MarshalIndent(body, "", "\t")
		if err != nil {
			app.errorJSON(w, err)
			return
		}
		log.Printf("Forwarding request to %s, body: %s", url, string(jsonData))
		reader = bytes.NewReader(jsonData)
	} else {
		reader = bytes.NewReader(nil)
	}

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer resp.Body.Close()

	var payload jsonResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		app.errorJSON(w, err)
		return
	}

	if payload.Error {
		log.Printf("Post service returned error: %s", payload.Message)
		statusCode := resp.StatusCode
		if statusCode == 0 {
			statusCode = http.StatusInternalServerError
		}
		app.errorJSON(w, errors.New(payload.Message), statusCode)
		return
	}

	app.writeJSON(w, resp.StatusCode, payload)
}

func (app *Config) sendMail(w http.ResponseWriter, msg MailPayload) {
	// Send mail via RabbitMQ (async)
	err := app.sendMailViaRabbit(msg.To, msg.Subject, msg.Message)
	if err != nil {
		log.Printf("Error sending mail via RabbitMQ: %v", err)
		app.errorJSON(w, err)
		return
	}

	// send back json
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Message queued for delivery to " + msg.To

	app.writeJSON(w, http.StatusAccepted, payload)
}

// sendMailViaRabbit sends an email through RabbitMQ
func (app *Config) sendMailViaRabbit(to, subject, message string) error {
	emitter, err := event.NewEventEmitter(app.Rabbit)
	if err != nil {
		return err
	}
	return emitter.SendMail(to, subject, message)
}

// logEventViaRabbit 通过 RabbitMQ 将日志事件发送给 logger-service。
// w: HTTP 响应对象，用于返回结果给调用者
// l: 日志内容，类型 LogPayload（包含 Name 和 Data 字段）
func (app *Config) logEventViaRabbit(w http.ResponseWriter, l LogPayload) {
	// 调用 pushToQueue 将日志发送到 RabbitMQ 队列
	err := app.pushToQueue(l.Name, l.Data)
	if err != nil {
		// 如果发送失败，返回 JSON 错误响应
		app.errorJSON(w, err)
		return
	}

	// 构造返回给客户端的 JSON 响应
	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged via RabbitMQ"

	// 写回 HTTP 响应，状态码 202 Accepted 表示请求已接收
	app.writeJSON(w, http.StatusAccepted, payload)
}

// pushToQueue 将消息发送到 RabbitMQ 队列
// name: 日志类型，例如 "log.INFO" 或 "log.ERROR"
// msg: 日志内容
func (app *Config) pushToQueue(name, msg string) error {
	// 创建一个新的 RabbitMQ Emitter，用于发送消息
	emitter, err := event.NewEventEmitter(app.Rabbit)
	if err != nil {
		return err // 如果创建失败直接返回错误
	}

	// 构造要发送的日志负载
	payload := LogPayload{
		Name: name, // 日志类型
		Data: msg,  // 日志内容
	}

	// 将 payload 转成 JSON 字符串
	// json.MarshalIndent 可以格式化输出，方便调试
	j, _ := json.MarshalIndent(&payload, "", "\t")

	// 调用 Emitter.Push 发布消息到交换机
	// 这里 routing key 固定为 "log.INFO"，可以根据需求修改
	err = emitter.Push(string(j), "log.INFO")
	if err != nil {
		return err
	}

	return nil
}

// =======================
// RESTful Auth API Handlers
// =======================

func (app *Config) LoginREST(w http.ResponseWriter, r *http.Request) {
	log.Println("RESTful: POST /auth/login")
	var payload AuthPayload
	if err := app.readJSON(w, r, &payload); err != nil {
		app.errorJSON(w, err)
		return
	}
	app.forwardToAuthService(w, nil, "POST", "http://authentication-service/auth/login", payload)
}

func (app *Config) LogoutREST(w http.ResponseWriter, r *http.Request) {
	log.Println("RESTful: POST /auth/logout")
	app.forwardToAuthService(w, r, "POST", "http://authentication-service/auth/logout", nil)
}

func (app *Config) RegisterREST(w http.ResponseWriter, r *http.Request) {
	log.Println("RESTful: POST /auth/register")
	var payload regPayload
	if err := app.readJSON(w, r, &payload); err != nil {
		app.errorJSON(w, err)
		return
	}
	app.forwardToAuthService(w, nil, "POST", "http://authentication-service/auth/register", payload)
}

func (app *Config) VerifyEmailREST(w http.ResponseWriter, r *http.Request) {
	log.Println("RESTful: POST /auth/verify-email")
	var payload VerifyCodePayload
	if err := app.readJSON(w, r, &payload); err != nil {
		app.errorJSON(w, err)
		return
	}
	app.forwardToAuthService(w, nil, "POST", "http://authentication-service/auth/verify-email", payload)
}

func (app *Config) ForgotPasswordREST(w http.ResponseWriter, r *http.Request) {
	log.Println("RESTful: POST /auth/forgot-password")
	var payload ForgotPasswordPaylod
	if err := app.readJSON(w, r, &payload); err != nil {
		app.errorJSON(w, err)
		return
	}
	app.forwardToAuthService(w, nil, "POST", "http://authentication-service/auth/forgot-password", payload)
}

func (app *Config) ResetPasswordREST(w http.ResponseWriter, r *http.Request) {
	log.Println("RESTful: POST /auth/reset-password")
	var payload ResetPasswordPayload
	if err := app.readJSON(w, r, &payload); err != nil {
		app.errorJSON(w, err)
		return
	}
	app.forwardToAuthService(w, nil, "POST", "http://authentication-service/auth/reset-password", payload)
}

func (app *Config) ProfileREST(w http.ResponseWriter, r *http.Request) {
	log.Println("RESTful: GET /auth/profile")
	app.forwardToAuthService(w, r, "GET", "http://authentication-service/auth/profile", nil)
}

// =======================
// RESTful Favorites API Handlers
// =======================

func (app *Config) GetUserFavoriteIDsREST(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	log.Printf("RESTful: GET favorites for user: %s", userId)
	url := "http://favourite-service/favorites/" + userId + "/ids"
	app.forwardToFavoriteService(w, "GET", url, nil)
}

func (app *Config) AddFavoriteREST(w http.ResponseWriter, r *http.Request) {
	log.Println("RESTful: POST /favorites")
	var payload struct {
		UserID int `json:"userId"`
		PostID int `json:"postId"`
	}
	if err := app.readJSON(w, r, &payload); err != nil {
		app.errorJSON(w, err)
		return
	}
	app.forwardToFavoriteService(w, "POST", "http://favourite-service/favorites", payload)
}

func (app *Config) RemoveFavoriteREST(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	postId := chi.URLParam(r, "postId")
	log.Printf("RESTful: DELETE favorite - user: %s, post: %s", userId, postId)
	url := "http://favourite-service/favorites/" + userId + "/" + postId
	app.forwardToFavoriteService(w, "DELETE", url, nil)
}

func (app *Config) SyncFavoritesREST(w http.ResponseWriter, r *http.Request) {
	log.Println("RESTful: POST /favorites/sync")
	var payload struct {
		UserID  int   `json:"userId"`
		PostIDs []int `json:"postIds"`
	}
	if err := app.readJSON(w, r, &payload); err != nil {
		app.errorJSON(w, err)
		return
	}
	app.forwardToFavoriteService(w, "POST", "http://favourite-service/favorites/sync", payload)
}

func (app *Config) forwardToFavoriteService(w http.ResponseWriter, method, url string, body any) {
	var reader *bytes.Reader
	if body != nil {
		jsonData, err := json.MarshalIndent(body, "", "\t")
		if err != nil {
			app.errorJSON(w, err)
			return
		}
		log.Printf("Forwarding request to %s, body: %s", url, string(jsonData))
		reader = bytes.NewReader(jsonData)
	} else {
		reader = bytes.NewReader(nil)
	}

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer resp.Body.Close()

	var payload jsonResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		app.errorJSON(w, err)
		return
	}

	if payload.Error {
		log.Printf("Favorite service returned error: %s", payload.Message)
		statusCode := resp.StatusCode
		if statusCode == 0 {
			statusCode = http.StatusInternalServerError
		}
		app.errorJSON(w, errors.New(payload.Message), statusCode)
		return
	}

	app.writeJSON(w, resp.StatusCode, payload)
}
