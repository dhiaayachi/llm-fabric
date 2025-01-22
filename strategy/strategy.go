package strategy

import (
	"context"
	agentinfo "github.com/dhiaayachi/llm-fabric/proto/gen/agent_info/v1"
)

type TaskAgent struct {
	Task   string
	Schema any
	Agent  *agentinfo.AgentInfo
	Node   *agentinfo.NodeInfo
}

type Strategy interface {
	Execute(task string, Agents []*agentinfo.AgentsNodeInfo, localLLM LocalLLM) []*TaskAgent
	Finalize(responses []string, localLLM LocalLLM) string
}

type LocalLLM interface {
	SubmitTask(ctx context.Context, task string, schema any) (string, error)
}
