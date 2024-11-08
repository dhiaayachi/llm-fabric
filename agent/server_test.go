package agent_test

import (
	"context"
	"errors"
	"google.golang.org/grpc/credentials/insecure"
	"testing"

	"github.com/dhiaayachi/llm-fabric/agent"
	"github.com/dhiaayachi/llm-fabric/llm"
	agentinfo "github.com/dhiaayachi/llm-fabric/proto/gen/agent_info/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"net"
)

// MockLlm is a mock of the llm.Llm interface generated by mockery.
type MockLlm struct {
	mock.Mock
}

func (m *MockLlm) SubmitTask(ctx context.Context, task string, opts ...*agentinfo.LlmOpt) (string, error) {
	args := m.Called(ctx, task, opts)
	return args.Get(0).(string), args.Error(1)
}

func (m *MockLlm) GetCapabilities() []agentinfo.Capability {
	args := m.Called()
	return args.Get(0).([]agentinfo.Capability)
}

func (m *MockLlm) GetTools() []agentinfo.Tool {
	args := m.Called()
	return args.Get(0).([]agentinfo.Tool)
}

func startTestServer(llmMock llm.Llm) (*grpc.ClientConn, func(), error) {
	// Set up a buffer connection that simulates a network for gRPC communication
	bufSize := 1024 * 1024
	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer()
	agentServer := agent.NewServer(llmMock, &agent.Config{})

	agentinfo.RegisterAgentServiceServer(s, agentServer)
	go func() {
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()

	// Dial a connection to the buffer connection
	conn, err := grpc.DialContext(context.Background(), "bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		s.Stop()
		_ = lis.Close()
	}
	return conn, cleanup, nil
}

func TestSubmitTask_Success(t *testing.T) {
	// Set up the mock and expected behavior
	llmMock := new(MockLlm)
	task := "Test Task"
	expectedResponse := "Task Response"

	llmMock.On("SubmitTask", mock.Anything, task, mock.Anything).Return(expectedResponse, nil)

	// Start the test gRPC server
	conn, cleanup, err := startTestServer(llmMock)
	defer cleanup()
	assert.NoError(t, err)

	// Create the gRPC client
	client := agentinfo.NewAgentServiceClient(conn)

	// Call SubmitTask
	req := &agentinfo.SubmitTaskRequest{Task: task}
	resp, err := client.SubmitTask(context.Background(), req)

	// Assert results
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, resp.Response)
	llmMock.AssertExpectations(t)
}

func TestSubmitTask_Error(t *testing.T) {
	// Set up the mock to return an error
	llmMock := new(MockLlm)
	task := "Test Task"
	expectedError := errors.New("failed to process task")

	llmMock.On("SubmitTask", mock.Anything, task, mock.Anything).Return("", expectedError)

	// Start the test gRPC server
	conn, cleanup, err := startTestServer(llmMock)
	defer cleanup()
	assert.NoError(t, err)

	// Create the gRPC client
	client := agentinfo.NewAgentServiceClient(conn)

	// Call SubmitTask and expect an error
	req := &agentinfo.SubmitTaskRequest{Task: task}
	resp, err := client.SubmitTask(context.Background(), req)

	// Assert the error and response
	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, expectedError.Error(), status.Convert(err).Message())
	llmMock.AssertExpectations(t)
}
