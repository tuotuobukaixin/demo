package models

import (
"time"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type UsageDetailRecordDAOHis struct {
}
type UsageDetailRecordHis struct {
	Id                   int       `orm:"auto;pk"`
	Billing_type         string    `orm:"size(2)"`   //1=应用计费，2=服务计费
	Resource_id          string    `orm:"size(50)"`  //资源实例ID---VM ID或容器 ID ---部署管理模块提供
	Order_id             string    `orm:"size(50)"`  //订单号 ---API Server提供给部署管理模块，部署管理模块记录到数据库
	App_id               string    `orm:"size(50)"`  //应用id---部署管理模块提供：VM ID或容器 ID
	App_type             string    `orm:"size(3)"`   //应用类型(cf=容器，dk=Docker，vm=VM) ---部署管理模块提供
	CloudserviceTypeCode string    `orm:"size(50)"`  //云资源ResourceID的云服务类型
	Resource_type_code   string    `orm:"size(50)"`  //资源类型编码--VM,Container ---部署管理模块提供
	Resource_spec_code   string    `orm:"size(100)"` //资源规格编码 ---部署管理模块提供
	Product_id           string    `orm:"size(50)"`  //产品ID ---API Server提供给部署管理模块，部署管理模块记录到数据库
	Domain_id            string    `orm:"size(50)"`  //租户ID ---API Server提供给部署管理模块，部署管理模块记录到数据库
	Project_id           string    `orm:"size(50)"`  //用户ProjectID ---API Server提供给部署管理模块，部署管理模块记录到数据库
	Region_code          string    `orm:"size(50)"`  //地区编码 ---API Server提供给部署管理模块，部署管理模块记录到数据库
	Az_code              string    `orm:"size(50)"`  //次级区域编码 ---API Server提供给部署管理模块，部署管理模块记录到数据库
	AccumulateFactorName string    `orm:"size(200)"` //约定的累积因子名
	AccumulateFactorVal  float64   `orm:"size(200)"` //AccumulateFactorName的累积值
	ExtendParams         string    `orm:"size(200)"` //多个扩展字段，以‘，’分隔
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
	orm.RegisterModel(new(UsageDetailRecordHis))
	//modelregister.RegisterModel(new(UsageDetailRecord))
}

//增加话单记录
func (this *UsageDetailRecordDAOHis) BatchAddUdrHis(usageDetailRecords []UsageDetailRecordHis) error {
	o := orm.NewOrm()
	_, err := o.InsertMulti(len(usageDetailRecords), usageDetailRecords)
	if err != nil {
		beego.Error("Batch add usage detail record failed. usageDetailRecords:", err)
	}
	return err
}
