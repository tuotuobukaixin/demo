package common

import (
    "bytes"
    "fmt"
    "net"
    "os"
    "os/exec"
    "runtime/debug"
    "strconv"
    "strings"
    "time"

    "github.com/astaxie/beego"
)

var local_ip string = ""

var wait_for_beego_log_finished = 5 //second
var QUIT_EXECPTION = 1

func Quit(code int) {
    beego.Info("Quit: Waiting for outputing beego info......")

    time.Sleep(time.Duration(wait_for_beego_log_finished) * time.Second)
    os.Exit(code)
}

func Execute_shell_cmd(cmdline string) (output string, err error) {
    var out bytes.Buffer
    err = nil

    curdir, _ := os.Getwd()
    beego.Debug("pwd: ", curdir, " ,excute:", cmdline)

    cmd := exec.Command("/bin/bash", "-c", cmdline)
    cmd.Stdout = &out
    cmd.Stderr = &out

    err = cmd.Run()
    if err != nil {
        beego.Debug("shell execute result: " + err.Error())
    }
    beego.Debug("output: ", out.String())

    output = out.String()
    return
}

func Get_cur_datetime_int() int64 {
    return Time2Int(time.Now())
}

func Itime2Time(sec int64) time.Time {
    var nsec int64 = 0
    return time.Unix(sec, nsec).UTC()
}

func Time2Int(tm time.Time) int64 {
    return tm.UTC().Unix()
}

func Get_cur_datetime() time.Time {
    return time.Now().UTC()
}

func Get_cur_datetime_string() string {
    now := Get_cur_datetime()
    tm := Time2String(now)
    return tm
}

func Time2String(tm time.Time) string {
    //    year,mon,day := tm.Date()
    //    hour,min,sec := tm.Clock()
    //    datetime := fmt.Sprintf("%d%02d%02d%02d%02d%02d", year,mon,day,hour,min,sec)
    datetime := tm.Format("20060102150405")

    return datetime
}

//将传入时间转化为yyyymmddhhmiss
func FormatTime_ymdhms(t time.Time) string {
    return t.Format("20060102150405")
}

//将传入时间转化为yyyymmddhhmiss_SSS
func FormatTime_ymdhms_ms(t time.Time) string {
    return fmt.Sprintf("%s_%03s", t.Format("20060102150405"), GetMillisInSec(t))
}

//获取在当前秒中的毫秒数
func GetMillisInSec(t time.Time) string {
    secs := t.Unix()
    nanos := t.UnixNano()
    millis := nanos / 1000000
    return fmt.Sprintf("%d", millis-secs*1000)
}

func Get_curpath() string {
    curdir, _ := os.Getwd()

    return curdir
}

func Assert(cond bool) {
    if !cond {
        err := EAssert
        time.Sleep(time.Duration(wait_for_beego_log_finished) * time.Second)
        panic(err)
    }
}

func Assert_int(act, expect int) {
    if expect == act {
    } else {
        beego.Error("expecting: ", expect, " , bug got: ", act)
        debug.PrintStack()
        Assert(expect == act)
    }
}

func Assert_string(act, expect string) {
    if expect == act {
    } else {
        beego.Error("expecting: ", expect, " , bug got: ", act)
        debug.PrintStack()
        Assert(expect == act)
    }
}

func Assert_time(act, expect time.Time) {
    y1, m1, d1 := act.Date()
    h1, mm1, s1 := act.Clock()
    y2, m2, d2 := expect.Date()
    h2, mm2, s2 := expect.Clock()

    if y1 != y2 || m1 != m2 || d1 != d2 ||
        h1 != h2 || mm1 != mm2 || s1 != s2 {
    } else {
        beego.Error("expecting: ", expect, " , bug got: ", act)
        debug.PrintStack()
    }
}

func Is_valid_ip(ip string) (valid bool) {
    ipaddr := net.ParseIP(ip)
    if ipaddr == nil {
        beego.Warn("invalid ip:", ip)
        valid = false
        return
    } else {
        beego.Debug("valid ip:", ip, " ,ipaddr:", ipaddr)
    }

    valid = true
    return
}

func Is_valid_endpoint(ep string) (valid bool) {
    subs1 := strings.Split(ep, "//")
    if len(subs1) != 2 {
        beego.Warn("invalid ep:", ep)
        valid = false
        return
    }

    if subs1[0] != "http:" && subs1[0] != "https:" {
        beego.Warn("invalid ep:", ep, " ,part:", subs1[0]+"//")
        valid = false
        return
    }

    subs2 := strings.Split(subs1[1], ":")
    if len(subs2) == 1 {
        valid = Is_valid_ip(subs2[0])
        if valid != true {
            beego.Warn("invalid ep:", ep, " ,part:", subs2[0])
            return
        }
    } else if len(subs2) == 2 {
        valid = Is_valid_ip(subs2[0])
        if valid != true {
            beego.Warn("invalid ep:", ep, " ,part:", subs2[0])
            return
        }

        port, err := strconv.Atoi(subs2[1])
        if err != nil {
            beego.Warn("invalid ep:", ep, " ,part:", subs2[1], ", error:", err)
            valid = false
            return
        }

        if port <= 0 || port >= 65536 {
            beego.Warn("invalid ep:", ep, ", port:", port)
            valid = false
            return
        }
    }

    valid = true
    return
}

func PrintSensitiveInfo(info string) string {
    str := ""

    if info == "" {
        str = "[]"
    } else {
        str = "[******]"
    }

    return str
}


func getLocalIP() string {
    addr, err := net.ResolveUDPAddr("udp", "1.2.3.4:1")
    if err != nil {
        beego.Error("net.ResolveUDPAddr failed, error:", err)
        return ""
    }

    conn, err := net.DialUDP("udp", nil, addr)
    if err != nil {
        beego.Error("net.DialUDP failed, error:", err, "address:", addr)
        return ""
    }

    defer conn.Close()

    ip, _, err := net.SplitHostPort(conn.LocalAddr().String())
    if err != nil {
        beego.Error("net.SplitHostPort failed, error:", err, "address:", conn.LocalAddr().String())
        return ""
    }

    return ip
}

func GetLocalIP() string {
    if local_ip == "" {
        beego.Debug("Get local IP by command again.")
        local_ip = getLocalIP()
    }
    return local_ip

}
