package grpc

import (
	"context"
	"github.com/dhiaayachi/llm-fabric/llm"
	"github.com/dhiaayachi/llm-fabric/proto/gen/agent_external/v1"
	agentinfo "github.com/dhiaayachi/llm-fabric/proto/gen/agent_info/v1"
	"github.com/dhiaayachi/llm-fabric/strategy"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"
	"net"
	"time"
)

type agent interface {
	DispatchTask(ctx context.Context, task string, schema any) (string, error)
	GetStrategies() []strategy.Strategy
	GetAgents() []*agentinfo.AgentsNodeInfo
	GetLocalLlm() llm.Llm
	SubmitTask(ctx context.Context, task string, schema *anypb.Any) (string, error)
}

type Server struct {
	agent_external.UnimplementedAgentServiceServer
	srv        *grpc.Server
	agent      agent
	logger     *logrus.Logger
	ListenAddr string
	listener   net.Listener
}

type Config struct {
	ListenAddr string
	Logger     *logrus.Logger
}

func (srv *Server) SubmitTask(ctx context.Context, request *agent_external.SubmitTaskRequest) (*agent_external.SubmitTaskResponse, error) {
	resp := &agent_external.SubmitTaskResponse{}
	response, err := srv.agent.SubmitTask(ctx, request.Task, request.Schema)
	if err != nil {
		return nil, err
	}
	resp.Response = response
	return resp, nil
}

func (srv *Server) DispatchTask(ctx context.Context, request *agent_external.DispatchTaskRequest) (*agent_external.DispatchTaskResponse, error) {
	resp := &agent_external.DispatchTaskResponse{}
	response, err := srv.agent.DispatchTask(ctx, request.Task, request.Schema)
	if err != nil {
		return nil, err
	}
	resp.Response = response
	return resp, nil
}

var _ agent_external.AgentServiceServer = &Server{}

func NewServer(agent agent, conf *Config) *Server {
	srv := Server{srv: grpc.NewServer(), agent: agent, logger: conf.Logger, ListenAddr: conf.ListenAddr}
	srv.srv.RegisterService(&agent_external.AgentService_ServiceDesc, &srv)
	return &srv
}

func (srv *Server) Start(ctx context.Context) {
	lis, err := net.Listen("tcp", srv.ListenAddr)
	if err != nil {
		srv.logger.WithError(err).Error("failed to listen")
	}
	srv.listener = lis
	go func() {
		for {
			err = srv.srv.Serve(srv.listener)
			if err != nil {
				srv.logger.WithError(err).Error("failed to serve")
			}
			after := time.After(time.Second)
			select {
			case <-after:
			case <-ctx.Done():
				return
			}
		}
	}()
}
