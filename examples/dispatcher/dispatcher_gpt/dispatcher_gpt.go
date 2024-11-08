package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dhiaayachi/llm-fabric/agent"
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
		err := fmt.Errorf("could not find capability to solve")
		d.logger.Fatal(err)
		return nil
	}
	return []*strategy.TaskAgent{{Agent: capabaleAgent, Task: task}}
}

func (d *Dispatcher) Finalize(responses []string, _ llm.Llm) string {
	r := ""
	for _, response := range responses {
		r += response
	}
	return r
}

func main() {

	agentInfo := agentv1.AgentInfo{
		Description: "Ollama agent_info",
		Capabilities: []*agentv1.Capability{
			{Id: "1", Description: "text summarization"},
			{Id: "2", Description: "image generation"},
			{Id: "3", Description: "text generation"},
		},
		Tools:   make([]*agentv1.Tool, 0),
		Id:      ulid.Make().String(),
		Address: "127.0.0.1:3442",
	}

	logger := logrus.New()

	// Create discoverer
	s := store.NewInMemoryStore()

	serfConf := serf.DefaultConfig()
	serfConf.NodeName = agentInfo.Id
	serfConf.MemberlistConfig.BindPort = 2222
	dicso, err := discoverer.NewSerfDiscoverer(serfConf, s, logger)
	if err != nil {
		logrus.Fatal(err)
	}
	err = dicso.Join(context.Background(), []string{"localhost:2222"}, &agentInfo)
	if err != nil {
		logrus.Fatal(err)
	}
	time.Sleep(1 * time.Second)

	// Create local llm
	l := llm.NewGPT(openai.DefaultConfig(os.Getenv("OPENAI_TOKEN")),
		logger,
		"gpt-4o-2024-08-06",
		"assistant",
		[]agentv1.Capability{{Id: "4", Description: "dispatch tasks to other agents"}},
		[]agentv1.Tool{})

	srv := agent.NewServer(l, &agent.Config{Logger: logger, ListenAddr: "0.0.0.0:3442"})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	srv.Start(ctx)

	// Create fabric
	f := fabric.NewFabric(dicso, &Dispatcher{logger: logger}, l)
	response, err := f.SubmitTask(context.Background(), "Can you summarize this text?: Johannes Gutenberg (1398 – 1468) "+
		"was a German goldsmith and publisher who introduced printing to Europe. His introduction of mechanical "+
		"movable type printing to Europe started the Printing Revolution and is widely regarded as the most important "+
		"event of the modern period. It played a key role in the scientific revolution and laid the basis for the "+
		"modern knowledge-based economy and the spread of learning to the masses.\\n\\nGutenberg many contributions "+
		"to printing are: the invention of a process for mass-producing movable type, the use of oil-based ink for "+
		"printing books, adjustable molds, and the use of a wooden printing press. His truly epochal invention was "+
		"the combination of these elements into a practical system that allowed the mass production of printed books "+
		"and was economically viable for printers and readers alike.\n\nIn Renaissance Europe, the arrival of mechanical "+
		"movable type printing introduced the era of mass communication which permanently altered the structure"+
		" of society. The relatively unrestricted circulation of information—including revolutionary ideas—transcended"+
		" borders, and captured the masses in the Reformation. The sharp increase in literacy broke the monopoly "+
		"of the literate elite on education and learning and bolstered the emerging middle class.")
	if err != nil {
		logrus.Fatal(err)
	}

	logger.WithField("response", response).Info("final response!")
}
