package engine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cloudwego/eino/compose"
)

// LLMConfig 给执行引擎传入的 LLM 调用配置
type LLMConfig struct {
	Provider     string
	BaseURL      string
	APIKey       string
	DefaultModel string
}

// Engine 工作流执行引擎，基于 Eino 的 compose.Chain 编排
type Engine struct {
	llm LLMConfig
}

func New(llm LLMConfig) *Engine {
	return &Engine{llm: llm}
}

// NodeLog 单个节点执行日志
type NodeLog struct {
	NodeID     string
	NodeType   string
	Status     string
	Input      string
	Output     string
	Error      string
	StartedAt  time.Time
	EndedAt    time.Time
	DurationMs int64
}

// ExecutionResult 整次执行结果
type ExecutionResult struct {
	Status     string
	OutputJSON string
	ErrorMsg   string
	StartedAt  time.Time
	EndedAt    time.Time
	DurationMs int64
	NodeLogs   []NodeLog
}

// flowState 执行过程中的内存状态，跨节点传递
type flowState struct {
	Vars     map[string]string
	NodeLogs []NodeLog
}

// Execute 执行 DSL，返回结构化执行结果
func (e *Engine) Execute(ctx context.Context, dslJSON, inputJSON string) (*ExecutionResult, error) {
	dsl, err := Parse(dslJSON)
	if err != nil {
		return nil, err
	}
	ordered, err := dsl.LinearOrder()
	if err != nil {
		return nil, err
	}

	chain := compose.NewChain[*flowState, *flowState]()

	for i := range ordered {
		node := ordered[i] // capture by value
		chain.AppendLambda(compose.InvokableLambda(func(ctx context.Context, st *flowState) (*flowState, error) {
			return e.runNode(ctx, st, node)
		}))
	}

	runnable, err := chain.Compile(ctx)
	if err != nil {
		return nil, fmt.Errorf("compile chain failed: %w", err)
	}

	initVars, err := parseInputVars(inputJSON)
	if err != nil {
		return nil, err
	}

	state := &flowState{Vars: initVars}
	startedAt := time.Now()

	out, runErr := runnable.Invoke(ctx, state)

	endedAt := time.Now()
	result := &ExecutionResult{
		StartedAt:  startedAt,
		EndedAt:    endedAt,
		DurationMs: endedAt.Sub(startedAt).Milliseconds(),
	}
	if out != nil {
		result.NodeLogs = out.NodeLogs
	}
	if runErr != nil {
		result.Status = "failed"
		result.ErrorMsg = runErr.Error()
		return result, runErr
	}

	finalOutput := ""
	if out != nil {
		finalOutput = out.Vars["__final_output"]
	}
	outBytes, _ := json.Marshal(map[string]string{"result": finalOutput})
	result.Status = "success"
	result.OutputJSON = string(outBytes)
	return result, nil
}

// runNode 执行单个 DSL 节点，并落一条节点日志
func (e *Engine) runNode(ctx context.Context, st *flowState, node Node) (*flowState, error) {
	startedAt := time.Now()
	logEntry := NodeLog{
		NodeID:    node.ID,
		NodeType:  node.Type,
		StartedAt: startedAt,
	}
	inputSnapshot, _ := json.Marshal(st.Vars)
	logEntry.Input = string(inputSnapshot)

	defer func() {
		logEntry.EndedAt = time.Now()
		logEntry.DurationMs = logEntry.EndedAt.Sub(startedAt).Milliseconds()
		st.NodeLogs = append(st.NodeLogs, logEntry)
	}()

	switch node.Type {
	case NodeStart:
		logEntry.Status = "success"
		logEntry.Output = string(inputSnapshot)
		return st, nil

	case NodePromptTemplate:
		template, _ := node.Config["template"].(string)
		if template == "" {
			err := errors.New("prompt_template requires non-empty config.template")
			logEntry.Status = "failed"
			logEntry.Error = err.Error()
			return st, err
		}
		rendered := renderTemplate(template, st.Vars)
		st.Vars["prompt"] = rendered
		logEntry.Status = "success"
		logEntry.Output = rendered
		return st, nil

	case NodeLLM:
		prompt := st.Vars["prompt"]
		if prompt == "" {
			err := errors.New("llm node requires upstream prompt")
			logEntry.Status = "failed"
			logEntry.Error = err.Error()
			return st, err
		}
		modelName, _ := node.Config["model"].(string)
		if modelName == "" {
			modelName = e.llm.DefaultModel
		}
		answer := callLLMStub(modelName, prompt)
		st.Vars["llm_output"] = answer
		st.Vars["__final_output"] = answer
		logEntry.Status = "success"
		logEntry.Output = answer
		return st, nil

	case NodeEnd:
		final := st.Vars["__final_output"]
		if final == "" {
			final = st.Vars["prompt"]
		}
		st.Vars["__final_output"] = final
		logEntry.Status = "success"
		logEntry.Output = final
		return st, nil

	default:
		err := fmt.Errorf("unsupported node type: %s", node.Type)
		logEntry.Status = "failed"
		logEntry.Error = err.Error()
		return st, err
	}
}

// renderTemplate 简单的 {{key}} 占位符替换
func renderTemplate(tpl string, vars map[string]string) string {
	out := tpl
	for k, v := range vars {
		out = strings.ReplaceAll(out, "{{"+k+"}}", v)
		out = strings.ReplaceAll(out, "{{ "+k+" }}", v)
	}
	return out
}

// callLLMStub 不依赖真实云端，用确定性回显，便于先打通链路
// 后续接入真实 OpenAI/通义/豆包 时只替换此处
func callLLMStub(model, prompt string) string {
	return fmt.Sprintf("[stub-llm:%s] %s", model, prompt)
}

// parseInputVars 把外部 inputJson 解析为执行变量表
func parseInputVars(raw string) (map[string]string, error) {
	vars := map[string]string{}
	if strings.TrimSpace(raw) == "" {
		return vars, nil
	}
	var decoded map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &decoded); err != nil {
		return nil, fmt.Errorf("invalid input json: %w", err)
	}
	for k, v := range decoded {
		switch x := v.(type) {
		case string:
			vars[k] = x
		case float64:
			vars[k] = fmt.Sprintf("%v", x)
		case bool:
			vars[k] = fmt.Sprintf("%v", x)
		default:
			b, _ := json.Marshal(x)
			vars[k] = string(b)
		}
	}
	return vars, nil
}
