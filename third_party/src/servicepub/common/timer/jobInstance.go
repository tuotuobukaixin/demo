package timer

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

const (
	JOB_STATUS_TOBE_EXECUTE = "1"
	JOB_STATUS_EXECUTING    = "2"
	JOB_STATUS_EXECUTED     = "3"
)

type JobInstance struct {
	Job_sche_id      int    `orm:"auto;pk"`
	Task_name        string `orm:"size(100);index"`
	Job_type         string `orm:"size(2)"`
	Create_time      time.Time
	Expect_exec_time time.Time
	Start_time       time.Time `orm:"null"`
	End_time         time.Time `orm:"null"`
	Process_ip       string    `orm:"null;size(30)"`
	Status           string    `orm:"default(1);size(2)"`
	retry_times      int       `orm:"null"`
}

type JobInstanceDAO struct {
}

func init() {
	orm.RegisterModel(new(JobInstance))
	//modelregister.RegisterModel(new(JobInstance))
}

//创建任务实例
func (this *JobInstanceDAO) AddJobInstance(jobInstance *JobInstance) error {
	o := orm.NewOrm()
	_, err := o.Insert(jobInstance)
	if err != nil {
		beego.Error("Add task instance failed. jobInstance:", jobInstance)
	}

	return err
}

//根据任务名称更新有效任务实例的处理IP
func (this *JobInstanceDAO) UpdateIpByTaskName(taskName string, ip string) int64 {
	o := orm.NewOrm()
	num, err := o.QueryTable("JobInstance").Filter("Task_name", taskName).Filter("Status", JOB_STATUS_TOBE_EXECUTE).Filter("Process_ip", "").Update(orm.Params{
		"Process_ip": ip, "Status": JOB_STATUS_EXECUTING,
	})

	if nil != err {
		beego.Error("Update valid task instance's ip by task name failed, error:", err, "task name:", taskName, "ip:", ip)
		return 0
	}
	return num
}

//根据任务名称查询有效的任务实例信息
func (this *JobInstanceDAO) QueryValidInstanceByName(taskName string) (JobInstance, error) {
	var jobInstance JobInstance
	o := orm.NewOrm()
	err := o.QueryTable("JobInstance").Filter("Task_name", taskName).Filter("Status", JOB_STATUS_TOBE_EXECUTE).One(&jobInstance)

	return jobInstance, err
}

//当任务执行完成时，调用该函数
func (this *JobInstanceDAO) FinishJob(taskName string, ip string, beginTime time.Time, endTime time.Time) {
	o := orm.NewOrm()
	_, err := o.QueryTable("JobInstance").Filter("Task_name", taskName).Filter("Process_ip", ip).Filter("Status", JOB_STATUS_EXECUTING).Update(orm.Params{
		"Start_time": beginTime, "End_time": endTime, "Status": JOB_STATUS_EXECUTED})

	if nil != err {
		beego.Error("Update job instance failed, error:", err, "task name:", taskName, "ip:", ip, "beginTime:", beginTime, "endTime:", endTime)
	}

}
