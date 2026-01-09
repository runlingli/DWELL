package main

import (
	"context"
	"database/sql" // Go 标准库，提供数据库操作能力
	"fmt"          // 用于字符串格式化
	"log"          // 提供日志功能
	"net/http"     // 提供 HTTP 服务能力
	"os"           // 用于读取环境变量
	"post-service/data"
	"time" // 提供时间相关功能

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
	DB     *sql.DB     // 数据库连接
	Models data.Models // 数据模型集合
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	log.Println("Starting post service")

	// 连接数据库
	pgConn := connectToPG()
	if pgConn == nil {
		log.Panic("Can't connect to Postgres!") // 如果连接失败，直接终止程序
	}

	app := Config{
		DB:     pgConn,
		Models: data.New(pgConn), // 初始化 Models，绑定数据库连接
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort), // 设置监听端口
		Handler: app.routes(),                // 指定路由处理器
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err) // 启动失败直接终止程序
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToPG() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready ...")
			counts++
		} else {
			log.Println("Connected to Postgres!")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for two seconds....")
		time.Sleep(2 * time.Second)
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
