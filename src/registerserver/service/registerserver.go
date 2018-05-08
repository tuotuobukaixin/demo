package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"registerserver/models"
	"registerserver/util"
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
	var rspgameserver []models.GameServerGet
	name := r.URL.Query().Get("name")
	if name == "" {

		gameservers, err := models.GetGameServers()
		if err != nil {
			errStr := "get gameserver failed"
			util.LOGGER.Error(errStr, err)
			ErrorResponse(500, errStr, w)
			return
		}
		for _, gameserver := range gameservers {

			var temp models.GameServerGet
			temp.Name = gameserver.Name
			temp.Status = gameserver.Status
			temp.ServiceAddr = gameserver.ServiceAddr
			rspgameserver = append(rspgameserver, temp)
		}
	} else {
		gameserver, err := models.GetGameServer(name)
		if err != nil {
			errStr := "get gameserver failed"
			util.LOGGER.Error(errStr, err)
			ErrorResponse(500, errStr, w)
			return
		}
		var temp models.GameServerGet
		temp.Name = gameserver.Name
		temp.Status = gameserver.Status
		temp.ServiceAddr = gameserver.ServiceAddr
		rspgameserver = append(rspgameserver, temp)
	}

	data, _ := json.Marshal(&rspgameserver)

	HttpResponse(200, data, w)
	util.LOGGER.Info("get gameserver success full")
}

// OpsGetQuota get user quota
// @Title get a quota info
// @Router /cce/quota/ops [get]
func UpdateGameserver(w http.ResponseWriter, r *http.Request) {

	v := r.Header.Get("Content-Type")
	v = strings.ToLower(v)
	v = strings.Replace(v, " ", "", -1)
	if v != "application/json" && v != "application/json;charset=utf-8" {
		errStr := "OpsUpdateQuota, content-Type error"
		util.LOGGER.Error(errStr, nil)
		ErrorResponse(400, errStr, w)
		return
	}

	gameserver := models.GameServerGet{}
	//get request
	_, err := TransferReqToInterface(r, &gameserver)
	if err != nil {
		errStr := "transfer request to interface err"
		util.LOGGER.Error(errStr, err)
		ErrorResponse(500, "transfer request to interface err", w)
		return
	}
	if gameserver.Name != "" {
		game, err := models.GetGameServer(gameserver.Name)
		if err != nil {
			errStr := "get gameserver failed"
			util.LOGGER.Error(errStr, err)
			ErrorResponse(500, errStr, w)
			return
		}
		game.Status = gameserver.Status
		game.ServiceAddr = gameserver.ServiceAddr
		game.TcpTest = gameserver.TcpTest
		game.FileSize = gameserver.FileSize
		game.FileTest = gameserver.FileTest
		game.TcpNum = gameserver.TcpNum
		err = models.UpdateGameServer(game)
		if err != nil {
			errStr := "Update gameserver failed"
			util.LOGGER.Error(errStr, err)
			ErrorResponse(500, errStr, w)
			return
		}
		HttpResponse(200, []byte(""), w)
		util.LOGGER.Info("update gameserver successfully")
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
func AddGameserver(w http.ResponseWriter, r *http.Request) {

	v := r.Header.Get("Content-Type")
	v = strings.ToLower(v)
	v = strings.Replace(v, " ", "", -1)
	if v != "application/json" && v != "application/json;charset=utf-8" {
		errStr := "OpsUpdateQuota, content-Type error"
		util.LOGGER.Error(errStr, nil)
		ErrorResponse(400, errStr, w)
		return
	}

	gameserver := models.GameServerGet{}
	//get request
	_, err := TransferReqToInterface(r, &gameserver)
	if err != nil {
		errStr := "transfer request to interface err"
		util.LOGGER.Error(errStr, err)
		ErrorResponse(500, "transfer request to interface err", w)
		return
	}
	game := models.GameServer{}
	game.ServiceAddr = gameserver.ServiceAddr
	game.Status = gameserver.Status
	game.Name = gameserver.Name
	eng, err := models.GetGameServer(game.Name)
	if eng != nil {
		game.ID = eng.ID
		game.FileSize = eng.FileSize
		game.FileTest = eng.FileTest
		game.TcpNum = eng.TcpNum
		game.TcpTest = eng.TcpTest
		err = models.UpdateGameServer(&game)
		if err != nil {
			errStr := "Add quota failed"
			util.LOGGER.Error(errStr, err)
			ErrorResponse(500, errStr, w)
			return
		}
	} else {
		err = models.AddGameServer(&game)
		if err != nil {
			errStr := "Add quota failed"
			util.LOGGER.Error(errStr, err)
			ErrorResponse(500, errStr, w)
			return
		}
	}
	HttpResponse(200, []byte(""), w)
	util.LOGGER.Info("Add gameserver successfully")

}
