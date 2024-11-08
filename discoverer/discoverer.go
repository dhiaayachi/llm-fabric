package discoverer

import (
	"context"
	agentinfo "github.com/dhiaayachi/llm-fabric/proto/gen/agent_info/v1"
)

// Discoverer allow agents to discover each other.
type Discoverer interface {
	Join(ctx context.Context, addresses []string, agent *agentinfo.AgentInfo) error
	GetAgents() []*agentinfo.AgentInfo
}
