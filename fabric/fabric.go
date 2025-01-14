package fabric

import (
	"context"
	"fmt"
	"github.com/dhiaayachi/llm-fabric/discoverer"
	"github.com/dhiaayachi/llm-fabric/fabric/internal/grpc"
	"github.com/dhiaayachi/llm-fabric/llm"
	"github.com/dhiaayachi/llm-fabric/proto/gen/agent_external/v1"
	"github.com/dhiaayachi/llm-fabric/strategy"
	"github.com/sirupsen/logrus"
)

type Fabric struct {
	localLlm   llm.Llm
	discoverer discoverer.Discoverer
	strategy   strategy.Strategy
}

func NewFabric(ctx context.Context, discoverer discoverer.Discoverer, strategy strategy.Strategy, llm llm.Llm, logger *logrus.Logger, grpcPort int) *Fabric {
	srv := grpc.NewServer(llm, &grpc.Config{Logger: logger, ListenAddr: fmt.Sprintf("0.0.0.0:%d", grpcPort)})
	srv.Start(ctx)
	return &Fabric{discoverer: discoverer, strategy: strategy, localLlm: llm}
}

func (f Fabric) SubmitTask(ctx context.Context, task string) (string, error) {
	taskAgents := f.strategy.Execute(task, f.discoverer.GetAgents(), f.localLlm)
	rsps := make([]string, 0)
	for _, taskAgent := range taskAgents {
		client, err := grpc.GetClient(taskAgent.Node.Address, taskAgent.Node.Port)
		if err != nil {
			return "", err
		}
		response, err := client.SubmitTask(ctx, &agent_external.SubmitTaskRequest{Task: taskAgent.Task})
		if err != nil {
			return "", err
		}
		rsps = append(rsps, response.Response)
	}
	return f.strategy.Finalize(rsps, f.localLlm), nil
}
