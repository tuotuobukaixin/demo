package common

import (
    "errors"
    "github.com/astaxie/beego"
    "io"
    "io/ioutil"
    "net/http"
    "net/http/httputil"
    "net/url"
    "crypto/x509"
    "crypto/tls"
    "strings"
    "strconv"
    "time"
)

type Rest interface {
    Set_timeout(timeout int) //second
    Http_post_json(url string, jsonbody string) (statuscode int, resbody string, err error)
    Http_post_ext(url string, sbody string, headers map[string]string) (statuscode int, resbody string, err error)
    
    Http_put_octectstream(url string, bodycontent []byte) (statuscode int, err error)
    Http_put_octectstream_ext(url string, bodycontent []byte, headers map[string]string) (statuscode int, err error)

    Http_put_json(url string, bodycontent string) (statuscode int, resbody string, err error)
    Http_put_json_ext(url string, bodycontent string, headers map[string]string) (statuscode int, resbody string, err error)

    Http_get(url string, jsonbody string) (statuscode int, resbody string, err error)
    Http_get_ext(url string, jsonbody string, headers map[string]string) (statuscode int, resbody string, err error)

    Http_delete(url string, jsonbody string) (statuscode int, resbody string, err error)
    Http_delete_ext(url string, jsonbody string, headers map[string]string) (statuscode int, resbody string, err error)

    ReverseProxy(host string, newurl string, query string, req *http.Request, respwriter http.ResponseWriter) (err error)
    Parse_URI(url string) (username string, userpass string, host string, port string, uri string, query string, err error) 
}
const (
    HEADER_BASICAUTH_USERNAME = "basicauth_username"
    HEADER_BASICAUTH_PASSWD   = "basicauth_passwd"
    HEADER_CONTENT_TYPE       = "Content-Type"
    HEADER_CONTENT_LENGTH     = "Content-Length"
    HEADER_X_AUTH_TOKEN       = "X-Auth-Token"
    HEADER_TENANT_ID          = "Tenant-id"
)

const (
    DEFAULT_HTTP_REQUEST_TIMEOUT = 30
)
type rest_impl struct {
    timeout   int   //timeout for a request
    certpool  *x509.CertPool
}

func New_Rest() Rest {
    impl := new(rest_impl)
    if nil == impl {
        return nil
    }

    impl.timeout = DEFAULT_HTTP_REQUEST_TIMEOUT
    return impl
}

func New_Rest_withcert(certfiles []string) Rest {
    impl := new(rest_impl)
    if nil == impl {
        return nil
    }
    
    certpool := init_cert_pool(certfiles)
    impl.certpool = certpool

    if impl.certpool == nil {
         beego.Warn("Rest init with no cert pool.")
    }
    
    return impl
}

func init_cert_pool(certfiles []string) *x509.CertPool{
    beego.Debug("init_cert_pool...")
    
    if len(certfiles) == 0 {
        return nil
    }
    
    certpool := x509.NewCertPool()
    for _, certfile := range certfiles {
        pemCert, err := ioutil.ReadFile(certfile)
        if err != nil {
            beego.Warn("Load certification fail. error: ", err)
            return nil
        }

        beego.Debug("init_cert_pool: add certfile ", certfile)
        certpool.AppendCertsFromPEM(pemCert)
    }

    return certpool
}

func (this *rest_impl) set_tls(url string, req *http.Request, client *http.Client) {
    beego.Debug("set tls, req: ", req.URL.Scheme, ", ", req.URL)
    
    if(req.URL.Scheme == "https") {
        tlsconfig := tls.Config{}
        if this.certpool != nil {
            beego.Debug("set tls, set cert agent.")
            tlsconfig.RootCAs = this.certpool           
        } else {
            beego.Debug("set tls, skip verify.")
        }

        tlsconfig.InsecureSkipVerify = true

        client.Transport = &http.Transport{
            TLSClientConfig: &tlsconfig,
        }
    }
}

