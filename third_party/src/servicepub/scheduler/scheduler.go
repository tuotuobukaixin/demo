package scheduler

import (
    "sort"
    "time"

    "github.com/astaxie/beego"
)

type Task interface {
    Run() (err error)
    GetTaskName() string
}
type TimingTaskMgr interface {
    Start()
    AddCycleTask(task Task, firstExecTime time.Time, interval time.Duration)
}

type timing_task_mgr_impl struct {
    jobList []*Job
    running bool
    newJob  chan *Job
}

func New_TimingTaskMgr() TimingTaskMgr {
    return &timing_task_mgr_impl{
        newJob:  make(chan *Job),
        running: false,
    }
}

func (this *timing_task_mgr_impl) AddCycleTask(task Task, firstExecTime time.Time, interval time.Duration) {
    beego.Info("Add a cycle task,", "task name:", task.GetTaskName(), "start time:", firstExecTime, "schedule interval:", interval)
    schedule := cycleScheduler{firstExecTime: firstExecTime, interval: interval}

    job := &Job{schedule: &schedule, nextExecTime: firstExecTime, task: task}
    if !this.running {
        beego.Info("Schedule executor has not started, the job will be executed when schedule executor start")
        this.jobList = append(this.jobList, job)
        return
    }
    beego.Info("Wait to add the job into runing ScheduledExecutor.")
    this.newJob <- job
}

type scheduler interface {   
    Next(time.Time) time.Time
}

type Job struct {
    schedule scheduler
    nextExecTime time.Time
    preExecTime time.Time
    task Task
}

func (this *timing_task_mgr_impl) Start() {
    this.running = true
    beego.Info("Start schedule executor.")
    go this.run()
}

func (this *timing_task_mgr_impl) run() {
    for {
        beego.Info("Timing task begin to execute..")
    
        sort.Sort(byTime(this.jobList))
        now := time.Now()
    
        var nextExecTime time.Time
        if len(this.jobList) == 0 || this.jobList[0].nextExecTime.IsZero() {
            beego.Debug("There is no job or execute time is 0, wait to add job.")
    
            nextExecTime = now.AddDate(10, 0, 0)
        } else {
            nextExecTime = this.jobList[0].nextExecTime
        }
        if nextExecTime.Before(now) {
            //the "nextExecTime" has passed for some reason, it can be trigered imediately
            nextExecTime = now
        }
        beego.Debug("Next execution time:", nextExecTime.Sub(now))

        select {
        case now = <-time.After(nextExecTime.Sub(now)):
            beego.Debug("It is time to execute jobs.")
                for _, job := range this.jobList {

                    if job.nextExecTime != nextExecTime {
                    break
                }
                beego.Debug("now execute task:", job.task.GetTaskName())
                go job.task.Run()
                
                job.preExecTime = nextExecTime
                job.nextExecTime = job.schedule.Next(now)
                beego.Debug("next execute task:", job.task.GetTaskName(), 
                    job.preExecTime, job.nextExecTime)
            }
            case newJob := <-this.newJob:
            beego.Info("Add a new job into runing schedule executor.")
            this.jobList = append(this.jobList, newJob)
            //newJob.nextExecTime = newJob.schedule.Next(now)
        }

        beego.Info("Timing task  finished.")
    }

}

type cycleScheduler struct {
    firstExecTime time.Time
    interval      time.Duration
}

func (this *cycleScheduler) Next(t time.Time) time.Time {
    return t.Add(this.interval - time.Duration(t.Nanosecond())*time.Nanosecond)
}

type byTime []*Job

func (s byTime) Len() int {
    return len(s)
}

func (s byTime) Less(i, j int) bool {
    return s[i].nextExecTime.Before(s[j].nextExecTime)
}

func (s byTime) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}
