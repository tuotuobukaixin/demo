package billing

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"pub/billing/models"
	"pub/common"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

const (
	HEADER_FLAG                     = "10"
	RECORD_FLAG                     = "20"
	FOOTER_FLAG                     = "90"
	FILE_VERSION                    = "01"
	FILE_TYPE                       = "text"
	HEADER_SIZE                     = 50
	FOOTER_SIZE                     = 50
	UDR_FILE_PERMISSION os.FileMode = 0644
)

var resTypesToBeUDR []string
var fileMaxSize = getFileMaxSize()
var backupFileName = "/file_zip_backup"

type UdrWriter struct {
	file                 common.File
	ipLast3Pos           string
	usageDetailRecordDAO models.UsageDetailRecordDAO
}

func init() {
	resTypesToBeUDR = getResTypeToBeUDR()
}

func NewUdrWriter() UdrWriter {
	udrWriter := new(UdrWriter)

	//取得IP后3位
	ipLast3Pos := strings.Split(common.GetLocalIP(), ".")[3]
	if len(ipLast3Pos) != 3 {
		ipLast3Pos = "0" + ipLast3Pos
	}

	udrWriter.ipLast3Pos = ipLast3Pos
	udrWriter.file = common.New_File()
	udrWriter.usageDetailRecordDAO = *new(models.UsageDetailRecordDAO)

	//初始化话单生成文件目录
	tempBillingFolder := getBillingFolder()
	finalBillingFolder := tempBillingFolder + "_zip"
	err := udrWriter.file.Dir_create(tempBillingFolder)
	if err != nil {
		beego.Error("Write UDR file failed because billing folder create failed. error: ", err, "billing_folder:", tempBillingFolder)
	}

	//初始化话单备份文件目录
	finalBillingBackupFolder := finalBillingFolder + backupFileName
	err = udrWriter.file.Dir_create(finalBillingBackupFolder)
	if err != nil {
		beego.Error("Write UDR file failed because billing folder create failed. error: ", err, "billing_folder:", finalBillingBackupFolder)
	}

	//初始化话单生成文件目录
	err = udrWriter.file.Dir_create(finalBillingFolder)
	if err != nil {
		beego.Error("Write UDR file failed because final billing folder create failed. error: ", err, "finalBillingFolder:", finalBillingFolder)
	}

	for _, resourceType := range resTypesToBeUDR {
		newDirName := getAbbreviationofResType(resourceType)
		//初始化话单生成文件目录
		FinalBillingFolder := finalBillingFolder + "/" + newDirName

		err = udrWriter.file.Dir_create(FinalBillingFolder)
		if err != nil {
			beego.Error("Write UDR file failed because final billing folder create failed. error: ", err, "finalBillingFolder:", FinalBillingFolder)
		}

		//初始化话单备份目录vm
		FinalBillingBackupFolder := finalBillingBackupFolder + "/" + newDirName

		err = udrWriter.file.Dir_create(FinalBillingBackupFolder)
		if err != nil {
			beego.Error("Write UDR file failed because final billing folder create failed. error: ", err, "finalBillingFolder:", FinalBillingBackupFolder)
		}
	}

	return *udrWriter

}

func (this *UdrWriter) WriterCycleUDR() {
	//不同的云资源类型生成到不同文件中
	for _, resourceType := range resTypesToBeUDR {
		//查询出待输出话单文件的话单信息
		udrs := this.usageDetailRecordDAO.QueryUDRbyResType_status(resourceType, models.UDR_STATUS_DB_READY)
		if len(udrs) != 0 {
			this.WriteUDR(udrs, resourceType)
		} else {
			//在整点时才生成空话单，此处考虑到调用外部接口获取耗时过长，所以在[0-2)分钟时都可以生成话单
			if time.Now().Minute() < 2 {
				//如果没有生成话单则生成空话单
				this.gen_blank_billing_file(resourceType)
			}
		}
	}
}

