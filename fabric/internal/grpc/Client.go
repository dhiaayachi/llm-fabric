package grpc

import (
	"fmt"
	"github.com/dhiaayachi/llm-fabric/proto/gen/agent_external/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func MakeClient(host string, port int32) (*Client, error) {
	addr := fmt.Sprintf("%s:%d", host, port)
	// Dial a connection to the buffer connection
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{agent_external.NewAgentServiceClient(conn), conn.Close}, nil
}

type Client struct {
	agent_external.AgentServiceClient
	Close func() error
}
