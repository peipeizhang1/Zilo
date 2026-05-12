package svc

import (
	"zilo/backend/api/internal/config"
	"zilo/backend/rpc/workflow/workflowrpc"

	"github.com/zeromicro/go-zero/zrpc"
)

// ServiceContext API 网关只持有 RPC 客户端，业务与 DB 交互全部下沉到 RPC 服务
type ServiceContext struct {
	Config      config.Config
	WorkflowRpc workflowrpc.WorkflowRpc
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:      c,
		WorkflowRpc: workflowrpc.NewWorkflowRpc(zrpc.MustNewClient(c.WorkflowRpc)),
	}
}
