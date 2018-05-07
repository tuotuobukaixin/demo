package models

type EventMessage struct {
	Id                    int    `orm:"auto;pk"`
	ObjectId              string `json:"objectId" validate:"nonzero"`
	ObjectClass           string `json:"objectClass" validate:"regexp=^(Application|Service|PaaS_Component|PaaS_Middleware)$"`
	SerialNumber          int64  `json:"serialNumber" validate:"nonzero"`
	EventTime             string `json:"eventTime" validate:"nonzero"`
	EventType             int64  `json:"eventType" validate:"min=1,max=21"`
	AlarmIndicator        string `json:"alarmIndicator"`
	EventId               int64  `json:"eventId" validate:"nonzero"`
	EventName             string `json:"eventName" validate:"nonzero"`
	PerceivedSeverity     string `json:"perceivedSeverity" validate:"regexp=^(Info)$"`
	AdditionalText        string `json:"additionalText"`
	AdditionalInformation string `json:"additionalInformation"`
	CreatedAt             string `json:"created_at"`
}

type EventMessageResult struct {
	EventMessage
	DetailedInformation string `json:"detailedInformation"`
}

type EventMessages struct {
	TotalResults float64              `json:"total_results"`
	Resources    []EventMessageResult `json:"resources"`
}
