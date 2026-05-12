package engine

import (
	"encoding/json"
	"errors"
	"fmt"
)

// 节点类型常量
const (
	NodeStart          = "start"
	NodeEnd            = "end"
	NodePromptTemplate = "prompt_template"
	NodeLLM            = "llm"
)

// DSL 定义：与前端 JSON 格式对齐
type DSL struct {
	Meta  Meta   `json:"meta"`
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

type Meta struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Mode    string `json:"mode"`
}

type Node struct {
	ID     string                 `json:"id"`
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
}

type Edge struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// Parse 解析并校验 DSL
func Parse(raw string) (*DSL, error) {
	if raw == "" {
		return nil, errors.New("dsl is empty")
	}
	var dsl DSL
	if err := json.Unmarshal([]byte(raw), &dsl); err != nil {
		return nil, fmt.Errorf("invalid dsl json: %w", err)
	}
	if len(dsl.Nodes) == 0 {
		return nil, errors.New("dsl has no nodes")
	}
	return &dsl, nil
}

// LinearOrder 把 DSL 拓扑排序为线性节点序列
// V1 仅支持线性 DAG（start -> ... -> end）
func (d *DSL) LinearOrder() ([]Node, error) {
	nodeMap := make(map[string]Node, len(d.Nodes))
	for _, n := range d.Nodes {
		if _, dup := nodeMap[n.ID]; dup {
			return nil, fmt.Errorf("duplicate node id: %s", n.ID)
		}
		nodeMap[n.ID] = n
	}

	indegree := make(map[string]int)
	successor := make(map[string]string)
	for _, e := range d.Edges {
		if _, ok := nodeMap[e.From]; !ok {
			return nil, fmt.Errorf("edge from unknown node: %s", e.From)
		}
		if _, ok := nodeMap[e.To]; !ok {
			return nil, fmt.Errorf("edge to unknown node: %s", e.To)
		}
		if _, exists := successor[e.From]; exists {
			return nil, fmt.Errorf("node %s has multiple successors (only linear flow supported in v1)", e.From)
		}
		successor[e.From] = e.To
		indegree[e.To]++
	}

	var startID string
	for _, n := range d.Nodes {
		if indegree[n.ID] == 0 {
			if startID != "" {
				return nil, errors.New("multiple start nodes detected; expected exactly one")
			}
			startID = n.ID
		}
	}
	if startID == "" {
		return nil, errors.New("no start node found (cycle?)")
	}

	ordered := make([]Node, 0, len(d.Nodes))
	visited := make(map[string]bool)
	cur := startID
	for {
		if visited[cur] {
			return nil, fmt.Errorf("cycle detected at node %s", cur)
		}
		visited[cur] = true
		ordered = append(ordered, nodeMap[cur])
		next, ok := successor[cur]
		if !ok {
			break
		}
		cur = next
	}
	if len(ordered) != len(d.Nodes) {
		return nil, errors.New("graph is not connected; some nodes are unreachable")
	}
	return ordered, nil
}
