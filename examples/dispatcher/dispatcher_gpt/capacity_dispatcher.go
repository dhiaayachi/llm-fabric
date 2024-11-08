package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dhiaayachi/llm-fabric/llm"
	agentinfo "github.com/dhiaayachi/llm-fabric/proto/gen/agent_info/v1"
	"github.com/dhiaayachi/llm-fabric/strategy"
	"github.com/sashabaranov/go-openai/jsonschema"
	"github.com/sirupsen/logrus"
)

type CapacityDispatcher struct {
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

func (d *CapacityDispatcher) Execute(task string, Agents []*agentinfo.AgentInfo, localLLM llm.Llm) []*strategy.TaskAgent {
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

	o := &agentinfo.LlmOpt{Typ: agentinfo.LlmOptType_LLM_OPT_TYPE_GPTResponseFormat}
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

	var capabaleAgent *agentinfo.AgentInfo
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

func (d *CapacityDispatcher) Finalize(responses []string, _ llm.Llm) string {
	r := ""
	for _, response := range responses {
		r += response
	}
	return r
}
