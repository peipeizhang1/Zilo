// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"context"

	"zilo/backend/api/internal/pkg/userctx"
	"zilo/backend/api/internal/svc"
	"zilo/backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"zilo/backend/rpc/workflow/workflowrpc"
)

type MeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MeLogic {
	return &MeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MeLogic) Me() (resp *types.MeResp, err error) {
	userID, err := userctx.GetUserIDFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	out, err := l.svcCtx.WorkflowRpc.Me(l.ctx, &workflowrpc.MeReq{UserId: userID})
	if err != nil {
		return nil, err
	}
	return &types.MeResp{
		UserId:   out.UserId,
		Email:    out.Email,
		Nickname: out.Nickname,
	}, nil
}
