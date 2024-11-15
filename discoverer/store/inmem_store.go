package store

import agentinfo "github.com/dhiaayachi/llm-fabric/proto/gen/agent_info/v1"

type dataType map[*agentinfo.NodeInfo]map[string]*agentinfo.AgentInfo

type InMemoryStore struct {
	data dataType
}

func (im InMemoryStore) Store(agent *agentinfo.AgentInfo, node *agentinfo.NodeInfo) error {
	if _, ok := im.data[node]; !ok {
		im.data[node] = make(map[string]*agentinfo.AgentInfo)
	}
	im.data[node][agent.Id] = agent
	return nil
}

func (im InMemoryStore) GetAll() []*agentinfo.AgentsNodeInfo {
	agentsNodes := make([]*agentinfo.AgentsNodeInfo, len(im.data))
	i := 0
	for n, agents := range im.data {
		agentsNodes[i] = &agentinfo.AgentsNodeInfo{}
		agentsNodes[i].Node = n
		for _, agent := range agents {
			agentsNodes[i].Agents = append(agentsNodes[i].Agents, agent)
		}
		i++
	}
	return agentsNodes
}

func NewInMemoryStore() Store {
	data := make(dataType)
	return &InMemoryStore{data: data}
}
