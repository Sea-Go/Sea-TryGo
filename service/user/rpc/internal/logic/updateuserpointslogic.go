package logic

import (
	"context"

	"sea-try-go/service/user/rpc/internal/model"
	"sea-try-go/service/user/rpc/internal/svc"
	"sea-try-go/service/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UpdateUserPointsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserPointsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserPointsLogic {
	return &UpdateUserPointsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserPointsLogic) UpdateUserPoints(in *pb.UpdateUserPointsReq) (*pb.UpdateUserPointsResp, error) {
	// Step 1: 查询用户
	user, err := l.svcCtx.UserModel.FindOneByUid(l.ctx, in.Uid)
	if err == gorm.ErrRecordNotFound {
		l.Logger.Errorf("查询用户失败: uid=%d, error=%v", in.Uid, err)
		return &pb.UpdateUserPointsResp{
			Success: false,
			Message: "用户不存在",
		}, nil
	}
	if err != nil {
		return nil, err
	}

	// Step 2: 计算新积分
	// TODO  协调积分系统和用户系统的积分的类型 是i32还是i64
	newScore := int32(int64(user.Score) + in.Points)

	// Step 3: 校验积分是否足够（减积分场景）
	if newScore < 0 {
		l.Logger.Infof("积分不足: uid=%d, 当前积分=%d, 变动=%d", in.Uid, user.Score, in.Points)
		return &pb.UpdateUserPointsResp{
			Success:  false,
			NewScore: user.Score,
			Message:  "积分不足",
		}, nil
	}

	// Step 4: 更新积分
	err = l.svcCtx.UserModel.UpdateUserById(l.ctx, in.Uid, &model.User{
		Score: newScore,
	})
	if err != nil {
		l.Logger.Errorf("更新积分失败: uid=%d, error=%v", in.Uid, err)
		return &pb.UpdateUserPointsResp{
			Success: false,
			Message: "更新失败",
		}, err
	}

	l.Logger.Infof("积分更新成功: uid=%d, 原积分=%d, 变动=%d, 新积分=%d",
		in.Uid, user.Score, in.Points, newScore)

	return &pb.UpdateUserPointsResp{
		Success:  true,
		NewScore: newScore,
		Message:  "更新成功",
	}, nil
}
