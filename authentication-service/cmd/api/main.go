package main

import (
	"authentication/data" // 包含 Models 和 User 数据模型
	"authentication/internal/mail"
	"authentication/internal/oauth"
	"authentication/internal/store"
	"authentication/internal/token"
	"context"
	"database/sql" // Go 标准库，提供数据库操作能力
	"fmt"          // 用于字符串格式化
	"log"          // 提供日志功能
	"net/http"     // 提供 HTTP 服务能力
	"os"           // 用于读取环境变量
	"time"         // 提供时间相关功能

	"github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	_ "github.com/jackc/pgconn"        // PostgreSQL 驱动依赖
	_ "github.com/jackc/pgx/v4"        // PostgreSQL 驱动依赖
	_ "github.com/jackc/pgx/v4/stdlib" // PostgreSQL 标准库兼容层
	//空白标识符：导入这个包，但不直接在代码里使用它的任何函数、类型或变量
)

const webPort = "80" // 服务监听端口

var counts int64 // 用来记录数据库重试次数

// Config 是整个应用的配置结构体
type Config struct {
	DB           *sql.DB             // 数据库连接池
	RefreshStore *store.RefreshStore //redis 存储 jti
	TokenService *token.Service      //生成/刷新token pairs
	OAuthService *oauth.OauthService // oauth 服务
	MailService  *mail.MailService   // 邮件服务
	Models       data.Models         // 数据模型集合
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// main 函数是应用入口
func main() {
	log.Println("Starting authentication service")

	// 连接数据库
	pgConn := connectToPG()
	if pgConn == nil {
		log.Panic("Can't connect to Postgres!") // 如果连接失败，直接终止程序
	}

	redisClient := connectToRedis()
	var refreshStore *store.RefreshStore
	if redisClient != nil {
		refreshStore = store.NewRefreshStore(redisClient)
		log.Println("redis connected successfully!")
	}

	accessSecret := os.Getenv("ACCESS_SECRET")
	refreshSecret := os.Getenv("REFRESH_SECRET")
	googleRedirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	var GoogleOauthConfig = &oauth2.Config{
		RedirectURL:  googleRedirectURL,
		ClientID:     googleClientID,
		ClientSecret: googleClientSecret,
		Scopes: []string{
			"openid",
			"profile",
			"email",
		},
		Endpoint: google.Endpoint,
	}

	tokenService := token.NewService(refreshStore, accessSecret, refreshSecret)
	oauthService := oauth.NewOauthService(refreshStore, GoogleOauthConfig)
	mailService := mail.NewMailService(refreshStore)

	// 配置应用
	app := Config{
		DB:           pgConn,
		RefreshStore: refreshStore,
		TokenService: tokenService,
		OAuthService: oauthService,
		MailService:  mailService,
		Models:       data.New(pgConn), // 初始化 Models，绑定数据库连接
	}

	// 创建 HTTP 服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort), // 设置监听端口
		Handler: app.routes(),                // 指定路由处理器
	}

	// 启动 HTTP 服务（阻塞调用）
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err) // 启动失败直接终止程序
	}
}

// ========================
// 打开数据库连接
// ========================
func openDB(dsn string) (*sql.DB, error) {
	// sql.Open 不会立即建立连接，而是返回一个连接池
	db, err := sql.Open("pgx", dsn) // pgx 驱动
	if err != nil {
		return nil, err
	}

	// Ping 确认数据库可以连接
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// ========================
// 连接数据库，带重试机制
// ========================
func connectToPG() *sql.DB {
	// DSN（Data Source Name）从环境变量获取
	// 例如："postgres://user:pass@localhost:5432/dbname?sslmode=disable"
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready ...")
			counts++ // 记录重试次数
		} else {
			log.Println("Connected to Postgres!")
			return connection // 成功连接返回
		}

		// 超过 10 次重试就放弃
		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for two seconds....")
		time.Sleep(2 * time.Second) // 等待 2 秒后再次尝试
		continue
	}
}

func connectToRedis() *redis.Client {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "redis:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := rdb.Ping(ctx).Err()
	if err != nil {
		log.Println("Redis not ready yet...")
		return nil
	}

	log.Println("Connected to Redis!")
	return rdb
}