func (this *UdrWriter) WriteUDR(udrs []models.UsageDetailRecord, resourceType string) error {

	leftSize := 0
	//创建话单文件时间
	var createFileTime time.Time
	var err error
	fileOpened := false
	udrsInFile := []models.UsageDetailRecord{}
	totalBillNum := 0
	fileName, midFileName := "", ""

	//统一错误处理
	defer func() {
		if err != nil {
			this.errorHandle(udrsInFile, err)

			//出错时删除生成的话单文件
			this.file.Close()
			os.Remove(midFileName)
			os.Remove(fileName)
		}
	}()

	for _, udr := range udrs {
		billNum := 0
		billString := ""

		billNum, billString = this.gen_bill(udr)
		billSize := bytes.Count([]byte(billString), nil)

		//beego.Debug("Current file object is", this.file)
		//判断文件是否写满，如果已写满，则写入尾记录内容
		if leftSize < billSize {
			if fileOpened {
				err = this.write_file_footer(createFileTime, totalBillNum)
				if err != nil {
					beego.Error("Writer billing file footer fail, file:", this.file)
					return err
				}

				this.file.Rename(fileName)
				fileOpened = false
				beego.Info("The UDRS write file suceess, filename:", fileName, "UDR:", udrsInFile)
				this.successHandle(udrsInFile, fileName, models.UDR_STATUS_FILE_READY)

				this.SendFileToBSS(udrsInFile, fileName, resourceType)
				udrsInFile = []models.UsageDetailRecord{}
			}

			//设定话单生成时间
			createFileTime = time.Now()

			//生成话单文件名
			fileName, midFileName = this.gen_bill_filename(createFileTime, udr.Resource_type_code)

			//创建话单文件
			err = this.file.Open(midFileName)
			if err != nil {
				beego.Error("Open billing file fail: ", midFileName)
				//this.errorHandle([]db.UsageDetailRecord{udr}, err)
				return err
			}
			defer this.file.Close()
			fileOpened = true
			leftSize = fileMaxSize - HEADER_SIZE - FOOTER_SIZE

			//新建文件后，重新实始化话单数量
			totalBillNum = 0

			err = this.write_file_header(createFileTime)
			if err != nil {
				beego.Error("Writer billing file header fail: ", this.file)
				//this.errorHandle([]db.UsageDetailRecord{udr}, err)
				return err
			}
		}
		//存放完成的话单
		udrsInFile = append(udrsInFile, udr)
		totalBillNum = totalBillNum + billNum

		beego.Debug("Write bill string to file, id:", udr.Id)
		err = this.file.Writestring(billString)
		if err != nil {
			beego.Debug("Write bill string to file failed, file:", this.file, "billString: ", billString)
			//this.errorHandle(udrsInFile, err)
			return err
		}
		beego.Debug("Write bill string to file success.")

		leftSize = leftSize - billSize
	}

	if fileOpened == true {
		err = this.write_file_footer(createFileTime, totalBillNum)
		if err != nil {
			beego.Error("Writer billing file footer fail, file", this.file)
			//this.errorHandle(udrsInFile, err)
			return err
		}

		this.file.Rename(fileName)
		fileOpened = false
		beego.Info("The UDRS write file suceess, filename:", fileName, "UDR:", udrsInFile)
		this.successHandle(udrsInFile, fileName, models.UDR_STATUS_FILE_READY)
		this.SendFileToBSS(udrsInFile, fileName, resourceType)
		udrsInFile = nil
	}

	beego.Debug("Billing proces finished.")

	return err
}

