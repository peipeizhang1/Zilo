# 后端开发文档（V0.1）

## 1. 目标与原则

- 后端目标：提供稳定、可扩展、可观测的工作流定义与执行平台。
- 固定技术栈：Go + go-zero + Eino。
- 架构铁律：严格 Handler -> Logic -> Model 分层，不跨层写业务。

---

## 2. 系统模块

- `workflow-api`：对前端提供 REST 接口（工作流管理/发布/执行/日志）。
- `workflow-engine`：DSL 校验、编译、执行（调用 Eino Workflow/Graph）。
- `workflow-model`：DB 交互层（goctl 生成 + 扩展）。
- `auth`：注册登录、JWT 签发、权限校验。

---

## 3. go-zero 工程结构建议

- `api/*.api`：接口定义。
- `internal/handler`：HTTP 入参绑定、响应封装。
- `internal/logic`：业务编排、事务控制、调用 model/engine。
- `internal/model`：数据库访问层（goctl 生成）。
- `internal/svc`：ServiceContext 注入依赖。
- `internal/types`：api 生成的请求响应结构体。

---

## 4. 数据库设计

## 4.1 用户与鉴权
- `users`
  - `id`, `email`, `password_hash`, `nickname`, `status`, `created_at`, `updated_at`
- `user_sessions`（可选）
  - `id`, `user_id`, `refresh_token_hash`, `expired_at`, `created_at`

## 4.2 工作流领域
- `workflows`
  - `id`, `owner_id`, `name`, `description`, `status`, `latest_version`, `created_at`, `updated_at`
- `workflow_versions`
  - `id`, `workflow_id`, `version`, `dsl_json`, `is_published`, `created_by`, `created_at`
- `workflow_executions`
  - `id`, `workflow_id`, `version_id`, `trigger_type`, `status`, `input_json`, `output_json`, `error_msg`, `started_at`, `ended_at`, `duration_ms`
- `workflow_execution_logs`
  - `id`, `execution_id`, `node_id`, `node_type`, `status`, `input_json`, `output_json`, `error_msg`, `started_at`, `ended_at`, `duration_ms`

## 4.3 索引建议
- `users(email)` 唯一索引。
- `workflow_versions(workflow_id, version)` 唯一索引。
- `workflow_executions(workflow_id, started_at)` 索引。
- `workflow_execution_logs(execution_id, node_id)` 复合索引。

---

## 5. API 设计（V1）

## 5.1 鉴权
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/refresh`
- `POST /api/v1/auth/logout`
- `GET /api/v1/auth/me`

## 5.2 工作流管理
- `POST /api/v1/workflows`
- `GET /api/v1/workflows`
- `GET /api/v1/workflows/{id}`
- `PUT /api/v1/workflows/{id}`
- `DELETE /api/v1/workflows/{id}`

## 5.3 草稿、发布、回滚
- `POST /api/v1/workflows/{id}/draft`
- `POST /api/v1/workflows/{id}/publish`
- `GET /api/v1/workflows/{id}/versions`
- `POST /api/v1/workflows/{id}/rollback/{version}`

## 5.4 执行与日志
- `POST /api/v1/workflows/{id}/execute`
- `POST /api/v1/workflows/{id}/debug`
- `GET /api/v1/executions/{executionId}`
- `GET /api/v1/executions/{executionId}/logs`

---

## 6. 执行引擎设计（Eino）

## 6.1 实际落地（已实现）
- 引擎包：`backend/internal/engine`，由 RPC 服务 `workflow-rpc` 注入使用。
- 调用入口：`Engine.Execute(ctx, dslJson, inputJson) (*ExecutionResult, error)`。
- 执行链路：
  1. `Parse` 解析并校验 DSL JSON。
  2. `LinearOrder` 拓扑排序为线性 DAG（V1 仅支持单出边线性流）。
  3. 用 `compose.NewChain[*flowState, *flowState]` 按顺序追加 Lambda 节点，再 `Compile` 为可运行体。
  4. `Invoke` 跑链路，每个节点写入 `NodeLog`，最终落 `workflow_executions` + `workflow_execution_logs`。
- 当前已支持节点：`start` / `prompt_template`（{{var}} 模板渲染）/ `llm`（V1 走 stub，便于无 Key 也可联调）/ `end`。

## 6.2 与 RPC 集成
- `rpc/workflow/internal/svc/servicecontext.go` 注入 `Engine` 实例。
- `executeworkflowlogic.go`：
  - 创建 `workflow_executions` 头记录（running）。
  - 调引擎执行。
  - 把所有节点日志逐条插入 `workflow_execution_logs`。
  - 用最终 `status / output / error / duration_ms` 写回 execution 终态。

## 6.3 节点类型（V1 路线图）
- 已完成：`start` / `prompt_template` / `llm`(stub) / `end`
- 下一阶段：`http` / `condition` / `parallel`
- 远期：`code_lambda`（沙箱执行）

## 6.4 后续工作
- 把 `llm` stub 替换为真实 OpenAI / 通义 / 豆包 适配器（按 `LLM.Provider` 路由）。
- 支持非线性 DAG：condition 分支、parallel 多出边、汇合节点。
- 节点级 timeout / 重试 / 错误分支。

---

## 7. 安全与稳定性

- JWT 中间件：保护 `/console` 相关接口。
- 密钥保护：供应商 API Key 仅服务端保存（加密存储）。
- 幂等控制：`execution_id` + `business_key` + Redis 锁。
- 超时控制：节点级 timeout + 全流程 deadline。
- 限流熔断：利用 go-zero 内置能力（限流、负载保护、熔断）。

---

## 8. 开发顺序（后端）

### 阶段 1：鉴权最小闭环
- 注册/登录/刷新 token。
- 用户态获取接口。

### 阶段 2：工作流定义最小闭环
- workflow CRUD。
- draft 保存 + 版本发布。

### 阶段 3：执行最小闭环（已完成）
- ✅ start/prompt_template/llm(stub)/end 节点跑通。
- ✅ execution + execution_logs 真实落库。
- ✅ 引擎单测覆盖：DSL 校验 + 真实 chain 执行（`internal/engine/engine_test.go`）。

### 阶段 4：增强能力（进行中）
- 接入真实 LLM Provider（替换 stub）。
- condition/parallel/http 节点。
- 失败重试、回滚（已完成）、审计增强。

---

## 9. 验收标准（Backend DoD）

- 鉴权链路可用（注册/登录/刷新）。
- 工作流可创建、保存草稿、发布版本。
- 已发布版本可执行并产生日志。
- 失败执行可定位具体节点错误。
- 代码符合 go-zero 三层规范，model 层不混业务逻辑。
