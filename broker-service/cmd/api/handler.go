package main

import (
	"broker/event"
	"bytes"         // 用于把 JSON 数据转成 io.Reader（HTTP 请求体需要）
	"encoding/json" // JSON 编码 / 解码
	"errors"        // 创建错误对象
	"log"
	"net/http" // HTTP 服务与客户端
	"time"
)

// =======================
// 前端 → Broker 的请求结构
// =======================

// RequestPayload 表示：
// 前端发送给 broker 的完整 JSON 请求结构
type RequestPayload struct {
	Action   string            `json:"action"`
	Register regPayload        `json:"register,omitempty"`
	Auth     AuthPayload       `json:"auth,omitempty"`
	Log      LogPayload        `json:"log,omitempty"`
	Mail     MailPayload       `json:"mail,omitempty"`
	Verify   VerifyCodePayload `json:"verify,omitempty"`
	Resource string            `json:"resource,omitempty"`
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
		// 登录认证 → 转发给 authentication-service
		app.register(w, requestPayload.Register)

	case "auth":
		// 登录认证 → 转发给 authentication-service
		app.authenticate(w, requestPayload.Auth)

	case "log":
		app.logEventViaRabbit(w, requestPayload.Log)

	case "mail":
		app.sendMail(w, requestPayload.Mail)

	case "verify":
		// 验证码验证 → 转发给 authentication-service
		app.verifyCode(w, requestPayload.Verify)

	case "resource":
		app.getResource(w, r, requestPayload.Resource)

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

func (app *Config) sendMail(w http.ResponseWriter, msg MailPayload) {
	jsonData, _ := json.MarshalIndent(msg, "", "\t")

	// call the mail service
	mailServiceURL := "http://mailer-service/send"

	// post to mail service
	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// make sure we get back the right status code
	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling mail service"))
		return
	}

	// send back json
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Message sent to " + msg.To

	app.writeJSON(w, http.StatusAccepted, payload)

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