//如果没有话单信息，则生成空话单文件
func (this *UdrWriter) gen_blank_billing_file(resType string) (err error) {
	beego.Info("It does not have available UDR, generate blank file, resourceType:", resType)

	//gen file name
	//blank_file_res_type := "container"

	//设定生成话单文件的时间
	createFileTime := time.Now()

	filename, tmpfilename := this.gen_bill_filename(createFileTime, resType)

	//open file for writing
	err = this.file.Open(tmpfilename)
	if err != nil {
		beego.Error("Open billing file fail: ", tmpfilename)
	}

	defer this.file.Close()
	err = this.write_file_header(createFileTime)
	if err != nil {
		beego.Error("Writer billing file header fail: ", this.file)
	}

	recourd_num := 0
	err = this.write_file_footer(createFileTime, recourd_num)
	if err != nil {
		beego.Error("Writer billing file footer fail: ", this.file)
	}

	this.file.Rename(filename)

	_, err = this.CompressFileAll(filename, resType)
	if err != nil {
		beego.Error("CompressFileAll fail: ", this.file)
	}

	return
}

func (this *UdrWriter) gen_bill_filename(createFileTime time.Time, restype string) (fileName string, tmpFileName string) {
	//conf := conf.Get_conf()
	//转成UTC时间作为文件名中的时间部分
	timestampStr := common.FormatTime_ymdhms_ms(createFileTime.UTC())

	//beego.Debug("gen billing time:", timestampStr)
	fileName = fmt.Sprintf("HWS_%s_%s_%s_%s.csv", getRegion(), getAbbreviationofResType(restype), timestampStr, this.ipLast3Pos)
	fileName = getBillingFolder() + "/" + fileName
	tmpFileName = fileName + ".mid"

	beego.Debug("Generate new udr file:", fileName)

	return fileName, tmpFileName
}

func (this *UdrWriter) gen_bill(udr models.UsageDetailRecord) (billnum int, bill string) {
	beego.Debug("Begin to assemble billing information, UDR:", udr)
	billnum = 0
	lastBillingTime := udr.Billing_start_time
	thisBillingTime := udr.Billing_end_time
	beego.Debug("ResId:", udr.Resource_id, "lastBillingTime:", lastBillingTime, "thisBillingTime", thisBillingTime)
	for lastBillingTime.Before(thisBillingTime) {
		tmpBillingTime := lastBillingTime.Add(1 * time.Hour)
		//第一次出话单或停止时话单开始时间和结束时间存在不是一个小时的情况
		if thisBillingTime.Before(tmpBillingTime) {
			tmpBillingTime = thisBillingTime
		} else {
			//如果上一次计费时间存在分钟情况，那么本次计费也进行小时取整处理
			tmpBillingTime = tmpBillingTime.Truncate(time.Hour)
		}
		beego.Debug("One udr is resId:", udr.Resource_id, "begin time:", lastBillingTime, "end time:", tmpBillingTime)

		//转化成UTC时间
		lasttime := common.FormatTime_ymdhms(lastBillingTime.UTC())
		thistime := common.FormatTime_ymdhms(tmpBillingTime.UTC())

		createRecTimeStr := common.FormatTime_ymdhms(time.Now().UTC())
		var duration float64

		//容器计费在这里计算时长
		if udr.Billing_type == common.UDR_BILLING_TYPE_APP {
			duration = tmpBillingTime.Sub(lastBillingTime).Seconds()
		} else {
			duration = udr.AccumulateFactorVal
		}

		billrec := fmt.Sprintf("%s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %f | %s",
			RECORD_FLAG, createRecTimeStr, udr.Project_id, udr.Region_code, udr.Az_code, udr.CloudserviceTypeCode, udr.Resource_type_code, udr.Resource_spec_code,
			udr.Resource_id, udr.BssParams, lasttime, thistime,
			udr.AccumulateFactorName, duration, udr.ExtendParams)

		billrec = billrec + "\r\n"

		bill = bill + billrec
		billnum++

		lastBillingTime = tmpBillingTime

		beego.Debug("Generate billing information:", billrec)
	}

	beego.Debug("Generate billing succeed. All billing:", bill, "Bill Num:", billnum)
	return billnum, bill
}