func (this *rest_impl) set_headers(req *http.Request, headers map[string]string) (statuscode int, resbody string, err error) {
    var basicauth_username = ""
    var basicauth_passwd = ""

    if headers == nil {
        return
    }
    
    for header_name, header_value :=range headers {
        if header_name == HEADER_BASICAUTH_USERNAME {
            basicauth_username = header_value
        } else if header_name == HEADER_BASICAUTH_PASSWD {
            basicauth_passwd = header_value
        } else {
            req.Header.Set(header_name, header_value)
        }
    }

    if basicauth_username != "" {
//      beego.Debug("Http set basic auth header: username/passwd", basicauth_username, " / ", basicauth_passwd)
        req.SetBasicAuth(basicauth_username, basicauth_passwd)
    }

    return
}

func (this *rest_impl) Set_timeout(timeout int) () {
    this.timeout = timeout
}

func (this *rest_impl) Http_post_json(url string, jsonbody string) (statuscode int, resbody string, err error) {
    var headers map[string]string
    headers = make(map[string]string)
    if jsonbody != "" {
        headers[HEADER_CONTENT_TYPE] = "application/json"
    }
    return this.Http_post_ext(url, jsonbody, headers)
}

func (this *rest_impl) new_client_req(method, url, sbody string, headers map[string]string) (client *http.Client, req *http.Request, err error) {
    beego.Debug("http_request: method = ", method, " ,url = ", url, " ,timeout = ", this.timeout)
//  beego.Debug("http_request: url = ", url , " ,\nbody", sbody)

    body := io.Reader(nil)
    if sbody != "" {
        body = strings.NewReader(sbody)
    }
    req, err = http.NewRequest(method, url, body)   
    if err != nil {
        beego.Warn("new_client_req create request failed. method = ", method, " ,url = ", url, " ,error:", err) 
//      beego.Warn("new_client_req create request failed. method = ", method, " ,url = ", url , ", error:", err, " ,\nbody", sbody) 
        return
    }
    
    this.set_headers(req, headers)

    req.ContentLength = int64( len(sbody) )

    //client with timeout settings
    timeout := time.Duration(time.Duration(this.timeout) * time.Second)
    client = &http.Client{
        Timeout: timeout}
    this.set_tls(url, req, client)

    return
}


func (this *rest_impl) Http_post_ext(url string, sbody string, headers map[string]string) (statuscode int, resbody string, err error) {
    beego.Debug("Http_post_ext: [url] = ", url)
//  beego.Debug("Http_post_ext: [url][params] = ", url, sbody)

    var client *http.Client
    var req    *http.Request
    client, req, err = this.new_client_req("POST", url, sbody, headers)
    if err != nil {
        beego.Warn("Http_post_ext: create client/req failed. error: ", err)
        return
    }
    
    var resp *http.Response
    beego.Info("Sent http post request, [url]=", url)
    resp, err = client.Do(req)
    beego.Info("3rd system response the request, error:", err)
    var resbodybuf []byte
    if resp != nil {
        defer resp.Body.Close()
        statuscode = resp.StatusCode
        resbodybuf, _ = ioutil.ReadAll(resp.Body)       
    } else {
        if err == nil {
            err = EHttpError
        }
    }
    resbody = string(resbodybuf)

    if err != nil || statuscode/100 != 2 {
        if err == nil {
            err = errors.New("http post failed, statuscode: " + strconv.Itoa(statuscode))
        }
        beego.Warn("Http_post failed: [url:] ", url, "[statuscode:]", statuscode, "[resp:]", resbody, " , error:", err)
        
        return
    }

    beego.Debug("Http_post succeed: [url] = %s", url)
    return
}

func (this *rest_impl) Http_put_octectstream(url string, bodycontent []byte) (statuscode int, err error) {
    var headers map[string]string
    headers = make(map[string]string)
    if len(bodycontent) > 0{
        headers[HEADER_CONTENT_TYPE] = "application/octet-stream"
    }
    return this.Http_put_octectstream_ext(url, bodycontent, headers)
}

