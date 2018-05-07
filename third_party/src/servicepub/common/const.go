// const
package common

const (
    SECOND_IN_MIN = 60
    MINUTE_IN_HOUR = 60
    SECOND_IN_HOUR = (SECOND_IN_MIN * MINUTE_IN_HOUR)
)

const (
    RES_TYPE_VM = "VM"
    RES_TYPE_CONTAINER = "Container"	
)

//billing_type
const (
    BILLING_TYPE_ON_DEMAND = "on-demand" 
    BILLING_TYPE_ON_PERIOD = "on-period" 
    BILLING_TYPE_ON_QUANTITY = "on-quantity"
)

//CloudserviceTypeCode
const (
	CLOUD_SERVICE_TYPE_CODE_CONTAINER = "hws.resource.type.container"
	CLOUD_SERVICE_TYPE_CODE_VM        = "hws.resource.type.vm"
)

const (
	//1=应用计费，2=服务计费
	UDR_BILLING_TYPE_APP     = "1"
	UDR_BILLING_TYPE_SERVICE = "2"
)

const (
	RESULT_CODE_SUCCESS = 0
	RESULT_CODE_FAIL    = 1
)