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

type RefreshTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshTokenLogic) RefreshToken(req *types.RefreshTokenReq) (resp *types.RefreshTokenResp, err error) {
	out, err := l.svcCtx.WorkflowRpc.RefreshToken(l.ctx, &workflowrpc.RefreshTokenReq{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		return nil, err
	}

	return &types.RefreshTokenResp{
		AccessToken: out.AccessToken,
		ExpireAt:    out.ExpireAt,
	}, nil
}
