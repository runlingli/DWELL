package main

// 导入标准库中的 net/http
// http 包提供了 HTTP 服务器和客户端的基础能力
import (
	"net/http"

	// chi 是一个轻量级的 HTTP 路由库
	// 用来管理路由（URL -> 处理函数）
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	// cors 用来处理跨域请求（CORS）
	// 当前端和后端不在同一个域名/端口时必须用
	"github.com/go-chi/cors"
)

// routes 是 Config 结构体的一个方法
// 它负责：
// 1. 创建路由器
// 2. 注册中间件
// 3. 绑定路由
// 4. 最终返回一个 http.Handler
func (app *Config) routes() http.Handler {
	// 创建一个新的 chi 路由器
	// mux 可以理解为“路由总管”
	mux := chi.NewRouter()

	// 使用 CORS 中间件
	// 作用：允许浏览器从其他域名访问你的后端 API
	mux.Use(cors.Handler(cors.Options{
		// 允许哪些来源访问你的后端
		// https://* 和 http://* 表示：允许所有 http / https 域名
		AllowedOrigins: []string{"https://*", "http://*"},

		// 允许哪些 HTTP 方法
		// GET    : 获取数据
		// POST   : 提交数据
		// PUT    : 更新数据
		// DELETE : 删除数据
		// OPTIONS: 浏览器在跨域时自动发送的“预检请求”
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},

		// 允许前端在请求中携带哪些 HTTP 头
		AllowedHeaders: []string{
			"Accept",        // 告诉服务器：客户端希望接收什么格式的数据（如 application/json）
			"Authorization", // 用来携带认证信息（如 JWT Token）
			"Content-Type",  // 告诉服务器：请求体的数据格式
			"X-CSRF-TOKEN",  // 防止 CSRF 攻击的安全令牌
		},

		// 允许前端“读取”的响应头
		// 默认情况下浏览器只能读到少量响应头
		ExposedHeaders: []string{"Link"},

		// 是否允许携带 Cookie / Authorization 等凭证
		// 如果你使用 session 或需要登录状态，这个通常要 true
		AllowCredentials: true,

		// 预检请求（OPTIONS）的缓存时间（单位：秒）
		// 300 秒内，浏览器不会重复发送 OPTIONS 请求
		MaxAge: 300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))

	mux.Post("/authenticate", app.Authenticate)

	mux.HandleFunc("/authenticate/google", app.GoogleLoginHandler)

	mux.HandleFunc("/oauth/google/callback", app.GoogleCallbackHandler)

	//mux.Post("/refresh", app.Refresh)

	mux.Post("/register", app.Register)

	mux.Post("/verify-code", app.VerifyEmail)

	mux.Get("/resource/profile", app.Profile)

	// 返回配置完成的路由器
	// mux 实现了 http.Handler 接口
	return mux
}
