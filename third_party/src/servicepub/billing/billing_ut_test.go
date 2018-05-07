package billing

import (
    . "pub/common"
    . "api/common"
    . "models/billing"
    . "models/db"
    . "srvbroker/test/ut/mock"
    . "strings"
    "testing"
    "time"
)

///mock data begin---
var testcase string = ""
var billing_contents [3]string
var write_count int = 0

///mock data end---

type test_DBVisitor struct {
    DBVisitor_mock
}

func (this *test_DBVisitor) Get_OrderTbl() OrderTbl {
    if testcase == "Test_ut_Billing_billingOnPeriod" {
        return new(test_OrderTbl2)
    }

    return new(test_OrderTbl)
}

func (this *test_DBVisitor) Get_ServiceInstanceTbl() ServiceInstanceTbl {
    if testcase == "Test_ut_Billing_billingOnPeriod" {
        return new(test_instanceTbl)
    }

    return new(test_instanceTbl)
}

func (this *test_DBVisitor) Get_OrderResources(orderid string) (billing_time time.Time, vms []VMInfo, err error) {
    billing_time = Get_cur_datetime()
    vms = make([]VMInfo, 2)
    vms[0].VMID = "vm123"
    vms[0].VMID = "vm456"
    err = nil

    return
}

type test_OrderTbl struct {
    OrderTbl_mock
}

func (this *test_OrderTbl) Get_all_orders() (orders []string, err error) {
    if testcase == "Test_ut_billing_basic" {
        orders = make([]string, 1)
        orders[0] = "order_abc1234"
        err = nil
        return
    } else if testcase == "Test_ut_Billing_no_order" {
        err = nil
        return
    }

    return
}

func (this *test_OrderTbl) Update_billingtime(order_id string, billingtime time.Time) (err error) {
    return
}

func (this *test_OrderTbl) Get_order_by_orderId(order_id string) (order *OrderEntry, err error) {
    order = new(OrderEntry)
    order.Order_id = order_id
    order.Billing_type = BILLING_TYPE_ON_QUANTITY
    order.User_id = "userid01"
    order.Region_id = "regionid01"
    order.Az_id = "azid01"
    order.Billing_time = Get_cur_datetime()

    err = nil
    return
}

type test_OrderTbl2 struct {
    OrderTbl_mock
}

func (this *test_OrderTbl2) Get_all_orders() (orders []string, err error) {
    orders = make([]string, 1)
    orders[0] = "order_abc1234"
    err = nil
    return

    return
}
func (this *test_OrderTbl2) Get_order_by_orderId(order_id string) (order *OrderEntry, err error) {
    order = new(OrderEntry)
    order.Order_id = order_id
    order.Billing_type = BILLING_TYPE_ON_PERIOD
    order.User_id = "userid01"
    order.Region_id = "regionid01"
    order.Az_id = "azid01"
    order.Billing_time = Get_cur_datetime()

    err = nil
    return
}

type test_instanceTbl struct {
    ServiceInstanceTbl_mock
}
func (this*test_instanceTbl)Find_inst_by_inst_id(guid string) (inst *ServiceInstance, err error) {
    inst = new(ServiceInstance)
    inst.Guid = guid
    inst.State = SRV_INSTANCE_RUNNING
    
    return
}

type test_Broker struct {
    Servicebroker_mock
}

func (this *test_Broker) Get_billing_info(instance_id string, orderid string, vms []VMInfo, lastbillingtime time.Time, thisbillingtime time.Time) (billinginfo []BillingInfo, err error) {
    billinginfo = make([]BillingInfo, 2)
    billinginfo[0].Resinfo.CloudserviceTypeCode = "123456789"
    billinginfo[0].Resinfo.ResourceTypeCode = "a123"
    billinginfo[0].Resinfo.ResourceSpecCode = "a456"
    billinginfo[0].Resinfo.ResourceId = "ar123456"
    billinginfo[0].Resinfo.ExtendParams = "resinfo_ext_param_01"
    billinginfo[0].Resextinfo.OrderId = orderid
    billinginfo[0].Resextinfo.ProductId = "prodid_01"
    billinginfo[0].AccumulateFactorName = "Factor_time"
    billinginfo[0].AccumulateFactorVal = 0.123456
    billinginfo[0].ExtendParams = "billing_ext_param_01"

    billinginfo[1].Resinfo.CloudserviceTypeCode = "987654321"
    billinginfo[1].Resinfo.ResourceTypeCode = "b123"
    billinginfo[1].Resinfo.ResourceSpecCode = "b456"
    billinginfo[1].Resinfo.ResourceId = "br123456"
    billinginfo[1].Resinfo.ExtendParams = "resinfo_ext_param_02"
    billinginfo[1].Resextinfo.OrderId = orderid
    billinginfo[1].Resextinfo.ProductId = "prodid_02"
    billinginfo[1].AccumulateFactorName = "Factor_length"
    billinginfo[1].AccumulateFactorVal = 0.654321
    billinginfo[1].ExtendParams = "billing_ext_param_02"

    err = nil
    return
}

type test_File struct {
    File_mock
}

func (this *test_File) StableWritestring(content string) (err error) {
    billing_contents[write_count] = content
    write_count = write_count + 1
    return
}

var interval = 10
var maxfilesize = 512
var billfolder = "/home/abc"
var region = "SC"
var restype = "VM"

