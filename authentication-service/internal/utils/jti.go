package utils // 定义 token 包，用于处理与令牌（token）相关的逻辑

import (
	"crypto/rand"     // 导入 crypto/rand，用于生成加密安全级别的随机数
	"encoding/base64" // 导入 base64，用于将二进制数据编码为字符串
	"fmt"
	"math/big"
	"time"
)

const RedisTimeout = 2 * time.Second
const MailTimeout = 5 * time.Second

// GenerateJTI 用于生成一个唯一的 JTI（JWT ID）
// JTI 通常用于标识一个 JWT，防止重放攻击或用于 token 黑名单
func GenerateJTI() (string, error) {
	// 创建一个长度为 32 字节的字节切片
	// 32 字节 = 256 位，随机性和安全性都足够高
	b := make([]byte, 32)

	// 使用加密安全的随机数生成器填充字节切片
	_, err := rand.Read(b)
	if err != nil {
		// 如果随机数生成失败，返回空字符串和错误
		return "", err
	}

	// 使用 URL-safe 的 Base64 编码方式将字节切片转换为字符串
	// RawURLEncoding 不会添加 "=" 填充符，适合在 URL、JWT 中使用
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func GenerateCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000)) // 0~999999
	if err != nil {
		// 如果随机数生成失败，返回空字符串和错误
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil // 保证 6 位
}
