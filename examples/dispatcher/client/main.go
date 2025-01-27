package main

import (
	"context"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dhiaayachi/llm-fabric/discoverer"
	"github.com/dhiaayachi/llm-fabric/discoverer/store"
	"github.com/dhiaayachi/llm-fabric/fabric"
	agentinfo "github.com/dhiaayachi/llm-fabric/proto/gen/agent_info/v1"
	"github.com/hashicorp/serf/serf"
	"github.com/sirupsen/logrus"
)

func main() {

	grpcPort, err := strconv.Atoi(os.Getenv("GRPC_PORT"))
	if err != nil {
		logrus.Fatalf("failed to parse GRPC_PORT as integer: %v", err)
	}

	serfPort, err := strconv.Atoi(os.Getenv("SERF_PORT"))
	if err != nil {
		logrus.Fatalf("failed to parse GRPC_PORT as integer: %v", err)
	}

	agentInfo := agentinfo.AgentsNodeInfo{
		Node: &agentinfo.NodeInfo{
			Port: int32(grpcPort),
		},
	}

	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	w := logger.WriterLevel(logrus.DebugLevel)

	// note that you are responsible for closing the writer
	defer w.Close()

	// Create discoverer
	s := store.NewInMemoryStore()

	serfConf := serf.DefaultConfig()
	serfConf.LogOutput = w

	serfConf.MemberlistConfig.BindAddr = "0.0.0.0"
	serfConf.MemberlistConfig.BindPort = serfPort
	dicso, err := discoverer.NewSerfDiscoverer(serfConf, s, logger)
	if err != nil {
		logrus.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	addrs := strings.Split(os.Getenv("SERF_JOIN_ADDRS"), " ")
	err = dicso.Join(ctx, addrs, &agentInfo)
	if err != nil {
		logrus.Fatal(err)
	}
	time.Sleep(30 * time.Second)

	// Create fabric
	c := fabric.NewClient(dicso, logger)
	logger.Info("sending request to dispatch!")
	response, err := c.SubmitTask(context.Background(), "Can you summarize this text?: Johannes Gutenberg (1398 – 1468) "+
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
	select {
	case <-ctx.Done():
	}
}
