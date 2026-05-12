import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import {
  createWorkflow,
  getExecution,
  getExecutionLogs,
  executeWorkflow,
  listWorkflows,
  publishWorkflow,
  saveDraft,
  WorkflowItem
} from "../services/workflow";

export function ConsolePage() {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [items, setItems] = useState<WorkflowItem[]>([]);
  const [creating, setCreating] = useState(false);
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [executeInput, setExecuteInput] = useState('{"name":"Zilo","topic":"工作流引擎"}');
  const [actionMsg, setActionMsg] = useState("");
  const [latestExecutionId, setLatestExecutionId] = useState<number | null>(null);
  const [executionDetail, setExecutionDetail] = useState<{
    status: string;
    durationMs: number;
    outputJson: string;
  } | null>(null);
  const [executionLogs, setExecutionLogs] = useState<
    Array<{ nodeId: string; nodeType: string; status: string; durationMs: number }>
  >([]);

  const fetchList = async () => {
    const accessToken = localStorage.getItem("zilo_access_token");
    if (!accessToken) {
      setError("未登录，请先登录");
      return;
    }
    const res = await listWorkflows();
    setItems(res.list || []);
  };

  useEffect(() => {
    fetchList()
      .catch((err) => setError(err instanceof Error ? err.message : "加载失败"))
      .finally(() => setLoading(false));
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  const handleCreate = async () => {
    if (!name.trim()) {
      setError("工作流名称不能为空");
      return;
    }
    setCreating(true);
    setError("");
    setActionMsg("");
    try {
      const created = await createWorkflow({ name: name.trim(), description: description.trim() });
      const dsl = {
        meta: { name: created.name, version: "draft", mode: "workflow" },
        nodes: [
          { id: "start_1", type: "start", config: {} },
          {
            id: "prompt_1",
            type: "prompt_template",
            config: { template: "你好 {{name}}, 请介绍一下 {{topic}}" }
          },
          { id: "llm_1", type: "llm", config: { model: "gpt-4o-mini" } },
          { id: "end_1", type: "end", config: {} }
        ],
        edges: [
          { from: "start_1", to: "prompt_1" },
          { from: "prompt_1", to: "llm_1" },
          { from: "llm_1", to: "end_1" }
        ]
      };
      await saveDraft(created.id, JSON.stringify(dsl));
      await publishWorkflow(created.id);
      await fetchList();
      setName("");
      setDescription("");
      setActionMsg("创建并发布初始版本成功");
    } catch (err) {
      setError(err instanceof Error ? err.message : "创建失败");
    } finally {
      setCreating(false);
    }
  };

  const handleExecute = async (workflowId: number) => {
    setError("");
    setActionMsg("");
    try {
      const inputPayload = executeInput.trim() || "{}";
      JSON.parse(inputPayload); // 提前校验 JSON 合法性
      const resp = await executeWorkflow(workflowId, inputPayload);
      setLatestExecutionId(resp.executionId);
      setExecutionDetail(null);
      setExecutionLogs([]);
      setActionMsg(`执行成功：executionId=${resp.executionId}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : "执行失败");
    }
  };

  const handleLoadExecution = async () => {
    if (!latestExecutionId) {
      setError("暂无 executionId，请先执行一次");
      return;
    }
    setError("");
    try {
      const [detail, logs] = await Promise.all([getExecution(latestExecutionId), getExecutionLogs(latestExecutionId)]);
      setExecutionDetail({
        status: detail.status,
        durationMs: detail.durationMs,
        outputJson: detail.outputJson
      });
      setExecutionLogs(
        (logs.list || []).map((item) => ({
          nodeId: item.nodeId,
          nodeType: item.nodeType,
          status: item.status,
          durationMs: item.durationMs
        }))
      );
      setActionMsg(`已加载 executionId=${latestExecutionId} 的详情和日志`);
    } catch (err) {
      setError(err instanceof Error ? err.message : "加载执行详情失败");
    }
  };

  return (
    <div style={{ padding: 28 }}>
      <h2 style={{ marginTop: 0 }}>工作流控制台</h2>
      <p style={{ color: "#9fb8ce" }}>登录后的最小可用列表页，下一步会接工作流创建和编辑器。</p>
      <div style={{ marginBottom: 16 }}>
        <Link to="/" className="btn btn-ghost">
          返回首页
        </Link>
      </div>
      <div
        style={{
          border: "1px solid rgba(79,209,255,0.35)",
          borderRadius: 10,
          padding: 12,
          marginBottom: 16,
          background: "rgba(8,18,31,0.65)"
        }}
      >
        <div style={{ fontWeight: 600, marginBottom: 8 }}>新建工作流</div>
        <div style={{ display: "grid", gap: 8, maxWidth: 520 }}>
          <input
            className="auth-input"
            placeholder="名称"
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
          <input
            className="auth-input"
            placeholder="描述"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
          />
          <button className="btn" style={{ width: 160 }} disabled={creating} onClick={handleCreate}>
            {creating ? "创建中..." : "创建并发布"}
          </button>
        </div>
      </div>
      {loading && <div>加载中...</div>}
      {!loading && error && <div className="auth-error">{error}</div>}
      {!loading && !error && actionMsg && <div style={{ color: "#8fffd1", marginBottom: 10 }}>{actionMsg}</div>}
      {!loading && (
        <div
          style={{
            border: "1px solid rgba(79,209,255,0.25)",
            borderRadius: 10,
            padding: 12,
            marginBottom: 14,
            background: "rgba(8,18,31,0.5)"
          }}
        >
          <div style={{ fontWeight: 600, marginBottom: 6 }}>执行参数 (JSON)</div>
          <textarea
            className="auth-input"
            style={{ width: "100%", minHeight: 70, fontFamily: "monospace" }}
            value={executeInput}
            onChange={(e) => setExecuteInput(e.target.value)}
          />
          <div style={{ color: "#9fb8ce", fontSize: 12, marginTop: 6 }}>
            将作为 prompt_template 节点变量注入，例如默认 DSL 中的 {"{{name}}"} 与 {"{{topic}}"}。
          </div>
          <div style={{ marginTop: 10 }}>
            <button className="btn btn-sm" onClick={handleLoadExecution}>
              查看最近执行详情
            </button>
            {latestExecutionId && (
              <span style={{ marginLeft: 10, color: "#8bbdd9" }}>executionId: {latestExecutionId}</span>
            )}
          </div>
        </div>
      )}
      {!loading && executionDetail && (
        <div
          style={{
            border: "1px solid rgba(79,209,255,0.35)",
            borderRadius: 10,
            padding: 12,
            marginBottom: 12,
            background: "rgba(8,18,31,0.65)"
          }}
        >
          <div style={{ fontWeight: 600, marginBottom: 6 }}>执行详情</div>
          <div style={{ color: "#9fb8ce", fontSize: 14 }}>
            状态: {executionDetail.status} | 耗时: {executionDetail.durationMs}ms
          </div>
          <pre style={{ whiteSpace: "pre-wrap", color: "#cbe9ff", fontSize: 12, marginTop: 8 }}>
            {executionDetail.outputJson || "(empty output)"}
          </pre>
        </div>
      )}
      {!loading && executionLogs.length > 0 && (
        <div
          style={{
            border: "1px solid rgba(79,209,255,0.35)",
            borderRadius: 10,
            padding: 12,
            marginBottom: 12,
            background: "rgba(8,18,31,0.65)"
          }}
        >
          <div style={{ fontWeight: 600, marginBottom: 6 }}>节点日志</div>
          {executionLogs.map((item) => (
            <div key={`${item.nodeId}-${item.nodeType}`} style={{ color: "#9fb8ce", fontSize: 13, marginBottom: 4 }}>
              {item.nodeId} ({item.nodeType}) - {item.status} - {item.durationMs}ms
            </div>
          ))}
        </div>
      )}
      {!loading && !error && (
        <div>
          <div style={{ marginBottom: 10 }}>共 {items.length} 条</div>
          <div style={{ display: "grid", gap: 10 }}>
            {items.map((item) => (
              <div
                key={item.id}
                style={{
                  border: "1px solid rgba(79,209,255,0.35)",
                  borderRadius: 10,
                  padding: 12,
                  background: "rgba(8,18,31,0.75)"
                }}
              >
                <div style={{ fontWeight: 600 }}>{item.name}</div>
                <div style={{ color: "#9fb8ce", fontSize: 14 }}>{item.description || "暂无描述"}</div>
                <div style={{ color: "#8bbdd9", fontSize: 12, marginTop: 6 }}>
                  版本: {item.latestVersion} | 更新时间: {item.updatedAt}
                </div>
                <div style={{ marginTop: 8 }}>
                  <button className="btn btn-sm" onClick={() => handleExecute(item.id)}>
                    执行
                  </button>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}

