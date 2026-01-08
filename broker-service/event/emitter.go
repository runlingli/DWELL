package event

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go" // RabbitMQ 官方 Go 客户端
)

// Emitter 封装了向 RabbitMQ 发布事件的逻辑
type Emitter struct {
	connection *amqp.Connection // RabbitMQ 连接对象
}

// setup 用于初始化 Emitter，例如声明交换机
func (e *Emitter) setup() error {
	// 创建一个 channel
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}

	// defer 确保函数结束时关闭 channel，避免资源泄露
	defer channel.Close()

	// 声明交换机（logs_topic），确保消息可以路由
	// 确保交换机存在，或者先创建交换机。listener和producer拿到同一个交换机。
	return declareExchange(channel)
}

// Push 向指定交换机发送一条消息
// event: 消息内容
// severity: routing key，例如 "log.INFO", "log.ERROR"
func (e *Emitter) Push(event string, severity string) error {
	// 每次发送消息都需要一个 channel
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close() // 发送完毕关闭 channel

	log.Println("Pushing to channel")

	// 发布消息到交换机
	err = channel.Publish(
		"logs_topic", // 交换机名称
		severity,     // routing key
		false,        // mandatory：如果 true 且没有匹配队列会返回消息
		false,        // immediate：如果 true 且没有消费者会返回消息
		amqp.Publishing{
			ContentType: "text/plain",  // 消息类型
			Body:        []byte(event), // 消息内容，需要转成 []byte
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// NewEventEmitter 创建一个新的 Emitter 对象，并执行初始化
func NewEventEmitter(conn *amqp.Connection) (Emitter, error) {
	emitter := Emitter{
		connection: conn, // 绑定连接
	}

	// 执行 setup 初始化交换机等
	err := emitter.setup()
	if err != nil {
		return Emitter{}, err
	}

	return emitter, nil
}
