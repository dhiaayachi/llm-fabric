package agent

import (
	"github.com/dhiaayachi/llm-fabric/proto/gen/agent_info/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetClient(addr string) (*Client, error) {
	// Dial a connection to the buffer connection
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{agent_info.NewAgentServiceClient(conn), conn.Close}, nil
}

type Client struct {
	agent_info.AgentServiceClient
	Close func() error
}
