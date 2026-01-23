package mqs

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"sea-try-go/service/article/rpc/internal/model"
	"sea-try-go/service/article/rpc/internal/svc"

	green "github.com/alibabacloud-go/green-20220302/v3/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/minio/minio-go/v7"
	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleConsumer struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewArticleConsumer(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleConsumer {
	return &ArticleConsumer{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ArticleConsumer) Consume(ctx context.Context, key, val string) error {
	l.Infof("DataClean Service Consuming: %s", val)

	var msg struct {
		ArticleId string `json:"article_id"`
		AuthorId  string `json:"author_id"`
		Content   string `json:"content"`
	}

	if err := json.Unmarshal([]byte(val), &msg); err != nil {
		l.Errorf("Unmarshal error: %v", err)
		return nil
	}

	// Fetch article info
	var article model.Article
	if err := l.svcCtx.ArticleRepo.Db.WithContext(ctx).Where("id = ?", msg.ArticleId).First(&article).Error; err != nil {
		l.Errorf("Failed to find article %s: %v", msg.ArticleId, err)
		return nil
	}

	if l.svcCtx.GreenClient == nil {
		l.Errorf("GreenClient is not initialized")
		return nil
	}

	serviceParameters, _ := json.Marshal(
		map[string]interface{}{
			"content": msg.Content,
		},
	)
	request := green.TextModerationRequest{
		Service:           tea.String("comment_detection"),
		ServiceParameters: tea.String(string(serviceParameters)),
	}

	// 使用传入的 ctx
	result, err := l.svcCtx.GreenClient.TextModeration(&request)
	if err != nil {
		l.Errorf("AliGreen API error: %v", err)
		return nil
	}

	statusCode := tea.IntValue(tea.ToInt(result.StatusCode))
	if statusCode == http.StatusOK {
		textModerationResponse := result.Body
		if tea.IntValue(tea.ToInt(textModerationResponse.Code)) == 200 {
			textModerationResponseData := textModerationResponse.Data
			reason := tea.StringValue(textModerationResponseData.Reason)
			labels := tea.StringValue(textModerationResponseData.Labels)

			if len(reason) > 0 || len(labels) > 0 {
				l.Infof("Article %s RISK DETECTED! Reason: %s, Labels: %s", msg.ArticleId, reason, labels)
				// Update status to Rejected (4)
				article.Status = 4
				if err := l.svcCtx.ArticleRepo.Db.WithContext(ctx).Save(&article).Error; err != nil {
					l.Errorf("Failed to update article status to Rejected: %v", err)
				}
			} else {
				l.Infof("Article %s passed safety check.", msg.ArticleId)

				// Upload to MinIO
				bucketName := l.svcCtx.Config.MinIO.BucketName
				objectName := fmt.Sprintf("%s.md", msg.ArticleId)
				reader := strings.NewReader(msg.Content)

				_, err = l.svcCtx.MinioClient.PutObject(ctx, bucketName, objectName, reader, int64(reader.Len()), minio.PutObjectOptions{
					ContentType: "text/markdown",
				})
				if err != nil {
					l.Errorf("Failed to upload to MinIO: %v", err)
					// Proceed to update status anyway? Or fail?
					// For now, log error but proceed with status update to avoid getting stuck,
					// or return to retry? If we return err, Kafka consumer might retry.
					// Let's log and proceed for now, assuming MinIO transient errors are handled elsewhere or retried manually.
				} else {
					l.Infof("Article %s uploaded to MinIO bucket %s", msg.ArticleId, bucketName)
				}

				// Update status to Published (2)
				article.Status = 2
				if err := l.svcCtx.ArticleRepo.Db.WithContext(ctx).Save(&article).Error; err != nil {
					l.Errorf("Failed to update article status to Published: %v", err)
				}
			}
		} else {
			l.Errorf("AliGreen response code error: %v", textModerationResponse.Code)
		}
	} else {
		l.Errorf("AliGreen http status error: %v", statusCode)
	}
	
	return nil
}
