package fabric

import (
	"context"
	"errors"
	"github.com/dhiaayachi/llm-fabric/discoverer"
	"github.com/dhiaayachi/llm-fabric/fabric/internal/grpc"
	"github.com/dhiaayachi/llm-fabric/proto/gen/agent_external/v1"
	"github.com/dhiaayachi/llm-fabric/proto/gen/agent_info/v1"
	"github.com/sirupsen/logrus"
)

type Client struct {
	discoverer discoverer.Discoverer
	logger     *logrus.Logger
}

func NewClient(discoverer discoverer.Discoverer, logger *logrus.Logger) *Client {
	return &Client{discoverer: discoverer, logger: logger}
}

func (c Client) SubmitTask(ctx context.Context, task string) (string, error) {

	agents := c.discoverer.GetAgents()

	dispatchers := make([]*agent_info.AgentsNodeInfo, 0)

	c.logger.WithField("num_agents", len(agents)).Debug("Submitting task")
	for _, agent := range agents {
		if agent.Agents[0].IsDispatcher {
			dispatchers = append(dispatchers, agent)
		}
	}

	c.logger.WithField("num_dispatchers", len(dispatchers)).Debug("Dispatching task")
	if len(dispatchers) == 0 {
		return "", errors.New("no dispatchers found")
	}

	client, err := grpc.MakeClient(dispatchers[0].Node.Address, dispatchers[0].Node.Port)
	if err != nil {
		return "", err
	}

	response, err := client.SubmitTask(ctx, &agent_external.SubmitTaskRequest{Task: task})
	if err != nil {
		return "", err
	}
	return response.GetResponse(), nil
}
