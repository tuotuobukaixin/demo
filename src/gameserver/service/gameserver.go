package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"gameserver/models"
	"gameserver/util"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	rt "registerserver/models"
	"strings"
	"time"
	"strconv"
)

//ErrorRsp used to struct the error response
type ErrorRsp struct {
	Code   int    `json:"code"`
	Desc   string `json:"description"`
	Detail string `json:"error_code"`
}

//ErrorResponse get the error response
func ErrorResponse(statusCode int, detail string, w http.ResponseWriter) {
	errRsp := ErrorRsp{Detail: detail}
	data, _ := json.Marshal(&errRsp)
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write(data)
}
func HttpResponse(statusCode int, data []byte, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write(data)
}

//TransferReqToInterface :Transfer Request to interface
func TransferReqToInterface(r *http.Request, reqMsg interface{}) ([]byte, error) {
	msg, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(msg, &reqMsg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
func RegisterGameserver() error {
	var temp = rt.GameServerGet{}
	temp.Name = util.Config.ServerName
	temp.ServiceAddr = util.Config.ServerAddr
	temp.Status = "Running"
	reqInfo, err := json.Marshal(temp)
	if err != nil {
		return err
	}
	data := bytes.NewReader([]byte(reqInfo))
	router := util.Config.Registerurl + "/api/v1/gameserver"
	_, status_code, _, _ := util.DoHttpRequest("POST", router, "application/json", data, "", "")
	if status_code != 200 {
		return errors.New("register failed")
	}
	return nil
}
func Writefile(size int) float64 {
	util.LOGGER.Info("begin write file")
	sb := bytes.Buffer{}
	for j := 0; j < 8000; j++ {
		sb.WriteString("test")
	}
	str := sb.String()
	start := time.Now()
	_ = os.Remove("file/" + util.Config.ServerName)
	file, err := os.OpenFile("file/"+util.Config.ServerName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		util.LOGGER.Error("open file failed.", err)
		return float64(0)
	}
	defer file.Close()

	for i := 0; i < size*32; i++ {
		file.WriteString(str)
	}
	cost := time.Since(start).Seconds()
	util.LOGGER.Info("finish write file")
	return float64(size) / cost

}

func Readfile() float64 {
	util.LOGGER.Info("begin read file")
	start := time.Now()
	file, err := os.OpenFile("file/"+util.Config.ServerName, os.O_RDWR, 0666)
	if err != nil {
		util.LOGGER.Error("open file failed.", err)
		return float64(0)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return float64(0)
	}

	var size = stat.Size()
	// define read block size = 2
	buf := make([]byte, 64000)
	for {
		_, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				util.LOGGER.Error("Read file failed.", err)
				return float64(0)
			}
		}
	}
	cost := time.Since(start).Seconds()
	util.LOGGER.Info("finish read file")
	return float64(size/1024/1024) / cost

}
func TestTCP(num int) (int, int, string) {
	gss, err := rt.GetGameServers()
	if err != nil {
		return 0, 0,""
	}
	total := len(gss)
	success := 0
	detail:=""
	for i := 0; i < len(gss); i++ {
		if gss[i].ServiceAddr != "" {

			router := fmt.Sprintf("http://%s.default.svc.cluster.local:8088/api/v1/gameserverhealth", gss[i].Name)
			rsp, status_code, _, _ := util.DoHttpRequest("GET", router, "application/json", nil, "", "")
			if status_code == 200 {
				success++
			}else {
				detail = strconv.Itoa(status_code) + string(rsp) +";" + detail
			}

		}

	}
	return success, total, detail
}
func Gothread() {
	for {
		time.Sleep(10 * time.Second)
		var tmp models.GameServerTestResult
		gs, _ := rt.GetGameServer(util.Config.ServerName)
		if gs.FileTest {
			writespeed := Writefile(gs.FileSize)
			tmp.FileWriteSpeed = int(writespeed)
			readspeed := Readfile()
			tmp.FileReadSpeed = int(readspeed)
		}
		if gs.TcpTest {
			success, totol,detail := TestTCP(gs.TcpNum)
			tmp.Total = totol
			tmp.Success = success
			tmp.Detail = detail
		}
		if gs.FileTest || gs.TcpTest {
			tmp.Name = gs.Name
			tmp.Time = time.Now().String()
			models.AddGameServerTestResult(&tmp)
		}
		results, err := models.GetGameServerResult(util.Config.ServerName)

		if err == nil && len(results) >= 500 {
			util.LOGGER.Info("begin write result")
			WriteFile(results)
			models.DeleteGameServer(util.Config.ServerName)
		}
	}
}
func WriteFile(results []models.GameServerTestResult) {
	file, err := os.OpenFile("file/result-"+util.Config.ServerName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		util.LOGGER.Error("open file failed.", err)
		return
	}
	defer file.Close()

	for i := 0; i < len(results); i++ {
		str := fmt.Sprintf("%s,%d,%d,%d/%d,%s,%s\n", results[i].Name, results[i].FileReadSpeed,
			results[i].FileWriteSpeed, results[i].Success, results[i].Total, results[i].Time, results[i].Detail)
		file.WriteString(str)
	}
}

// OpsGetQuota get user quota
// @Title get a quota info
// @Router /cce/quota/ops [get]
func GetGameserverDetail(w http.ResponseWriter, r *http.Request) {

	v := r.Header.Get("Content-Type")
	v = strings.ToLower(v)
	v = strings.Replace(v, " ", "", -1)
	if v != "application/json" && v != "application/json;charset=utf-8" {
		errStr := "get gameserver, content-Type error"
		util.LOGGER.Error(errStr, nil)
		ErrorResponse(400, errStr, w)
		return
	}
	var rspgameserver []models.GameServerTestResultGet

	gameservers, err := models.GetGameServersResult()
	if err != nil {
		errStr := "get gameserver failed"
		util.LOGGER.Error(errStr, err)
		ErrorResponse(500, errStr, w)
		return
	}
	for _, gameserver := range gameservers {

		var temp models.GameServerTestResultGet
		temp.Time = gameserver.Time
		temp.FileReadSpeed = gameserver.FileReadSpeed
		temp.FileWriteSpeed = gameserver.FileWriteSpeed
		temp.Success = gameserver.Success
		temp.Total = gameserver.Total
		rspgameserver = append(rspgameserver, temp)
	}

	data, _ := json.Marshal(&rspgameserver)

	HttpResponse(200, data, w)
	util.LOGGER.Info("get gameserver success full")
}

// OpsGetQuota get user quota
// @Title get a quota info
// @Router /cce/quota/ops [get]
func Health(w http.ResponseWriter, r *http.Request) {

	v := r.Header.Get("Content-Type")
	v = strings.ToLower(v)
	v = strings.Replace(v, " ", "", -1)
	if v != "application/json" && v != "application/json;charset=utf-8" {
		errStr := "get gameserver, content-Type error"
		util.LOGGER.Error(errStr, nil)
		ErrorResponse(400, errStr, w)
		return
	}

	HttpResponse(200, []byte(""), w)
}
