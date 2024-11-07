package fabric

import (
	"context"
	"github.com/dhiaayachi/llm-fabric/discoverer"
	"github.com/dhiaayachi/llm-fabric/llm"
	agentv1 "github.com/dhiaayachi/llm-fabric/proto/gen/agent/v1"
	"github.com/dhiaayachi/llm-fabric/strategy"
	"google.golang.org/grpc"
)

type Fabric struct {
	localLlm   llm.Llm
	discoverer discoverer.Discoverer
	strategy   strategy.Strategy
}

func NewFabric(discoverer discoverer.Discoverer, strategy strategy.Strategy, llm llm.Llm) *Fabric {
	return &Fabric{discoverer: discoverer, strategy: strategy, localLlm: llm}
}

func (f Fabric) SubmitTask(ctx context.Context, task string) (string, error) {
	taskAgents := f.strategy.Execute(task, f.discoverer.GetAgents(), f.localLlm)
	rsps := make([]string, 0)
	for _, taskAgent := range taskAgents {
		conn, err := grpc.NewClient(taskAgent.Agent.Address)
		if err != nil {
			return "", err
		}
		client := agentv1.NewAgentServiceClient(conn)
		response, err := client.SubmitTask(ctx, &agentv1.SubmitTaskRequest{Task: taskAgent.Task})
		if err != nil {
			return "", err
		}
		rsps = append(rsps, response.Response)
	}
	return f.strategy.Finalize(rsps), nil
}