func (this *UdrWriter) write_file_header(createFileTime time.Time) (err error) {
	beego.Debug("Write file header...")
	//将话单头记录中的时间格式化为UTC时间
	datatimeStr := common.FormatTime_ymdhms(createFileTime.UTC())
	content := fmt.Sprintf("%s | %s | %s | %s \r\n", HEADER_FLAG, datatimeStr, FILE_VERSION, FILE_TYPE)
	err = this.file.Writestring(content)
	if err != nil {
		beego.Warn("Writer billing file header fail: ", this.file, content)
		return
	}

	beego.Debug("Write file header succeed.")
	return
}

func (this *UdrWriter) write_file_footer(createFileTime time.Time, record_num int) (err error) {
	beego.Debug("Write file footer...")
	//将其转化为UTC时间再格式化
	datetimeStr := common.FormatTime_ymdhms(createFileTime.UTC())
	content := fmt.Sprintf("%s | %s | %d \r\n", FOOTER_FLAG, datetimeStr, record_num)
	err = this.file.Writestring(content)
	if err != nil {
		beego.Warn("Writer billing file footer fail: ", this.file, content)
		return
	}

	beego.Debug("Write file footer succeed.")
	return
}

func (this *UdrWriter) SendFileToBSS(udrs []models.UsageDetailRecord, fullFileName, resourceType string) error {
	beego.Info("Begin to compress and sent udr to BSS, fileName", fullFileName)
	var err error
	var finalZipFileName string

	//统一错误处理
	defer func() {
		if err != nil {
			this.errorHandle(udrs, err)
		}
	}()

	if fullFileName == "" {
		beego.Error("CSV file is blank, so udrs is invalid data, udrs:", udrs)
		err = common.EInvalidBillingData
		return err
	}

	finalZipFileName, err = this.CompressFileAll(fullFileName, resourceType)
	if err != nil {
		return err
	}

	//更新话单处理状态
	this.successHandle(udrs, finalZipFileName, models.UDR_STATUS_BSS_READY)

	return err
}

func (this *UdrWriter) CompressFileAll(fullFileName, resourceType string) (string, error) {
	finalZipFileName, err := this.CompressFile(fullFileName, resourceType, false)
	if err != nil {
		return "", err
	}

	beego.Info("Compress and sent udr to BSS successfully, zip file name", finalZipFileName)

	finalZipBackFileName, err := this.CompressFile(fullFileName, resourceType, true)
	if err != nil {
		os.Remove(finalZipFileName)
		return "", err
	}
	beego.Info("Compress and sent udr to BAK FILE successfully, zip file name", finalZipBackFileName)

	//成功后删除临时目录下生成CSV文件
	os.Remove(fullFileName)

	return finalZipFileName, nil
}

