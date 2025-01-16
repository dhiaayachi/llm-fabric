package main

import (
	"context"
	"github.com/dhiaayachi/llm-fabric/discoverer"
	"github.com/dhiaayachi/llm-fabric/discoverer/store"
	"github.com/dhiaayachi/llm-fabric/fabric"
	"github.com/dhiaayachi/llm-fabric/llm"
	agentinfo "github.com/dhiaayachi/llm-fabric/proto/gen/agent_info/v1"
	"github.com/dhiaayachi/llm-fabric/strategy"
	"github.com/hashicorp/serf/serf"
	"github.com/oklog/ulid/v2"
	"github.com/sirupsen/logrus"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
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
		Agents: []*agentinfo.AgentInfo{
			{
				Description: "Ollama agent_info",
				Capabilities: []*agentinfo.Capability{
					{Description: "text summarization"},
					{Description: "text generation"},
				},
				Tools: make([]*agentinfo.Tool, 0),
				Id:    ulid.Make().String(),

				Cost: 1,
			},
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

	serfConf.Logger = log.New(w, "", 0)

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
	time.Sleep(1 * time.Second)
	// Create local llm
	parse, err := url.Parse(os.Getenv("OLLAMA_URL"))
	if err != nil {
		logrus.Fatal(err)
	}
	l := llm.NewOllama(parse.String(), logger, "llama3.2", "dispatcher")

	_ = fabric.NewAgent(ctx, dicso, []strategy.Strategy{}, l, logger, grpcPort)
	select {
	case <-ctx.Done():
	}
}
