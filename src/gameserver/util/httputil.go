package util

import (
	"net/http"

	"net/url"
	"sync"
	"time"
	"io"
	"io/ioutil"
	"crypto/tls"
)

var CONFIG = make(map[string]string)
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
func DoHttpRequest(method string, requrl string, contentType string, body io.Reader, token string, subjecttoken string) (data []byte, statusCode int, header http.Header, err error) {

	req, err := http.NewRequest(method, requrl, body)
	if err != nil {
		return nil, 500, nil, err
	}

	req.Header.Set("Content-Type", contentType)

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