package svc

import (
	"zilo/backend/model"
	"zilo/backend/rpc/workflow/internal/config"
	"zilo/backend/rpc/workflow/internal/engine"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config               config.Config
	UserModel            model.UsersModel
	WorkflowModel        model.WorkflowsModel
	WorkflowVersionModel model.WorkflowVersionsModel
	ExecutionModel       model.WorkflowExecutionsModel
	ExecutionLogModel    model.WorkflowExecutionLogsModel
	Engine               *engine.Engine
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:               c,
		UserModel:            model.NewUsersModel(conn),
		WorkflowModel:        model.NewWorkflowsModel(conn),
		WorkflowVersionModel: model.NewWorkflowVersionsModel(conn),
		ExecutionModel:       model.NewWorkflowExecutionsModel(conn),
		ExecutionLogModel:    model.NewWorkflowExecutionLogsModel(conn),
		Engine: engine.New(engine.LLMConfig{
			Provider:     c.LLM.Provider,
			BaseURL:      c.LLM.BaseURL,
			APIKey:       c.LLM.APIKey,
			DefaultModel: c.LLM.DefaultModel,
		}),
	}
}
