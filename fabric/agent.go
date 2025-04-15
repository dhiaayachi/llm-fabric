package fabric

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dhiaayachi/llm-fabric/discoverer"
	"github.com/dhiaayachi/llm-fabric/fabric/internal/grpc"
	"github.com/dhiaayachi/llm-fabric/llm"
	"github.com/dhiaayachi/llm-fabric/proto/gen/agent_external/v1"
	agentinfo "github.com/dhiaayachi/llm-fabric/proto/gen/agent_info/v1"
	"github.com/dhiaayachi/llm-fabric/strategy"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
)

type Agent struct {
	localLlm           llm.Llm
	discoverer         discoverer.Discoverer
	DispatchStrategies []strategy.Strategy
	cancel             context.CancelFunc
	logger             *logrus.Logger
}

func (a *Agent) SubmitTask(ctx context.Context, task string, schema *anypb.Any) (string, error) {
	s, err := fromAnyPB(schema)
	response, err := a.localLlm.SubmitTask(ctx, task, s)
	if err != nil {
		return "", err
	}
	return response, nil
}

func fromAnyPB(schema *anypb.Any) (any, error) {
	m := make(map[string]any)

	bytes := schema.GetValue()
	if len(bytes) == 0 {
		return nil, nil
	}

	err := json.Unmarshal(bytes, &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (a *Agent) GetLocalLlm() llm.Llm {
	return a.localLlm
}

func (a *Agent) GetAgents() []*agentinfo.AgentsNodeInfo {
	return a.discoverer.GetAgents()
}

func (a *Agent) GetStrategies() []strategy.Strategy {
	return a.DispatchStrategies
}

func NewAgent(ctx context.Context, discoverer discoverer.Discoverer, strategies []strategy.Strategy, llm llm.Llm, logger *logrus.Logger, grpcPort int) *Agent {
	agent := Agent{discoverer: discoverer, DispatchStrategies: strategies, localLlm: llm, logger: logger}
	srv := grpc.NewServer(&agent, &grpc.Config{Logger: logger, ListenAddr: fmt.Sprintf("0.0.0.0:%d", grpcPort)})
	srv.Start(ctx)
	return &agent
}

func (a *Agent) DispatchTask(ctx context.Context, task string, schema any) (string, error) {
	taskAgents := a.DispatchStrategies[0].Execute(task, a.discoverer.GetAgents(), a.localLlm)
	rsps := make([]string, 0)
	for _, taskAgent := range taskAgents {
		client, err := grpc.MakeClient(taskAgent.Node.Address, taskAgent.Node.Port)
		if err != nil {
			return "", err
		}
		var s *anypb.Any
		if taskAgent.Schema != nil {
			a.logger.Infof("dispatch task agent with schema: %s", taskAgent.Schema)
			s, err = toAnyPB(taskAgent.Schema)
			if err != nil {
				return "", err
			}
		}

		response, err := client.SubmitTask(ctx, &agent_external.SubmitTaskRequest{Task: taskAgent.Task, Schema: s})
		if err != nil {
			return "", err
		}
		rsps = append(rsps, response.Response)
	}
	return a.DispatchStrategies[0].Finalize(rsps, a.localLlm), nil

}

func toAnyPB(schema any) (*anypb.Any, error) {
	m := make(map[string]any)

	bytes, err := json.Marshal(schema)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &m)
	if err != nil {
		return nil, err
	}

	newStruct, err := structpb.NewStruct(m)

	s, err := anypb.New(newStruct)
	return s, nil
}
