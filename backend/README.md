# Backend Quick Start

## 目录结构（go-zero 标准布局）

```text
backend/
├── api/                       # API 网关服务
│   ├── workflow.api           # API 描述文件（生成 handler/types 的源）
│   ├── workflow.go            # 入口
│   ├── etc/workflow-api.yaml  # API 网关配置
│   └── internal/{config,handler,logic,svc,types,pkg/userctx}
├── rpc/workflow/              # 工作流 RPC 服务（承载所有业务逻辑）
│   ├── workflow.proto         # gRPC IDL
│   ├── workflow.go            # 入口
│   ├── etc/workflow.yaml      # RPC 服务配置（DB / JWT / LLM）
│   ├── workflowpb/, workflowrpc/
│   └── internal/{config,server,svc,logic,auth,engine}
├── model/                     # 共享 DB 模型（goctl model 生成 + 自定义扩展）
│   └── sql/schema.sql
├── deployments/, scripts/, README.md, go.mod, go.sum
```

设计要点：
- `api/` 只做 HTTP 入参/出参 + 调下游 RPC，不直接读写 DB。
- `rpc/workflow/` 持有所有业务逻辑，并通过 `internal/engine`（基于 Eino）执行 DSL，通过 `internal/auth` 签发/校验 JWT。
- `model/` 作为同层级共享包，被 RPC 服务复用；后续若新增 RPC 服务（例如 user-rpc / billing-rpc）也直接复用同一个 model。

## 1) 准备配置

```bash
cp .env.example .env
```

按需修改：
- `api/etc/workflow-api.yaml`：HTTP 端口 + RPC 客户端 endpoint + JWT Secret。
- `rpc/workflow/etc/workflow.yaml`：MySQL DataSource、JWT Secret、LLM Provider/Key。

## 2) 初始化数据库

```bash
mysql -u root -p < model/sql/schema.sql
```

## 3) 重新生成 go-zero API 代码（修改了 workflow.api 之后才需要）

```bash
./scripts/gen-api.sh
go mod tidy
```

## 4) 启动服务（先 RPC 后 API）

```bash
# 业务层
cd rpc/workflow && go run . -f etc/workflow.yaml

# 网关层
cd api && go run . -f etc/workflow-api.yaml
```

## 5) 跑引擎单测（无需外部依赖）

```bash
go test ./rpc/workflow/internal/engine/... -v
```
