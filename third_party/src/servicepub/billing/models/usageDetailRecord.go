package models

import (
	"pub/common"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

const (
	UDR_STATUS_DB_READY        = 1
	UDR_STATUS_FILE_READY      = 2
	UDR_STATUS_BSS_READY       = 3
	PROCESS_TIMEOUT_END_MINUTE = -30
	//超期一天不处理
	PROCESS_TIMEOUT_BEGIN_MINUTE = -1440
)

type UsageDetailRecord struct {
	Id                   int       `orm:"auto;pk"`
	Billing_type         string    `orm:"size(2)"`        //1=应用计费，2=服务计费
	Resource_id          string    `orm:"size(50)"`       //资源实例ID---VM ID或容器 ID ---部署管理模块提供
	Order_id             string    `orm:"size(50)"`       //订单号 ---API Server提供给部署管理模块，部署管理模块记录到数据库
	App_id               string    `orm:"size(50)"`       //应用id---部署管理模块提供：VM ID或容器 ID
	App_type             string    `orm:"size(3)"`        //应用类型(cf=容器，dk=Docker，vm=VM) ---部署管理模块提供
	CloudserviceTypeCode string    `orm:"size(50)"`       //云资源ResourceID的云服务类型
	Resource_type_code   string    `orm:"size(50)"`       //资源类型编码--VM,Container ---部署管理模块提供
	Resource_spec_code   string    `orm:"size(100)"`      //资源规格编码 ---部署管理模块提供
	Product_id           string    `orm:"size(50)"`       //产品ID ---API Server提供给部署管理模块，部署管理模块记录到数据库
	Domain_id            string    `orm:"size(50)"`       //租户ID ---API Server提供给部署管理模块，部署管理模块记录到数据库
	Project_id           string    `orm:"size(50)"`       //用户ProjectID ---API Server提供给部署管理模块，部署管理模块记录到数据库
	Region_code          string    `orm:"size(50)"`       //地区编码 ---API Server提供给部署管理模块，部署管理模块记录到数据库
	Az_code              string    `orm:"size(50)"`       //次级区域编码 ---API Server提供给部署管理模块，部署管理模块记录到数据库
	AccumulateFactorName string    `orm:"size(200)"`      //约定的累积因子名
	AccumulateFactorVal  float64   `orm:"size(200)"`      //AccumulateFactorName的累积值
	ExtendParams         string    `orm:"size(200)"`      //多个扩展字段，以‘，’分隔
	BssParams            string    `orm:"size(200);null"` //扩展的运营参数
	Create_time          time.Time //生成话单时间
	Billing_start_time   time.Time //计费周期开始时间
	Billing_end_time     time.Time //计费周期结束时间
	Status               int       `orm:"size(2)"`         //话单处理状态  1: 计费初始化  2: 生成话单   3: 已上传到BSS
	Result_code          int       `orm:"size(2);null"`    //处理结果   0: 成功     1: 失败
	Result_message       string    `orm:"size(2000);null"` //处理结果描述
	File_name            string    `orm:"size(100);null"`  //话单文件名称
	Update_time          time.Time `orm:"null"`            //更新时间
	Process_ip           string    `orm:"size(30);null"`   //计费主机IP
}

func init() {
	orm.RegisterModel(new(UsageDetailRecord))
	//modelregister.RegisterModel(new(UsageDetailRecord))
}

type UsageDetailRecordDAO struct {
}

//增加话单记录
func (this *UsageDetailRecordDAO) BatchAddUdr(usageDetailRecords []UsageDetailRecord) error {
	o := orm.NewOrm()
	_, err := o.InsertMulti(len(usageDetailRecords), usageDetailRecords)
	if err != nil {
		beego.Error("Batch add usage detail record failed. usageDetailRecords:", usageDetailRecords)
	}
	return err
}

//增加话单记录，接入外部传入的Ormer，保证事务一致
func (this *UsageDetailRecordDAO) AddUdrWithOrmer(ormer orm.Ormer, usageDetailRecord *UsageDetailRecord) (int64, error) {
	id, err := ormer.Insert(usageDetailRecord)
	if err != nil {
		beego.Error("Add usage detail record failed. usageDetailRecord:", usageDetailRecord)
	}
	return id, err
}

//增加话单记录，接入外部传入的Ormer，保证事务一致
func (this *UsageDetailRecordDAO) BatchAddUdrWithOrmer(ormer orm.Ormer, usageDetailRecords []UsageDetailRecord) error {
	_, err := ormer.InsertMulti(len(usageDetailRecords), usageDetailRecords)
	if err != nil {
		beego.Error("Batch add usage detail record failed. error:", err, "usageDetailRecords:", usageDetailRecords)
	}
	return err
}

//增加话单记录
func (this *UsageDetailRecordDAO) AddUsageDetailRecord(usageDetailRecord *UsageDetailRecord) error {
	o := orm.NewOrm()
	_, err := o.Insert(usageDetailRecord)
	if err != nil {
		beego.Error("Add usage detail record failed. usageDetailRecord:", usageDetailRecord)
	}
	return err
}

func (this *UsageDetailRecordDAO) QueryUDRbyStatus(status int) []UsageDetailRecord {
	var udrs []UsageDetailRecord
	_, err := orm.NewOrm().QueryTable("UsageDetailRecord").Filter("Status", status).All(&udrs)
	if nil != err {
		beego.Error("Query UDR by status failed, err: ", err, "status", status)
	}
	return udrs
}

func (this *UsageDetailRecordDAO) QueryUDRbyResType_status(resType string, status int) []UsageDetailRecord {
	var udrs []UsageDetailRecord
	timeoutEndTime := time.Now().Add(PROCESS_TIMEOUT_END_MINUTE * time.Minute)

	//只查询出成功的话单，失败由失败任务来处理
	_, err := orm.NewOrm().QueryTable("UsageDetailRecord").Filter("Status", status).Filter("Resource_type_code", resType).Filter("Result_code", common.RESULT_CODE_SUCCESS).Filter("Create_time__gte", timeoutEndTime).All(&udrs)
	if nil != err {
		beego.Error("Query UDR by status failed, err: ", err, "status:", status, "Resource_type_code:", resType)
	}
	return udrs
}

func (this *UsageDetailRecordDAO) UpdateUDR(usageDetailRecord *UsageDetailRecord) {
	_, err := orm.NewOrm().Update(usageDetailRecord)
	if err != nil {
		beego.Error("Update usage detail record failed. error:", err, "usageDetailRecord:", usageDetailRecord)
	}

}

//查询出指定状态失败或超期的话单
func (this *UsageDetailRecordDAO) QueryFailedUDRByType_Status(resType string, status int) []UsageDetailRecord {
	var udrs []UsageDetailRecord

	//超时时间
	timeoutBeginTime := time.Now().Add(PROCESS_TIMEOUT_BEGIN_MINUTE * time.Minute)
	timeoutEndTime := time.Now().Add(PROCESS_TIMEOUT_END_MINUTE * time.Minute)

	condition := orm.NewCondition()
	//失败条件，但不处理超期一天的失败话单
	failedCond := condition.And("Result_code", common.RESULT_CODE_FAIL).And("Create_time__gte", timeoutBeginTime)

	//处理超期条件
	timeOutCond := condition.And("Result_code", common.RESULT_CODE_SUCCESS).And("Create_time__gte", timeoutBeginTime).And("Create_time__lt", timeoutEndTime)

	filterCond := condition.And("Status", status).And("Resource_type_code", resType).AndCond(failedCond.OrCond(timeOutCond))

	_, err := orm.NewOrm().QueryTable("UsageDetailRecord").SetCond(filterCond).All(&udrs)
	if nil != err {
		beego.Error("Query failed UDR by status failed, err: ", err, "status:", status, "Resource_type_code:", resType)
	}
	return udrs
}

func (this *UsageDetailRecordDAO) UpdateUdrAfterError(id int, resultCode int, resultMessage string, ip string, updateTime time.Time) {
	_, err := orm.NewOrm().QueryTable("UsageDetailRecord").Filter("Id", id).Update(orm.Params{
		"Result_code":    resultCode,
		"Result_message": resultMessage,
		"Process_ip":     ip,
		"Update_time":    updateTime,
	})

	if err != nil {
		beego.Error("Update UDR after operation error, err: ", err, "Id:", id, "Result_code:", resultCode, "Result_message:", resultMessage)
	}
}

func (this *UsageDetailRecordDAO) UpdateUdrAfterSuccess(id int, resultCode int, resultMessage string, ip string, updateTime time.Time, status int, fileName string) {
	_, err := orm.NewOrm().QueryTable("UsageDetailRecord").Filter("Id", id).Update(orm.Params{
		"Result_code":    resultCode,
		"Result_message": resultMessage,
		"Process_ip":     ip,
		"Update_time":    updateTime,
		"Status":         status,
		"File_name":      fileName,
	})

	if err != nil {
		beego.Error("Update UDR after operation error, err: ", err, "Id:", id, "Result_code:", resultCode, "Result_message:", resultMessage)
	}
}

func (this *UsageDetailRecordDAO) QueryTimeOutUdr() []UsageDetailRecord {
	var udrs []UsageDetailRecord

	//超时时间
	timeoutBeginTime := time.Now().Add(PROCESS_TIMEOUT_BEGIN_MINUTE * time.Minute)

	condition := orm.NewCondition()

	//处理超期条件
	timeOutCond := condition.And("Create_time__lt", timeoutBeginTime)

	filterCond := condition.AndCond(timeOutCond)

	_, err := orm.NewOrm().QueryTable("UsageDetailRecord").SetCond(filterCond).All(&udrs)
	if nil != err {
		beego.Error("Query failed UDR by status failed, err: ", err)
	}
	return udrs
}

func (this *UsageDetailRecordDAO) DelAll(usageDetailRecords []UsageDetailRecord) error {
	o := orm.NewOrm()
	for _, udr := range usageDetailRecords {
		_, err := o.Delete(&udr)
		if err == nil {
			beego.Info("delete a record from Record table success, Resource_id:", udr.Resource_id, "App_id", udr.App_id)
		} else {
			beego.Error("del a record from  Record table fail, will Rollback: ", err)
			return err
		}
	}

	return nil
}
