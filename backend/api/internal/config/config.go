package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

// Config API 网关配置：只关心 HTTP 服务自身 + 下游 RPC 客户端 + JWT 校验密钥
// DB / Redis / LLM 等下沉到 RPC 服务（rpc/workflow）配置中
type Config struct {
	rest.RestConf

	WorkflowRpc zrpc.RpcClientConf

	Auth struct {
		AccessSecret  string
		AccessExpire  int64
		RefreshExpire int64
	}
}
