package grpc

import (
	"context"
	"github.com/dhiaayachi/llm-fabric/llm"
	"github.com/dhiaayachi/llm-fabric/proto/gen/agent_external/v1"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"time"
)

type Server struct {
	agent_external.UnimplementedAgentServiceServer
	srv        *grpc.Server
	llm        llm.Llm
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
	response, err := srv.llm.SubmitTask(ctx, request.Task, request.Opts...)
	if err != nil {
		return nil, err
	}
	resp.Response = response
	return resp, nil
}

var _ agent_external.AgentServiceServer = &Server{}

func NewServer(llm llm.Llm, conf *Config) *Server {
	srv := Server{srv: grpc.NewServer(), llm: llm, logger: conf.Logger, ListenAddr: conf.ListenAddr}
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
