package common

import "time"


type UnifiedBillingData struct {
	BillingInfo
	BssParams            string
	BillingType          string     //1=应用计费，2=服务计费
	ProjectId            string
	RegionCode           string
	AzCode               string
	CreateTime           time.Time
	BillingStartTime     time.Time
	BillingEndTime       time.Time
	AppId                string
	AppType              string
	DomainId             string
}
