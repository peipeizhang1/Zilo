import { httpRequest } from "./http";

export type WorkflowItem = {
  id: number;
  name: string;
  description: string;
  latestVersion: number;
  updatedAt: string;
};

export type ListWorkflowResp = {
  total: number;
  list: WorkflowItem[];
};

function authHeader() {
  const token = localStorage.getItem("zilo_access_token") || "";
  return {
    Authorization: `Bearer ${token}`
  };
}

export function listWorkflows() {
  return httpRequest<ListWorkflowResp>("/api/v1/workflows?page=1&pageSize=50", {
    headers: authHeader()
  });
}

export function createWorkflow(payload: { name: string; description?: string }) {
  return httpRequest<WorkflowItem>("/api/v1/workflows", {
    method: "POST",
    headers: authHeader(),
    body: JSON.stringify(payload)
  });
}

export function saveDraft(workflowId: number, dslJson: string) {
  return httpRequest<void>(`/api/v1/workflows/${workflowId}/draft`, {
    method: "POST",
    headers: authHeader(),
    body: JSON.stringify({ dslJson })
  });
}

export function publishWorkflow(workflowId: number) {
  return httpRequest<{ version: number }>(`/api/v1/workflows/${workflowId}/publish`, {
    method: "POST",
    headers: authHeader(),
    body: JSON.stringify({})
  });
}

export function executeWorkflow(workflowId: number, inputJson = "") {
  return httpRequest<{ executionId: number; status: string }>(`/api/v1/workflows/${workflowId}/execute`, {
    method: "POST",
    headers: authHeader(),
    body: JSON.stringify({ inputJson })
  });
}

export function getExecution(executionId: number) {
  return httpRequest<{
    id: number;
    workflowId: number;
    versionId: number;
    status: string;
    inputJson: string;
    outputJson: string;
    errorMsg: string;
    startedAt: string;
    endedAt: string;
    durationMs: number;
  }>(`/api/v1/executions/${executionId}`, {
    headers: authHeader()
  });
}

export function getExecutionLogs(executionId: number) {
  return httpRequest<{
    list: Array<{
      nodeId: string;
      nodeType: string;
      status: string;
      inputJson: string;
      outputJson: string;
      errorMsg: string;
      startedAt: string;
      endedAt: string;
      durationMs: number;
    }>;
  }>(`/api/v1/executions/${executionId}/logs`, {
    headers: authHeader()
  });
}

