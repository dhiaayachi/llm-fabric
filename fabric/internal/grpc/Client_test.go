package grpc_test

import (
	grpc2 "github.com/dhiaayachi/llm-fabric/fabric/internal/grpc"
	"github.com/dhiaayachi/llm-fabric/proto/gen/agent_external/v1"
	"net"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

// Mock server implementation for AgentService
type mockAgentServiceServer struct {
	agent_external.UnimplementedAgentServiceServer
}

func startMockServer(t *testing.T) (string, func()) {
	lis, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err, "failed to start TCP listener")

	grpcServer := grpc.NewServer()
	agent_external.RegisterAgentServiceServer(grpcServer, &mockAgentServiceServer{})

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

	u := strings.Split(serverAddr, ":")
	require.Len(t, u, 2)

	portNum, err := strconv.Atoi(u[1])
	require.NoError(t, err)
	// Get a client using the GetClient function
	client, err := grpc2.GetClient(u[0], int32(portNum))
	require.NoError(t, err, "GetClient failed")
	require.NotNil(t, client, "client should not be nil")

	// Test the Close function
	err = client.Close()
	assert.NoError(t, err, "client.Close() should not return an error")
}

func TestGetClient_InvalidAddress(t *testing.T) {
	// Try to connect to an invalid address
	client, err := grpc2.GetClient("%%%invalid/", 333)
	require.Error(t, err, "GetClient should return an error")
	require.Nil(t, client, "client should be nil")
}
