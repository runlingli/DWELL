package main

import (
	"fmt"
	"listener/event" // 自定义包，用于封装 RabbitMQ 消费者逻辑
	"log"
	"math" // 用于计算指数退避的 backoff 时间
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go" // RabbitMQ 官方 Go 客户端
)

func main() {
	// 尝试连接 RabbitMQ
	rabbitConn, err := connect()
	if err != nil {
		// 连接失败，打印错误并退出程序
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close() // 程序结束时关闭连接

	// 程序已连接，开始监听消息
	log.Println("Listening for and consuming RabbitMQ messages...")

	// 创建消费者实例（封装了 queue/binding/exchange 的逻辑）
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err) // 这里直接 panic，因为消费者无法创建，程序无法运行
	}

	// 指定要监听的 routing keys
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err) // 监听过程中出错打印日志
	}
}

// connect 尝试连接 RabbitMQ，并带有指数退避（exponential backoff）fixed->exponential
// 返回一个 amqp.Connection 对象
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
