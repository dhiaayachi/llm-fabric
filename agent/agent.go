package agent

import (
	"github.com/dhiaayachi/llm-fabric/discoverer"
)

type Agent struct {
	discoverer discoverer.Discoverer
}

func NewAgent(d discoverer.Discoverer) *Agent {
	return &Agent{discoverer: d}
}
