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

func availableCapabilities(Agents []*agentinfo.AgentInfo) []*agentinfo.Capability {
	res := make([]*agentinfo.Capability, 0)
	capa := make(map[string]*agentinfo.Capability)
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

func (d *CapabilityDispatcher) Execute(task string, Agents []*agentinfo.AgentInfo, localLLM llm.Llm) []*strategy.TaskAgent {
	capa := availableCapabilities(Agents)
	prompt := "select the best capabilities to answer the following task:\\n\\n"
	prompt = prompt + fmt.Sprintf("%s\\n\\n", task)
	prompt = prompt + "the available capabilities are:\\n"
	marchal, err := json.Marshal(capa)
	if err != nil {
		d.logger.Fatal(err)
	}
	prompt = prompt + fmt.Sprintf("%s\\n\\n", marchal)
	prompt = prompt + "\\n select a sub set of capabilities that an AI agent_info should have to solve this task, return a subset of capabilities that are needed to solve this task. Only select from the provided capabilities" +
		"(minimum 1 and maximum 3)"

	o := &llmoptions.LlmOpt{Typ: llmoptions.LlmOptType_LLM_OPT_TYPE_OLLAMA_RESPONSE_SCHEMA}
	type result struct {
		Capabilities []struct {
			Id          string `json:"id"`
			Description string `json:"description"`
		} `json:"capabilities"`
	}

	v := result{Capabilities: []struct {
		Id          string `json:"id"`
		Description string `json:"description"`
	}{
		{Id: "Id of the selected capability",
			Description: "Description of the selected capability"},
	}}
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

	var capabaleAgents []*agentinfo.AgentInfo
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
			capabaleAgents = append(capabaleAgents, a)
		}
	}

	if len(capabaleAgents) == 0 {
		err := fmt.Errorf("could not find capability to solve")
		d.logger.Fatal(err)
		return nil
	}
	d.logger.WithField("capabaleAgents", capabaleAgents).Info("found capable agents")
	var selectedAgent = capabaleAgents[0]
	for _, a := range capabaleAgents {
		if selectedAgent.Cost > a.Cost {
			selectedAgent = a
		}
	}
	d.logger.WithField("selectedAgent", selectedAgent).Info("selected agent")
	return []*strategy.TaskAgent{{Agent: selectedAgent, Task: task}}
}

func (d *CapabilityDispatcher) Finalize(responses []string, _ llm.Llm) string {
	r := ""
	for _, response := range responses {
		r += response
	}
	return r
}