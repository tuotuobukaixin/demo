package common

import (
    "time"
)

type ResInfo struct {
    CloudserviceTypeCode string
    ResourceTypeCode     string
    ResourceSpecCode     string
    ResourceId           string
    ExtendParams         string
    ResExentParam        string
}

type ResExtInfo struct {
    OrderId   string
    ProductId string
}

type BillingInfo struct {
    Resinfo              ResInfo
    Resextinfo           ResExtInfo
    AccumulateFactorName string
    AccumulateFactorVal  float64
    ExtendParams         string
}

//后续废弃
type BillingRecord struct {
    Billinginfo     []BillingInfo 
    Lastbillingtime time.Time 
    Thisbillingtime time.Time 
    Userid          string 
    Regionid        string 
    Azid            string
}
type BillingDataGetter interface {
    Begin() (err error)
    Get_data() (records []BillingRecord, err error)
    End() (err error)
    /*
        usage:
        1) Begin
        2) Get_data can be called many times, return ENoMoreData when there is no more billing data
        3) End
    */
}

type ServerInfo struct {
    Endpoint string `json:"endpoint"`
    Account  string `json:"account"`
    Passwd   string `json:"passwd"`
}

//bss
type Resource struct {
    TenantId      string         `json:"tenantId"` 
    ResourceInfos []ResourceInfo `json:"resourceInfos"` 
}

type ResourceInfo struct {
    ResourceId   string `json:"resourceId"` 
    ResourceType int    `json:"resourceType"` 
}


