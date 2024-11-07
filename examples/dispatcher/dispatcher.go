package main

import (
	"context"
	"encoding/json"
	"fmt"
	discoverer "github.com/dhiaayachi/llm-fabric/discoverer"
	"github.com/dhiaayachi/llm-fabric/discoverer/store"
	"github.com/dhiaayachi/llm-fabric/fabric"
	"github.com/dhiaayachi/llm-fabric/llm"
	agentv1 "github.com/dhiaayachi/llm-fabric/proto/gen/agent/v1"
	strategy "github.com/dhiaayachi/llm-fabric/strategy"
	"github.com/hashicorp/serf/serf"
	"github.com/oklog/ulid/v2"
	"github.com/ollama/ollama/api"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"time"
)

type Dispatcher struct {
	logger logrus.Logger
}

func availableCapabilities(Agents []*agentv1.Agent) []*agentv1.Capability {
	res := make([]*agentv1.Capability, 0)
	capa := make(map[string]*agentv1.Capability)
	for _, a := range Agents {
		for _, c := range a.Capabilities {
			capa[c.Id] = c
		}
	}
	for _, c := range capa {
		res = append(res, c)
	}
	return res
}

func (d *Dispatcher) Execute(task string, Agents []*agentv1.Agent, localLLM llm.Llm) []*strategy.TaskAgent {
	capa := availableCapabilities(Agents)
	prompt := "select the best capabilities to answer the following task:\\n\\n"
	prompt = prompt + fmt.Sprintf("%s\\n\\n", task)
	prompt = prompt + "the available capabilities are:\\n"
	marchal, err := json.Marshal(capa)
	if err != nil {
		d.logger.Fatal(err)
	}
	prompt = prompt + fmt.Sprintf("%s\\n\\n", marchal)
	prompt = prompt + "\\n select a set of capabilities that an AI agent should have to solve this task, return a subset of capabilities that are needed to solve this task" +
		"(minimum 1 and maximum 3)"

	response, err := localLLM.SubmitTask(context.Background(), prompt, "json")
	if err != nil {
		d.logger.Fatal(err)
	}

	c := make([]agentv1.Capability, 0)
	err = json.Unmarshal([]byte(response), &c)
	if err != nil {
		d.logger.Fatal(err)
	}

	d.logger.WithFields(logrus.Fields{"response": response}).Info("got a response!")

	return nil
}

func (d *Dispatcher) Finalize(_ []string) string {
	return ""
}

func main() {

	agent := agentv1.Agent{
		Description: "Ollama agent",
		Capabilities: []*agentv1.Capability{
			{Id: "1", Description: "text summarization"},
			{Id: "2", Description: "image generation"},
			{Id: "3", Description: "text generation"},
		},
		Tools: make([]*agentv1.Tool, 0),
		Id:    ulid.Make().String(),
	}

	logger := logrus.New()

	// Create discoverer
	s := store.NewInMemoryStore()

	serfConf := serf.DefaultConfig()
	serfConf.NodeName = agent.Id
	serfConf.MemberlistConfig.BindPort = 2222
	dicso, err := discoverer.NewSerfDiscoverer(serfConf, s, logger)
	if err != nil {
		logrus.Fatal(err)
	}
	err = dicso.Join(context.Background(), []string{"localhost:2222"}, &agent)
	if err != nil {
		logrus.Fatal(err)
	}

	time.Sleep(10 * time.Second)
	// Create local llm
	parse, err := url.Parse("http://localhost:11434")
	if err != nil {
		logrus.Fatal(err)
	}
	var ollama = api.NewClient(parse, http.DefaultClient)
	l := llm.NewOllama(ollama,
		logger,
		"llama3.2",
		"dispatcher",
		[]agentv1.Capability{{Id: "4", Description: "dispatch tasks to other agents"}},
		[]agentv1.Tool{})

	// Create fabric

	f := fabric.NewFabric(dicso, &Dispatcher{}, l)
	_, err = f.SubmitTask(context.Background(), "Can you summarize this text?")
	if err != nil {
		logrus.Fatal(err)
	}
}
