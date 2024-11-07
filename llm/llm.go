package llm

import (
	"context"
	agentv1 "github.com/dhiaayachi/llm-fabric/proto/gen/agent/v1"
)

type Llm interface {
	// SubmitTask TODO: task need to be changed to proto eventually
	SubmitTask(ctx context.Context, task string, respFormat string) (response string, err error)
	GetCapabilities() []agentv1.Capability // Abilities or features the llm supports
	GetTools() []agentv1.Tool              // Abilities or features the llm supports
}