//压缩文件，并传输到指定的NFS上
func (this *UdrWriter) CompressFile(fullFileName, resourceType string, isBakFile bool) (string, error) {
	var err error
	var udrZipFileName string

	//将文件压缩到文件所在目录下
	udrZipFileName, err = common.ZipFile(fullFileName, "")

	//不管成功或失败都删除该中间文件
	defer os.Remove(udrZipFileName)

	if err != nil {
		beego.Error("Compress file failed, fullFileName:", fullFileName)
		return "", err
	}

	_, fileNameNoSuffix, _ := common.SplitFullFileName(fullFileName)

	//设置最终ZIP文件存放目录
	destPath := getBillingFolder() + "_zip"

	var finalZipFileName string

	if isBakFile == true {
		finalZipFileName = destPath + backupFileName + "/" + getAbbreviationofResType(resourceType) +
			"/" + fileNameNoSuffix + ".zip"
	} else {
		finalZipFileName = destPath + "/" + getAbbreviationofResType(resourceType) +
			"/" + fileNameNoSuffix + ".zip"
	}

	fw, err := os.Create(finalZipFileName)

	//如果出错，则删除zip文件,以免出现空文件
	defer func() {
		fw.Close()
		if err != nil {
			os.Remove(finalZipFileName)
		}
	}()

	if err != nil {
		beego.Error("Create final zip file failed, error:", err)
		return "", err
	}

	w := zip.NewWriter(fw)
	defer w.Close()

	//写UDR的ZIP文件
	srcFile, err := os.Open(udrZipFileName)
	defer srcFile.Close()
	if err != nil {
		beego.Error("Open UDR zip file failed, error:", err)
		return "", err
	}

	header := &zip.FileHeader{Name: fileNameNoSuffix + ".zip"}

	//由于SetModTime中会转化成UTC时间，导致文件修改时间显示不对，所以直接将当前时间的字面时间作为UTC时间
	header.SetModTime(common.TimeLiteralToUTC(time.Now()))
	header.Method = zip.Deflate

	destFile, err := w.CreateHeader(header)
	if err != nil {
		beego.Error("Create file in zip failed, error:", err)
		return "", err
	}

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		beego.Error("Write zip file failed, error:", err)
		return "", err
	}

	header = &zip.FileHeader{Name: fileNameNoSuffix + ".zip.md5"}
	header.Method = zip.Deflate

	//由于SetModTime中会转化成UTC时间，导致文件修改时间显示不对，所以直接将当前时间的字面时间作为UTC时间
	header.SetModTime(common.TimeLiteralToUTC(time.Now()))
	//写MD5文件
	destFile, err = w.CreateHeader(header)

	if err != nil {
		beego.Error("Create md5 file in zip failed, error:", err)
		return "", err
	}
	hashSum := common.Md5sumFile(udrZipFileName)

	_, err = destFile.Write([]byte(hashSum))
	if err != nil {
		beego.Error("Write md5 infor failed, error:", err)
		return "", err
	}

	err = w.Flush()

	if err != nil {
		beego.Error("Flush final zip file failed, error:", err)
		return "", err
	}

	beego.Debug("Compress file success, fullFileName:", fullFileName, "finalZipFileName:", finalZipFileName)

	//提前将引用udrZipFileName的地方关闭
	srcFile.Close()

	//修改生成文件的权限，让BSS SFTP用户可以访问
	err = os.Chmod(finalZipFileName, UDR_FILE_PERMISSION)
	if err != nil {
		beego.Error("Change file permission failed, error:", err)
		return "", err
	}

	return finalZipFileName, nil
}

//操作成功后的处理
func (this *UdrWriter) successHandle(udrs []models.UsageDetailRecord, fullFileName string, status int) {
	beego.Debug("Update UDR status to success status, udrs:", udrs)
	for _, udr := range udrs {
		udr.Process_ip = common.GetLocalIP()
		udr.Result_code = common.RESULT_CODE_SUCCESS
		udr.Result_message = "Success"
		udr.Update_time = time.Now()
		udr.Status = status
		udr.File_name = fullFileName
		//避免对象传递过程中出错，所以只更新指定字段
		this.usageDetailRecordDAO.UpdateUdrAfterSuccess(udr.Id, udr.Result_code, udr.Result_message, udr.Process_ip, udr.Update_time, udr.Status, udr.File_name)
	}
}

//操作失败后统一处理
func (this *UdrWriter) errorHandle(udrs []models.UsageDetailRecord, err error) {
	beego.Debug("Operation failed, Update UDR status to failed status, udrs:", udrs)
	for _, udr := range udrs {
		udr.Process_ip = common.GetLocalIP()
		udr.Result_code = common.RESULT_CODE_FAIL
		udr.Result_message = err.Error()
		udr.Update_time = time.Now()

		//避免对象传递过程中出错，所以只更新指定字段
		this.usageDetailRecordDAO.UpdateUdrAfterError(udr.Id, udr.Result_code, udr.Result_message, udr.Process_ip, udr.Update_time)
	}

	//调用上报告警
}

