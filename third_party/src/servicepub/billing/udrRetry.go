package billing

import (
	"pub/common"
	"pub/common/timer"
	"pub/billing/models"
	"time"

	"github.com/astaxie/beego"
)

type UdrRetryTimer struct {
	udrWriter            UdrWriter
	usageDetailRecordDAO models.UsageDetailRecordDAO
}

//构造函数
func NewUdrRetryTimer() UdrRetryTimer {
	udrRetryTimer := new(UdrRetryTimer)
	udrRetryTimer.usageDetailRecordDAO = *new(models.UsageDetailRecordDAO)
	udrRetryTimer.udrWriter = NewUdrWriter()

	return *udrRetryTimer
}

//启动任务
func (this *UdrRetryTimer) StartUdrRetryTimer() {
	//加入生成话单的任务
	timer := timer.GetScheduledExecutor()

	//服务启动时延迟一段时间执行
	timer.AddSimpleSingletonJob(this, time.Now().Add(1*time.Minute), 10*time.Minute)

}

//返回该任务的名称，用于timer中使用
func (this *UdrRetryTimer) GetTaskName() string {
	return "UDR_RETRY_TIMER"
}

//实现定时任务的接口
func (this *UdrRetryTimer) Run() error {
	var err error
	beego.Info("Begin to execute Udr Retry Timer")
	//查询状态为计费初始化的失败或超时的话单
	for _, resourceType := range resTypesToBeUDR {
		udrs := this.usageDetailRecordDAO.QueryFailedUDRByType_Status(resourceType, models.UDR_STATUS_DB_READY)
		if len(udrs) > 0 {
			beego.Info("Retry UDRs which status is UDR_STATUS_DB_READY, udrs:", udrs)
			err = this.udrWriter.WriteUDR(udrs, resourceType)
		}
	}

	//查询状态为话单已生成的失败或超时的话单
	for _, resourceType := range resTypesToBeUDR {
		udrs := this.usageDetailRecordDAO.QueryFailedUDRByType_Status(resourceType, models.UDR_STATUS_FILE_READY)
		//将话单按文件名存放
		udrMap := map[string][]models.UsageDetailRecord{}
		for _, udr := range udrs {
			tmpUdrs := udrMap[udr.File_name]
			tmpUdrs = append(tmpUdrs, udr)
			udrMap[udr.File_name] = tmpUdrs
		}

		for retryFileName, retryUdrs := range udrMap {
			if len(retryUdrs) > 0 {
				beego.Info("Retry UDRs which status is UDR_STATUS_FILE_READY, fileName:", retryFileName, "udrs:", retryUdrs)
				err = nil
				if retryFileName == "" {
					beego.Error("retryFileName is null, please check data in usage_detail_record, udrs:", retryUdrs)
					continue
				}
				//判断文件是否存在
				if !common.Exist(retryFileName) {
					err = this.udrWriter.RewriteForNoFileUDR(retryUdrs, retryFileName)
				}
				if err == nil {
					err = this.udrWriter.SendFileToBSS(retryUdrs, retryFileName, resourceType)
				} else {
					beego.Error("Rewrite UDR file failed, udrs:", retryUdrs, "fileName:", retryFileName)
				}
			}
		}

	}

	beego.Info("Finish to execute Udr Retry Timer.")

	return err
}
