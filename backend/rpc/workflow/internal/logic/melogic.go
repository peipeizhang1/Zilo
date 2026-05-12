package logic

import (
	"context"
	"errors"

	"zilo/backend/model"
	"zilo/backend/rpc/workflow/internal/svc"
	"zilo/backend/rpc/workflow/workflowpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type MeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MeLogic {
	return &MeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MeLogic) Me(in *workflowpb.MeReq) (*workflowpb.MeResp, error) {
	if in.UserId <= 0 {
		return nil, errors.New("invalid user id")
	}
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, uint64(in.UserId))
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &workflowpb.MeResp{
		UserId:   int64(user.Id),
		Email:    user.Email,
		Nickname: user.Nickname,
	}, nil
}
