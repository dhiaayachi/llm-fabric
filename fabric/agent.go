package fabric

import (
	"context"
	"fmt"
	"github.com/dhiaayachi/llm-fabric/discoverer"
	"github.com/dhiaayachi/llm-fabric/fabric/internal/grpc"
	"github.com/dhiaayachi/llm-fabric/llm"
	"github.com/dhiaayachi/llm-fabric/proto/gen/agent_external/v1"
	agentinfo "github.com/dhiaayachi/llm-fabric/proto/gen/agent_info/v1"
	llmoptions "github.com/dhiaayachi/llm-fabric/proto/gen/llm_options/v1"
	"github.com/dhiaayachi/llm-fabric/strategy"
	"github.com/sirupsen/logrus"
)

type Agent struct {
	localLlm           llm.Llm
	discoverer         discoverer.Discoverer
	DispatchStrategies []strategy.Strategy
	cancel             context.CancelFunc
}

func (a Agent) GetLocalLlm() llm.Llm {
	return a.localLlm
}

func (a Agent) GetAgents() []*agentinfo.AgentsNodeInfo {
	return a.discoverer.GetAgents()
}

func (a Agent) GetStrategies() []strategy.Strategy {
	return a.DispatchStrategies
}

func NewAgent(ctx context.Context, discoverer discoverer.Discoverer, strategies []strategy.Strategy, llm llm.Llm, logger *logrus.Logger, grpcPort int) *Agent {

	agent := Agent{discoverer: discoverer, DispatchStrategies: strategies, localLlm: llm}
	srv := grpc.NewServer(&agent, &grpc.Config{Logger: logger, ListenAddr: fmt.Sprintf("0.0.0.0:%d", grpcPort)})
	srv.Start(ctx)
	return &agent
}

func (a Agent) SubmitTask(ctx context.Context, task string, opts []*llmoptions.LlmOpt) (string, error) {
	if len(a.DispatchStrategies) > 0 {
		taskAgents := a.DispatchStrategies[0].Execute(task, a.discoverer.GetAgents(), a.localLlm)

		rsps := make([]string, 0)
		for _, taskAgent := range taskAgents {
			client, err := grpc.GetClient(taskAgent.Node.Address, taskAgent.Node.Port)
			if err != nil {
				return "", err
			}
			response, err := client.SubmitTask(ctx, &agent_external.SubmitTaskRequest{Task: taskAgent.Task, Opts: opts})
			if err != nil {
				return "", err
			}
			rsps = append(rsps, response.Response)
		}
		return a.DispatchStrategies[0].Finalize(rsps, a.localLlm), nil
	} else {
		response, err := a.localLlm.SubmitTask(ctx, task, opts...)
		if err != nil {
			return "", err
		}
		return response, nil
	}
}
