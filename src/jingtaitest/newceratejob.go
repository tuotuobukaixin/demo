package main

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
	"crypto/tls"
	"encoding/json"
	"strconv"
)
var DELETE = false
var JOBNUM = 20000




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
		req.Header.Set("Authorization", "Bearer " +token)
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


func createjob(urlstring string,sfs_name string, token string,job_name string)  {
	job :=`{	"apiVersion": "batch/v1",	"kind": "Job",	"metadata": {		"name": "`+job_name+`",
		"namespace": "default"
	},
	"spec": {
		"template": {
			"metadata": {
				"name": "`+job_name+`"
			},
			"spec": {
				"containers": [{
					"env": [{
							"name": "TIMEOUT",
							"value": "3600"
						},
						{
							"name": "JOB_NAME",
							"value": "`+job_name+`"
						}
					],
					"image": "swr.cn-north-1.myhuaweicloud.com/jingtai/logtest:latest",
					"imagePullPolicy": "Always",
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
					"terminationMessagePolicy": "File",
					"volumeMounts": [{
						"mountPath": "/root/file",
						"name": "`+sfs_name+`"
					}]
				}],
				"dnsPolicy": "ClusterFirst",
				"imagePullSecrets": [{
					"name": "default-secret"
				}],
				"restartPolicy": "OnFailure",
				"schedulerName": "default-scheduler",
				"securityContext": {},
				"terminationGracePeriodSeconds": 30,
				"volumes": [{
					"name": "`+sfs_name+`",
					"persistentVolumeClaim": {
						"claimName": "`+sfs_name+`"
					}
				}]
			}
		}
	}
}
`
	data := bytes.NewReader([]byte(job))
	router :=  urlstring + "/apis/batch/v1/namespaces/default/jobs"
	_, status_code, _, _ := DoHttpRequest("POST", router, "application/json;charset=utf8", data, token, "")
	if status_code != 201 {
		return
	}
}

