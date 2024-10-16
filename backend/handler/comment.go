package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/kingwrcy/moments/db"
	"github.com/kingwrcy/moments/vo"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type CommentHandler struct {
	base BaseHandler
}

func NewCommentHandler(injector do.Injector) *CommentHandler {
	return &CommentHandler{do.MustInvoke[BaseHandler](injector)}
}

// RemoveComment godoc
//
//	@Tags		Comment
//	@Summary	删除评论
//	@Accept		json
//	@Produce	json
//	@Param		id			query	int		true	"评论ID"
//	@Param		x-api-token	header	string	true	"登录TOKEN"
//	@Success	200
//	@Router		/api/comment/remove [post]
func (c CommentHandler) RemoveComment(ctx echo.Context) error {
	context := ctx.(CustomContext)
	currentUser := context.CurrentUser()
	id, err := strconv.Atoi(ctx.QueryParam("id"))
	if err != nil {
		return FailResp(ctx, ParamError)
	}
	var (
		comment db.Comment
		memo    db.Memo
	)
	if err = c.base.db.First(&comment, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return FailResp(ctx, ParamError)
	}
	if err = c.base.db.First(&memo, comment.MemoId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return FailResp(ctx, ParamError)
	}

	if currentUser.Id != memo.UserId && currentUser.Id != 1 {
		return FailRespWithMsg(ctx, Fail, "没有权限")
	}
	if c.base.db.Delete(&comment).RowsAffected != 1 {
		return FailRespWithMsg(ctx, Fail, "删除失败")
	}
	return SuccessResp(ctx, h{})
}

func checkGoogleRecaptcha(logger zerolog.Logger, sysConfigVO vo.FullSysConfigVO, token string) error {
	if sysConfigVO.EnableGoogleRecaptcha {
		if token == "" {
			return errors.New("token必填")
		}
		params := url.Values{}
		params.Set("secret", sysConfigVO.GoogleSecretKey)
		params.Set("response", token)

		response, err := http.Post("https://recaptcha.net/recaptcha/api/siteverify?"+params.Encode(), "", nil)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			return errors.New("google验证服务无法正常返回")
		}
		resp, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}
		logger.Info().Str("Action", "评论").Msgf("google resp: %s", resp)

		var result map[string]interface{}
		err = json.Unmarshal(resp, &result)
		if err != nil {
			return err
		}
		if success, ok := result["success"].(bool); ok {
			if success {
				if score, ok := result["score"].(float64); ok {
					if score > 0.5 {
						return nil
					}
				}
			}
		}
		return errors.New("人机校验不通过")
	}
	return nil
}

// AddComment godoc
//
//	@Tags		Comment
//	@Summary	添加评论
//	@Accept		json
//	@Produce	json
//	@Param		object	body	vo.AddCommentReq	true	"添加评论"
//	@Success	200
//	@Router		/api/comment/add [post]
func (c CommentHandler) AddComment(ctx echo.Context) error {
	var (
		req         vo.AddCommentReq
		comment     db.Comment
		now         = time.Now()
		sysConfig   db.SysConfig
		sysConfigVO vo.FullSysConfigVO
	)
	err := ctx.Bind(&req)
	if err != nil {
		c.base.log.Error().Msgf("发表评论时参数校验失败,原因:%s", err)
		return FailResp(ctx, ParamError)
	}
	c.base.db.First(&sysConfig)
	_ = json.Unmarshal([]byte(sysConfig.Content), &sysConfigVO)

	if !sysConfigVO.EnableComment {
		return FailRespWithMsg(ctx, Fail, "评论未开启")
	}

	if err := checkGoogleRecaptcha(c.base.log, sysConfigVO, req.Token); err != nil {
		return FailRespWithMsg(ctx, Fail, err.Error())
	}
	if context, ok := ctx.(CustomContext); ok {
		currentUser := context.CurrentUser()
		if currentUser == nil {
			comment.Username = req.Username
		} else {
			comment.Username = currentUser.Nickname
			comment.Author = fmt.Sprintf("%d", currentUser.Id)
		}
	}

	comment.Content = req.Content
	comment.Email = req.Email
	comment.CreatedAt = &now
	comment.UpdatedAt = &now
	comment.ReplyTo = req.ReplyTo
	comment.Website = req.Website
	comment.MemoId = req.MemoID

	if err = c.base.db.Save(&comment).Error; err == nil {
		// 调用 handleNewComment 函数
		handleNewComment(comment)
		return SuccessResp(ctx, h{})
	}
	return FailRespWithMsg(ctx, Fail, "发表评论失败")
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

func handleNewComment(comment db.Comment) {
	webhookURL := os.Getenv("FEISHU_WEBHOOK_URL")
	if webhookURL == "" {
		log.Println("FEISHU_WEBHOOK_URL environment variable is not set")
		return
	}

	message := fmt.Sprintf("你的Memo ID: %d 下有新的评论:\n %s给 %s 回复: %s",
		comment.MemoId, comment.Username, comment.ReplyTo, comment.Content)

	err := SendFeishuWebhook(webhookURL, message)
	if err != nil {
		log.Printf("Failed to send webhook: %v", err)
	}
}

// FeishuWebhookPayload 是飞书 Webhook 的请求体结构
type FeishuWebhookPayload struct {
	MsgType string `json:"msg_type"`
	Content struct {
		Text string `json:"text"`
	} `json:"content"`
}
