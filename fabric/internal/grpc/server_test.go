package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/dhiaayachi/llm-fabric/proto/gen/agent_external/v1"
	"github.com/dhiaayachi/llm-fabric/strategy"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/anypb"
	"testing"

	"github.com/dhiaayachi/llm-fabric/llm"
	agentinfo "github.com/dhiaayachi/llm-fabric/proto/gen/agent_info/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type agentMock struct {
	mock.Mock
}

func (a *agentMock) SubmitTask(ctx context.Context, task string, schema *anypb.Any) (string, error) {
	args := a.Called(ctx, task, schema)
	return args.Get(0).(string), args.Error(1)
}

func (a *agentMock) DispatchTask(ctx context.Context, task string, schema any) (string, error) {
	args := a.Called(ctx, task, schema)
	return args.Get(0).(string), args.Error(1)
}

func (a *agentMock) GetStrategies() []strategy.Strategy {
	args := a.Called()
	return args.Get(0).([]strategy.Strategy)
}

func (a *agentMock) GetAgents() []*agentinfo.AgentsNodeInfo {
	args := a.Called()
	return args.Get(0).([]*agentinfo.AgentsNodeInfo)
}

func (a *agentMock) GetLocalLlm() llm.Llm {
	args := a.Called()
	return args.Get(0).(llm.Llm)
}

func startTestServer(llmMock agent) (*grpc.ClientConn, func(), error) {

	// Start the gRPC server
	s := grpc.NewServer()
	agentServer := NewServer(llmMock, &Config{Logger: logrus.New(), ListenAddr: "localhost:0"})
	agent_external.RegisterAgentServiceServer(s, agentServer)

	ctx, cancel := context.WithCancel(context.Background())
	agentServer.Start(ctx)

	// Retrieve the actual port used by the listener
	serverAddr := agentServer.listener.Addr().String()

	// Connect to the server using gRPC Dial
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		cancel()
		return nil, nil, fmt.Errorf("failed to dial gRPC server: %w", err)
	}
	return conn, cancel, nil
}

func TestSubmitTask_Success(t *testing.T) {
	// Set up the mock and expected behavior
	llmMock := new(agentMock)
	task := "Test Task"
	expectedResponse := "Task Response"

	llmMock.On("SubmitTask", mock.Anything, task, mock.Anything).Return(expectedResponse, nil)

	// Start the test gRPC server
	conn, cleanup, err := startTestServer(llmMock)
	defer cleanup()
	assert.NoError(t, err)

	// Create the gRPC client
	client := agent_external.NewAgentServiceClient(conn)

	// Call DispatchTask
	req := &agent_external.SubmitTaskRequest{Task: task}
	resp, err := client.SubmitTask(context.Background(), req)

	// Assert results
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, resp.Response)
	llmMock.AssertExpectations(t)
}

func TestSubmitTask_Error(t *testing.T) {
	// Set up the mock to return an error
	llmMock := new(agentMock)
	task := "Test Task"
	expectedError := errors.New("failed to process task")

	llmMock.On("SubmitTask", mock.Anything, task, mock.Anything).Return("", expectedError)

	// Start the test gRPC server
	conn, cleanup, err := startTestServer(llmMock)
	defer cleanup()
	assert.NoError(t, err)

	// Create the gRPC client
	client := agent_external.NewAgentServiceClient(conn)

	// Call DispatchTask and expect an error
	req := &agent_external.SubmitTaskRequest{Task: task}
	resp, err := client.SubmitTask(context.Background(), req)

	// Assert the error and response
	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, expectedError.Error(), status.Convert(err).Message())
	llmMock.AssertExpectations(t)
}
