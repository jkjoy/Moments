//go:build !prod

package main

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kingwrcy/moments/db"
	_ "github.com/kingwrcy/moments/docs"
	"github.com/kingwrcy/moments/handler"
	"github.com/kingwrcy/moments/log"
	"github.com/kingwrcy/moments/middleware"
	"github.com/kingwrcy/moments/vo"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/samber/do/v2"
	_ "github.com/swaggo/echo-swagger"
	"gorm.io/gorm"
)

var gitCommitID string

func newEchoEngine(_ do.Injector) (*echo.Echo, error) {
	e := echo.New()
	return e, nil
}

// @title		Moments API
// @version	0.2.1
func main() {

	injector := do.New()
	var cfg vo.AppConfig

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		fmt.Printf("读取配置文件异常:%s", err)
		return
	}

	do.ProvideValue(injector, &cfg)
	do.Provide(injector, log.NewLogger)

	myLogger := do.MustInvoke[zerolog.Logger](injector)
	if gitCommitID != "" {
		myLogger.Info().Msgf("git commit id = %s", gitCommitID)
	}
	handleEmptyConfig(myLogger, &cfg)

	do.Provide(injector, db.NewDB)
	do.Provide(injector, newEchoEngine)
	do.Provide(injector, handler.NewBaseHandler)

	tx := do.MustInvoke[*gorm.DB](injector)

	e := do.MustInvoke[*echo.Echo](injector)
	e.Use(middleware.Auth(injector))

	setupRouter(injector)

	migrateTo3(tx, myLogger)
	myLogger.Info().Msgf("服务端启动成功,监听:%d端口...", cfg.Port)
	e.HideBanner = true
	err = e.Start(fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		myLogger.Fatal().Msgf("服务启动失败,错误原因:%s", err)
	}
}

// FeishuWebhookPayload 是飞书 Webhook 的请求体结构
type FeishuWebhookPayload struct {
	MsgType string `json:"msg_type"`
	Content struct {
		Text string `json:"text"`
	} `json:"content"`
}

// SendFeishuWebhook 发送飞书 Webhook
func SendFeishuWebhook(webhookURL string, message string) error {
	client := resty.New()

	payload := FeishuWebhookPayload{
		MsgType: "text",
		Content: struct {
			Text string `json:"text"`
		}{
			Text: message,
		},
	}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Post(webhookURL)

	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to send webhook, status code: %d", resp.StatusCode())
	}

	log.Println("Webhook sent successfully")
	return nil
}
