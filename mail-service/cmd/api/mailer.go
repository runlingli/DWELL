package main

import (
	"bytes"         // 用于在内存中构建模板渲染结果
	"html/template" // HTML 模板引擎（自动转义，防 XSS）
	"log"           // 日志输出
	"time"          // 超时控制

	"github.com/vanng822/go-premailer/premailer" // 用于把 CSS inline 到 HTML（兼容邮件客户端）
	mail "github.com/xhit/go-simple-mail/v2"     // SMTP 客户端库
	// mail为包别名
)

// =======================
// Mail：邮件发送器配置
// =======================

// Mail 描述的是“如何连接并发送邮件”
// 这是一个基础设施层（infrastructure）的配置对象
type Mail struct {
	Domain      string // 邮件域名（某些 SMTP 服务会用）
	Host        string // SMTP 服务器地址
	Port        int    // SMTP 端口
	Username    string // SMTP 用户名
	Password    string // SMTP 密码
	Encryption  string // 加密方式：tls / ssl / none
	FromAddress string // 默认发件人邮箱
	FromName    string // 默认发件人名称
}

// =======================
// Message：业务层邮件对象
// =======================

// Message 描述“一封要发送的邮件”
// 注意：它不关心 SMTP，只关心邮件内容
type Message struct {
	From        string         // 发件人邮箱
	FromName    string         // 发件人名称
	To          string         // 收件人邮箱
	Subject     string         // 邮件标题
	Attachments []string       // 附件路径列表
	Data        any            // 原始业务数据（通常是字符串或结构体）
	DataMap     map[string]any // 专供模板渲染使用的数据
}

//所有字段都会有一个零值，如果没有设置就是零值

// =====================================
// SendSMTPMessage：对外的“发送邮件”接口
// =====================================

func (m *Mail) SendSMTPMessage(msg Message) error {

	msg.From = m.FromAddress
	msg.FromName = m.FromName

	// 构造模板渲染所需的数据
	// 这里统一放入 DataMap，避免模板直接依赖 Message 结构
	data := map[string]any{
		"message": msg.Data,
	}

	msg.DataMap = data

	// 构建 HTML 邮件内容（模板 + CSS inline）
	formattedMessage, err := m.buildHTMLMessage(msg)
	if err != nil {
		return err
	}

	// 构建纯文本邮件内容（用于不支持 HTML 的客户端）
	plainMessage, err := m.buildPlainTextMessage(msg)
	if err != nil {
		return err
	}

	// =======================
	// SMTP 客户端配置
	// =======================

	server := mail.NewSMTPClient()
	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = m.getEncryption(m.Encryption)

	// 不保持长连接（适合微服务场景）
	server.KeepAlive = false

	// 连接和发送超时，防止请求卡死
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	// 建立 SMTP 连接
	smtpClient, err := server.Connect()
	if err != nil {
		log.Println(err)
		return err
	}

	// =======================
	// 构建邮件内容
	// =======================

	email := mail.NewMSG()
	email.SetFrom(msg.From).
		AddTo(msg.To).
		SetSubject(msg.Subject)

	// 设置纯文本正文（fallback）
	email.SetBody(mail.TextPlain, plainMessage)

	// 添加 HTML 正文（multipart/alternative）
	email.AddAlternative(mail.TextHTML, formattedMessage)

	// 添加附件（如果有）
	if len(msg.Attachments) > 0 {
		for _, x := range msg.Attachments {
			email.AddAttachment(x)
		}
	}

	// 发送邮件
	err = email.Send(smtpClient)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// =====================================
// buildHTMLMessage：生成 HTML 邮件内容
// =====================================

func (m *Mail) buildHTMLMessage(msg Message) (string, error) {

	// HTML 邮件模板路径
	templateToRender := "./templates/mail.html.gohtml"

	// 解析模板文件
	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	// 使用 bytes.Buffer 接收渲染结果
	var tpl bytes.Buffer

	// 执行模板，传入 DataMap
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	// 获取渲染后的 HTML
	formattedMessage := tpl.String()

	// 将 <style> 中的 CSS 转为 inline style
	// 这是为了兼容 Gmail / Outlook 等邮件客户端
	formattedMessage, err = m.inlineCSS(formattedMessage)
	if err != nil {
		return "", err
	}

	return formattedMessage, nil
}

// =====================================
// buildPlainTextMessage：生成纯文本邮件
// =====================================

func (m *Mail) buildPlainTextMessage(msg Message) (string, error) {

	// 纯文本模板路径
	templateToRender := "./templates/mail.plain.gohtml"

	// 解析模板
	t, err := template.New("email-plain").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer

	// 渲染模板
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	plainMessage := tpl.String()

	return plainMessage, nil
}

// =====================================
// inlineCSS：CSS 内联处理
// =====================================

func (m *Mail) inlineCSS(s string) (string, error) {

	// premailer 配置选项
	options := premailer.Options{
		RemoveClasses:     false, // 保留 class
		CssToAttributes:   false, // 不强制转为 HTML 属性
		KeepBangImportant: true,  // 保留 !important
	}

	// 创建 premailer 实例
	prem, err := premailer.NewPremailerFromString(s, &options)
	// 已经拿到一个准备好的处理器，但它还没有做实际的转换工作。
	if err != nil {
		return "", err
	}

	// 执行 CSS inline 转换
	html, err := prem.Transform()
	if err != nil {
		return "", err
	}

	return html, nil
}

// =====================================
// getEncryption：字符串 → SMTP 加密枚举
// =====================================

func (m *Mail) getEncryption(s string) mail.Encryption {
	switch s {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSLTLS
	case "none", "":
		return mail.EncryptionNone
	default:
		// 默认使用 STARTTLS（相对安全）
		return mail.EncryptionSTARTTLS
	}
}
