package billing

import (
    "bytes"
    "fmt"
    "github.com/astaxie/beego"
    . "pub/common"
    "time"
  . "pub/scheduler"
)

type Billing interface {
    Start() (err error)
    Stop() (err error)
    Run() (err error) //can be called by timer or in thread multi times
    GetTaskName() string
}

type Billing_cfg struct {
    Billfolder  string
    Maxfilesize int
    Region      string   //used in file name
    Res_type    string   //used in file name
}

func New_Billing(cfg *Billing_cfg, billing_data BillingDataGetter, file File, timingTaskMgr TimingTaskMgr) Billing {
    billing := new(billing_impl)

    billing.billing_data = billing_data
    
    billing.maxfilesize = cfg.Maxfilesize
    billing.billfolder = cfg.Billfolder
    billing.region = cfg.Region
    billing.res_type = cfg.Res_type
    
    billing.timingTaskMgr = timingTaskMgr
    billing.file = file

    return billing
}

type billing_impl struct {
    billing_data BillingDataGetter

    billfolder  string
    maxfilesize int
    region      string
    res_type    string
    
    //
    file  File
    timingTaskMgr TimingTaskMgr
}

func (this *billing_impl) Start() (err error) {
    beego.Info("Billing starting...")

    err = this.file.Dir_create(this.billfolder)
    if err != nil {
        beego.Error("Billing started failed, billing folder creat fail. error: ", err)
        return
    }

    this.timingTaskMgr.Start()    
    this.timingTaskMgr.AddCycleTask(this, Get_cur_datetime().Add(1 * time.Minute), 1 * time.Minute) //1*time.Hour)
    
    beego.Info("Billing started. now:", Get_cur_datetime())

    return
}
func (this *billing_impl) GetTaskName() string {
    return "service_billing"
}

var filever string = "01"
var filefmt string = "text"

var header_flag string = "10"
var record_flag string = "20"
var footer_flag string = "90"

func (this *billing_impl) gen_blank_billing_file() (err error) {
    datetime := Get_cur_datetime_string()

    filename := ""
    tmpfilename := ""
    
    //gen file name
    filename, tmpfilename = this.gen_bill_filename(this.billfolder, datetime, this.res_type)

    //open file for writing
    err = this.file.Open(tmpfilename)
    if err != nil {
        beego.Error("Open billing file fail: ", tmpfilename, " ,error:", err)
        //to do: writing failure handling
    }

    err = this.write_file_header(header_flag, datetime, filever, filefmt)
    if err != nil {
        beego.Error("Writer billing file header fail: ", this.file, " ,error:", err)
        //to do: writing failure handling
    }

    recourd_num := 0
    err = this.write_file_footer(footer_flag, datetime, recourd_num)
    if err != nil {
        beego.Error("Writer billing file footer fail: ", this.file, " ,error:", err)
        //to do: writing failure handling
    }
    this.file.Close()
    this.file.Rename(filename)
    
    return
}

