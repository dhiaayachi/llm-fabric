package discoverer

import (
	"context"
	"github.com/dhiaayachi/llm-fabric/discoverer/store"
	agentinfo "github.com/dhiaayachi/llm-fabric/proto/gen/agent_info/v1"
	"github.com/hashicorp/serf/serf"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"time"
)

const moduleLog = "discoverer.serf"

type Serf interface {
	Join(existing []string, ignoreOld bool) (int, error)
	LocalMember() serf.Member
	UserEvent(name string, payload []byte, coalesce bool) error
}

type SerfDiscoverer struct {
	serf   Serf
	evtCh  chan serf.Event
	cancel context.CancelFunc
	store  store.Store
	logger *logrus.Logger
	agent  *agentinfo.AgentsNodeInfo
}

func (s *SerfDiscoverer) GetAgents() []*agentinfo.AgentsNodeInfo {
	return s.store.GetAll()
}

func (s *SerfDiscoverer) GetDispatchers() []*agentinfo.AgentsNodeInfo {
	agents := s.store.GetAll()
	dispatchers := make([]*agentinfo.AgentsNodeInfo, 0)
	for _, agent := range agents {
		if agent.Agents[0].IsDispatcher {
			dispatchers = append(dispatchers, agent)
		}
	}
	return dispatchers
}

func (s *SerfDiscoverer) Join(ctx context.Context, addresses []string, agent *agentinfo.AgentsNodeInfo) error {
	_, err := s.serf.Join(addresses, true)
	if err != nil {
		s.logger.WithError(err).WithField("module", moduleLog).Error("failed to join Serf cluster")
		return err
	}
	s.logger.WithFields(logrus.Fields{
		"module":    moduleLog,
		"addresses": addresses,
	}).Info("joined Serf cluster successfully")
	s.agent = agent
	go s.run(ctx, s.evtCh, 5*time.Second)
	return nil
}

func (s *SerfDiscoverer) run(ctx context.Context, ch chan serf.Event, tickDelay time.Duration) {
	after := time.After(1 * time.Millisecond)

	for {
		select {
		case <-after:
			s.agent.Node.Address = s.serf.LocalMember().Addr.String()
			marshal, err := proto.Marshal(s.agent)
			s.logger.WithField("module", moduleLog).Debug("sending agent_info info")
			if err != nil {
				s.logger.WithError(err).Error("failed to marshal agent_info info")
				return
			}
			err = s.serf.UserEvent("agent_info broadcast", marshal, true)
			if err != nil {
				s.logger.WithError(err).Error("failed to broadcast agent_info info")
			}
			after = time.After(tickDelay)
		case <-ctx.Done():
			s.logger.WithField("module", moduleLog).Info("context cancelled, stopping event consumption")
			return

		case evt := <-ch:
			s.logger.WithFields(logrus.Fields{
				"module":     moduleLog,
				"event_type": evt.EventType().String(),
			}).Debug("received event")
			switch evt.(type) {
			case serf.UserEvent:
				ue := evt.(serf.UserEvent)
				agentsNode := agentinfo.AgentsNodeInfo{}

				err := proto.Unmarshal(ue.Payload, &agentsNode)
				if err != nil {
					s.logger.WithError(err).WithField("module", moduleLog).Error("error unmarshalling llm")
					continue
				}

				for _, agent := range agentsNode.Agents {

					s.logger.WithFields(logrus.Fields{
						"module":   moduleLog,
						"agent_id": agent.Id,
						"address":  agentsNode.Node.Address,
					}).Debug("discovered a new llm")

					err = s.store.Store(agent, agentsNode.Node)
					if err != nil {
						s.logger.WithError(err).WithField("module", moduleLog).Error("error storing llm")
						continue
					}
					s.logger.WithFields(logrus.Fields{
						"module":   moduleLog,
						"agent_id": agent.Id,
					}).Debug("llm stored successfully")
				}
			}
		}
	}
}

func NewSerfDiscoverer(conf *serf.Config, store store.Store, logger *logrus.Logger) (Discoverer, error) {
	e := make(chan serf.Event)
	conf.EventCh = e
	s, err := serf.Create(conf)
	if err != nil {
		logger.WithError(err).WithField("module", moduleLog).Error("failed to create Serf instance")
		return nil, err
	}

	logger.WithField("module", moduleLog).Info("created SerfDiscoverer successfully")
	return &SerfDiscoverer{serf: s, evtCh: e, store: store, logger: logger}, nil
}
