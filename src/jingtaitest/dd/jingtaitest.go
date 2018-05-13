package dd

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"
	"errors"
	"bytes"
	"encoding/json"
	"crypto/tls"
)

var REGION = "cn-north-1"
var CLUSTER_UUID="b8c96ab0-5461-11e8-941c-0255ac101f3e"
var project_id="dc98b5ef5f434e99b4ade7a870d3a503"
var resov_node="35f1eafd-54cd-11e8-9481-0255ac101f3e"



var URLINFO = make(map[int]string)
var connectionPool = struct {
	sync.RWMutex
	pool map[string]*http.Client
}{pool: make(map[string]*http.Client)}

//NewConnection :create a new connection
func NewConnection(requestPath string) (httpClient *http.Client, err error) {
	requestURL, err := url.Parse(requestPath)
	if err != nil {
		return
	}
	requestHost := requestURL.Host
	//Create a new security connection and update the cache
	httpClient, err = createConnection()
	if err == nil {
		connectionPool.Lock()
		connectionPool.pool[requestHost] = httpClient
		connectionPool.Unlock()
	}
	return
}

func createConnection() (httpClient *http.Client, err error) {

	timeout := time.Duration(10 * time.Second)
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true},}
	httpClient = &http.Client{Timeout: timeout,Transport: tr}
	return
}

//GetConnection :get a existing connection
func GetConnection(requestPath string) (httpClient *http.Client, ok bool) {
	requestURL, err := url.Parse(requestPath)
	if err != nil {
		return
	}
	requestHost := requestURL.Host
	connectionPool.RLock()
	existingConn, ok := connectionPool.pool[requestHost]
	connectionPool.RUnlock()
	if ok { //Get the existing connection in cache
		return existingConn, true
	}
	ok = false
	return
}

// DoHttpRequest send request to ops server
func DoHttpRequest(method string, requrl string, contentType string, body io.Reader, token string, clusterid string) (data []byte, statusCode int, header http.Header, err error) {

	req, err := http.NewRequest(method, requrl, body)
	if err != nil {
		return nil, 500, nil, err
	}

	req.Header.Set("Content-Type", contentType)
	if token != "" {
		req.Header.Set("X-Auth-Token", token)
	}
	if clusterid != "" {
		req.Header.Set("X-Cluster-ID", clusterid)
	}
	requestURL, err := url.Parse(requrl)
	if err != nil {
		return
	}
	requestHost := requestURL.Host

	var httpClient *http.Client
	c, ok := GetConnection(requrl)
	if ok { // The connection existing in cache
		httpClient = c
	} else { //Have to create a new connection
		httpClient, err = NewConnection(requestURL.Scheme + "://" + requestHost)
		if err != nil {
			return nil, 500, nil, err
		}
	}

	resp, err := httpClient.Do(req)

	if err != nil {
		httpClient, err = NewConnection(requestURL.Scheme + "://" + requestHost)
		if err != nil { //Try to refresh the cache and try again in case the error caused by the cache incorrect
			return nil, 500, nil, err
		}
		resp, err = httpClient.Do(req)
		if err != nil { //Try to refresh the cache and try again in case the error caused by the cache incorrect
			return nil, 500, nil, err
		}
	}

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 500, nil, err
	}

	defer resp.Body.Close()
	return data, resp.StatusCode, resp.Header, nil
}

