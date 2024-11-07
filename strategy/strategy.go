package strategy

import (
	"github.com/dhiaayachi/llm-fabric/llm"
	agentv1 "github.com/dhiaayachi/llm-fabric/proto/gen/agent/v1"
)

type TaskAgent struct {
	Task  string
	Agent *agentv1.Agent
}

type Strategy interface {
	Execute(task string, Agents []*agentv1.Agent, localLLM llm.Llm) []*TaskAgent
	Finalize(responses []string) string
}
