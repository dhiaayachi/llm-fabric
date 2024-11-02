package discoverer

import (
	"context"
	agentv1 "github.com/dhiaayachi/llm-fabric/proto/gen/agent/v1"
	"github.com/hashicorp/serf/serf"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

const moduleLog = "discoverer.serf"

type Store interface {
	Store(agent *agentv1.Agent) error
	// GetByID(id string) (agentv1.Agent, error)
	// GetByTool(name string) (agentv1.Agent, error)
	// GetByCapability(capability agentv1.Capability) (agentv1.Agent, error)
}

type Serf interface {
	Join(existing []string, ignoreOld bool) (int, error)
	LocalMember() serf.Member
}

type SerfDiscoverer struct {
	serf   Serf
	evtCh  chan serf.Event
	cancel context.CancelFunc
	store  Store
	logger *logrus.Logger
}

func (s *SerfDiscoverer) Join(ctx context.Context, addresses []string) error {
	_, err := s.serf.Join(addresses, true)
	if err != nil {
		s.logger.WithError(err).WithField("module", moduleLog).Error("failed to join Serf cluster")
		return err
	}
	ctx, cancel := context.WithCancel(ctx)
	s.cancel = cancel
	s.logger.WithFields(logrus.Fields{
		"module":    moduleLog,
		"addresses": addresses,
	}).Info("joined Serf cluster successfully")
	go s.consumeEvts(ctx, s.evtCh, s.serf.LocalMember().Name)
	return nil
}

func (s *SerfDiscoverer) consumeEvts(ctx context.Context, ch chan serf.Event, name string) {
	for {
		select {
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
				agent := agentv1.Agent{}

				err := proto.Unmarshal(ue.Payload, &agent)
				if err != nil {
					s.logger.WithError(err).WithField("module", moduleLog).Error("error unmarshalling agent")
					continue
				}
				if agent.Id == name {
					s.logger.WithFields(logrus.Fields{
						"module":   moduleLog,
						"agent_id": agent.Id,
						"address":  agent.Address,
					}).Warn("received notification for self")
					continue
				}
				s.logger.WithFields(logrus.Fields{
					"module":   moduleLog,
					"agent_id": agent.Id,
					"address":  agent.Address,
				}).Info("discovered a new agent")

				err = s.store.Store(&agent)
				if err != nil {
					s.logger.WithError(err).WithField("module", moduleLog).Error("error storing agent")
					continue
				}
				s.logger.WithFields(logrus.Fields{
					"module":   moduleLog,
					"agent_id": agent.Id,
				}).Debug("agent stored successfully")
			}
		}
	}
}

func NewSerfDiscoverer(conf *serf.Config, store Store, logger *logrus.Logger) (Discoverer, error) {
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