//timer task
func (this *billing_impl) Run() (err error) {
    beego.Debug("Billing procesing begins...")

    header_len := 50
    footer_len := 50
    recsizemax := this.maxfilesize - header_len - footer_len

    timestamp := Get_cur_datetime_string()
    sieral_number := 1

    billing_file_num := 0
    file_opened := false
    filesize := 0

    err = this.billing_data.Begin()
    if err == ENotBillingTime {
        beego.Debug("Not billing time, will stop billing.")
        return
    } else if err != nil {
        beego.Error("Billing data get fail at the beining. error: ", err)
        return
    }

    filename := ""
    tmpfilename := ""
    for err != ENoMoreData {
        timestamp = Get_cur_datetime_string()

        var records []BillingRecord
        records, err = this.billing_data.Get_data()
        if err == ENoMoreData {
            beego.Debug("Not more bill data.bill generation will finish.")
            break
        }

        beego.Debug("Got bill data:", records)
        //gen bill for this order
        bill := ""
        billnum := 0
        billnum, bill, err = this.gen_bill(records, record_flag, timestamp)
        if err == ENeedNoBill {
            continue
        } else if err != nil {
            beego.Error("Gen bill fail: ", records, " ,error:", err)
            //to do: writing failure handling
        }
        sieral_number = sieral_number + billnum
        rec_len := bytes.Count([]byte(bill), nil)

        //close this file if it's closed to full, or write to it
        if filesize < rec_len {
            if file_opened == true {

                err = this.write_file_footer(footer_flag, timestamp, sieral_number-1)
                if err != nil {
                    beego.Error("Writer billing file footer fail: ", this.file, " ,error:", err)
                    //to do: writing failure handling
                }
                this.file.Close()
                this.file.Rename(filename)
                file_opened = false
            }

            //gen file name
            timestamp = Get_cur_datetime_string()
            filename, tmpfilename = this.gen_bill_filename(this.billfolder, timestamp, this.res_type)

            Assert(this.billfolder != "")

            //open file for writing
            err = this.file.Open(tmpfilename)
            if err != nil {
                beego.Error("Open billing file fail: ", tmpfilename, " ,error:", err)
                //to do: writing failure handling
            }
            file_opened = true
            billing_file_num++
            filesize = recsizemax
            sieral_number = 1

            err = this.write_file_header(header_flag, timestamp, filever, filefmt)
            if err != nil {
                beego.Error("Writer billing file header fail: ", this.file, " ,error:", err)
                //to do: writing failure handling
            }
        }

        beego.Debug("billing write bill record...")
        err = this.file.StableWritestring(bill)
        if err != nil {
            beego.Debug("billing write bill record failed", this.file, ", ", bill, " ,error:", err)
            return
        }
        beego.Debug("billing write bill record succeed")

        filesize = filesize - rec_len
    }
    this.billing_data.End()

    if file_opened == true {
        err = this.write_file_footer(footer_flag, timestamp, sieral_number-1)
        if err != nil {
            beego.Error("Writer billing file footer fail: ", this.file, " ,error:", err)
            //to do: writing failure handling
        }
        this.file.Close()       
        this.file.Rename(filename)      
        file_opened = false
    }

    if billing_file_num == 0 {
        err = this.gen_blank_billing_file()
        if err != nil {
            beego.Error("gen blank billing file fail, error: ", err)
            //to do: writing failure handling
        }
    }

    err = this.billing_data.End()
    if err != nil {
        beego.Error("Billing data get fail at the end. error: ", err)
        return
    }
    
    beego.Debug("Billing process finished.")
    return
}

func (this *billing_impl) write_file_header(header_flag string, timestamp string, filever string, filefmt string) (err error) {
    beego.Debug("billing write_file_header...")

    content := fmt.Sprintf("%s | %s | %s | %s \r\n", header_flag, timestamp, filever, filefmt)
    err = this.file.StableWritestring(content)
    if err != nil {
        beego.Warn("Writer billing file header fail: ", this.file, content, " ,error:", err)
        return
    }

    beego.Debug("billing write_file_header succeed.")
    return
}

func (this *billing_impl) write_file_footer(footer_flag string, timestamp string, record_num int) (err error) {
    beego.Debug("billing write_file_footer...")

    content := fmt.Sprintf("%s | %s | %d \r\n", footer_flag, timestamp, record_num)
    err = this.file.StableWritestring(content)
    if err != nil {
        beego.Warn("billing Writer billing file footer fail: ", this.file, content, " ,error:", err)
        return
    }

    beego.Debug("billing write_file_footer succeed.")
    return
}

func (this *billing_impl) Stop() (err error) {
    beego.Info("Billing stopping...")

    beego.Info("Billing stopped.")
    return
}

func (this *billing_impl) gen_bill(records []BillingRecord, record_flag string, 
    timestamp string) (billnum int, bill string, err error) {
    beego.Debug("gen_bill: ")

    billnum = 0
    for _, record := range records {
        for _, resrec := range record.Billinginfo {
            if resrec.AccumulateFactorVal != 0.0 {
                lasttime := Time2String(record.Lastbillingtime)
                thistime := Time2String(record.Thisbillingtime)
    
                bssparams := resrec.Resextinfo.OrderId + ":" + resrec.Resextinfo.ProductId
                billrec := fmt.Sprintf("%s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %f | %s",
                    record_flag, timestamp, record.Userid, record.Regionid, record.Azid, 
                    resrec.Resinfo.CloudserviceTypeCode, resrec.Resinfo.ResourceTypeCode, resrec.Resinfo.ResourceSpecCode,
                    resrec.Resinfo.ResourceId, bssparams, lasttime, thistime,
                    resrec.AccumulateFactorName, resrec.AccumulateFactorVal, resrec.ExtendParams)

                billrec = billrec + "\r\n"

                bill = bill + billrec
                billnum++

                beego.Debug("bill generated: billrec")
            } else {
                beego.Debug("No bill will be generated: accumulator factor val is 0.0, AccumulateFactorVal", resrec.AccumulateFactorVal)
            }
        }
    }
    

    beego.Debug("gen_bill succeed. num: ", billnum)
    return
}

func (this *billing_impl) gen_bill_filename(folder string, datetime string, restype string) (filename string, tmpfilename string) {
    filename = fmt.Sprintf("HWS_%s_%s_%s.csv", this.region, restype, datetime)
    filename = folder + "/" + filename
    tmpfilename = filename + ".mid"
    
    return
}

