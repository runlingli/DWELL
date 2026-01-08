package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go" // RabbitMQ 官方 Go 客户端
)

// Consumer 封装了 RabbitMQ 消费者的基本信息
type Consumer struct {
	conn      *amqp.Connection // RabbitMQ 连接对象
	queueName string           // 队列名称（如果是随机队列会在 setup 中生成）
}

// NewConsumer 创建一个新的消费者对象
func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn, // 绑定连接
	}

	// 调用 setup 函数初始化交换机等配置
	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

// setup 初始化 RabbitMQ 消费者需要的资源（channel、exchange 等）
func (consumer *Consumer) setup() error {
	// 从连接创建 channel
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	// 声明交换机
	return declareExchange(channel)
}

// Payload 表示从队列接收到的消息结构
type Payload struct {
	Name string `json:"name"` // 消息类型，例如 "log"、"auth"
	Data string `json:"data"` // 消息内容
}

// Listen 开始监听指定的 topics（routing keys）
func (consumer *Consumer) Listen(topics []string) error {
	// 每次监听都要创建一个 channel
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close() // 函数结束时关闭 channel

	// 声明一个随机队列，用于临时接收消息
	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	// 将队列绑定到交换机，并绑定指定的 topics
	for _, s := range topics {
		err = ch.QueueBind(
			q.Name,       // 队列名
			s,            // routing key
			"logs_topic", // 交换机名
			false,        // 不等待服务器确认
			nil,          // 额外参数
		)

		if err != nil {
			return err
		}
	}

	// 从队列中消费消息
	// 此处只有一个消费者
	messages, err := ch.Consume(
		q.Name, // 队列名
		"",     // 消费者名，空表示自动生成
		true,   // 自动 ack
		false,  // 排他队列
		false,  // no-local，不支持
		false,  // no-wait
		nil,    // 额外参数
	)
	if err != nil {
		return err
	}

	// forever 通道用来阻塞主线程，让 goroutine 持续监听消息
	forever := make(chan bool)

	// 启动 goroutine 处理消息
	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload) // 将 JSON 数据解码到结构体

			// 每条消息开启一个 goroutine 异步处理
			go handlePayload(payload)
		}
	}()

	fmt.Printf("Waiting for message [Exchange, Queue] [logs_topic, %s]\n", q.Name)
	<-forever // 阻塞，保证程序不退出

	return nil
}

// handlePayload 根据 Payload.Name 决定如何处理消息
func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		// 打日志
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}

	case "auth":
		// TODO: 可以在这里处理认证相关消息

	// 可以根据业务需要增加更多 case
	default:
		// 默认情况也记录日志
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	}
}

// logEvent 将消息发送到日志微服务
func logEvent(entry Payload) error {
	// 将结构体编码为 JSON
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log" // 日志微服务 URL

	// 构造 HTTP POST 请求
	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json") // 设置请求头

	client := &http.Client{}

	// 发送请求
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// 如果日志服务返回不是 Accepted (202)，则认为失败
	if response.StatusCode != http.StatusAccepted {
		return err
	}

	return nil
}
