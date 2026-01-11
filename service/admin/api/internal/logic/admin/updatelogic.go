// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"
	"encoding/json"
	"errors"

	"sea-try-go/service/admin/api/internal/model"
	"sea-try-go/service/admin/api/internal/svc"
	"sea-try-go/service/admin/api/internal/types"
	"sea-try-go/service/common/cryptx"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.UpdateSelfReq) (resp *types.UpdateSelfResp, err error) {
	userId := l.ctx.Value("userId").(json.Number)
	id, er := userId.Int64()
	if er != nil {
		return nil, er
	}
	updates := make(map[string]interface{})
	if len(req.Username) > 0 {
		updates["username"] = req.Username
	}

	if len(req.Password) > 0 {
		newPassword, e := cryptx.PasswordEncrypt(req.Password)
		if e != nil {
			return nil, e
		}
		updates["password"] = newPassword
	}

	if len(req.Email) > 0 {
		updates["email"] = req.Email
	}
	if req.Extrainfo != nil {
		updates["extra_info"] = req.Extrainfo
	}

	if len(updates) > 0 {
		err = l.svcCtx.DB.Model(&model.Admin{}).Where("id = ?", id).Updates(updates).Error
		if err != nil {
			return nil, errors.New("更新失败:" + err.Error())
		}
	}
	var newAdmin model.Admin
	err = l.svcCtx.DB.Model(&model.Admin{}).Where("id = ?", id).First(&newAdmin).Error
	if err != nil {
		return nil, err
	}

	return &types.UpdateSelfResp{
		Success: true,
		Admin: types.AdminInfo{
			Id:        newAdmin.Id,
			Username:  newAdmin.Username,
			Email:     newAdmin.Email,
			Extrainfo: newAdmin.ExtraInfo,
		},
	}, nil

}
