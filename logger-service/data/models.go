package data

import (
	"context" // 用于控制数据库操作的超时 / 取消
	"log"     // 日志输出
	"time"    // 时间戳

	"go.mongodb.org/mongo-driver/bson"           // MongoDB 的 BSON 构造工具
	"go.mongodb.org/mongo-driver/bson/primitive" // MongoDB 特有类型（如 ObjectID）
	"go.mongodb.org/mongo-driver/mongo"          // Mongo 客户端核心类型
	"go.mongodb.org/mongo-driver/mongo/options"  // Mongo 查询 / 连接配置
)

// =======================
// Mongo 客户端（全局）
// =======================

// client 是 MongoDB 的客户端连接池
// 整个 data 包都会复用它
var client *mongo.Client

// =======================
// 初始化 data 层
// =======================

// New 用于初始化 Models
// 在 main 中调用：data.New(mongoClient)
func New(mongo *mongo.Client) Models {

	// 把 main 传进来的 Mongo client 保存到包级变量
	client = mongo

	// 返回 Models，供上层使用
	return Models{
		LogEntry: LogEntry{},
	}
}

// =======================
// Models 容器
// =======================

// Models 是 data 层的“统一出口”
// 以后如果你有 User、Audit、Metric，都加在这里
type Models struct {
	LogEntry LogEntry
}

// =======================
// LogEntry 数据模型
// =======================

// LogEntry 表示 MongoDB 中 logs 集合的一条文档
type LogEntry struct {
	ID string `bson:"_id,omitempty" json:"id,omitempty"`
	// Mongo 中的 _id
	// omitempty 表示：插入时如果为空，Mongo 自动生成

	Name string `bson:"name" json:"name"`
	// 日志来源（哪个服务）

	Data string `bson:"data" json:"data"`
	// 日志内容

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	// 创建时间

	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
	// 更新时间
}

// =======================
// 插入一条日志
// =======================

func (l *LogEntry) Insert(entry LogEntry) error {

	// 获取数据库 logs 中的 logs 集合
	collection := client.Database("logs").Collection("logs")

	// InsertOne 插入一条文档
	_, err := collection.InsertOne(
		context.TODO(), // 这里没有设置超时（简单写法）
		LogEntry{
			Name:      entry.Name,
			Data:      entry.Data,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	)

	if err != nil {
		log.Println("Error inserting into logs:", err)
		return err
	}

	return nil
}

// =======================
// 查询所有日志
// =======================

func (l *LogEntry) All() ([]*LogEntry, error) {

	// 创建一个 15 秒超时的 context
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	// 创建查询选项
	opts := options.Find()

	// 设置排序：按 created_at 倒序（最新的在前）
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	// 执行查询
	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println("Finding all docs error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	// 用来存放结果
	var logs []*LogEntry

	// 遍历游标
	for cursor.Next(ctx) {
		var item LogEntry

		// 把 BSON 解码成 Go struct
		err := cursor.Decode(&item)
		if err != nil {
			log.Print("Error decoding log into slice:", err)
			return nil, err
		}

		logs = append(logs, &item)
	}

	return logs, nil
}

// =======================
// 根据 ID 查询一条日志
// =======================

func (l *LogEntry) GetOne(id string) (*LogEntry, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	// Mongo 的 _id 是 ObjectID，不是字符串
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var entry LogEntry

	// 查找单条文档
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

// =======================
// 删除整个集合（危险操作）
// =======================

func (l *LogEntry) DropCollection() error {

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	// 删除整个集合
	if err := collection.Drop(ctx); err != nil {
		return err
	}

	return nil
}

// =======================
// 更新一条日志
// =======================

func (l *LogEntry) Update() (*mongo.UpdateResult, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	// 把字符串 ID 转成 ObjectID
	docID, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil {
		return nil, err
	}

	// UpdateOne：
	// - 第一个参数：过滤条件
	// - 第二个参数：更新内容（$set）
	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": docID},
		bson.D{
			{"$set", bson.D{
				{"name", l.Name},
				{"data", l.Data},
				{"updated_at", time.Now()},
			}},
		},
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}
