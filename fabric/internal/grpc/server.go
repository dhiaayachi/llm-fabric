package grpc

import (
	"context"
	"github.com/dhiaayachi/llm-fabric/llm"
	"github.com/dhiaayachi/llm-fabric/proto/gen/agent_external/v1"
	agentinfo "github.com/dhiaayachi/llm-fabric/proto/gen/agent_info/v1"
	llmoptions "github.com/dhiaayachi/llm-fabric/proto/gen/llm_options/v1"
	"github.com/dhiaayachi/llm-fabric/strategy"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"time"
)

type agent interface {
	SubmitTask(ctx context.Context, task string, opts []*llmoptions.LlmOpt) (string, error)
	GetStrategies() []strategy.Strategy
	GetAgents() []*agentinfo.AgentsNodeInfo
	GetLocalLlm() llm.Llm
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
	response, err := srv.agent.SubmitTask(ctx, request.Task, request.Opts)
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

func (srv *Server) submitTask(ctx context.Context, task string, opts []*llmoptions.LlmOpt) (string, error) {
	strategies := srv.agent.GetStrategies()
	if len(strategies) > 0 {
		taskAgents := strategies[0].Execute(task, srv.agent.GetAgents(), srv.agent.GetLocalLlm())

		rsps := make([]string, 0)
		for _, taskAgent := range taskAgents {
			client, err := GetClient(taskAgent.Node.Address, taskAgent.Node.Port)
			if err != nil {
				return "", err
			}
			response, err := client.SubmitTask(ctx, &agent_external.SubmitTaskRequest{Task: taskAgent.Task, Opts: opts})
			if err != nil {
				return "", err
			}
			rsps = append(rsps, response.Response)
		}
		return strategies[0].Finalize(rsps, srv.agent.GetLocalLlm()), nil
	} else {
		response, err := srv.agent.GetLocalLlm().SubmitTask(ctx, task, opts...)
		if err != nil {
			return "", err
		}
		return response, nil
	}
}
