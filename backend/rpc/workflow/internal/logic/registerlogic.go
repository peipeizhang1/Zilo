package logic

import (
	"context"
	"errors"
	"strings"

	"zilo/backend/model"
	"zilo/backend/rpc/workflow/internal/svc"
	"zilo/backend/rpc/workflow/workflowpb"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *workflowpb.RegisterReq) (*workflowpb.RegisterResp, error) {
	email := strings.TrimSpace(strings.ToLower(in.Email))
	password := strings.TrimSpace(in.Password)
	nickname := strings.TrimSpace(in.Nickname)
	if email == "" || password == "" {
		return nil, errors.New("email and password are required")
	}
	if len(password) < 8 {
		return nil, errors.New("password must be at least 8 chars")
	}
	if nickname == "" {
		nickname = "new_user"
	}

	_, err := l.svcCtx.UserModel.FindOneByEmail(l.ctx, email)
	if err == nil {
		return nil, errors.New("email already registered")
	}
	if !errors.Is(err, model.ErrNotFound) {
		return nil, err
	}

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	result, err := l.svcCtx.UserModel.Insert(l.ctx, &model.Users{
		Email:        email,
		PasswordHash: string(hashBytes),
		Nickname:     nickname,
		Status:       1,
	})
	if err != nil {
		return nil, err
	}
	userID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &workflowpb.RegisterResp{
		UserId: userID,
		Email:  email,
	}, nil
}
