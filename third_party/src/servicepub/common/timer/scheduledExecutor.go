package timer

import (
	"pub/common"
	"errors"
	"sort"
	"time"

	"github.com/astaxie/beego"
)

const (
	JOB_TYPE_SINGLETON = "2"
	Jobabc             = "2"
)

var scheduledExecutor *ScheduledExecutor = nil

//任务调度执行器
type ScheduledExecutor struct {
	jobList []*Job
	running bool
	newJob  chan *Job
}

//调度工作
type Job struct {
	schedule Schedule

	nextExecTime time.Time

	preExecTime time.Time

	task Task

	jobType string
}

//调度器，负责管理任务的执行时间，需要实现Next方法
type Schedule interface {
	//计算下一次任务执行时间
	Next(time.Time) time.Time
}

//任务接口，由具体的任务来实现
type Task interface {
	Run() error

	//任务做为唯一键，不允许重名
	GetTaskName() string
}

//初始化调试执行器
func newScheduledExecutor() *ScheduledExecutor {
	return &ScheduledExecutor{
		newJob:  make(chan *Job),
		running: false,
	}
}

//单例模式
func GetScheduledExecutor() *ScheduledExecutor {
	if scheduledExecutor == nil {
		scheduledExecutor = newScheduledExecutor()
	}
	return scheduledExecutor
}

//增加简单调度任务
func (this *ScheduledExecutor) AddSimpleJob(task Task, firstExecTime time.Time, interval time.Duration) {
	this.addJob(task, firstExecTime, interval, "")
}

//增加单例任务
func (this *ScheduledExecutor) AddSimpleSingletonJob(task Task, firstExecTime time.Time, interval time.Duration) {
	this.addJob(task, firstExecTime, interval, JOB_TYPE_SINGLETON)
}

func (this *ScheduledExecutor) addJob(task Task, firstExecTime time.Time, interval time.Duration, jobType string) error {
	beego.Info("Add a simple job,", "task name:", task.GetTaskName(), "start time:", firstExecTime, "schedule interval:", interval)
	schedule := SimpleSchedule{firstExecTime: firstExecTime, interval: interval}

	//判断任务是否重名
	for _, j := range this.jobList {
		if task.GetTaskName() == j.task.GetTaskName() {
			beego.Error("Not allowed to add a repeated task, task name:", task.GetTaskName())
			return errors.New("Repeated task")
		}
	}

	job := &Job{schedule: &schedule, nextExecTime: firstExecTime, task: task, jobType: jobType}

	if jobType == JOB_TYPE_SINGLETON {
		//对于周期性单例任务通过数据库来进行任务分配
		this.createInstanceIfNotExist(*job)
	}

	if !this.running {
		beego.Info("Schedule executor has not started, the job will be executed when schedule executor start")
		this.jobList = append(this.jobList, job)
		return nil
	}
	beego.Info("Wait to add the job into runing ScheduledExecutor.")
	this.newJob <- job

	return nil
}

//启动任务调度框架
func (this *ScheduledExecutor) Start() {
	this.running = true
	beego.Info("Start schedule executor.")
	go this.run()
}

//循环调度所有的任务
func (this *ScheduledExecutor) run() {
	for {
		beego.Info("Begin to find next execution time.")
		//对工作按执行时间先后排序
		sort.Sort(byTime(this.jobList))
		now := time.Now()
		//获取第一个可执行工作的执行时间
		var firstJobExecTime time.Time
		if len(this.jobList) == 0 || this.jobList[0].nextExecTime.IsZero() {
			beego.Info("There is no job or execute time is 0, wait to add job.")
			//如果获取不到任务或任务无执行时间，则进行等待
			firstJobExecTime = now.AddDate(10, 0, 0)
		} else {
			firstJobExecTime = this.jobList[0].nextExecTime
		}
		beego.Info("Next execution time:", firstJobExecTime.Sub(now))

		select {
		case now = <-time.After(firstJobExecTime.Sub(now)):
			beego.Info("It is time to execute jobs.")
			//看哪些工作时间与firstJobExecTime相同，如果相同则执行该工作
			for _, job := range this.jobList {

				//如果任务与执行时间
				if job.nextExecTime != firstJobExecTime {
					break
				}
				nextExecTime := job.schedule.Next(firstJobExecTime)
				go this.runJob(*job, nextExecTime)
				job.preExecTime = firstJobExecTime
				job.nextExecTime = nextExecTime

			}
		//新任务加入，触发执行该条件
		case newJob := <-this.newJob:
			beego.Info("Add a new job into runing schedule executor.")
			this.jobList = append(this.jobList, newJob)
			//newJob.nextExecTime = newJob.schedule.Next(now)
		}

		beego.Info("Finish to execute a number of jobs.")
	}

}

