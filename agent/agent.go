package agent

import agentv1 "github.com/dhiaayachi/llm-fabric/proto/gen/agent/v1"

type Agent interface {
	// SubmitTask TODO: task need to be changed to proto eventually
	SubmitTask(task string) (response string, err error)
	GetCapabilities() ([]agentv1.Capability, error) // Abilities or features the agent supports
	GetTools() ([]agentv1.Tool, error)              // Abilities or features the agent supports
}
