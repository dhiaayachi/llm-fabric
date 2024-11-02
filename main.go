package main

import (
	"flag"
	agentv1 "github.com/dhiaayachi/llm-fabric/proto/gen/agent/v1"
	"github.com/hashicorp/serf/serf"
	"github.com/oklog/ulid/v2"
	"google.golang.org/protobuf/proto"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	bind := flag.String("bind", "", "")
	join := flag.String("join", "", "")
	flag.Parse()
	dconf := serf.DefaultConfig()
	eventCh := make(chan serf.Event)
	dconf.NodeName = ulid.Make().String()
	dconf.EventCh = eventCh
	if bind == nil || *bind == "" {
		log.Printf("--bind is mandatory")
		os.Exit(-1)
	}

	bindAddr, err := net.ResolveUDPAddr("udp", *bind)
	if err != nil {
		log.Printf("invalid bind addr %s", err.Error())
		os.Exit(-1)
	}
	dconf.MemberlistConfig.BindAddr = bindAddr.IP.String()
	dconf.MemberlistConfig.BindPort = bindAddr.Port
	dconf.Logger = log.New(os.Stdout, "serf", 0)
	s, err := serf.Create(dconf)
	if err != nil {
		log.Printf("error creating serf instance %s", err.Error())
		os.Exit(-1)
	}

	if join != nil && *join != "" {
		_, err = s.Join([]string{*join}, false)
		if err != nil {
			log.Printf("error joining serf pool %s", err.Error())
			os.Exit(-1)
		}
	}

	// add local agent to the fabric
	svc := agentv1.Agent{Id: dconf.NodeName, Address: bindAddr.String()}
	marshal, err := proto.Marshal(&svc)
	if err != nil {
		log.Printf("error marshalling json %s", err.Error())
	}
	err = s.UserEvent("new service", marshal, false)
	if err != nil {
		log.Printf("error sending user event %s", err.Error())
	}

	t := time.Tick(5000 * time.Millisecond)
	for {
		select {
		case <-t:
			log.Printf("num members %d", s.NumNodes())

		case evt := <-eventCh:
			log.Printf("got event %s", evt.EventType().String())
			switch evt.(type) {
			case serf.UserEvent:
				ue := evt.(serf.UserEvent)
				agent := agentv1.Agent{}

				err = proto.Unmarshal(ue.Payload, &agent)
				if err != nil {
					log.Printf("error unmarshalling agent %s", err.Error())
				}
				if agent.Id == dconf.NodeName {
					log.Printf("got a notification for self %s from %s", agent.Id, agent.Address)
				} else {
					log.Printf("got a new agent %s from %s", agent.Id, agent.Address)
				}
			}
		}
	}
}
