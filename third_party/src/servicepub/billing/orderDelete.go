/*
	话单定时任务的创建和删除过期的备份话单（话单备份有效时间为7天）
*/

package billing

import (
	"github.com/astaxie/beego"
	"io/ioutil"
	"os"
	"pub/common/timer"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	backup_file_type = ".zip" // 备份文件类型
)

var file_valid_days = getBackupFileTimer()                           // 备份文件有效期（天）
var backup_path string = getBillingFolder() + "_zip/file_zip_backup" // //备份文件存放路径
var vm_backup_path string = backup_path + "/vm/"
var container_backup_path string = backup_path + "/container/"

type OrderDelete struct {
}

//构造函数
func NewOrderDelete() OrderDelete {

	OrderDelete := new(OrderDelete)
	return *OrderDelete
}

//启动任务
func (this *OrderDelete) StartOrderDeleteTimer() {

	//加入生成话单备份的任务
	timer := timer.GetScheduledExecutor()

	//timer.AddSimpleSingletonJob(this, time.Now().Add(1*time.Minute).Truncate(time.Hour), 10*time.Minute)
	//添加定时任务
	//arg1:首次执行时间，arg2：间隔时间
	//每隔10分钟检测一次
	intervalStr := getBackupDeleteTaskInterval()
	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		beego.Error("Wrong config backupDelete_interval(like: 24h(24 hours),24m(24 minutes), 24s(24 seconds) in app.conf),system will use default value(24h)")
		interval, _ = time.ParseDuration("24h")
	}
	timer.AddSimpleSingletonJob(this, time.Now().Add(10*time.Minute), interval)
}

//返回该任务的名称，用于timer中使用
func (this *OrderDelete) GetTaskName() string {

	return "Delete Order"
}

//实现定时任务的接口
func (this *OrderDelete) Run() error {

	err := deleteOrder()
	if err != nil {
		beego.Error("delete order ERROR: ", err.Error())
	}
	return err
}

func deleteOrder() error {

	//保留有效期内的备份文件，删除过期的备份
	//删除vm备份
	err := walk(vm_backup_path, backup_file_type)
	if err != nil {
		return err
	}
	//删除container备份
	err = walk(container_backup_path, backup_file_type)
	if err != nil {
		return err
	}
	return nil

}

//遍历目录
func walk(dirPth, suffix string) error {

	beego.Info("Start clean backup order ", dirPth)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return err
	}

	suffix = strings.ToUpper(suffix)
	for _, fi := range dir {
		if fi.IsDir() {
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			beego.Debug("judge file ", fi.Name(), " is not a week ago file")
			//删除指定的文件
			err = delete_order_before_aweek(dirPth + fi.Name())
			if err != nil {
				return err
			}
		}
	}
	return nil

}

//删除符合要求的文件
func delete_order_before_aweek(filename string) error {

	fileInfo, err := os.Stat(filename)
	if err != nil {
		return err
	}
	now := time.Now()
	//获取文件创建时间
	create_time_unix := reflect.ValueOf(fileInfo.Sys()).Elem().FieldByName("Ctim").Field(0).Int()
	file_create_time := time.Unix(create_time_unix, 0)

	//获取7天之前的时间
	file_timer := 24 * file_valid_days
	days, _ := time.ParseDuration("-" + strconv.Itoa(file_timer) + "h")
	before_time := now.Add(days)
	beego.Debug("file_name: ", filename, " file_create_time: ", file_create_time, " now_time: ", now, " before_time: ", before_time)

	//删除7天之前备份话单
	if file_create_time.Before(before_time) {
		beego.Info("Delete ", filename, " Create at ", file_create_time.Format("07-01-02-03 04:05:06"))
		err = os.Remove(filename)
		if err != nil {
			return err
		}
	}
	return nil

}
