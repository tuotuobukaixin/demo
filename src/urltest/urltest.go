package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
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
	httpClient = &http.Client{Timeout: timeout}
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
func config() {
	fi, err := os.Open("config")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		if strings.Contains(string(a), "=") {
			CONFIG[strings.Split(string(a), "=")[0]] = strings.Split(string(a), "=")[1]
		}
	}
}
func urlconfig() {
	fi, err := os.Open(CONFIG["urlfile"])
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()
	num, _ := strconv.Atoi(CONFIG["threadnum"])
	br := bufio.NewReader(fi)
	index := 0
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}

		URLINFO[index] = fmt.Sprintf("%s %s", URLINFO[index], string(a))
		index++
		if index == num {
			index = 0
		}
	}
	return
}
func theardfuc(num int) {
	for {
		urlip := strings.Fields(URLINFO[num])
		for i := 0; i < len(urlip); i++ {
			if urlip[i] != "" {
				router := fmt.Sprintf("http://%s", urlip[i])
				_, status_code, _, _ := DoHttpRequest("GET", router, "application/json", nil, "", "")
				fmt.Println(fmt.Sprintf("%d %s %d", num, urlip[i], status_code))
			}

		}
	}
}
func main() {
	config()
	fmt.Println(CONFIG)
	urlconfig()
	fmt.Println(URLINFO)
	num, _ := strconv.Atoi(CONFIG["threadnum"])
	for a := 0; a < num; a++ {
		go theardfuc(a)
	}
	time.Sleep(360000 * time.Hour)
}
