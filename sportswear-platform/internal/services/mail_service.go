package services

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"

	"sportswear-platform/internal/utils"
)

// MailConfig 邮件配置
type MailConfig struct {
	SMTPHost  string
	SMTPPort  string
	Username  string
	Password  string
	FromName  string
	FromEmail string
}

// MailService 邮件服务
type MailService struct {
	config  *MailConfig
	enabled bool
}

// NewMailService 创建邮件服务
func NewMailService() *MailService {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	fromName := os.Getenv("SMTP_FROM_NAME")
	fromEmail := os.Getenv("SMTP_FROM_EMAIL")

	enabled := host != "" && username != "" && password != ""

	if fromName == "" {
		fromName = "Sportswear Platform"
	}
	if fromEmail == "" {
		fromEmail = username
	}
	if fromEmail == "" {
		fromEmail = "noreply@sportswear-platform.com"
	}

	return &MailService{
		config: &MailConfig{
			SMTPHost:  host,
			SMTPPort:  port,
			Username:  username,
			Password:  password,
			FromName:  fromName,
			FromEmail: fromEmail,
		},
		enabled: enabled,
	}
}

// IsEnabled 邮件服务是否可用
func (s *MailService) IsEnabled() bool {
	return s.enabled
}

// SendLeadNotification 发送新询盘通知邮件
func (s *MailService) SendLeadNotification(to []string, data *LeadNotificationData) error {
	if !s.enabled {
		utils.Logger.Infow("邮件服务未配置，跳过邮件发送", "lead_name", data.Name)
		return nil
	}

	subject := fmt.Sprintf("新询盘通知 - %s (%s)", data.Name, data.Company)
	body, err := s.renderTemplate("lead_notification", data)
	if err != nil {
		return err
	}

	return s.send(to, subject, body)
}

// LeadNotificationData 询盘通知模板数据
type LeadNotificationData struct {
	Name        string
	Email       string
	Phone       string
	Company     string
	Country     string
	Message     string
	Score       int
	ScoreLevel  string
	ProjectType string
	Source      string
	Time        string
	AdminURL    string
}

// renderTemplate 渲染邮件模板
func (s *MailService) renderTemplate(name string, data interface{}) (string, error) {
	tmpl := template.Must(template.New(name).Parse(s.getTemplate(name)))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// getTemplate 获取邮件模板
func (s *MailService) getTemplate(name string) string {
	switch name {
	case "lead_notification":
		return `<!DOCTYPE html>
<html>
<head><meta charset="utf-8"></head>
<body style="font-family: Arial, sans-serif; background: #f5f5f5; padding: 20px;">
<div style="max-width: 600px; margin: 0 auto; background: white; border-radius: 8px; padding: 30px;">
<h2 style="color: #1e3a8a; margin-top: 0;">新询盘通知</h2>
<table style="width: 100%; border-collapse: collapse;">
<tr><td style="padding: 8px 0; color: #666; width: 100px;">姓名</td><td style="padding: 8px 0;"><strong>{{.Name}}</strong></td></tr>
<tr><td style="padding: 8px 0; color: #666;">邮箱</td><td style="padding: 8px 0;"><a href="mailto:{{.Email}}">{{.Email}}</a></td></tr>
{{if .Phone}}<tr><td style="padding: 8px 0; color: #666;">电话</td><td style="padding: 8px 0;">{{.Phone}}</td></tr>{{end}}
{{if .Company}}<tr><td style="padding: 8px 0; color: #666;">公司</td><td style="padding: 8px 0;">{{.Company}}</td></tr>{{end}}
{{if .Country}}<tr><td style="padding: 8px 0; color: #666;">国家</td><td style="padding: 8px 0;">{{.Country}}</td></tr>{{end}}
{{if .ProjectType}}<tr><td style="padding: 8px 0; color: #666;">项目类型</td><td style="padding: 8px 0;">{{.ProjectType}}</td></tr>{{end}}
{{if .Score}}<tr><td style="padding: 8px 0; color: #666;">评分</td><td style="padding: 8px 0;">{{.Score}} ({{.ScoreLevel}})</td></tr>{{end}}
{{if .Message}}<tr><td style="padding: 8px 0; color: #666;">留言</td><td style="padding: 8px 0; font-style: italic;">{{.Message}}</td></tr>{{end}}
</table>
{{if .AdminURL}}<p style="margin-top: 20px;"><a href="{{.AdminURL}}" style="display: inline-block; background: #1e3a8a; color: white; padding: 10px 20px; border-radius: 4px; text-decoration: none;">查看详情 →</a></p>{{end}}
<p style="color: #999; font-size: 12px; margin-top: 20px; border-top: 1px solid #eee; padding-top: 10px;">收到时间: {{.Time}}</p>
</div>
</body>
</html>`
	default:
		return ""
	}
}

// send 发送邮件
func (s *MailService) send(to []string, subject, body string) error {
	addr := fmt.Sprintf("%s:%s", s.config.SMTPHost, s.config.SMTPPort)

	// 构建邮件头
	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromEmail)
	headers["To"] = to[0] // 简化：仅支持单收件人
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	var msg bytes.Buffer
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n")
	msg.WriteString(body)

	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.SMTPHost)
	return smtp.SendMail(addr, auth, s.config.FromEmail, to, msg.Bytes())
}
