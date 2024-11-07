package store

import agentv1 "github.com/dhiaayachi/llm-fabric/proto/gen/agent/v1"

type Store interface {
	Store(agent *agentv1.Agent) error
	GetAll() []*agentv1.Agent
}