func (this *rest_impl) Http_put_octectstream_ext(url string, bodycontent []byte, headers map[string]string) (statuscode int, err error) {
    beego.Debug("Http_put_octectstream: [url]= ", url)
    //beego.Debug("Http_put_octectstream: [url][params]= ", url, bodycontent)

    if len(bodycontent) > 0{
        headers[HEADER_CONTENT_TYPE] = "application/octet-stream"
    }
    
    var client *http.Client
    var req    *http.Request
    client, req, err = this.new_client_req("PUT", url, string(bodycontent), headers)
    if err != nil {
        beego.Warn("Http_put_octectstream_ext: create client/req failed. error: ", err)
        return
    }
    
/*  
    body := strings.NewReader(string(bodycontent))
    req, _ := http.NewRequest("PUT", url, body)

    req.Header.Set("Content-Type", "application/octet-stream")
    this.set_headers(req, headers)

    client := &http.Client{}
    this.set_tls(url, req, client)
*/
    var resp *http.Response
    beego.Info("Sent http put request, [url]=", url)
    resp, err = client.Do(req)
    beego.Info("3rd system response the request, error:", err)

    var resbodybuf []byte
    if resp != nil {
        defer resp.Body.Close()
        statuscode = resp.StatusCode
        resbodybuf, _ = ioutil.ReadAll(resp.Body)       
    } else {
        if err == nil {
            err = EHttpError
        }
    }

    resbody := string(resbodybuf)

    if err != nil || statuscode/100 != 2 {
        if err == nil {
            err = errors.New("http put octect failed, statuscode: " + strconv.Itoa(statuscode))
        }
        beego.Warn("Http_put octect failed: [url:]", url, "[statuscode:]", statuscode, "[resp:]", resbody, ", error:", err)
        
        return
    }

    beego.Debug("Http_put succeed: [url:]", url, "[statuscode:]", statuscode)
    return
}

func (this *rest_impl) Http_put_json(url string, bodycontent string) (statuscode int, resbody string, err error) {
    var headers map[string]string
    headers = make(map[string]string)
    if len(bodycontent) > 0{
        headers[HEADER_CONTENT_TYPE] = "application/json"
    }
    return this.Http_put_json_ext(url, bodycontent, headers)
}
func (this *rest_impl) Http_put_json_ext(url string, bodycontent string, headers map[string]string) (statuscode int, resbody string, err error) {
    beego.Debug("Http_put_json: [url]= ", url)
//  beego.Debug("Http_put_json: [url][params]= ", url, bodycontent)

/*  
    body := strings.NewReader(string(bodycontent))
    req, _ := http.NewRequest("PUT", url, body)

    req.Header.Set("Content-Type", "application/json")
    this.set_headers(req, headers)

    client := &http.Client{}
    this.set_tls(url, req, client)
*/
    if bodycontent != "" {
        headers[HEADER_CONTENT_TYPE] = "application/json"
    }
    
    var client *http.Client
    var req    *http.Request
    client, req, err = this.new_client_req("PUT", url, bodycontent, headers)
    if err != nil {
        beego.Warn("Http_put_json_ext: create client/req failed. error: ", err)
        return
    }
    
    var resp *http.Response
    beego.Info("Sent http put request, url:", url)
    resp, err = client.Do(req)
    beego.Info("3rd system response the request, error:", err)

    var resbodybuf []byte
    if resp != nil {
        defer resp.Body.Close()
        statuscode = resp.StatusCode
        resbodybuf, _ = ioutil.ReadAll(resp.Body)       
    } else {
        if err == nil {
            err = EHttpError
        }
        beego.Warn("Http_put_json, resp:", resp)
    }
    resbody = string(resbodybuf)

    if err != nil || statuscode/100 != 2 {
        if err == nil {
            err = errors.New("http put json failed, statuscode: " + strconv.Itoa(statuscode))
        }
        beego.Warn("Http_put_json failed: [url:] ", url, "[statuscode:]", statuscode, "[resp:]", resbody, ", error:", err)
        
        return
    }

    beego.Debug("Http_put_json succeed: [url] = %s", url)
    return
}
func (this *rest_impl) Http_get(url string, jsonbody string) (statuscode int, resbody string, err error) {
    var headers map[string]string
    headers = make(map[string]string)
    if len(jsonbody) > 0{
        headers[HEADER_CONTENT_TYPE] = "application/json"
    }
    return this.Http_get_ext(url, jsonbody, headers)
}