func (this *UdrWriter) RewriteForNoFileUDR(udrs []models.UsageDetailRecord, fileName string) error {
	//创建话单文件时间
	beego.Debug("Begin to rewrite UDR file. fileName:", fileName, "udrs:", udrs)

	var createFileTime time.Time
	var err error
	totalBillNum := 0

	midFileName := fileName + ".mid"

	//统一错误处理
	defer func() {
		if err != nil {
			this.errorHandle(udrs, err)

			//出错时删除生成的话单文件
			this.file.Close()
			os.Remove(midFileName)
			os.Remove(fileName)
		}
	}()

	//设定话单生成时间
	createFileTime = time.Now()

	//创建话单文件
	err = this.file.Open(midFileName)
	if err != nil {
		beego.Error("Open billing file fail: ", midFileName)
		return err
	}
	defer this.file.Close()
	err = this.write_file_header(createFileTime)
	if err != nil {
		beego.Error("Writer billing file header fail: ", this.file)
		return err
	}

	for _, udr := range udrs {
		billNum, billString := this.gen_bill(udr)

		totalBillNum = totalBillNum + billNum

		beego.Debug("Write bill string to file, id:", udr.Id)
		err = this.file.Writestring(billString)
		if err != nil {
			beego.Debug("Write bill string to file failed, file:", this.file, "billString: ", billString)
			return err
		}
	}
	beego.Debug("Write bill string to file success.")

	err = this.write_file_footer(createFileTime, totalBillNum)

	if err != nil {
		beego.Error("Writer billing file footer fail, file", this.file)
		return err
	}

	this.file.Rename(fileName)
	//this.successHandle(udrsInFile, fileName, db.UDR_STATUS_FILE_READY)

	beego.Debug("Rewrite UDR file finished.")

	return err
}

//获取话单文件最大容量
func getFileMaxSize() int {
	var fileMaxSizeStr = beego.AppConfig.String("billing_file_maxsize")

	fileMaxSize, err := strconv.Atoi(fileMaxSizeStr)
	if err != nil {
		beego.Error("billing_file_maxsize is wrong, please check app.conf, system will use default value(10M).")
		fileMaxSize = 10
	}

	return fileMaxSize * 1024 * 1024
}

//获取CSV文件存放目录
func getBillingFolder() string {
	return beego.AppConfig.String("billing_folder")
}

//获取备份删除定时任务间隔时间
func getBackupDeleteTaskInterval() string {
	taskInterval := beego.AppConfig.String("backupDelete_task_interval")
	if taskInterval == "" {
		beego.Error("empty delete backup file task interval time,please check app.conf system will use default value(24h) ")
		return "24h"
	}
	return taskInterval

}

//获取备份文件有效期
func getBackupFileTimer() int {
	var backupTimerStr = beego.AppConfig.String("order_backup_timer")
	//验证是否是一个正整数（合法的备份有效期）
	backupTimer, err := strconv.Atoi(backupTimerStr)
	if err != nil {
		beego.Error("order backup file timer wrong,please check app.conf, system will use default value(7 days).")
		return 7
	}
	return backupTimer
}

//获取区域编码
func getRegion() string {
	return beego.AppConfig.String("regioncode")
}

//获取资源类型的简写：全称小数点最后几位为简写
func getAbbreviationofResType(resType string) string {
	strs := strings.Split(resType, ".")
	return strs[len(strs)-1]
}

func copyFile(src, dst string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer dstFile.Close()

	return io.Copy(dstFile, srcFile)
}

//获取资源类型的简写：全称小数点最后几位为简写
func getResTypeToBeUDR() []string {
	resTypeStr := beego.AppConfig.String("billing_res_type")
	if resTypeStr != "" {
		return strings.Split(resTypeStr, "|")
	} else {
		beego.Emergency("billing_res_type is null, please check app.conf.")
		return nil
	}
}
