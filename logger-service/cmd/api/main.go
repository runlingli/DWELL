package main

import (
	"context"          // Go 官方的上下文包：用于控制“超时 / 取消 / 生命周期”
	"fmt"              // 字符串格式化，例如 fmt.Sprintf
	"log"              // 标准日志输出
	"log-service/data" // 你自己的数据层（Mongo 的封装）
	"net/http"         // HTTP 服务器
	"time"             // 时间相关（超时、sleep 等）

	"go.mongodb.org/mongo-driver/mongo"         // MongoDB 官方 Go 驱动
	"go.mongodb.org/mongo-driver/mongo/options" // Mongo 连接配置
)

// =======================
// 常量定义
// =======================

const (
	webPort  = "80"                    // HTTP 服务端口
	rpcPort  = "5001"                  // 预留：RPC 服务端口（当前未使用）
	mongoURL = "mongodb://mongo:27017" // MongoDB 地址（Docker service name）
	gRpcPort = "50001"                 // 预留：gRPC 服务端口
)

// =======================
// 全局 Mongo 客户端
// =======================

// client 是 MongoDB 的连接客户端
// 之所以定义为全局变量：
// - 整个服务生命周期内只需要一个连接池
// - Models 内部会复用它
var client *mongo.Client

// =======================
// 应用配置结构体
// =======================

// Config 是整个服务的“依赖容器”
// 你后面所有 handler 都通过 app *Config 访问共享资源
type Config struct {
	Models data.Models // 数据访问层（对 Mongo 的封装）
}

// =======================
// main：程序入口
// =======================

func main() {

	// context.WithTimeout 的作用：
	// 给后续操作一个“最长存活时间”
	// 超过 15 秒会自动取消
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	// defer cancel()：
	// - main 结束时，通知所有使用该 ctx 的操作“该结束了”
	defer cancel()
	// =======================
	// 连接 MongoDB
	// =======================

	mongoClient, err := connectToMongo(ctx)
	if err != nil {
		// 如果数据库连不上，服务没有存在意义，直接崩
		log.Panic(err)
	}

	// 保存到全局变量
	client = mongoClient

	// =======================
	// 确保 Mongo 连接被关闭
	// =======================

	// defer + 匿名函数：
	// main 退出时自动执行
	// 优雅断开 Mongo 连接
	//但最多只等 ctx 允许的时间。
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			// 如果断开失败，说明资源泄漏风险
			panic(err)
		}
	}()

	// =======================
	// 构建应用配置
	// =======================

	// data.New(client)：
	// - 把 Mongo client 注入数据层
	// - 返回 Models（Repository 层）
	app := Config{
		Models: data.New(client),
	}

	// =======================
	// 启动 HTTP 服务
	// =======================

	log.Println("Starting service on port", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort), // 监听端口
		Handler: app.routes(),                // 路由入口
	}

	// 阻塞运行 HTTP 服务
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic()
	}
}

// =======================
// Mongo 连接函数
// =======================

func connectToMongo(ctx context.Context) (*mongo.Client, error) {

	// =======================
	// 1. 构建 Mongo 连接配置
	// =======================

	// options.Client() 创建一个 Mongo 客户端配置对象
	// ApplyURI 指定 Mongo 地址
	clientOptions := options.Client().ApplyURI(mongoURL)

	// 设置认证信息
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// =======================
	// 2. 建立连接
	// =======================

	// mongo.Connect：
	// - 不会立刻报错所有问题
	// - 返回的是一个“连接池客户端”
	c, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Println("Error connecting:", err)
		return nil, err
	}

	log.Println("Connected to mongo!")

	return c, nil
}
