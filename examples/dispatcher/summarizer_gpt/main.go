package main

import (
	"context"
	"fmt"
	"github.com/dhiaayachi/llm-fabric/agent"
	"github.com/dhiaayachi/llm-fabric/discoverer"
	"github.com/dhiaayachi/llm-fabric/discoverer/store"
	"github.com/dhiaayachi/llm-fabric/llm"
	agentinfo "github.com/dhiaayachi/llm-fabric/proto/gen/agent_info/v1"
	"github.com/hashicorp/serf/serf"
	"github.com/oklog/ulid/v2"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
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

	agentInfo := agentinfo.AgentInfo{
		Description: "Ollama agent_info",
		Capabilities: []*agentinfo.Capability{
			{Id: "1", Description: "text summarization"},
			{Id: "3", Description: "text generation"},
		},
		Tools: make([]*agentinfo.Tool, 0),
		Id:    ulid.Make().String(),
		Port:  int32(grpcPort),
		Cost:  10, // GPT is high cost
	}

	logger := logrus.New()

	// Create discoverer
	s := store.NewInMemoryStore()

	serfConf := serf.DefaultConfig()

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
	l := llm.NewGPT(openai.DefaultConfig(os.Getenv("OPENAI_TOKEN")), logger, "gpt-4o-2024-08-06", "assistant")

	srv := agent.NewServer(l, &agent.Config{Logger: logger, ListenAddr: fmt.Sprintf("0.0.0.0:%d", grpcPort)})
	srv.Start(ctx)
	select {
	case <-ctx.Done():
	}
}
