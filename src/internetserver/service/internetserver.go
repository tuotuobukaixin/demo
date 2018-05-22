package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"internetserver/util"
	"strings"
	"github.com/garyburd/redigo/redis"
	"fmt"
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
func RedisFunc() string{
	c, err := redis.Dial("tcp", util.Config.Redisserver)
	if err != nil {
		util.LOGGER.Error("Connect to redis error", err)
		fmt.Println("Connect to redis error", err)
		return ""
	}
	defer c.Close()
	is_key_exit, err := redis.Bool(c.Do("EXISTS", "mykey1"))
	if err != nil {
		util.LOGGER.Error("exit", err)
		fmt.Println("exit", err)
	}
	if !is_key_exit {
		_, err = c.Do("SET", "mykey", "superWang")
		if err != nil {
			util.LOGGER.Error("redis set failed:", err)
			fmt.Println("redis set failed:", err)
		}
	}
	username, err := redis.String(c.Do("GET", "mykey"))
	if err != nil {
		util.LOGGER.Error("redis get failed:", err)
		fmt.Println("redis set failed:", err)
	}
	return username

}
type Result struct {
	Success        int    `json:"success"`
	Total          int    `json:"total"`
}
func GetInfo(w http.ResponseWriter, r *http.Request) {

	v := r.Header.Get("Content-Type")
	v = strings.ToLower(v)
	v = strings.Replace(v, " ", "", -1)
	if v != "application/json" && v != "application/json;charset=utf-8" {
		errStr := "get info, content-Type error"
		util.LOGGER.Error(errStr, nil)
		ErrorResponse(400, errStr, w)
		return
	}
	var temp = Result{}
	temp.Total = util.Config.Num
	temp.Success = 0

	for i:=0; i<util.Config.Num; i++{
		result := RedisFunc()
		if result == "superWang"{
			temp.Success ++
		}
	}
	data, _ := json.Marshal(&temp)
	HttpResponse(200, data, w)
}

