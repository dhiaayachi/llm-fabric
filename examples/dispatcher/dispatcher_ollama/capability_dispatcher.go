package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dhiaayachi/llm-fabric/llm"
	agentinfo "github.com/dhiaayachi/llm-fabric/proto/gen/agent_info/v1"
	llmoptions "github.com/dhiaayachi/llm-fabric/proto/gen/llm_options/v1"
	"github.com/dhiaayachi/llm-fabric/strategy"
	"github.com/sirupsen/logrus"
)

type CapabilityDispatcher struct {
	logger *logrus.Logger
}
type Agent struct {
	Id           string                  `json:"id"`
	Capabilities []*agentinfo.Capability `json:"capabilities"`
	Tools        []*agentinfo.Tool       `json:"tools"`
}

type AgentNode struct {
	AgentDesc Agent `json:"agent_description"`
	node      *agentinfo.NodeInfo
	agent     *agentinfo.AgentInfo
}

func availableAgents(agentsNodes []*agentinfo.AgentsNodeInfo) map[string]*AgentNode {

	res := make(map[string]*AgentNode)
	for _, agentN := range agentsNodes {
		for _, a := range agentN.Agents {
			res[a.Id] = &AgentNode{AgentDesc: Agent{Id: a.Id, Capabilities: a.Capabilities, Tools: a.Tools}, node: agentN.Node, agent: a}
		}
	}
	return res
}

func (d *CapabilityDispatcher) Execute(task string, agentsNodes []*agentinfo.AgentsNodeInfo, localLLM llm.Llm) []*strategy.TaskAgent {
	agents := availableAgents(agentsNodes)
	prompt := "select the best agents (1 to 3 agents) to answer the following task based on its capabilities and available tools:\\n\\n"
	prompt = prompt + fmt.Sprintf("%s\\n\\n", task)
	prompt = prompt + "the available agents are:\\n"
	marchal, err := json.Marshal(agents)
	if err != nil {
		d.logger.Fatal(err)
	}
	prompt = prompt + fmt.Sprintf("%s\\n\\n", marchal)

	o := &llmoptions.LlmOpt{Typ: llmoptions.LlmOptType_LLM_OPT_TYPE_OLLAMA_RESPONSE_SCHEMA}
	type result struct {
		Agents []struct {
			Id string `json:"id"`
		} `json:"agents"`
	}

	v := result{
		Agents: []struct {
			Id string `json:"id"`
		}{
			{Id: "Id of the selected agent"},
		},
	}

	schema, err := json.Marshal(v)
	if err != nil {
		d.logger.Fatal(err)
		return nil
	}

	err = llm.FromVal[string](o, string(schema))
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

	if len(res.Agents) == 0 {
		err := fmt.Errorf("could not find agent to solve")
		d.logger.Fatal(err)
		return nil
	}
	d.logger.WithField("capabaleAgents", res.Agents).Info("found capable agents")

	var selectedAgent = agents[res.Agents[0].Id]
	leastCost := selectedAgent.agent.Cost

	for _, agent := range res.Agents {
		if leastCost > agents[agent.Id].agent.Cost {
			leastCost = agents[agent.Id].agent.Cost
			selectedAgent = agents[agent.Id]
		}
	}
	d.logger.WithField("selectedAgent", selectedAgent).Info("selected agent")
	return []*strategy.TaskAgent{{Agent: selectedAgent.agent, Task: task, Node: selectedAgent.node}}
}

func (d *CapabilityDispatcher) Finalize(responses []string, _ llm.Llm) string {
	r := ""
	for _, response := range responses {
		r += response
	}
	return r
}