func GetToken(usename string, domain string,password string, project_name string) (error, string) {
	body := `{	"auth": {	"identity": {	"methods": ["password"],	"password": {	"user": {	"name": "` + usename +`",	"password": "` + password +`",	"domain": {	"name": "` + domain +`"	}	}	}	},	"scope": {	"project": {	"name": "` + project_name +`"	}	}	}}`
	jsonBlob := []byte(body)

	data := bytes.NewReader([]byte(jsonBlob))
	router := "https://iam." + REGION + ".myhuaweicloud.com/v3/auth/tokens"
	rsp, status_code, reshead, _ := DoHttpRequest("POST", router, "application/json;charset=utf8", data, "", "")
	if status_code != 201 {
		fmt.Println(string(rsp))
		return errors.New("register failed"), ""
	}
	return nil, reshead["X-Subject-Token"][0]
}
func GetNode(clusterid string, token string,project_id string) []byte{
	router := "https://cce." + REGION + ".myhuaweicloud.com/api/v3/projects/" + project_id + "/clusters/" + clusterid + "/nodes"
	rsp, status_code, _, _ := DoHttpRequest("GET", router, "application/json", nil, token, "")
	if status_code != 200 {
		fmt.Println(string(rsp))
		return nil
	}
	fmt.Println(string(rsp))
	return rsp

}
func Getpod(clusterid string, token string,hostname string) bool{
	router := "https://"+clusterid+".cce." + REGION + ".myhuaweicloud.com/api/v1/namespaces/default/pods?fieldSelector=spec.nodeName=" + hostname
	rsp, status_code, _, _ := DoHttpRequest("GET", router, "application/json", nil, token, clusterid)
	if status_code != 200 {
		fmt.Println(string(rsp))
		return false
	}
	var dat map[string]interface{}
	json.Unmarshal(rsp, &dat)
	for _,intem:=range dat["items"].([]interface{}) {
		_,flag1:= intem.(map[string]interface{})["metadata"]
		if flag1 {
			_,flag2:= intem.(map[string]interface{})["metadata"].(map[string]interface{})["labels"]
			if flag2 {
				_,flag3:= intem.(map[string]interface{})["metadata"].(map[string]interface{})["labels"].(map[string]interface{})["app"]
				if flag3 {
					if intem.(map[string]interface{})["metadata"].(map[string]interface{})["labels"].(map[string]interface{})["app"].(string) != "icagent" {
						return false
					}
				}
			}
		}
	}
	fmt.Println(hostname + "has no pod now")
	return true

}


func DeleteNode(clusterid string, token string,project_id string,node_id string) error{
	if node_id == resov_node {
		return nil
	}
	fmt.Println("begion to delete node " + node_id )
	router := "https://cce." + REGION + ".myhuaweicloud.com/api/v3/projects/"+project_id+"/clusters/"+clusterid+"/nodes/" + node_id
	rsp, status_code, _, _ := DoHttpRequest("DELETE", router, "application/json", nil, token, "")
	if status_code != 200 {
		fmt.Println(string(rsp))
		return errors.New("delete failed")
	}
	for a := 0; a < 10; a++ {
		time.Sleep(10*time.Second)

		_, status_code, _, _ := DoHttpRequest("GET", router, "application/json", nil, token, "")
		if status_code == 404 {
			break
		}
	}
	return nil

}


func main() {
	for {
		time.Sleep(1* time.Hour)
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		fmt.Println("start to delete node" )
		err, token := GetToken("shenhongjun-test", "shenhongjun-test", "qwer1234!", "cn-north-1")
		if err != nil {
			fmt.Println(err)
			return
		}
		var dat map[string]interface{}
		rsp := GetNode("b8c96ab0-5461-11e8-941c-0255ac101f3e", token, "dc98b5ef5f434e99b4ade7a870d3a503")

		json.Unmarshal(rsp, &dat)
		var hostname []string
		for _, intem := range dat["items"].([]interface{}) {
			_, flag1 := intem.(map[string]interface{})["status"]
			if flag1 {
				_, flag2 := intem.(map[string]interface{})["status"].(map[string]interface{})["privateIP"]
				if flag2 {
					hostname = append(hostname, intem.(map[string]interface{})["status"].(map[string]interface{})["privateIP"].(string))
					flag := Getpod("b8c96ab0-5461-11e8-941c-0255ac101f3e", token, intem.(map[string]interface{})["status"].(map[string]interface{})["privateIP"].(string))
					if flag {
						DeleteNode("b8c96ab0-5461-11e8-941c-0255ac101f3e", token, "dc98b5ef5f434e99b4ade7a870d3a503", intem.(map[string]interface{})["metadata"].(map[string]interface{})["uid"].(string))
					}
				}
			}
		}

	}


}
