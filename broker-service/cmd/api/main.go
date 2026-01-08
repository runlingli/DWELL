package main

import (
	"fmt" // 用来做字符串格式化（例如 fmt.Sprintf）
	"log" // 用来打印日志（比 fmt.Println 更适合服务端程序）
	"math"
	"net/http" // Go 标准库中的 HTTP 服务器实现
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go" // RabbitMQ 官方 Go 客户端
)

// webPort 定义服务监听的端口号
// 这里是字符串类型，是因为后面会和 ":" 拼接
const webPort = "80"

// Config 是一个配置结构体
// 目前是空的，但后续可以在这里放：
// - 数据库连接
// - 配置项
// - 依赖对象
type Config struct {
	Rabbit *amqp.Connection
}

// main 是 Go 程序的入口函数
// 程序从这里开始执行
func main() {
	// 尝试连接 RabbitMQ
	rabbitConn, err := connect()
	if err != nil {
		// 连接失败，打印错误并退出程序
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close() // 程序结束时关闭连接

	// 创建一个 Config 实例
	// app 会作为整个应用的“上下文”
	app := Config{
		Rabbit: rabbitConn,
	}

	// 打印一条启动日志
	// 提示服务正在监听哪个端口
	log.Printf("Starting broker service on port %s\n", webPort)

	// 创建一个 HTTP 服务器实例
	srv := &http.Server{
		// Addr 指定服务器监听的地址
		// ":80" 表示监听本机所有网卡的 80 端口
		Addr: fmt.Sprintf(":%s", webPort),

		// Handler 是请求的处理器
		// app.routes() 返回的是一个 http.Handler
		// 所有进来的 HTTP 请求都会先交给它处理
		Handler: app.routes(),
	}

	// 启动 HTTP 服务器
	// ListenAndServe 是一个阻塞调用
	// 一旦启动成功，程序会一直运行在这里
	err = srv.ListenAndServe()

	// 如果服务器启动失败（比如端口被占用）
	// err 就不为 nil
	if err != nil {
		// Panic 会打印错误并直接终止程序
		log.Panic(err)
	}
}

func connect() (*amqp.Connection, error) {
	var counts int64              // 记录尝试次数
	var backOff = 1 * time.Second // 初始退避时间
	var connection *amqp.Connection

	// 无限循环直到连接成功或超过重试次数
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq") // 默认 guest 用户
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...") // 服务未启动
			counts++
		} else {
			log.Println("Connected to RabbitMQ!") // 成功连接
			connection = c
			break // 跳出循环
		}

		// 超过最大重试次数，返回错误
		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		// 指数退避：等待时间 = 尝试次数的平方秒
		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...") // 打印等待信息
		time.Sleep(backOff)           // 暂停一段时间再重试
		continue                      // 循环下一次尝试
	}

	return connection, nil
}
