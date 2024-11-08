package store

import agentinfo "github.com/dhiaayachi/llm-fabric/proto/gen/agent_info/v1"

type Store interface {
	Store(agent *agentinfo.AgentInfo) error
	GetAll() []*agentinfo.AgentInfo
}
