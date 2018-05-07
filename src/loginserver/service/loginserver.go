package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"loginserver/models"
	"loginserver/util"
	kubrt "registerserver/models"
	"strings"
	"bytes"
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



// OpsGetQuota get user quota
// @Title get a quota info
// @Router /cce/quota/ops [get]
func GetGameserver(w http.ResponseWriter, r *http.Request) {

	v := r.Header.Get("Content-Type")
	v = strings.ToLower(v)
	v = strings.Replace(v, " ", "", -1)
	if v != "application/json" && v != "application/json;charset=utf-8" {
		errStr := "get gameserver, content-Type error"
		util.LOGGER.Error(errStr, nil)
		ErrorResponse(400, errStr, w)
		return
	}
	var rspgameserver []kubrt.GameServerGet
	name := r.URL.Query().Get("name")
	if name == "" {

		gameservers, err := kubrt.GetGameServers()
		if err != nil {
			errStr := "get gameserver failed"
			util.LOGGER.Error(errStr, err)
			ErrorResponse(500, errStr, w)
			return
		}
		for _, gameserver := range gameservers {

			var temp kubrt.GameServerGet
			temp.Name = gameserver.Name
			temp.Status = gameserver.Status
			temp.ServiceAddr = gameserver.ServiceAddr
			rspgameserver = append(rspgameserver, temp)
		}
	} else {
		gameserver, err := kubrt.GetGameServer(name)
		if err != nil {
			errStr := "get gameserver failed"
			util.LOGGER.Error(errStr, err)
			ErrorResponse(500, errStr, w)
			return
		}
		var temp kubrt.GameServerGet
		temp.Name = gameserver.Name
		temp.Status = gameserver.Status
		temp.ServiceAddr = gameserver.ServiceAddr
		rspgameserver = append(rspgameserver, temp)
	}

	data, _ := json.Marshal(&rspgameserver)

	HttpResponse(200, data, w)
	util.LOGGER.Info("get gameserver success full")
}

