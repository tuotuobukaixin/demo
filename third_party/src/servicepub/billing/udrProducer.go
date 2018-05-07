package billing

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"pub/billing/models"
	. "pub/common"
	"pub/common/timer"
	"strconv"
	"time"
)

type UdrProducer struct {
	billingDataReader BillingDataReader
}

//计费数据获取接口
type BillingDataReader interface {
	//传入Ormer对象，保证事务一致
	GetBillingDatas(ormer orm.Ormer, maxRow int) ([]UnifiedBillingData, error)
	StopApplication(ormer orm.Ormer, resId string) (UnifiedBillingData, error)
}

//构造函数
func NewUdrProducer(billingDataReader BillingDataReader) UdrProducer {
	udrProducer := new(UdrProducer)

	udrProducer.billingDataReader = billingDataReader
	return *udrProducer
}

//启动任务
func (this *UdrProducer) StartUdrProucer() {
	if this.billingDataReader == nil {
		beego.Error("Start UdrProucer failed because BillingDataReader is nil.")
		return
	}

	//加入生成话单的任务
	timer := timer.GetScheduledExecutor()
	//timer.AddSimpleSingletonJob(this, time.Now().Add(1*time.Hour).Truncate(time.Hour), 1*time.Hour)

	now := time.Now()
	nextMin := (now.Minute()/10 + 1) * 10
	timer.AddSimpleSingletonJob(this, now.Truncate(time.Hour).Add(time.Duration(nextMin)*time.Minute), 10*time.Minute)

}

//返回该任务的名称，用于timer中使用
func (this *UdrProducer) GetTaskName() string {
	return "Generate_Cycle_UDR"
}

func getReadOrderNumOnce() int {
	var numStr = beego.AppConfig.String("read_oder_num_once")
	//验证是否是一个正整数（合法的备份有效期）
	num, err := strconv.Atoi(numStr)
	if err != nil {
		beego.Error("read_oder_num_once wrong,please check app.conf, system will use default value(30).")
		return 30
	}
	return num
}

var MaxRow = getReadOrderNumOnce()

//实现定时任务的接口
func (this *UdrProducer) Run() error {
	//每次读取话单数量
	maxRow := MaxRow

	ormer := orm.NewOrm()
	usageDetailRecordDAO := new(models.UsageDetailRecordDAO)
	var err error = nil
	var unifiedBillingDatas []UnifiedBillingData

	//如果出错，则回滚事务
	defer func() {
		//beego.Info("error:", err)
		if err != nil {
			ormer.Rollback()
			beego.Error("Rollback UDR producer transaction because of the error:", err)
		} else {
			ormer.Commit()
		}
	}()

	//循环计费数据获取接口中读取数据
	for {
		//开启事务
		ormer.Begin()
		unifiedBillingDatas, err = this.billingDataReader.GetBillingDatas(ormer, maxRow)
		beego.Debug("Finish to get billing info from billingDataReader once, billingDatas:", unifiedBillingDatas)

		//如果数据获取完成，则跳出循环
		if err == ENoMoreData || len(unifiedBillingDatas) == 0 {
			ormer.Commit()
			break
		}
		udrs := billingDatasToUdrs(unifiedBillingDatas)

		err = usageDetailRecordDAO.BatchAddUdrWithOrmer(ormer, udrs)
		beego.Debug("Batch add udr into table, udrs:", udrs)
		if err != nil {
			return err
		}
		//每取一次提交
		ormer.Commit()
	}
	beego.Info("Finish to write billing data into udr table from billingDataReader.")
	writer := NewUdrWriter()
	writer.WriterCycleUDR()

	return err
}

func (this *UdrProducer) GenerateUdrAfterStop(resId string) error {
	ormer := orm.NewOrm()
	var err error
	var id int64
	var unifiedBillingData UnifiedBillingData
	//如果出错，则回滚事务
	defer func() {
		if err != nil {
			ormer.Rollback()
			beego.Error("Rollback UDR producer transaction because of the error:", err)
		} else {
			ormer.Commit()
		}

	}()
	//开启事务
	ormer.Begin()
	unifiedBillingData, err = this.billingDataReader.StopApplication(ormer, resId)
	if err != nil {
		return err
	}

	usageDetailRecord := billingDataToUdr(unifiedBillingData)
	usageDetailRecordDAO := new(models.UsageDetailRecordDAO)
	id, err = usageDetailRecordDAO.AddUdrWithOrmer(ormer, &usageDetailRecord)

	beego.Debug("Add UDR into DB, usageDetailRecord:", usageDetailRecord)
	if err != nil {
		return err
	}
	//每取一次提交
	ormer.Commit()
	//返回的ID回写
	usageDetailRecord.Id = int(id)
	beego.Info("Finish to write UDR to DB.")

	writer := NewUdrWriter()
	//写话单, 目前只有容器应用调到此函数，暂时参数写死
	writer.WriteUDR([]models.UsageDetailRecord{usageDetailRecord}, "hws.resource.type.container")
	return nil
}

func billingDatasToUdrs(unifiedBillingDatas []UnifiedBillingData) []models.UsageDetailRecord {
	udrs := []models.UsageDetailRecord{}
	for _, billingData := range unifiedBillingDatas {
		usageDetailRecord := billingDataToUdr(billingData)
		udrs = append(udrs, usageDetailRecord)
	}

	return udrs
}

//将billing转成话单对象
func billingDataToUdr(billingData UnifiedBillingData) models.UsageDetailRecord {
	usageDetailRecord := new(models.UsageDetailRecord)
	usageDetailRecord.App_id = billingData.AppId
	usageDetailRecord.App_type = billingData.AppType
	usageDetailRecord.Az_code = billingData.AzCode
	usageDetailRecord.Billing_start_time = billingData.BillingStartTime
	usageDetailRecord.Billing_end_time = billingData.BillingEndTime
	usageDetailRecord.Billing_type = billingData.BillingType
	if billingData.CreateTime.IsZero() {
		usageDetailRecord.Create_time = time.Now()
	} else {
		usageDetailRecord.Create_time = billingData.CreateTime
	}
	usageDetailRecord.Domain_id = billingData.DomainId
	usageDetailRecord.Order_id = billingData.Resextinfo.OrderId
	usageDetailRecord.Process_ip = GetLocalIP()
	usageDetailRecord.Product_id = billingData.Resextinfo.ProductId
	usageDetailRecord.Project_id = billingData.ProjectId
	usageDetailRecord.Region_code = billingData.RegionCode
	usageDetailRecord.Resource_id = billingData.Resinfo.ResourceId
	usageDetailRecord.CloudserviceTypeCode = billingData.Resinfo.CloudserviceTypeCode
	usageDetailRecord.Resource_spec_code = billingData.Resinfo.ResourceSpecCode
	usageDetailRecord.Resource_type_code = billingData.Resinfo.ResourceTypeCode
	usageDetailRecord.Status = models.UDR_STATUS_DB_READY
	usageDetailRecord.ExtendParams = billingData.ExtendParams
	usageDetailRecord.AccumulateFactorName = billingData.AccumulateFactorName
	usageDetailRecord.AccumulateFactorVal = billingData.AccumulateFactorVal
	usageDetailRecord.BssParams = billingData.BssParams

	return *usageDetailRecord
}
