package logic

import (
	"context"
	"errors"
	"strings"

	"zilo/backend/model"
	authx "zilo/backend/rpc/workflow/internal/auth"
	"zilo/backend/rpc/workflow/internal/svc"
	"zilo/backend/rpc/workflow/workflowpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RefreshTokenLogic) RefreshToken(in *workflowpb.RefreshTokenReq) (*workflowpb.RefreshTokenResp, error) {
	token := strings.TrimSpace(in.RefreshToken)
	if token == "" {
		return nil, errors.New("refreshToken is required")
	}
	claims, err := authx.ParseToken(l.svcCtx.Config.Auth.AccessSecret, token)
	if err != nil {
		return nil, errors.New("invalid refreshToken")
	}
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	accessToken, expireAt, err := authx.GenerateToken(
		l.svcCtx.Config.Auth.AccessSecret,
		user.Id,
		user.Email,
		l.svcCtx.Config.Auth.AccessExpire,
	)
	if err != nil {
		return nil, err
	}

	return &workflowpb.RefreshTokenResp{
		AccessToken: accessToken,
		ExpireAt:    expireAt,
	}, nil
}
