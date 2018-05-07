package legacycollectorregistrar

import (
	"encoding/json"
	"testing"

	"github.com/apcera/nats"
	"github.com/cloudfoundry/loggregatorlib/cfcomponent"
	"github.com/cloudfoundry/loggregatorlib/loggertesthelper"
	"github.com/cloudfoundry/yagnats/fakeyagnats"
	"github.com/stretchr/testify/assert"
)

func TestAnnounceComponent(t *testing.T) {
	cfc := getTestCFComponent()
	mbus := fakeyagnats.Connect()

	called := make(chan *nats.Msg, 10)
	mbus.Subscribe(AnnounceComponentMessageSubject, func(response *nats.Msg) {
		called <- response
	})

	registrar := NewCollectorRegistrar(mbus, loggertesthelper.Logger())
	registrar.announceComponent(cfc)

	expectedJson, _ := createYagnatsMessage(t, AnnounceComponentMessageSubject)

	payloadBytes := (<-called).Data
	assert.Equal(t, expectedJson, payloadBytes)
}

func TestSubscribeToComponentDiscover(t *testing.T) {
	cfc := getTestCFComponent()
	mbus := fakeyagnats.Connect()

	called := make(chan *nats.Msg, 10)
	mbus.Subscribe(DiscoverComponentMessageSubject, func(response *nats.Msg) {
		called <- response
	})

	registrar := NewCollectorRegistrar(mbus, loggertesthelper.Logger())
	registrar.subscribeToComponentDiscover(cfc)

	expectedJson, _ := createYagnatsMessage(t, DiscoverComponentMessageSubject)
	mbus.PublishRequest(DiscoverComponentMessageSubject, "unused-reply", expectedJson)

	payloadBytes := (<-called).Data
	assert.Equal(t, expectedJson, payloadBytes)
}

func createYagnatsMessage(t *testing.T, subject string) ([]byte, *nats.Msg) {

	expected := &AnnounceComponentMessage{
		Type:        "Loggregator Server",
		Index:       0,
		Host:        "1.2.3.4:5678",
		UUID:        "0-abc123",
		Credentials: []string{"user", "pass"},
	}

	expectedJson, err := json.Marshal(expected)
	assert.NoError(t, err)

	yagnatsMsg := &nats.Msg{
		Subject: subject,
		Reply:   "reply_to",
		Data:    expectedJson,
	}

	return expectedJson, yagnatsMsg
}

func getTestCFComponent() cfcomponent.Component {
	return cfcomponent.Component{
		IpAddress:         "1.2.3.4",
		Type:              "Loggregator Server",
		Index:             0,
		StatusPort:        5678,
		StatusCredentials: []string{"user", "pass"},
		UUID:              "abc123",
	}
}
