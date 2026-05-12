package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf

	Mysql struct {
		DataSource string
	}

	Auth struct {
		AccessSecret  string
		AccessExpire  int64
		RefreshExpire int64
	}

	LLM struct {
		Provider     string
		BaseURL      string
		APIKey       string
		DefaultModel string
	}
}
