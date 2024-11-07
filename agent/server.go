package agent

import (
	"context"
	"github.com/dhiaayachi/llm-fabric/llm"
	agentv1 "github.com/dhiaayachi/llm-fabric/proto/gen/agent/v1"
	"google.golang.org/grpc"
)

type Server struct {
	agentv1.UnimplementedAgentServiceServer
	srv *grpc.Server
	llm llm.Llm
}

func (s Server) SubmitTask(ctx context.Context, request *agentv1.SubmitTaskRequest) (*agentv1.SubmitTaskResponse, error) {
	resp := &agentv1.SubmitTaskResponse{}
	response, err := s.llm.SubmitTask(ctx, request.Task, "json")
	if err != nil {
		return nil, err
	}
	resp.Response = response
	return resp, nil
}

var _ agentv1.AgentServiceServer = &Server{}

func NewServer(llm llm.Llm) *Server {
	srv := Server{srv: grpc.NewServer(), llm: llm}
	srv.srv.RegisterService(&agentv1.AgentService_ServiceDesc, srv)
	return &srv
}
