package agent_test

import (
	"net"
	"testing"

	"github.com/dhiaayachi/llm-fabric/agent"
	"github.com/dhiaayachi/llm-fabric/proto/gen/agent_info/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

// Mock server implementation for AgentService
type mockAgentServiceServer struct {
	agent_info.UnimplementedAgentServiceServer
}

func startMockServer(t *testing.T) (string, func()) {
	lis, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err, "failed to start TCP listener")

	grpcServer := grpc.NewServer()
	agent_info.RegisterAgentServiceServer(grpcServer, &mockAgentServiceServer{})

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			panic(err)
		}
	}()

	cleanup := func() {
		grpcServer.Stop()
		_ = lis.Close()
	}
	return lis.Addr().String(), cleanup
}

func TestGetClient(t *testing.T) {
	// Start the mock server
	serverAddr, cleanup := startMockServer(t)
	defer cleanup()

	// Get a client using the GetClient function
	client, err := agent.GetClient(serverAddr)
	require.NoError(t, err, "GetClient failed")
	require.NotNil(t, client, "client should not be nil")

	// Test the Close function
	err = client.Close()
	assert.NoError(t, err, "client.Close() should not return an error")
}

func TestGetClient_InvalidAddress(t *testing.T) {
	// Try to connect to an invalid address
	client, err := agent.GetClient("%%%invalid/")
	require.Error(t, err, "GetClient should return an error")
	require.Nil(t, client, "client should be nil")
}
