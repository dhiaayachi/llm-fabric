package main

import (
	"context"
	"github.com/dhiaayachi/llm-fabric/agent"
	"github.com/dhiaayachi/llm-fabric/discoverer"
	"github.com/dhiaayachi/llm-fabric/discoverer/store"
	"github.com/dhiaayachi/llm-fabric/fabric"
	"github.com/dhiaayachi/llm-fabric/llm"
	agentv1 "github.com/dhiaayachi/llm-fabric/proto/gen/agent_info/v1"
	"github.com/hashicorp/serf/serf"
	"github.com/oklog/ulid/v2"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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
	srv.Start(ctx)

	// Create fabric
	f := fabric.NewFabric(dicso, &CapacityDispatcher{logger: logger}, l)
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
