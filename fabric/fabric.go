package fabric

import (
	"context"
	"github.com/dhiaayachi/llm-fabric/agent"
	"github.com/dhiaayachi/llm-fabric/discoverer"
	"github.com/dhiaayachi/llm-fabric/llm"
	"github.com/dhiaayachi/llm-fabric/proto/gen/agent_external/v1"
	"github.com/dhiaayachi/llm-fabric/strategy"
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
		client, err := agent.GetClient(taskAgent.Agent.Address, taskAgent.Agent.Port)
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
