package fabric

import (
	"context"
	"fmt"
	"github.com/dhiaayachi/llm-fabric/agent"
	"github.com/dhiaayachi/llm-fabric/discoverer"
	"github.com/dhiaayachi/llm-fabric/llm"
	agentinfo "github.com/dhiaayachi/llm-fabric/proto/gen/agent_info/v1"
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
		client, err := agent.GetClient(fmt.Sprintf("%s:%d", taskAgent.Agent.Address, taskAgent.Agent.Port))
		if err != nil {
			return "", err
		}
		response, err := client.SubmitTask(ctx, &agentinfo.SubmitTaskRequest{Task: taskAgent.Task})
		if err != nil {
			return "", err
		}
		rsps = append(rsps, response.Response)
	}
	return f.strategy.Finalize(rsps, f.localLlm), nil
}