func (this *ScheduledExecutor) runJob(job Job, nextExecTime time.Time) {
	if job.jobType == JOB_TYPE_SINGLETON {
		//如果本机抢到了任务的锁，那就本机来执行
		if this.lockJob(job.task.GetTaskName()) {
			beego.Info("The server lock the job, begin to execute the task:", job.task.GetTaskName())
			execBeginTime := time.Now()
			job.task.Run()
			execEndTime := time.Now()
			//更新执行时间和状态
			this.finishJob(job, execBeginTime, execEndTime)
			//创建下一个应用实例
			job.nextExecTime = nextExecTime
			this.createNextInstance(job)
		} else {
			beego.Info("The server did not get job lock, task:", job.task.GetTaskName(), "ip:", common.GetLocalIP())
		}

	} else {
		beego.Info("Begin to execute the task:", job.task.GetTaskName())
		job.task.Run()
	}
}

func (this *ScheduledExecutor) createInstanceIfNotExist(job Job) {
	taskName := job.task.GetTaskName()
	dao := new(JobInstanceDAO)
	tempJob, _ := dao.QueryValidInstanceByName(taskName)
	if tempJob.Task_name == "" {
		beego.Info("The task is not exist in DB, taskName:", taskName)
		this.createNextInstance(job)
	}

}

//向数据库中创建任务实例
func (this *ScheduledExecutor) createNextInstance(job Job) {
	jobInstance := new(JobInstance)
	jobInstance.Job_type = JOB_TYPE_SINGLETON
	jobInstance.Create_time = time.Now()
	jobInstance.Expect_exec_time = job.nextExecTime
	taskName := job.task.GetTaskName()
	jobInstance.Task_name = taskName
	jobInstance.Status = JOB_STATUS_TOBE_EXECUTE
	jobInstance.Process_ip = ""
	dao := new(JobInstanceDAO)
	dao.AddJobInstance(jobInstance)
}

//向数据库来申请并锁定任务
func (this *ScheduledExecutor) lockJob(taskName string) bool {
	localIP := common.GetLocalIP()
	if localIP == "" {
		beego.Error("Can't lock job, because it can't get local ip, taskName", taskName)
		return false
	}
	beego.Debug("current IP is", localIP)
	dao := new(JobInstanceDAO)
	updates := dao.UpdateIpByTaskName(taskName, localIP)
	if updates > 0 {
		return true
	}
	return false
}

func (this *ScheduledExecutor) finishJob(job Job, beginTime time.Time, endTime time.Time) {
	dao := new(JobInstanceDAO)
	dao.FinishJob(job.task.GetTaskName(), common.GetLocalIP(), beginTime, endTime)
}

//实现Job按执行时间先后顺序排序
type byTime []*Job

func (s byTime) Len() int {
	return len(s)
}

//返回元素i是否排在元素j的前面
func (s byTime) Less(i, j int) bool {
	return s[i].nextExecTime.Before(s[j].nextExecTime)
}

//交换元素i,j的位置
func (s byTime) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

//任务调度器的一种实现：简单任务调度器，从一个时间点开始，周期性（interval）进行调度
type SimpleSchedule struct {
	firstExecTime time.Time
	interval      time.Duration
}

func (this *SimpleSchedule) Next(t time.Time) time.Time {
	//将执行时间精确到秒级
	return t.Add(this.interval - time.Duration(t.Nanosecond())*time.Nanosecond)
}