func (this *rest_impl) Http_get_ext(url string, jsonbody string, headers map[string]string) (statuscode int, resbody string, err error) {
    beego.Debug("Http_get: [url] =  ", url)
    //beego.Debug("Http_get: [url][params] =  ", url, jsonbody)
/*  
    body := io.Reader(nil)
    if jsonbody != "" {
        body = strings.NewReader(jsonbody)
    }
    req, _ := http.NewRequest("GET", url, body)
    if jsonbody != "" {
        req.Header.Set("Content-Type", "application/json")
    }
    this.set_headers(req, headers)

    client := &http.Client{}
    this.set_tls(url, req, client)
*/

    var client *http.Client
    var req    *http.Request
    client, req, err = this.new_client_req("GET", url, jsonbody, headers)
    if err != nil {
        beego.Warn("Http_get_ext: create client/req failed. url:", url, " ,error: ", err)
        return
    }
    
    var resp *http.Response
    beego.Info("Sent http get request, [url]=", url)
    resp, err = client.Do(req)
    beego.Info("3rd system response the request, error:", err)
    var resbodybuf []byte
    if resp != nil {
        defer resp.Body.Close()
        statuscode = resp.StatusCode
        resbodybuf, _ = ioutil.ReadAll(resp.Body)       
    } else {
        if err == nil {
            err = EHttpError
        }
        beego.Warn("Http_get, resp:", resp)
    }
    resbody = string(resbodybuf)

    if err != nil || statuscode/100 != 2 {
        if statuscode == 404 {
            err = ENotFound
        } else if err == nil {
            err = errors.New("http get failed, statuscode: " + strconv.Itoa(statuscode))
        }
        beego.Warn("Http_get failed: [url:] ", url, "[statuscode:]", statuscode, "[resp:]", resbody, ", error: ", err)      
        
        return
    } else {
        beego.Debug("Http_get succeed: [url] = %s", url)
    }

    return
}

func (this *rest_impl) Http_delete(url string, jsonbody string) (statuscode int, resbody string, err error) {
    var headers map[string]string
    headers = make(map[string]string)
    if len(jsonbody) > 0{
        headers[HEADER_CONTENT_TYPE] = "application/json"
    }
    return this.Http_delete_ext(url, jsonbody, headers)
}
func (this *rest_impl) Http_delete_ext(url string, jsonbody string, headers map[string]string) (statuscode int, resbody string, err error) {
    beego.Debug("Http_delete: [url] = ", url)
    //beego.Debug("Http_delete: [url][params] = ", url, jsonbody)
/*  
    body := io.Reader(nil)
    if jsonbody != "" {
        body = strings.NewReader(jsonbody)
    }
    req, _ := http.NewRequest("DELETE", url, body)
    if jsonbody != "" {
        req.Header.Set("Content-Type", "application/json")
    }
    this.set_headers(req, headers)
    
    client := &http.Client{}
    this.set_tls(url, req, client)
*/
    if jsonbody != "" {
        headers[HEADER_CONTENT_TYPE] = "application/json"
    }
    var client *http.Client
    var req    *http.Request
    client, req, err = this.new_client_req("DELETE", url, jsonbody, headers)
    if err != nil {
        beego.Warn("Http_get_ext: create client/req failed. error: ", err)
        return
    }
    
    var resp *http.Response
    beego.Info("Sent http delete request, [url]=", url)
    resp, err = client.Do(req)
    beego.Info("3rd system response the request, error:", err)

    var resbodybuf []byte
    if resp != nil {
        defer resp.Body.Close()
        statuscode = resp.StatusCode
        resbodybuf, _ = ioutil.ReadAll(resp.Body)       
    } else {
        if err == nil {
            err = EHttpError
        }
        beego.Warn("Http_delete, resp:", resp)
    }
    resbody = string(resbodybuf)

    if err != nil || statuscode/100 != 2 {
        if err == nil {
            err = errors.New("http delete failed, statuscode: " + strconv.Itoa(statuscode))
        }
        beego.Warn("Http_delete failed: [url:] ", url, "[statuscode:]", statuscode, "[resp:]", resbody, ", error:", err)
        
        return
    }

    beego.Debug("Http_delete succeed: [url]= ", url)

    return
}

