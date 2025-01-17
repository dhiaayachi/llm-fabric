package store

import (
	agentinfo "github.com/dhiaayachi/llm-fabric/proto/gen/agent_info/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInMemoryStore_Store(t *testing.T) {
	// Arrange
	s := NewInMemoryStore().(*InMemoryStore)
	node := &agentinfo.NodeInfo{}
	agent := &agentinfo.AgentInfo{Id: "agent1"}

	// Act
	err := s.Store(agent, node)

	// Assert
	assert.Nil(t, err)
	assert.Len(t, s.data, 1)       // Check data contains one node
	assert.Len(t, s.data[node], 1) // Check node contains one agent

	// Additional test: Store again with same node and agent
	err = s.Store(agent, node)
	assert.Nil(t, err)
	// Assert data remains unchanged (no duplicates)
	assert.Len(t, s.data, 1)
	assert.Len(t, s.data[node], 1)
}

func TestInMemoryStore_GetAll(t *testing.T) {
	// Arrange
	store := NewInMemoryStore().(*InMemoryStore)
	node := &agentinfo.NodeInfo{}
	agent := &agentinfo.AgentInfo{Id: "agent1"}
	require.NoError(t, store.Store(agent, node))

	// Act
	agentsNodes := store.GetAll()

	// Assert
	assert.Len(t, agentsNodes, 1)
	assert.Equal(t, agentsNodes[0].Node, node)
	assert.Len(t, agentsNodes[0].Agents, 1)
	assert.Equal(t, agentsNodes[0].Agents[0].Id, agent.Id)

	// Additional test: Empty store
	store = NewInMemoryStore().(*InMemoryStore)
	agentsNodes = store.GetAll()
	assert.Empty(t, agentsNodes)
}

func TestNewInMemoryStore(t *testing.T) {
	// Arrange & Act
	store := NewInMemoryStore().(*InMemoryStore)

	// Assert
	assert.NotNil(t, store)
	assert.NotNil(t, store.data)
	assert.Empty(t, store.data)
}
