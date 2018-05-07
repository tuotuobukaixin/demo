package common

import (
      "github.com/astaxie/beego"
      "time"
)

type Timer_task interface {
    Run() (err error)
}

type Timer interface{
    Start(interval int, task Timer_task) (err error) //sec
    Stop()
}

type timer_impl struct {
    timer   *time.Timer
    interval int
    task     Timer_task

    terminate bool  //flag
}

func New_timer() Timer {
    timer := new(timer_impl)
    timer.terminate = true
    
    return timer
}

func (this *timer_impl) Start(interval int, task Timer_task) (err error){
    this.interval = interval
    this.task = task
    this.terminate = false
    
    timer := time.NewTimer(time.Duration(interval) * time.Second)
    this.timer = timer  
    go this.timer_process()  

    err = nil   
    return
}

func (this *timer_impl) Stop() {
    this.terminate = true
}

func (this *timer_impl) restart_timer() {
    this.timer.Reset(time.Duration(this.interval) * time.Second)
}

func (this *timer_impl) timer_process(){
    beego.Info("timer procesing...")

    for{
        select {
            case <-this.timer.C:
                beego.Info("timer expired. now: ", Get_cur_datetime())  
                err := this.task.Run()
                if err != nil {
                    beego.Warn("timer task run failed. Task: ", this.task, ", error:", err)
                }
                
                if ! this.terminate {
                    this.restart_timer()
                }
        }
    }
    
    beego.Info("timer proces quit.")
}