var billing_cfg Billing_cfg = Billing_cfg{billfolder, interval, maxfilesize, region, restype} 
func Test_ut_billing_basic(t *testing.T) {
    write_count = 0
    testcase = "Test_ut_billing_basic"

    db := new(test_DBVisitor)
    
    timer := new(Timer_mock)
    broker := new(test_Broker)
    file := new(test_File)

    billingdatagetter := New_OrderBillingData(db, broker)
    billing := New_Billing(&billing_cfg, billingdatagetter, file, timer)

    err := billing.Start()
    if err != nil {
        t.Error("billing stop failed, error: ", err)
        t.FailNow()
    }

    err = billing.Run()
    if err != nil {
        t.Error("billing failed, error: ", err)
    }

    err = billing.Stop()
    if err != nil {
        t.Error("billing start failed, error: ", err)
        t.FailNow()
    }

    t.Log("billing info:\n",
        "line 0:", billing_contents[0],
        "line 1:", billing_contents[1],
        "line 2:", billing_contents[2])

    if write_count != 3 {
        t.Error("bill gened is not right!")
    }

    if Contains(billing_contents[0], "10") && Contains(billing_contents[0], "text") {
        t.Log("billing file header is right.")
    } else {
        t.Error("billing file header is not right.")
    }

    if Contains(billing_contents[1], "20 |") &&
        Contains(billing_contents[1], "123456789") &&
        Contains(billing_contents[1], "userid01") &&
        Contains(billing_contents[1], "regionid01") &&
        Contains(billing_contents[1], "azid01") &&
        Contains(billing_contents[1], "a456") &&
        Contains(billing_contents[1], "ar123456") &&
        Contains(billing_contents[1], "order_abc1234:prodid_01") &&
        Contains(billing_contents[1], "Factor_time") &&
        Contains(billing_contents[1], "0.123456") &&
        Contains(billing_contents[1], "billing_ext_param_01") {
        t.Log("billing file record 1 is right.")
    } else {
        t.Error("billing file record 1 is not right.")
    }

    if Contains(billing_contents[1], "20 |") &&
        Contains(billing_contents[1], "987654321") &&
        Contains(billing_contents[1], "userid01") &&
        Contains(billing_contents[1], "regionid01") &&
        Contains(billing_contents[1], "azid01") &&
        Contains(billing_contents[1], "b456") &&
        Contains(billing_contents[1], "br123456") &&
        Contains(billing_contents[1], "order_abc1234:prodid_02") &&
        Contains(billing_contents[1], "Factor_length") &&
        Contains(billing_contents[1], "0.654321") &&
        Contains(billing_contents[1], "billing_ext_param_02") {
        t.Log("billing file record 2 is right.")
    } else {
        t.Error("billing file record 2 is not right.")
    }

    if Contains(billing_contents[2], "90") {
        t.Log("billing file footer is right.")
    } else {
        t.Error("billing file footer is not right.")
    }
}

func Test_ut_Billing_no_order(t *testing.T) {
    write_count = 0
    testcase = "Test_ut_Billing_no_order"

    db := new(test_DBVisitor)
    timer := new(Timer_mock)
    broker := new(test_Broker)
    file := new(test_File)

    billingdatagetter := New_OrderBillingData(db, broker)
    billing := New_Billing(&billing_cfg, billingdatagetter, file, timer)

    err := billing.Start()
    if err != nil {
        t.Error("billing stop failed, error: ", err)
        t.FailNow()
    }

    err = billing.Run()
    if err != nil {
        t.Error("billing failed, error: ", err)
    }

    err = billing.Stop()
    if err != nil {
        t.Error("billing start failed, error: ", err)
        t.FailNow()
    }

    t.Log("billing info:\n",
        "line 0:", billing_contents[0],
        "line 1:", billing_contents[1])

    if write_count != 2 {
        t.Error("bill gened is not right!")
    }

    if Contains(billing_contents[0], "10") && Contains(billing_contents[0], "text") {
        t.Log("billing file header is right.")
    } else {
        t.Error("billing file header is not right.")
    }

    if Contains(billing_contents[1], "90") {
        t.Log("billing file footer is right.")
    } else {
        t.Error("billing file footer is not right.")
    }
}

func Test_ut_Billing_billingOnPeriod(t *testing.T) {
    write_count = 0
    testcase = "Test_ut_Billing_billingOnPeriod"

    db := new(test_DBVisitor)
    timer := new(Timer_mock)
    broker := new(test_Broker)
    file := new(test_File)

    billingdatagetter := New_OrderBillingData(db, broker)
    billing := New_Billing(&billing_cfg, billingdatagetter, file, timer)

    err := billing.Start()
    if err != nil {
        t.Error("billing stop failed, error: ", err)
        t.FailNow()
    }

    err = billing.Run()
    if err != nil {
        t.Error("billing failed, error: ", err)
    }

    err = billing.Stop()
    if err != nil {
        t.Error("billing start failed, error: ", err)
        t.FailNow()
    }

    t.Log("billing info:\n",
        "line 0:", billing_contents[0],
        "line 1:", billing_contents[1])

    if write_count != 2 {
        t.Error("bill gened is not right!")
    }

    if Contains(billing_contents[0], "10") && Contains(billing_contents[0], "text") {
        t.Log("billing file header is right.")
    } else {
        t.Error("billing file header is not right.")
    }

    if Contains(billing_contents[1], "90") {
        t.Log("billing file footer is right.")
    } else {
        t.Error("billing file footer is not right.")
    }
}
