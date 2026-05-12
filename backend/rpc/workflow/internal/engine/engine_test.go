package engine

import (
	"context"
	"strings"
	"testing"
)

func TestEngine_Execute_LinearChain(t *testing.T) {
	dsl := `{
		"meta": {"name": "demo", "version": "draft", "mode": "workflow"},
		"nodes": [
			{"id": "start_1", "type": "start", "config": {}},
			{"id": "prompt_1", "type": "prompt_template", "config": {"template": "你好 {{name}}, 请介绍一下 {{topic}}"}},
			{"id": "llm_1", "type": "llm", "config": {"model": "gpt-4o-mini"}},
			{"id": "end_1", "type": "end", "config": {}}
		],
		"edges": [
			{"from": "start_1", "to": "prompt_1"},
			{"from": "prompt_1", "to": "llm_1"},
			{"from": "llm_1", "to": "end_1"}
		]
	}`

	input := `{"name": "Zilo", "topic": "工作流引擎"}`

	eng := New(LLMConfig{DefaultModel: "test-model"})
	result, err := eng.Execute(context.Background(), dsl, input)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if result.Status != "success" {
		t.Fatalf("expected status=success, got=%s err=%s", result.Status, result.ErrorMsg)
	}
	if len(result.NodeLogs) != 4 {
		t.Fatalf("expected 4 node logs, got %d", len(result.NodeLogs))
	}
	if !strings.Contains(result.OutputJSON, "Zilo") || !strings.Contains(result.OutputJSON, "工作流引擎") {
		t.Fatalf("expected output to contain template variables, got=%s", result.OutputJSON)
	}
	if !strings.Contains(result.OutputJSON, "stub-llm") {
		t.Fatalf("expected stub LLM marker in output, got=%s", result.OutputJSON)
	}
}

func TestEngine_Execute_DSLValidation(t *testing.T) {
	cases := []struct {
		name string
		dsl  string
	}{
		{"empty", ""},
		{"missing_nodes", `{"meta":{},"edges":[]}`},
		{"unknown_node_type", `{"meta":{},"nodes":[{"id":"n1","type":"unknown"}],"edges":[]}`},
	}
	eng := New(LLMConfig{})
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := eng.Execute(context.Background(), c.dsl, "")
			if err == nil {
				t.Fatalf("expected error for case %s", c.name)
			}
		})
	}
}
