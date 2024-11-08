package strategy

import (
	"github.com/dhiaayachi/llm-fabric/llm"
	agentinfo "github.com/dhiaayachi/llm-fabric/proto/gen/agent_info/v1"
)

type TaskAgent struct {
	Task  string
	Agent *agentinfo.AgentInfo
}

type Strategy interface {
	Execute(task string, Agents []*agentinfo.AgentInfo, localLLM llm.Llm) []*TaskAgent
	Finalize(responses []string) string
}