func Getjob(urlstring string, token string) []byte{
	router := urlstring + "/apis/batch/v1/namespaces/default/jobs"
	rsp, status_code, _, _ := DoHttpRequest("GET", router, "application/json", nil, token, "")
	if status_code != 200 {
		fmt.Println(string(rsp))
		return nil
	}
	fmt.Println(string(rsp))
	return rsp

}
func DeleteJob(urlstring string, token string,job_name string) error{
	fmt.Println("begin to delete job " + job_name)
	body := `{
		"kind": "DeleteOptions",
		"apiVersion": "v1",
		"propagationPolicy": "Foreground"
	}`
	jsonBlob := []byte(body)

	data := bytes.NewReader([]byte(jsonBlob))
	router := urlstring +"/apis/batch/v1/namespaces/default/jobs/" + job_name
	rsp, status_code, _, _ := DoHttpRequest("DELETE", router, "application/json", data, token, "")
	if status_code != 200 {
		fmt.Println(string(rsp))
		return errors.New("delete failed")
	}
	return nil

}
func main() {
	endpoint :=[]string{"https://114.116.27.140:5443","https://119.3.2.42:5443","https://139.159.163.53:5443","https://139.159.161.40:5443"}
	token := []string{"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImRlZmF1bHQtdG9rZW4tNHZmYmMiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoiZGVmYXVsdCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjcxNjE2MmIzLTU0NjItMTFlOC1hMGJlLWZhMTYzZWMyZWJlNiIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OmRlZmF1bHQifQ.WBZ79zAxElI1xgXQXRhUuTnfPxBnqfT08pLpjE5AwbVD7-TWbPzGfoqOcbwoZfnsmgfFlFUCksX0cXjRR6qI1Oiiaajh__dQJdfPWbXSDCavinveYF_pVss5AKXP5_FwcpiobNCLxtMU1HsYP4niiNny91hoK5Is37JZsMiLFzDpdzy56UxeNxZpdfdWLOGCNM8Pqb6rVPqgFGMZnXdskir52wlhP52dFEmMyQX8s3Chfaxzde2N_rrZ_8lDynvKABR_jrbwvvB2q-SIv8vRQeNM3BR_FtyUtBMCCSBbGPeg-APwLt5flbNlZCYqFdezer-MIfq08brn7N0VolAMBg",
		"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImRlZmF1bHQtdG9rZW4tMjRiN3MiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoiZGVmYXVsdCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6Ijk4NmYxNzFiLTU1ZDQtMTFlOC1iODZkLWZhMTYzZTM4ODA0MSIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OmRlZmF1bHQifQ.PTHrKsXyuYEjeIacGL3yATEK_t8W5ZlsNhRcX5Z3QIOi4AuAtLXU1IxIITf0keeZI83YZyN-IMl_kaK_pKTgtt1v02wi6ofpglIjBmC7udeNO44J2FjwfKAyfl7cgXavTc3HHmaVuuZDkE674g_Nj602m74nl_q2hGYlKzgnrqzMM4HVNGvMMWYAc50G2Zmw7PZFZ7f_pHfXu-gW5S_H7lQXuj7G3TY1jt9HvxH5l9lwOihRv8AMUshRsHdcRg9KT_8OgIikY2HXtgjt2oa4OpBhyGo2f3yAir3V32TUJFwFneNP5aHd-nLcj4OruGDLps9ZAbjn5gj-AzY44tQZEg",
		"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImRlZmF1bHQtdG9rZW4tYnFyZmIiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoiZGVmYXVsdCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6ImU3MzU3YmRkLTU1ZDMtMTFlOC04OTEyLWZhMTYzZWU2MGYzYiIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OmRlZmF1bHQifQ.YTc-GAWlSlEushwxo8CnUoMGAZrenhwiQkZ3FtRFftczzVhsU4W9d_NnQun0XNJ3Lb9Q5r_qTKhZ_fBcQN98dMSDsk4sLBwCaPX3BJaMrjoHpBs0wReGTjdmxY7TIhHSydAvfNzd15lnN4Tq_caMU67E8QJPNsvkAOnBljd2cUqMS7qr1F-VoL795udDUHc4yD92hvUeVbn1jPKu6mdSn2rl5-SfJuqDvSsxz14zLn9gahmrD_LMwo68RJGRZ5lvGlROtMSzWIAWjnlG3YNQm263PKdiUQHj_UJSMlhXPApdCklD8rvSo2ep0sie84ZoW4KcKK6svlVykqsDKc9mDA",
		"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImRlZmF1bHQtdG9rZW4tdmhmdHEiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoiZGVmYXVsdCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6ImYyZjI4YzQyLTU3NDItMTFlOC05N2RkLWZhMTYzZWE3OGRhNiIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OmRlZmF1bHQifQ.t0aC0sBU1_cT73hY7baV-WxfUrMmYliPIhPnYLI7dtawAe615Z3INFyhaLPfy5xG1YKs54LirUB8dDIzD-DsUgAyEKu86nS8AQMBoTVx0qGN8xL1U6TOPvIcAzotDkUnOZRFC98-GOlwau3SwBMExuLyTqJaKqgbSKZ4Tf_gRyFd0l4QziVQRQVqsVGxD4NByBneo8_08zDlSCEulnx3OrnRpF1ZiQuxwaEZDQgNz4VfXteooNpM4xPLTGKu5RojOOsXh8pnCajOcCzhZKP_g0QumwhDhdOdvDr4TfBSGlhhwv114SfGBKtUbeR6v8g2HGVcKY7OpSXvYO2Pan7_rA"}
	pvc := []string{"cce-sfs-jh1q6m9p-iycq","cce-sfs-jh3aoacy-7fm6","cce-sfs-jh3apgcw-szn0","cce-sfs-jh64l1s2-h86m"}


	//for {
	var dat map[string]interface{}
	for index := 3 ; index < 4 ; index++ {
		rsp := Getjob(endpoint[index], token[index])
		json.Unmarshal(rsp, &dat)
		if dat["items"] == nil || len(dat["items"].([]interface{})) == 0 {
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
			fmt.Println("start to create job")
			for i := 0; i < JOBNUM; i++ {
				if !DELETE {
					createjob(endpoint[index], pvc[index], token[index], "jobtest"+strconv.Itoa(i))
				}
			}
			rsp = Getjob(endpoint[index], token[index])
			json.Unmarshal(rsp, &dat)
		}
		for i := 0; i < JOBNUM; i++ {
			if DELETE {
				DeleteJob(endpoint[index], token[index], "jobtest"+strconv.Itoa(i))
			}
		}
	}
	//for _, intem := range dat["items"].([]interface{}) {
	//	_, flag1 := intem.(map[string]interface{})["status"]
	//	if flag1 {
	//		_, flag2 := intem.(map[string]interface{})["status"].(map[string]interface{})["conditions"]
	//		if flag2 {
	//			_, flag3 := intem.(map[string]interface{})["status"].(map[string]interface{})["conditions"].([]interface{})[0].(map[string]interface{})["type"]
	//			if flag3 {
	//				if intem.(map[string]interface{})["status"].(map[string]interface{})["conditions"].([]interface{})[0].(map[string]interface{})["type"].(string) == "Complete" {
	//					DeleteJob("b8c96ab0-5461-11e8-941c-0255ac101f3e", token, intem.(map[string]interface{})["metadata"].(map[string]interface{})["name"].(string))
	//				}
	//			}
	//		}
	//	}
	//}
	time.Sleep(10*time.Second)
	//}



}
