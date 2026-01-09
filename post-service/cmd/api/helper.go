package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// jsonResponse 用来定义接口返回给前端的 JSON 结构
// 这是一个“统一响应格式”，方便前端处理
type jsonResponse struct {
	// Error 首字母大写，表示这是一个“导出字段”
	// 只有导出字段才能被 json.Marshal 访问到
	Error bool `json:"error"`

	// Message 用来给前端返回一段可读的信息
	Message string `json:"message"`

	// Data 用来承载任意类型的数据
	// any 是 interface{} 的别名（Go 1.18+）
	// omitempty 表示：如果 Data 是 nil，就不会出现在 JSON 中
	Data any `json:"data,omitempty"`
}

// readJSON tries to read the body of a request and converts it into JSON
// readJSON 的作用是：
// 从 HTTP 请求体中读取 JSON，并解析到 data 中
func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	// 限制请求体的最大大小（1 MB）
	// 防止客户端发送超大 body，导致内存被耗尽（一种常见攻击方式）
	maxBytes := 1048576 // one megabyte

	// 使用 MaxBytesReader 包装原始的 r.Body
	// 如果请求体超过 maxBytes，会直接报错
	// w传入：用于在读取超限时，自动返回 413（Request Entity Too Large）
	// 第三个参数格式必须为int64
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// 创建一个 JSON 解码器
	// Decoder 是“流式”的，适合从请求体中逐步读取数据
	dec := json.NewDecoder(r.Body)

	// 将 JSON 解码到 data 中
	// data 通常是一个结构体指针
	err := dec.Decode(data)
	if err != nil {
		// 如果 JSON 格式错误、字段不匹配等，会在这里返回
		return err
	}

	// 再尝试解码一次
	// 目的：确保请求体中“只有一个 JSON 对象”
	// 要求剩余的是&struct{}{}空的
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		// 如果不是 EOF，说明 body 里还有多余内容
		// 这是非法的 JSON 请求
		return errors.New("body must have only a single JSON value")
	}

	return nil
}

// writeJSON takes a response status code and arbitrary data
// and writes a json response to the client
// writeJSON 的作用是：
// 把任意 Go 数据转换成 JSON，并写入 HTTP 响应
func (app *Config) writeJSON(
	w http.ResponseWriter,
	status int,
	data any,
	headers ...http.Header, // 可选参数，用来额外设置响应头
) error {

	// 将 Go 数据序列化为 JSON
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// 如果调用方传入了额外的 header
	// 就把它们写入响应头
	if len(headers) > 0 {
		for key, value := range headers[0] {
			// Header()返回type Header map[string][]string
			w.Header()[key] = value
		}
	}

	// 明确告诉客户端：返回的是 JSON
	w.Header().Set("Content-Type", "application/json")

	// 写入 HTTP 状态码
	w.WriteHeader(status)

	// 写入 JSON 响应体
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

// errorJSON takes an error, and optionally a response status code,
// and generates and sends a json error response
// errorJSON 的作用是：
// 把 Go 的 error 转换成统一格式的 JSON 错误响应
func (app *Config) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	// 默认使用 400 Bad Request
	statusCode := http.StatusBadRequest

	// 如果调用方传入了状态码，就使用它
	if len(status) > 0 {
		statusCode = status[0]
	}

	// 构造统一的错误响应结构
	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()

	// 复用 writeJSON，把错误返回给客户端
	return app.writeJSON(w, statusCode, payload)
}
