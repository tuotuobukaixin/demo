package eventclient

import (
	"errors"
	"eventclient/models"
	"eventclient/sinks"
	"kafkaclient/dataproducer"
	"paas_lager/lager"
)

//define the struct of Clients

type EventClient struct {
	Clients []Client
	logger  lager.Logger
}

func NewEmptyClient(logger lager.Logger) *EventClient {
	eventClient := EventClient{
		Clients: []Client{},
		logger:  logger,
	}
	return &eventClient
}

func NewDefaultClient(
	logger lager.Logger,
	producer *dataproducer.Producer,
	systemID int64,
) *EventClient {
	//check the config number
	sysConfig := models.LoadSystemIDFromPath(systemID)
	if sysConfig == nil {
		logger.Error("LoadSystemIDFromPath sysConfig Failed!", errors.New("load sysConfig Fail!"))
		return nil
	}
	eventClient := NewEmptyClient(logger)

	if producer == nil {
		logger.Warn("Develop mod has been opened", nil)
		return eventClient
	}
	client := sinks.NewKafkaClient(logger, producer, sysConfig)
	eventClient.RegisterSink(client)

	return eventClient
}

//the interface of Client
type Client interface {
	EventPublish(event models.EventMessage) error
	AlarmPublish(alarm models.AlarmMessage) error
	Close()
}

//Clients Register method
func (e *EventClient) RegisterSink(client Client) {
	e.Clients = append(e.Clients, client)
}

func (e *EventClient) EventPublish(event models.EventMessage) error {
	logger := e.logger
	if len(e.Clients) == 0 {
		logger.Debug("The alarm message is :", lager.Data{
			"data": event,
		})
		return nil
	}
	for _, eventClient := range e.Clients {
		err := eventClient.EventPublish(event)
		if err != nil {
			logger.Error("eventClient EventPublish Failed!", err)
			return err
		}
	}
	return nil
}

func (e *EventClient) AlarmPublish(alarm models.AlarmMessage) error {
	logger := e.logger
	if len(e.Clients) == 0 {
		logger.Debug("The alarm message is :", lager.Data{
			"data": alarm,
		})
		return nil
	}
	for _, eventClient := range e.Clients {
		err := eventClient.AlarmPublish(alarm)
		if err != nil {
			logger.Error("eventClient AlarmPublish Failed!", err)
			return err
		}
	}
	return nil
}

func (e *EventClient) Close() {
	for _, eventClient := range e.Clients {
		eventClient.Close()
	}
}
