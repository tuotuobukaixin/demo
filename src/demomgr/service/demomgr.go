package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"demomgr/models"
	"demomgr/util"
	"strings"
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


func GetDemoTest(w http.ResponseWriter, r *http.Request) {

	v := r.Header.Get("Content-Type")
	v = strings.ToLower(v)
	v = strings.Replace(v, " ", "", -1)
	if v != "application/json" && v != "application/json;charset=utf-8" {
		errStr := "get demotest, content-Type error"
		util.LOGGER.Error(errStr, nil)
		ErrorResponse(400, errStr, w)
		return
	}
	var rspdemotest []models.DemoTestGet
	name := r.URL.Query().Get("name")
	if name == "" {

		demotests, err := models.GetDemoTests()
		if err != nil {
			errStr := "get demotest failed"
			util.LOGGER.Error(errStr, err)
			ErrorResponse(500, errStr, w)
			return
		}
		for _, demotest := range demotests {

			var temp models.DemoTestGet
			temp.Name = demotest.Name
			temp.Status = demotest.Status
			temp.ServiceAddr = demotest.ServiceAddr
			temp.TcpTest = demotest.TcpTest
			temp.FileTest = demotest.FileTest
			temp.FileSize = demotest.FileSize
			temp.DownFile = demotest.DownFile
			temp.DownFileSum = demotest.DownFileSum
			temp.DownFileUrl = demotest.DownFileUrl
			rspdemotest = append(rspdemotest, temp)
		}
	} else {
		demotest, err := models.GetDemoTest(name)
		if err != nil {
			errStr := "get demotest failed"
			util.LOGGER.Error(errStr, err)
			ErrorResponse(500, errStr, w)
			return
		}
		var temp models.DemoTestGet
		temp.Name = demotest.Name
		temp.Status = demotest.Status
		temp.ServiceAddr = demotest.ServiceAddr
		temp.TcpTest = demotest.TcpTest
		temp.FileTest = demotest.FileTest
		temp.FileSize = demotest.FileSize
		temp.DownFile = demotest.DownFile
		temp.DownFileSum = demotest.DownFileSum
		temp.DownFileUrl = demotest.DownFileUrl
		rspdemotest = append(rspdemotest, temp)
	}

	data, _ := json.Marshal(&rspdemotest)

	HttpResponse(200, data, w)
	util.LOGGER.Info("get demotest success full")
}

// OpsGetQuota get user quota
// @Title get a quota info
// @Router /cce/quota/ops [get]
func UpdateDemoTest(w http.ResponseWriter, r *http.Request) {

	v := r.Header.Get("Content-Type")
	v = strings.ToLower(v)
	v = strings.Replace(v, " ", "", -1)
	if v != "application/json" && v != "application/json;charset=utf-8" {
		errStr := "OpsUpdateQuota, content-Type error"
		util.LOGGER.Error(errStr, nil)
		ErrorResponse(400, errStr, w)
		return
	}

	demotest := models.DemoTestGet{}
	//get request
	_, err := TransferReqToInterface(r, &demotest)
	if err != nil {
		errStr := "transfer request to interface err"
		util.LOGGER.Error(errStr, err)
		ErrorResponse(500, "transfer request to interface err", w)
		return
	}
	if demotest.Name != "" {
		game, err := models.GetDemoTest(demotest.Name)
		if err != nil {
			errStr := "get demotest failed"
			util.LOGGER.Error(errStr, err)
			ErrorResponse(500, errStr, w)
			return
		}
		game.Status = demotest.Status
		game.ServiceAddr = demotest.ServiceAddr
		game.TcpTest = demotest.TcpTest
		game.FileSize = demotest.FileSize
		game.FileTest = demotest.FileTest
		game.DownFile = demotest.DownFile
		game.DownFileSum = demotest.DownFileSum
		game.DownFileUrl = demotest.DownFileUrl
		err = models.UpdateDemoTest(game)
		if err != nil {
			errStr := "Update demotest failed"
			util.LOGGER.Error(errStr, err)
			ErrorResponse(500, errStr, w)
			return
		}
		HttpResponse(200, []byte(""), w)
		util.LOGGER.Info("update demotest successfully")
	} else {
		errStr := "name cannot be change"
		util.LOGGER.Error(errStr, err)
		ErrorResponse(400, errStr, w)
		return
	}

}

// OpsGetQuota get user quota
// @Title get a quota info
// @Router /cce/quota/ops [get]
func AddDemoTest(w http.ResponseWriter, r *http.Request) {

	v := r.Header.Get("Content-Type")
	v = strings.ToLower(v)
	v = strings.Replace(v, " ", "", -1)
	if v != "application/json" && v != "application/json;charset=utf-8" {
		errStr := "OpsUpdateQuota, content-Type error"
		util.LOGGER.Error(errStr, nil)
		ErrorResponse(400, errStr, w)
		return
	}

	demotest := models.DemoTestGet{}
	//get request
	_, err := TransferReqToInterface(r, &demotest)
	if err != nil {
		errStr := "transfer request to interface err"
		util.LOGGER.Error(errStr, err)
		ErrorResponse(500, "transfer request to interface err", w)
		return
	}
	game := models.DemoTest{}
	game.ServiceAddr = demotest.ServiceAddr
	game.Status = demotest.Status
	game.Name = demotest.Name
	game.Podip = demotest.Podip
	eng, err := models.GetDemoTest(game.Name)
	if eng != nil {
		game.ID = eng.ID
		game.FileSize = eng.FileSize
		game.FileTest = eng.FileTest
		game.DownFile = eng.DownFile
		game.TcpTest = eng.TcpTest
		game.DownFileSum = eng.DownFileSum
		game.DownFileUrl = eng.DownFileUrl
		err = models.UpdateDemoTest(&game)
		if err != nil {
			errStr := "Add quota failed"
			util.LOGGER.Error(errStr, err)
			ErrorResponse(500, errStr, w)
			return
		}
	} else {
		err = models.AddDemoTest(&game)
		if err != nil {
			errStr := "Add quota failed"
			util.LOGGER.Error(errStr, err)
			ErrorResponse(500, errStr, w)
			return
		}
	}
	HttpResponse(200, []byte(""), w)
	util.LOGGER.Info("Add demotest successfully")

}
