package main

import (
	"context"
	"github.com/dhiaayachi/llm-fabric/agent"
	"github.com/dhiaayachi/llm-fabric/discoverer"
	"github.com/dhiaayachi/llm-fabric/discoverer/store"
	"github.com/dhiaayachi/llm-fabric/llm"
	agentinfo "github.com/dhiaayachi/llm-fabric/proto/gen/agent_info/v1"
	"github.com/hashicorp/serf/serf"
	"github.com/oklog/ulid/v2"
	"github.com/ollama/ollama/api"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"time"
)

func main() {
	agentInfo := agentinfo.AgentInfo{
		Description: "Ollama agent_info",
		Capabilities: []*agentinfo.Capability{
			{Id: "1", Description: "text summarization"},
			{Id: "3", Description: "text generation"},
		},
		Tools: make([]*agentinfo.Tool, 0),
		Id:    ulid.Make().String(),
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
	err = dicso.Join(ctx, []string{"0.0.0.0:2222"}, &agentInfo)
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
		[]agentinfo.Capability{{Id: "4", Description: "dispatch tasks to other agents"}},
		[]agentinfo.Tool{})

	srv := agent.NewServer(l, &agent.Config{Logger: logger, ListenAddr: "0.0.0.0:3442"})
	srv.Start(ctx)
	select {
	case <-ctx.Done():
	}
}
