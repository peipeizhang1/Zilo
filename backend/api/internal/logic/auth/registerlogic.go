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

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	out, err := l.svcCtx.WorkflowRpc.Register(l.ctx, &workflowrpc.RegisterReq{
		Email:    req.Email,
		Password: req.Password,
		Nickname: req.Nickname,
	})
	if err != nil {
		return nil, err
	}
	return &types.RegisterResp{
		UserId: out.UserId,
		Email:  out.Email,
	}, nil
}
