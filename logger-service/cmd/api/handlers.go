package main

import (
	"log-service/data" // 你自己的数据层，用来操作 MongoDB
	"net/http"         // HTTP 协议相关（handler、status code 等）
)

// =======================
// 前端 / 上游服务发送的 JSON 结构
// =======================

// JSONPayload 用来“接收请求体里的 JSON 数据”
//
// 例如前端或 broker 发送：
//
//	{
//	  "name": "authentication",
//	  "data": "user admin@example.com logged in"
//	}
type JSONPayload struct {
	Name string `json:"name"` // 日志来源（哪个服务）
	Data string `json:"data"` // 日志内容
}

// =======================
// 写日志的 HTTP Handler
// =======================

// WriteLog 是一个 HTTP handler
// 当路由匹配到它时，这个函数会被调用
func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {

	var requestPayload struct {
		Name string `json:"name"` // 用户邮箱
		Data string `json:"data"` // 用户密码（明文）
	}
	// =======================
	// 1. 读取请求体中的 JSON
	// =======================

	// requestPayload 用来接收解析后的 JSON 数据
	//var requestPayload JSONPayload
	//HTTP 请求体是字节流, Go 需要一个 结构体模板 才能把 JSON 映射进来
	//JSONPayload 就是 “JSON → Go 的翻译模板”

	// app.readJSON：
	// - 从 r.Body 中读取数据
	// - 解析 JSON
	// - 填充到 requestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// =======================
	// 2. 构造数据库层需要的结构体
	// =======================

	// data.LogEntry 是“数据库模型”
	// 它代表 MongoDB 中的一条日志记录
	event := data.LogEntry{
		Name: requestPayload.Name, // 日志来源
		Data: requestPayload.Data, // 日志内容
	}

	// =======================
	// 3. 写入 MongoDB
	// =======================

	// Insert 是你在 data 层封装的数据库操作
	err = app.Models.LogEntry.Insert(event)
	if err != nil {
		// 如果写入失败，返回统一 JSON 错误响应
		app.errorJSON(w, err)
		return
	}

	// =======================
	// 4. 返回成功响应
	// =======================

	resp := jsonResponse{
		Error:   false,    // 没有错误
		Message: "logged", // 告诉调用方：日志已记录
	}

	// 返回 HTTP 202（Accepted）
	app.writeJSON(w, http.StatusAccepted, resp)
}
