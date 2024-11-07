package discoverer

import (
	"context"
	agentv1 "github.com/dhiaayachi/llm-fabric/proto/gen/agent/v1"
)

// Discoverer allow agents to discover each other.
type Discoverer interface {
	Join(ctx context.Context, addresses []string, agent *agentv1.Agent) error
	GetAgents() []*agentv1.Agent
}
