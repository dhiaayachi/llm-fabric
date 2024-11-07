package store

import agentv1 "github.com/dhiaayachi/llm-fabric/proto/gen/agent/v1"

type InMemoryStore struct {
	data map[string]*agentv1.Agent
}

func (im InMemoryStore) Store(agent *agentv1.Agent) error {
	im.data[agent.Id] = agent
	return nil
}

func (im InMemoryStore) GetAll() []*agentv1.Agent {
	agents := make([]*agentv1.Agent, len(im.data))
	i := 0
	for _, a := range im.data {
		agents[i] = a
		i++
	}
	return agents
}

func NewInMemoryStore() Store {
	data := make(map[string]*agentv1.Agent, 0)
	return &InMemoryStore{data: data}
}
