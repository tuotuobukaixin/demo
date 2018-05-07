package models

type AlarmMessage struct {
	Id                        int    `orm:"auto;pk"`
	ObjectId                  string `json:"objectId" validate:"nonzero"`
	ObjectClass               string `json:"objectClass" validate:"regexp=^(Application|Service|PaaS_Component|PaaS_Middleware)$"`
	SerialNumber              int64  `json:"serialNumber" validate:"nonzero"`
	EventTime                 string `json:"eventTime" validate:"nonzero"`
	EventType                 int64  `json:"eventType" validate:"min=1,max=21"`
	ClearedType               string `json:"clearedType" validate:"regexp=^(ADAC|ADMC|NA)$"`
	EventId                   int64  `json:"eventId" validate:"nonzero"`
	EventName                 string `json:"eventName" validate:"nonzero"`
	PerceivedSeverity         string `json:"perceivedSeverity" validate:"regexp=^(Critical|Major|Minor|Warning|Info|indeterminate|cleared)$"`
	ServiceAffectingIndicator string `json:"serviceAffectingIndicator"`
	RootCaseAlarmIndicator    string `json:"rootCaseAlarmIndicator"`
	RootCaseAlarmSN           string `json:"rootCaseAlarmSN"`
	ThresholdInfo             string `json:"thresholdInfo"`
	AdditionalText            string `json:"additionalText"`
	AdditionalInformation     string `json:"additionalInformation"`
	CreatedAt                 string `json:"created_at"`
}

type AlarmMessageResult struct {
	AlarmMessage
	DetailedInformation   string `json:"detailedInformation"`
	ProposedRepairActions string `json:"proposedRepairActions"`
}

type AlarmMessages struct {
	TotalResults float64              `json:"total_results"`
	Resources    []AlarmMessageResult `json:"resources"`
}
