package models

type Event struct {
	Id                    int
	EventType             int64  `json:"eventType" validate:"min=1,max=21"`
	ClearedType           string `json:"clearedType" validate:"regexp=^(ADAC|ADMC|NA)$"`
	EventId               int64  `json:"eventId" validate:"nonzero"`
	EventName             string `json:"eventName" validate:"nonzero"`
	PerceivedSeverity     string `json:"perceivedSeverity" validate:"regexp=^(Critical|Major|Minor|Warning|Info|indeterminate|cleared)$"`
	DetailedInformation   string `json:"detailedInformation"`
	ProposedRepairActions string `json:"proposedRepairActions"`
	Createdat             string `json:"created_at"`
	Updatedat             string `json:"updated_at"`
	BIsSystemDefine       bool   `json:"bIsSystemDefine"`
}

type Events struct {
	TotalResults float64 `json:"total_results"`
	Resources    []Event `json:"resources"`
}

type MessageSearchInfo struct {
	ObjectId          string
	ObjectClass       string
	EventName         string
	PerceivedSeverity string
	SerialNumber      int64
	EventType         int64
	EventId           int64
	LimitNum          int64
	Offset            int64
	EventTimeStart    string
	EventTimeEnd      string
}
