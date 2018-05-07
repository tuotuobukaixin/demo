package billing

import (
	"pub/common/timer"
	"pub/billing/models"
	"time"

	"github.com/astaxie/beego"
)

type TableHistory struct {
	beego.Controller
}

//构造函数
func NewTableHistory() TableHistory {
	TableHistory := new(TableHistory)

	return *TableHistory
}

//启动任务
func (this *TableHistory) StartTableHistory() {
	//加入生成话单的任务
	timer := timer.GetScheduledExecutor()
//timer.AddSimpleSingletonJob(this, time.Now().Add(1*time.Minute).Truncate(time.Hour), 10*time.Minute)
	timer.AddSimpleSingletonJob(this, time.Now().Add(10*time.Minute), 24*time.Hour)
}

//返回该任务的名称，用于timer中使用
func (this *TableHistory) GetTaskName() string {
	return "Generate_Table_History"
}

//实现定时任务的接口
func (this *TableHistory) Run() error {

	usageDetailRecordDAO := new(models.UsageDetailRecordDAO)

	udrs := usageDetailRecordDAO.QueryTimeOutUdr()
	beego.Debug("QueryTimeOutUdr, udrs:", udrs)

	UsageDetailRecordDAOHis := new(models.UsageDetailRecordDAOHis)

	udrsHis := recordDataToHis(udrs)

	err := UsageDetailRecordDAOHis.BatchAddUdrHis(udrsHis)
	beego.Debug("Batch add udr into table, udrs:", udrs)
	if err != nil {
		return err
	}

	err = usageDetailRecordDAO.DelAll(udrs)
	return err
}

func recordDataToHis(usageDetailRecords []models.UsageDetailRecord) []models.UsageDetailRecordHis {
	udrsHis := []models.UsageDetailRecordHis{}
	for _, udr := range usageDetailRecords {
		usageDetailRecordHis := new(models.UsageDetailRecordHis)
		usageDetailRecordHis.Billing_type = udr.Billing_type
		usageDetailRecordHis.Resource_id = udr.Resource_id
		usageDetailRecordHis.Order_id = udr.Order_id
		usageDetailRecordHis.App_id = udr.App_id
		usageDetailRecordHis.CloudserviceTypeCode = udr.CloudserviceTypeCode
		usageDetailRecordHis.Resource_type_code = udr.Resource_type_code
		usageDetailRecordHis.Resource_spec_code = udr.Resource_spec_code
		usageDetailRecordHis.Product_id = udr.Product_id
		usageDetailRecordHis.Domain_id = udr.Domain_id
		usageDetailRecordHis.Project_id = udr.Project_id
		usageDetailRecordHis.Region_code = udr.Region_code
		usageDetailRecordHis.Az_code = udr.Az_code
		usageDetailRecordHis.AccumulateFactorName = udr.AccumulateFactorName
		usageDetailRecordHis.AccumulateFactorVal = udr.AccumulateFactorVal
		usageDetailRecordHis.ExtendParams = udr.ExtendParams
		usageDetailRecordHis.BssParams = udr.BssParams
		usageDetailRecordHis.Create_time = udr.Create_time
		usageDetailRecordHis.Billing_start_time = udr.Billing_start_time
		usageDetailRecordHis.Billing_end_time = udr.Billing_end_time
		usageDetailRecordHis.Status = udr.Status
		usageDetailRecordHis.Result_code = udr.Result_code
		usageDetailRecordHis.Result_message = udr.Result_message
		usageDetailRecordHis.File_name = udr.File_name
		usageDetailRecordHis.Update_time = udr.Update_time
		usageDetailRecordHis.Process_ip = udr.Process_ip

		udrsHis = append(udrsHis, *usageDetailRecordHis)
	}
	return udrsHis
}
