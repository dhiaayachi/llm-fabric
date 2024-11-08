package agent

import (
	"context"
	"github.com/dhiaayachi/llm-fabric/llm"
	agentinfo "github.com/dhiaayachi/llm-fabric/proto/gen/agent_info/v1"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"time"
)

type Server struct {
	agentinfo.UnimplementedAgentServiceServer
	srv        *grpc.Server
	llm        llm.Llm
	logger     *logrus.Logger
	ListenAddr string
}

type Config struct {
	ListenAddr string
	Logger     *logrus.Logger
}

func (srv *Server) SubmitTask(ctx context.Context, request *agentinfo.SubmitTaskRequest) (*agentinfo.SubmitTaskResponse, error) {
	resp := &agentinfo.SubmitTaskResponse{}
	response, err := srv.llm.SubmitTask(ctx, request.Task, request.Opts...)
	if err != nil {
		return nil, err
	}
	resp.Response = response
	return resp, nil
}

var _ agentinfo.AgentServiceServer = &Server{}

func NewServer(llm llm.Llm, conf *Config) *Server {
	srv := Server{srv: grpc.NewServer(), llm: llm, logger: conf.Logger, ListenAddr: conf.ListenAddr}
	srv.srv.RegisterService(&agentinfo.AgentService_ServiceDesc, &srv)
	return &srv
}

func (srv *Server) Start(ctx context.Context) {
	go func() {
		for {
			lis, err := net.Listen("tcp", srv.ListenAddr)
			if err != nil {
				srv.logger.WithError(err).Error("failed to listen")
			}
			err = srv.srv.Serve(lis)
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
