// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"errors"
	"time"

	"sea-try-go/service/common/cryptx"
	"sea-try-go/service/common/jwt"
	"sea-try-go/service/user/api/internal/model"
	"sea-try-go/service/user/api/internal/svc"
	"sea-try-go/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	username := req.Username
	password := req.Password

	user := model.User{}

	//未找到和输入密码错误都显示用户名或密码错误,未找到不能提示找不到用户,否则存在安全隐患

	err = l.svcCtx.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	correct := cryptx.CheckPassword(user.Password, password)
	if correct != true {
		return nil, errors.New("用户名或密码错误")
	}
	now := time.Now().Unix()
	accessSecret := l.svcCtx.Config.Auth.AccessSecret
	accessExpire := l.svcCtx.Config.Auth.AccessExpire

	token, er := jwt.GetToken(accessSecret, now, accessExpire, int64(user.Id))
	if er != nil {
		return nil, er
	}
	return &types.LoginResp{
		Token: token,
	}, nil
}
