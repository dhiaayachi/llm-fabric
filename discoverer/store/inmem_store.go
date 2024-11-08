package store

import agentinfo "github.com/dhiaayachi/llm-fabric/proto/gen/agent_info/v1"

type InMemoryStore struct {
	data map[string]*agentinfo.AgentInfo
}

func (im InMemoryStore) Store(agent *agentinfo.AgentInfo) error {
	im.data[agent.Id] = agent
	return nil
}

func (im InMemoryStore) GetAll() []*agentinfo.AgentInfo {
	agents := make([]*agentinfo.AgentInfo, len(im.data))
	i := 0
	for _, a := range im.data {
		agents[i] = a
		i++
	}
	return agents
}

func NewInMemoryStore() Store {
	data := make(map[string]*agentinfo.AgentInfo, 0)
	return &InMemoryStore{data: data}
}
