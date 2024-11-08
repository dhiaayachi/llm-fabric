package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dhiaayachi/llm-fabric/discoverer"
	"github.com/dhiaayachi/llm-fabric/discoverer/store"
	"github.com/dhiaayachi/llm-fabric/fabric"
	"github.com/dhiaayachi/llm-fabric/llm"
	agentv1 "github.com/dhiaayachi/llm-fabric/proto/gen/agent_info/v1"
	"github.com/dhiaayachi/llm-fabric/strategy"
	"github.com/hashicorp/serf/serf"
	"github.com/oklog/ulid/v2"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type Dispatcher struct {
	logger *logrus.Logger
}

func availableCapabilities(Agents []*agentv1.AgentInfo) []*agentv1.Capability {
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

func (d *Dispatcher) Execute(task string, Agents []*agentv1.AgentInfo, localLLM llm.Llm) []*strategy.TaskAgent {
	capa := availableCapabilities(Agents)
	prompt := "select the best capabilities to answer the following task:\\n\\n"
	prompt = prompt + fmt.Sprintf("%s\\n\\n", task)
	prompt = prompt + "the available capabilities are:\\n"
	marchal, err := json.Marshal(capa)
	if err != nil {
		d.logger.Fatal(err)
	}
	prompt = prompt + fmt.Sprintf("%s\\n\\n", marchal)
	prompt = prompt + "\\n select a set of capabilities that an AI agent_info should have to solve this task, return a subset of capabilities that are needed to solve this task" +
		"(minimum 1 and maximum 3)"

	o := &agentv1.LlmOpt{Typ: agentv1.LlmOptType_LLM_OPT_TYPE_GPTResponseFormat}
	type result struct {
		Capabilities []struct {
			Id          string `json:"id"`
			Description string `json:"description"`
		} `json:"capabilities"`
	}

	schema, err := jsonschema.GenerateSchemaForType(result{})
	if err != nil {
		d.logger.Fatal(err)
		return nil
	}

	err = llm.FromVal[*jsonschema.Definition](o, schema)
	if err != nil {
		d.logger.Fatal(err)
	}
	response, err := localLLM.SubmitTask(context.Background(), prompt, o)
	if err != nil {
		d.logger.Fatal(err)
	}

	res := result{}
	err = json.Unmarshal([]byte(response), &res)
	if err != nil {
		d.logger.Fatal(err)
	}

	d.logger.WithFields(logrus.Fields{"response": res}).Info("got a response!")

	var capabaleAgent *agentv1.AgentInfo
	for _, a := range Agents {
		foundCap := false
		for _, c := range res.Capabilities {
			foundCap = false
			for _, ca := range a.Capabilities {
				if ca.Id == c.Id {
					//found it
					foundCap = true
					break
				}
			}
			if !foundCap {
				break
			}
		}
		if foundCap {
			capabaleAgent = a
			break
		}
	}
	if capabaleAgent == nil {
		d.logger.Fatal(fmt.Errorf("could not find capability to solve"))
	}

	return nil
}

func (d *Dispatcher) Finalize(_ []string) string {
	return ""
}

func main() {

	agent := agentv1.AgentInfo{
		Description: "Ollama agent_info",
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
	l := llm.NewGPT(openai.DefaultConfig(os.Getenv("OPENAI_TOKEN")),
		logger,
		"gpt-4o-2024-08-06",
		"assistant",
		[]agentv1.Capability{{Id: "4", Description: "dispatch tasks to other agents"}},
		[]agentv1.Tool{})

	// Create fabric

	f := fabric.NewFabric(dicso, &Dispatcher{logger: logger}, l)
	_, err = f.SubmitTask(context.Background(), "Can you summarize this text?")
	if err != nil {
		logrus.Fatal(err)
	}
}
