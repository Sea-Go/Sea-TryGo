package logic

import (
	"context"
	"fmt"
	"sea-try-go/service/article/rpc/internal/svc"
	"sea-try-go/service/article/rpc/model"
	"sea-try-go/service/article/rpc/pb"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type LikeArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLikeArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeArticleLogic {
	return &LikeArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LikeArticleLogic) LikeArticle(in *__.LikeArticleReq) (*__.LikeArticleResp, error) {
	// todo: add your logic here and delete this line
	// 用于去重：Set 结构，Key=article:liked:文章ID, Value=用户ID
	cacheLikeSetKey := fmt.Sprintf("article:liked:%s", in.ArticleId)
	// 用于计数：Hash 结构或 String 结构，Key=article:count:文章ID
	cacheCountKey := fmt.Sprintf("article:count:%s", in.ArticleId)

	// 2. 区分操作类型 (1: 点赞, 2: 取消)
	switch in.Type {
	case 1: // === 点赞 ===
		// A. 写入数据库 (利用数据库唯一索引做兜底去重)
		newLike := &model.ArticleLikes{
			ArticleId:  in.ArticleId,
			UserId:     in.UserId,
			CreateTime: time.Now(), // model 里如果是 time.Time 类型
		}

		// InsertResult 是 go-zero 生成代码的方法，通常返回 (sql.Result, error)
		_, err := l.svcCtx.ArticleLikesModel.Insert(l.ctx, newLike)
		if err != nil {
			// 如果错误是“唯一键冲突”，说明已经点过赞了
			// 这里的 sqlx.ErrDupEntry 是 go-zero 封装的常见数据库错误
			// 注意：Postgres 的 duplicate error 有时需要特殊判断，但 go-zero 通常会拦截
			return nil, fmt.Errorf("您已经点赞过了或系统繁忙: %v", err)
		}

		// B. 数据库写入成功，更新 Redis 缓存
		// 1. 将用户ID加入 Set
		l.svcCtx.RedisClient.SaddCtx(l.ctx, cacheLikeSetKey, in.UserId)
		// 2. 点赞数 +1 (原子操作)
		l.svcCtx.RedisClient.IncrCtx(l.ctx, cacheCountKey)
		// 3. (可选) 给 Key 设置一个过期时间，防止死数据，比如 1 周
		l.svcCtx.RedisClient.ExpireCtx(l.ctx, cacheCountKey, 3600*24*7)

	case 2: // === 取消点赞 ===
		// A. 删除数据库记录
		// 注意：生成的 Model 默认可能只有 Delete (by ID)。
		// 你可能需要在 articlelikesmodel.go 里手动加一个 DeleteByArticleIdAndUserId。
		// 这里假设你用的是基础 SQL 逻辑或者自定义方法：

		// 临时方案：如果没有自定义 Model 方法，我们这里只能模拟，
		// **你需要看下面的【缺失部分 2】来完善这行代码**
		err := l.svcCtx.ArticleLikesModel.DeleteByArticleUser(l.ctx, in.ArticleId, in.UserId)
		if err != nil {
			if err == sqlx.ErrNotFound {
				return nil, fmt.Errorf("您还没有点赞，无法取消")
			}
			return nil, err
		}

		// B. 更新 Redis
		l.svcCtx.RedisClient.SremCtx(l.ctx, cacheLikeSetKey, in.UserId)
		l.svcCtx.RedisClient.DecrCtx(l.ctx, cacheCountKey)

	default:
		return nil, fmt.Errorf("不支持的操作类型: %d", in.Type)
	}

	// 3. 获取最新点赞数返回
	// 直接读 Redis (速度快)
	val, err := l.svcCtx.RedisClient.GetCtx(l.ctx, cacheCountKey)
	if err != nil {
		// 如果 Redis 读失败，可以降级去读数据库 Count，或者暂时返回 0
		return &pb.LikeArticleResponse{LikeCount: 0}, nil
	}

	// 简单的 string 转 int32 转换
	var count int32
	fmt.Sscanf(val, "%d", &count)

	return &pb.LikeArticleResponse{
		LikeCount: count,
	}, nil
}
