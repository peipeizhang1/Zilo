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
	"golang.org/x/crypto/bcrypt"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *workflowpb.LoginReq) (*workflowpb.LoginResp, error) {
	email := strings.TrimSpace(strings.ToLower(in.Email))
	password := strings.TrimSpace(in.Password)
	if email == "" || password == "" {
		return nil, errors.New("email and password are required")
	}

	user, err := l.svcCtx.UserModel.FindOneByEmail(l.ctx, email)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}
	if user.Status != 1 || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return nil, errors.New("invalid email or password")
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
	refreshToken, _, err := authx.GenerateToken(
		l.svcCtx.Config.Auth.AccessSecret,
		user.Id,
		user.Email,
		l.svcCtx.Config.Auth.RefreshExpire,
	)
	if err != nil {
		return nil, err
	}

	return &workflowpb.LoginResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpireAt:     expireAt,
	}, nil
}
