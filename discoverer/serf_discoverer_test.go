package discoverer

import (
	"context"
	"errors"
	"testing"
	"time"

	agentv1 "github.com/dhiaayachi/llm-fabric/proto/gen/agent/v1"
	"github.com/hashicorp/serf/serf"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/proto"
)

// MockSerf implements the Serf interface for testing purposes.
type MockSerf struct {
	mock.Mock
}

func (m *MockSerf) Join(existing []string, ignoreOld bool) (int, error) {
	args := m.Called(existing, ignoreOld)
	return args.Int(0), args.Error(1)
}

func (m *MockSerf) LocalMember() serf.Member {
	args := m.Called()
	return args.Get(0).(serf.Member)
}

// MockStore implements the Store interface for testing purposes.
type MockStore struct {
	mock.Mock
}

func (m *MockStore) Store(agent *agentv1.Agent) error {
	args := m.Called(agent)
	return args.Error(0)
}

func setupLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	return logger
}

func TestJoin_Success(t *testing.T) {
	mockSerf := new(MockSerf)
	mockStore := new(MockStore)
	logger := setupLogger()

	mockSerf.On("Join", []string{"127.0.0.1"}, true).Return(1, nil)
	mockSerf.On("LocalMember").Return(serf.Member{Name: "local-llm"})

	discoverer := &SerfDiscoverer{
		serf:   mockSerf,
		evtCh:  make(chan serf.Event, 1),
		store:  mockStore,
		logger: logger,
	}

	ctx := context.Background()
	err := discoverer.Join(ctx, []string{"127.0.0.1"})

	assert.NoError(t, err, "Join should not return an error")
	mockSerf.AssertExpectations(t)
}

func TestJoin_Failure(t *testing.T) {
	mockSerf := new(MockSerf)
	mockStore := new(MockStore)
	logger := setupLogger()

	mockSerf.On("Join", []string{"127.0.0.1"}, true).Return(0, errors.New("join error"))

	discoverer := &SerfDiscoverer{
		serf:   mockSerf,
		evtCh:  make(chan serf.Event, 1),
		store:  mockStore,
		logger: logger,
	}

	ctx := context.Background()
	err := discoverer.Join(ctx, []string{"127.0.0.1"})

	assert.Error(t, err, "Join should return an error on failure")
	mockSerf.AssertExpectations(t)
}

func TestConsumeEvts_ProcessUserEvent(t *testing.T) {
	mockSerf := new(MockSerf)
	mockStore := new(MockStore)
	logger := setupLogger()

	agent := &agentv1.Agent{Id: "test-llm", Address: "127.0.0.1"}
	payload, _ := proto.Marshal(agent)
	mockEvent := serf.UserEvent{Payload: payload}
	mockStore.On("Store", mock.MatchedBy(func(a interface{}) bool {
		a1, ok := a.(*agentv1.Agent)
		if !ok {
			return false
		}
		return proto.Equal(a1, agent)
	})).Return(nil)

	discoverer := &SerfDiscoverer{
		serf:   mockSerf,
		evtCh:  make(chan serf.Event, 1),
		store:  mockStore,
		logger: logger,
	}

	discoverer.evtCh <- mockEvent
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		time.Sleep(500 * time.Millisecond)
		cancel() // Stop the goroutine after some time to prevent an infinite loop in tests
	}()

	discoverer.consumeEvts(ctx, discoverer.evtCh, "other-llm")
}

func TestConsumeEvts_SkipSelfNotification(t *testing.T) {
	mockSerf := new(MockSerf)
	mockStore := new(MockStore)
	logger := setupLogger()

	agent := &agentv1.Agent{Id: "local-llm", Address: "127.0.0.1"}
	payload, _ := proto.Marshal(agent)
	mockEvent := serf.UserEvent{Payload: payload}

	discoverer := &SerfDiscoverer{
		serf:   mockSerf,
		evtCh:  make(chan serf.Event, 1),
		store:  mockStore,
		logger: logger,
	}

	discoverer.evtCh <- mockEvent
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		time.Sleep(500 * time.Millisecond)
		cancel() // Stop the goroutine after some time to prevent an infinite loop in tests
	}()

	discoverer.consumeEvts(ctx, discoverer.evtCh, "local-llm")

	mockStore.AssertNotCalled(t, "Store", agent)
}

func TestConsumeEvts_UnmarshalError(t *testing.T) {
	mockSerf := new(MockSerf)
	mockStore := new(MockStore)
	logger := setupLogger()

	invalidPayload := []byte("invalid data")
	mockEvent := serf.UserEvent{Payload: invalidPayload}

	discoverer := &SerfDiscoverer{
		serf:   mockSerf,
		evtCh:  make(chan serf.Event, 1),
		store:  mockStore,
		logger: logger,
	}

	discoverer.evtCh <- mockEvent
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		time.Sleep(500 * time.Millisecond)
		cancel()
	}()

	discoverer.consumeEvts(ctx, discoverer.evtCh, "other-llm")

	mockStore.AssertNotCalled(t, "Store", mock.Anything)
}

func TestNewSerfDiscoverer_ErrorCreatingSerfInstance(t *testing.T) {
	logger := setupLogger()

	conf := &serf.Config{}
	_, err := NewSerfDiscoverer(conf, nil, logger)

	assert.Error(t, err, "NewSerfDiscoverer should return an error if Serf creation fails")
}

func TestNewSerfDiscoverer_Success(t *testing.T) {
	logger := setupLogger()
	mockStore := new(MockStore)

	conf := serf.DefaultConfig()
	e := make(chan serf.Event)
	conf.EventCh = e

	mockSerf := new(MockSerf)
	mockSerf.On("LocalMember").Return(serf.Member{Name: "local-llm"})

	// Assuming the serf.Create call has been mocked appropriately in actual tests or integration
	discoverer, err := NewSerfDiscoverer(conf, mockStore, logger)

	assert.NoError(t, err, "NewSerfDiscoverer should not return an error on success")
	assert.NotNil(t, discoverer, "discoverer should not be nil")
}
