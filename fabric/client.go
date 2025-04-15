package fabric

import (
	"context"
	"errors"
	"github.com/dhiaayachi/llm-fabric/discoverer"
	"github.com/dhiaayachi/llm-fabric/fabric/internal/grpc"
	"github.com/dhiaayachi/llm-fabric/proto/gen/agent_external/v1"
	"github.com/sirupsen/logrus"
)

type Client struct {
	discoverer discoverer.Discoverer
	logger     *logrus.Logger
}

func NewClient(discoverer discoverer.Discoverer, logger *logrus.Logger) *Client {
	return &Client{discoverer: discoverer, logger: logger}
}

func (c Client) DispatchTask(ctx context.Context, task string) (string, error) {

	dispatchers := c.discoverer.GetDispatchers()

	c.logger.WithField("num_dispatchers", len(dispatchers)).Debug("Dispatching task")
	if len(dispatchers) == 0 {
		return "", errors.New("no dispatchers found")
	}

	client, err := grpc.MakeClient(dispatchers[0].Node.Address, dispatchers[0].Node.Port)
	if err != nil {
		return "", err
	}

	response, err := client.DispatchTask(ctx, &agent_external.DispatchTaskRequest{Task: task})
	if err != nil {
		return "", err
	}
	return response.GetResponse(), nil
}
