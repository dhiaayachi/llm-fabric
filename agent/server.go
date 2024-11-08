package agent

import (
	"context"
	"github.com/dhiaayachi/llm-fabric/llm"
	agentinfo "github.com/dhiaayachi/llm-fabric/proto/gen/agent_info/v1"
	"google.golang.org/grpc"
)

type Server struct {
	agentinfo.UnimplementedAgentServiceServer
	srv *grpc.Server
	llm llm.Llm
}

func (s *Server) SubmitTask(ctx context.Context, request *agentinfo.SubmitTaskRequest) (*agentinfo.SubmitTaskResponse, error) {
	resp := &agentinfo.SubmitTaskResponse{}
	response, err := s.llm.SubmitTask(ctx, request.Task, request.Opts...)
	if err != nil {
		return nil, err
	}
	resp.Response = response
	return resp, nil
}

var _ agentinfo.AgentServiceServer = &Server{}

func NewServer(llm llm.Llm) *Server {
	srv := Server{srv: grpc.NewServer(), llm: llm}
	srv.srv.RegisterService(&agentinfo.AgentService_ServiceDesc, &srv)
	return &srv
}
