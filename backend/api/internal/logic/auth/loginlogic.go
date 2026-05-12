// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"context"

	"zilo/backend/api/internal/svc"
	"zilo/backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"zilo/backend/rpc/workflow/workflowrpc"
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
	out, err := l.svcCtx.WorkflowRpc.Login(l.ctx, &workflowrpc.LoginReq{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}
	return &types.LoginResp{
		AccessToken:  out.AccessToken,
		RefreshToken: out.RefreshToken,
		ExpireAt:     out.ExpireAt,
	}, nil
}
