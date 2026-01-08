package main

import (
	"log"      // 用于在服务端打印日志（调试 / 错误）
	"net/http" // HTTP handler、状态码等
)

// =======================
// 发送邮件的 HTTP Handler
// =======================

// SendMail 是一个 HTTP handler
// 当前端 / broker 请求发送邮件时，会调用这个函数
func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {

	// =======================
	// 1. 定义“接收 JSON 的结构体”
	// =======================

	// mailMessage 用来接收 HTTP 请求体里的 JSON
	// 注意：这个结构体只用于“接收请求”
	type mailMessage struct {
		To      string `json:"to"`      // 收件人邮箱
		Subject string `json:"subject"` // 邮件标题
		Message string `json:"message"` // 邮件正文
	}

	// requestPayload 用来存放解析后的 JSON 数据
	var requestPayload mailMessage

	// =======================
	// 2. 读取并解析 JSON 请求体
	// =======================

	// app.readJSON 会做：
	// - 从 r.Body 读取数据
	// - 解析 JSON
	// - 填充到 requestPayload
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		// JSON 格式错误 / 缺字段 / body 过大等
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	// =======================
	// 3. 构造“业务层的邮件对象”
	// =======================

	// Message 是你系统内部真正用来“发邮件”的结构体
	// 它和 HTTP JSON 没有直接关系
	msg := Message{
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	// =======================
	// 4. 调用 Mailer 发送邮件
	// =======================

	// app.Mailer 是一个已经配置好的邮件发送器
	// SendSMTPMessage 会真正通过 SMTP 把邮件发出去
	err = app.Mailer.SendSMTPMessage(msg)
	if err != nil {
		// SMTP 连接失败 / 认证失败 / 网络错误等
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	// =======================
	// 5. 返回成功响应
	// =======================

	payload := jsonResponse{
		Error:   false,                          // 没有错误
		Message: "sent to " + requestPayload.To, // 给调用方一个明确反馈
	}

	// 返回 HTTP 202（Accepted）
	app.writeJSON(w, http.StatusAccepted, payload)
}
