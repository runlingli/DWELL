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

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {

	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {

	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	log.Printf("Received request for action: %s", requestPayload.Action)
	switch requestPayload.Action {

	case "register":
		app.register(w, requestPayload.Register)

	case "auth":
		app.authenticate(w, requestPayload.Auth)

	case "logout":
		app.logout(w, r)

	case "log":
		app.logEventViaRabbit(w, requestPayload.Log)

	case "mail":
		app.sendMail(w, requestPayload.Mail)

	case "verify":
		app.verifyCode(w, requestPayload.Verify)

	case "resource":
		app.getResource(w, r, requestPayload.Resource)

	case "forgot-password":
		app.forgotPassword(w, requestPayload.ForgotPassword)

	case "reset-password":
		app.resetPassword(w, requestPayload.ResetPassword)

	case "get-posts":
		app.getAllPosts(w)

	case "create-post":
		app.createPost(w, requestPayload.Post)

	case "update-post":
		app.updatePost(w, requestPayload.Post)

	case "delete-post":
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

	app.forwardToAuthService(w, r, "GET", url, nil)
}

func (app *Config) oauthGoogleLogin(w http.ResponseWriter, r *http.Request) {
	log.Printf("Forwarding browser to authentication-service Google OAuth login")

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get("http://authentication-service/authenticate/google")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTemporaryRedirect || resp.StatusCode == http.StatusFound {
		location := resp.Header.Get("Location")
		log.Printf("Redirecting browser to %s", location)
		http.Redirect(w, r, location, http.StatusFound)
		return
	}

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

func (app *Config) forwardToAuthService(
	w http.ResponseWriter,
	r *http.Request,
	method string,
	url string,
	body any,
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

	if r != nil {
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

	for _, c := range resp.Cookies() {
		http.SetCookie(w, c)
	}

	// 解析 JSON 响应
	var payload jsonResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		app.errorJSON(w, err)
		return
	}

	if payload.Error {
		log.Printf("Authentication service returned error: %s", payload.Message)
		statusCode := resp.StatusCode
		if statusCode == 0 {
			statusCode = http.StatusInternalServerError
		}
		app.errorJSON(w, errors.New(payload.Message), statusCode)
		return
	}

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

func (app *Config) logEventViaRabbit(w http.ResponseWriter, l LogPayload) {
	err := app.pushToQueue(l.Name, l.Data)
	if err != nil {

		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged via RabbitMQ"

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) pushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(app.Rabbit)
	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	j, _ := json.MarshalIndent(&payload, "", "\t")

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
