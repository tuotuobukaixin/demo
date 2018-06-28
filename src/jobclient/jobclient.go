package main

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"
	//"strconv"
	"jobclient/util"
	//"jobclient/models"
	//"encoding/json"
	"github.com/garyburd/redigo/redis"
	"strconv"
)

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
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	httpClient = &http.Client{Timeout: timeout, Transport: tr}
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
		req.Header.Set("Authorization", "Bearer "+token)
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

func createjob(urlstring string, sfs_name string, token string, job_name string) {
	job := `{	"apiVersion": "batch/v1",	"kind": "Job",	"metadata": {		"name": "` + job_name + `",
		"namespace": "default"
	},
	"spec": {
		"template": {
			"metadata": {
				"name": "` + job_name + `"
			},
			"spec": {
				"containers": [{
					"env": [
						{
							"name": "redis_url",
							"value": "192.168.0.216"
						},
						{
							"name": "redis_port",
							"value": "6379"
						},
						{
							"name": "timeout",
							"value": "600"
						},
						{
							"name": "jobname",
							"value": "` + job_name + `"
						}
					],
					"image": "swr.cn-north-1.myhuaweicloud.com/cce-demo/jobtest:latest",
					"imagePullPolicy": "IfNotPresent",
					"lifecycle": {},
					"name": "container-0",
					"resources": {
						"limits": {
							"cpu": "100m",
							"memory": "0.1Gi"
						},
						"requests": {
							"cpu": "100m",
							"memory": "0.1Gi"
						}
					},
					"terminationMessagePath": "/dev/termination-log",
					"terminationMessagePolicy": "File"
				}],
				"dnsPolicy": "ClusterFirst",
				"imagePullSecrets": [{
					"name": "default-secret"
				}],
				"restartPolicy": "OnFailure",
				"schedulerName": "default-scheduler",
				"securityContext": {},
				"terminationGracePeriodSeconds": 30
			}
		}
	}
}
`
	data := bytes.NewReader([]byte(job))
	router := urlstring + "/apis/batch/v1/namespaces/default/jobs"
	rsp, status_code, _, _ := DoHttpRequest("POST", router, "application/json;charset=utf8", data, token, "")
	if status_code != 201 {
		fmt.Println(string(rsp))
		return
	}
}

func Getjob(urlstring string, token string) []byte {
	router := urlstring + "/apis/batch/v1/namespaces/default/jobs"
	rsp, status_code, _, _ := DoHttpRequest("GET", router, "application/json", nil, token, "")
	if status_code != 200 {
		fmt.Println(string(rsp))
		return nil
	}
	return rsp

}
func DeleteJob(urlstring string, token string, job_name string) error {
	body := `{
		"kind": "DeleteOptions",
		"apiVersion": "v1",
		"propagationPolicy": "Foreground"
	}`
	jsonBlob := []byte(body)

	data := bytes.NewReader([]byte(jsonBlob))
	router := urlstring + "/apis/batch/v1/namespaces/default/jobs/" + job_name
	rsp, status_code, _, _ := DoHttpRequest("DELETE", router, "application/json", data, token, "")
	if status_code != 200 {
		fmt.Println(string(rsp))
		return errors.New("delete failed")
	}
	return nil

}

func redis_get() []string {
	var value []string
	c, err := redis.Dial("tcp", util.Config.Redisserver)
	if err != nil {
		util.LOGGER.Error("Connect to redis error", err)
		fmt.Println("Connect to redis error", err)
		return value
	}
	defer c.Close()
	len, err := redis.Int(c.Do("llen", "joblist"))
	if err != nil {
		util.LOGGER.Error("redis get failed:", err)
		fmt.Println("redis set failed:", err)
		return value
	}
	if len > 0 {

		for i := 0; i < len; i++ {
			tmp, err := redis.String(c.Do("lpop", "joblist"))
			if err != nil {
				util.LOGGER.Error("redis get failed:", err)
				fmt.Println("redis set failed:", err)
			} else {

				value = append(value, tmp)
			}
		}
	}
	return value
}
func main() {
	endpoint := "https://192.168.0.167:5443"
	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImRlZmF1bHQtdG9rZW4tNHRjcHoiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoiZGVmYXVsdCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjM4Y2M2MGQ5LTVmZTUtMTFlOC05N2Q4LWZhMTYzZWY4NzVlMCIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OmRlZmF1bHQifQ.cmtex0EHZWx6zlssyu2SHYxrgeoIMVLjBJhce0jiUFQlUdhhaGuyacJzAASrhOSH7xSIHe59naFEKYkk_9iYAgO9qa9wyxLurgsCRsiXo4RafuSsevMSfv43ggXa-TSrZ64dnVVIaVW45ITKs2HBiI7gzU4mZtGthhtwfwc5ouKdB4BxojvHxdb57IH1OU9QrkY7_PCWfNb-RpxwOqSpR5VaJqzT5tgS1YsZ-Cy8kAYMBA45GHRikgYj-hnOLq7euEtIzCCuk0Lr_QcElEpmFE6_tOFepKDNl4vnXj9D2kvsoVYnGWfnn0GOvczVFb7lkiDYvgYHXBmZ9dpFm5P1bA"

	flag := false
	jobnum := 0
	for {
		time.Sleep(1 * time.Second)
		if flag && jobnum != util.Config.Num {
			finish_job := redis_get()
			for _, job_tmp := range finish_job {
				DeleteJob(endpoint, token, job_tmp)
				jobnum++
				time.Sleep(1 * time.Second)
			}
		} else {
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
			for i := 0; i < util.Config.Num; i++ {
				createjob(endpoint, "", token, "jobtest"+strconv.Itoa(i))
			}
			flag = true
			jobnum = 0
		}

	}

}