func (this *rest_impl) ReverseProxy(host string, newurl string, query string, req *http.Request, respwriter http.ResponseWriter) (err error) {  
    serverurl := "http://" + host

    var nurl *url.URL
    nurl, err = url.Parse(serverurl)
    if err != nil {
        beego.Warn("parse server url failed, url: ", serverurl, ", err: ", err)
        return
    }
    
    proxy := httputil.NewSingleHostReverseProxy(nurl)
    if proxy == nil {
        beego.Warn("crete http reverse proxy failed: host: ", host)
        return
    }
    req.URL.Path = newurl
    if query != "" {
        req.URL.RawQuery = query
    }
    
    beego.Debug("Proxying for url: ", req.URL.String())


    //muitipart special handling
    content_type := req.Header.Get("Content-Type")
    if strings.Contains(content_type, "avoidpart/form-data"){
        content_type = strings.Replace(content_type, "avoidpart/form-data", "multipart/form-data", -1) 
        req.Header.Set("Content-Type", content_type)
        beego.Debug("restore_multipart.")
    } 
    
    proxy.ServeHTTP(respwriter, req)

    return
}

func (this *rest_impl)is_at(c rune) bool {
    if c == '@' {
            return true
    } else {
            return false
    }
}
func (this *rest_impl)is_colon(c rune) bool {
    if c == ':' {
            return true
    } else {
            return false
    }
}
func (this *rest_impl)is_question(c rune) bool {
    if c == '?' {
            return true
    } else {
            return false
    }
}
func (this *rest_impl)is_slash(c rune) bool {
    if c == '/' {
            return true
    } else {
            return false
    }
}

func (this *rest_impl)Parse_URI(url string) (username string, userpass string, host string, port string, uri string, query string, err error) {
    left := url

    var aa []string
    var bb []string
    //find username, pass
    if strings.Contains(left, "@") {
        aa = strings.FieldsFunc(left, this.is_at)
        if len(aa) == 2 {
                bb = strings.FieldsFunc(aa[0], this.is_colon)
                bb_len := len(bb)
                if bb_len == 2 {
                        username = bb[0]
                        userpass = bb[1]
                } else {
                        err = EInvalidParam
                        return
                }

                left = aa[1]
        } else {
                username = ""
                userpass = ""
                left  = url
        }
    } else {
        username = ""
        userpass = ""
    }
    

    //find host, port
    aa = strings.Split(left, "/")
    if len(aa) > 0 && aa[0] != "" {
        aa[0] = strings.Trim(aa[0], " ")
        if aa[0] == "" {
            err = EInvalidParam
            return
        }
        
        hostinfo := ""
        bb = strings.Split(aa[0], ":")
        bb_len := len(bb)
        if bb_len == 1 {
            host = bb[0]
            port = ""
            hostinfo = host
        } else if bb_len == 2 {
            host = bb[0]
            port = bb[1]
            hostinfo = host + ":" + port
        } else {
            err = EInvalidParam
            return
        }
        
        bb = strings.Split(left, hostinfo)
        Assert(len(bb) > 0)
        left = bb[1]
    } else {
        host = ""
        port = ""
        left  = url
    }
    
    //find uri, query
    if strings.Contains(left, "?") {
        aa = strings.FieldsFunc(left, this.is_question)
        aa_len := len(aa)
        if aa_len == 1 {
                uri = aa[0]
                query = ""
        } else {
                uri = aa[0]
                query = aa[1]
        }
    } else {
        uri = left
        query = ""
    }
    
    return
}

