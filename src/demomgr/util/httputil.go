package util

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"

	"net/url"
	"sync"
	"time"
)

var connectionPool = struct {
	sync.RWMutex
	pool map[string]*http.Client
}{pool: make(map[string]*http.Client)}

//InsecurityConnection define
var InsecurityConnection = &http.Client{
	Timeout: time.Duration(30 * time.Second),
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		//DisableKeepAlives: true,
	},
}

//NewConnection :create a new connection
func NewConnection(requestPath string, caKey []byte, crt []byte, key []byte) (httpClient *http.Client, err error) {
	requestURL, err := url.Parse(requestPath)
	if err != nil {
		return
	}
	requestHost := requestURL.Host
	if caKey == nil || len(caKey) == 0 { //This is a insecurity connection
		httpClient = InsecurityConnection
		connectionPool.Lock()
		connectionPool.pool[requestHost] = InsecurityConnection
		connectionPool.Unlock()
		return
	}

	//Create a new security connection and update the cache
	httpClient, err = createConnection(caKey, crt, key)
	if err == nil {
		connectionPool.Lock()
		connectionPool.pool[requestHost] = httpClient
		connectionPool.Unlock()
	}
	return
}

func createConnection(caKey []byte, crt []byte, key []byte) (httpClient *http.Client, err error) {
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caKey)

	clientCrt, err := tls.X509KeyPair(crt, key)
	if err != nil {
		LOGGER.Error("X509KeyPair err:", err)
		return
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:      pool,
			Certificates: []tls.Certificate{clientCrt},
		},
	}
	timeout := time.Duration(10 * time.Second)
	httpClient = &http.Client{Transport: tr, Timeout: timeout}
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


