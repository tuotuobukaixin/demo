package common

import (
    "github.com/astaxie/beego"
    "encoding/xml"
    "net/http"
    "strings"
    "strconv"
    "errors"
)

type HttpRedirector interface {
    Redirect(req *http.Request, respwriter http.ResponseWriter) (resp string, err error)
}

func New_http_redirecter(apis_confs []string, rest Rest) HttpRedirector{
    redirector := new(http_redirector_impl)

    for _, conf := range apis_confs {
        var vapis apis
        //beego.Trace("parsing api conf: \n", conf)
        err := xml.Unmarshal([]byte(conf), &vapis)
        if err != nil {
            beego.Error("parse redirector table fail, err: ", err)
            return nil
        }
        
        redirector.apis_confs = append(redirector.apis_confs, vapis)
    }
    
    redirector.rest = rest

    redirector.build_search_tbl()
    return redirector
}

type apis struct {
    XMLName     xml.Name     `xml:"apis"`
    Properties  properties   `xml:"properties"`
    Apis        []api        `xml:"api"`
}

type properties struct {
    XMLName     xml.Name   `xml:"properties"`
    Properties  []property `xml:"property"`
}

type property struct {
    XMLName   xml.Name `xml:"property"`
    Name      string   `xml:"name"`
    Value     string   `xml:"value"`
}

type api struct {
    XMLName   xml.Name `xml:"api"`
    Version   string   `xml:"version"`
    Uri       string   `xml:"uri"`
    Protocol  string   `xml:"protocol"`
    Dests     destinations `xml:"destinations"`
}

type destinations struct {
    XMLName         xml.Name `xml:"destinations"`
    Hosts           hosts    `xml:"hosts"`
    Destinations    []destination   `xml:"destination"`
}

type hosts struct {
    XMLName         xml.Name `xml:"hosts"`
    Hosts           []string   `xml:"host"`    
}

type destination struct {
    XMLName   xml.Name `xml:"destination"`
    Protocol  string   `xml:"protocol"`
    Method    string   `xml:"method"`
    Uri       string   `xml:"uri"`
}

type act_slow_api struct {
    method   string

    prefix   string
    pattern  string
    postfix  string
    
    actapi  act_api
}
type act_slow2_api struct {
    method   string

    prefix    string
    pattern1  string
    middle    string
    pattern2  string
    postfix   string
    
    actapi  act_api
}

type act_api struct {
    host     string
    uri      string
}

type http_redirector_impl struct {
    apis_confs      []apis
    api_fast_tbl     map[string]*act_api
    api_slow_tbl    []*act_slow_api
    api_slow2_tbl   []*act_slow2_api
    rest            Rest
}

func (this *http_redirector_impl)get_api_default_host(ap *apis) string {
    for _, prop := range ap.Properties.Properties {
        if prop.Name == "host" {
            return prop.Value
        }
    }

    return ""
}
func (this *http_redirector_impl)build_search_tbl() {
    this.api_fast_tbl = make(map[string]*act_api)
    
    for _, vapis := range this.apis_confs {
        defaulthost := this.get_api_default_host(&vapis)
            
        for _, api := range vapis.Apis {
            //if not host specified, use the property defined
            if len(api.Dests.Hosts.Hosts) == 0{
                api.Dests.Hosts.Hosts = append(api.Dests.Hosts.Hosts, defaulthost)
            } else if api.Dests.Hosts.Hosts[0] == "" {
                api.Dests.Hosts.Hosts[0] = defaulthost
            }
            
            for _, dest := range api.Dests.Destinations {
                api.Uri = strings.TrimSpace(api.Uri)
                beego.Debug("preprocessing uri: ", api.Uri)
                if strings.Contains(api.Uri, "{") && strings.Contains(api.Uri, "}") {
                    aa := strings.Split(api.Uri, "{")
                    if len(aa) == 2 {
                        prefix := aa[0]

                        bb := strings.Split(api.Uri, prefix)
                        Assert_int(len(bb), 2)
                        left := bb[1]
                        
                        aa = strings.Split(left, "}")
                        Assert_int(len(aa), 2)
                        pattern := aa[0]+"}"
                        postfix := aa[1]
                        
                        slow_actapi := act_slow_api{dest.Method, prefix, pattern, postfix, act_api{api.Dests.Hosts.Hosts[0], dest.Uri}}
                        this.api_slow_tbl = append(this.api_slow_tbl, &slow_actapi)
                        beego.Debug("add to api slow search tbl: ", dest.Method, ", ", prefix, " |", pattern, "| ", postfix, ", dsturl: ", slow_actapi.actapi.host + slow_actapi.actapi.uri)
                    } else if len(aa) == 3 {
                        prefix := aa[0]

                        bb := strings.Split(api.Uri, prefix)
                        Assert_int(len(bb), 2)
                        left := bb[1]

                        aa = strings.Split(left, "}")
                        Assert_int(len(aa), 3)
                        pattern1 := aa[0]+"}"

                        bb = strings.Split(left, pattern1)
                        Assert_int(len(bb), 2)
                        left = bb[1]
                        beego.Debug("got now: ", prefix, pattern1, left)
                        
                        bb = strings.Split(left, "{")
                        Assert_int(len(bb), 2)
                        middle := bb[0]
                        left = "{" + bb[1]

                        aa = strings.Split(left, "}")
                        Assert_int(len(aa), 2)
                        pattern2 := aa[0]+"}"
                        postfix := aa[1]

                        slow2_actapi := act_slow2_api{dest.Method, prefix, pattern1, middle, pattern2, postfix, act_api{api.Dests.Hosts.Hosts[0], dest.Uri}}
                        this.api_slow2_tbl = append(this.api_slow2_tbl, &slow2_actapi)
                        beego.Debug("add to api slow2 search tbl: ", dest.Method, ", ", prefix, " |", pattern1, "| ", middle, " | ", pattern2, " | ", postfix, ", dsturl: ", slow2_actapi.actapi.host + slow2_actapi.actapi.uri)
                    }
                    
                } else {
                    key := dest.Method + " " + api.Protocol + "://" + api.Uri
                    actapi := act_api{api.Dests.Hosts.Hosts[0], dest.Uri}                   
                    this.api_fast_tbl[key] = &actapi
                    beego.Debug("add to api fast search tbl, key: ", key, ", dsturl: ", actapi.host + actapi.uri)
                }
            }
        }   
    }
}

func (this *http_redirector_impl)get_act_uri(method string, uri string) (actapi *act_api, err error){
    //step 1: search fast tbl
    key := method+ " " + "http://" + uri
    actapi = this.api_fast_tbl[key]
    if actapi != nil {      
        return
    }

    //step 2: search slow table
    for _, slow_actapi := range this.api_slow_tbl {
        if method == slow_actapi.method {
            prefix_pos := strings.Index(uri, slow_actapi.prefix)
            prefix_end := prefix_pos + len(slow_actapi.prefix)
            if slow_actapi.postfix != "" {
                postfix_pos := strings.Index(uri, slow_actapi.postfix)
                if prefix_pos != -1 && postfix_pos != -1  && postfix_pos > prefix_end {
                    aa := strings.Split(uri, slow_actapi.prefix)
                    Assert_int(len(aa), 2)

                    bb := strings.Split(aa[1], slow_actapi.postfix)
                    Assert_int(len(bb), 2)
                    generic_matched := bb[0]

                    actapi1 := slow_actapi.actapi
                    actapi1.uri = strings.Replace(actapi1.uri, slow_actapi.pattern, generic_matched, 1)
                    actapi = &actapi1
                    return
                }   
            } else {
                if prefix_pos != -1 {
                    aa := strings.Split(uri, slow_actapi.prefix)
                    Assert_int(len(aa), 2)
                    generic_matched := aa[1]

                    actapi1 := slow_actapi.actapi
                    actapi1.uri = strings.Replace(actapi1.uri, slow_actapi.pattern, generic_matched, 1)
                    actapi = &actapi1
                    return
                }   
            }
            
        }
    }

    //step 3: search slow2 table
    for _, slow2_actapi := range this.api_slow2_tbl {
        if method == slow2_actapi.method {
            prefix_pos := strings.Index(uri, slow2_actapi.prefix)
            prefix_end := prefix_pos + len(slow2_actapi.prefix)

            middle_pos := strings.Index(uri, slow2_actapi.middle)
            middle_end := middle_pos + len(slow2_actapi.middle)

            beego.Debug("pre_pos:", prefix_pos, prefix_end, "mid_pos:", middle_pos, middle_end)
            beego.Debug("postfix: [" + slow2_actapi.postfix + "]")
            if slow2_actapi.postfix == "" {
                if prefix_pos != -1 && middle_pos != -1  && middle_pos > prefix_end  {                  
                    aa := strings.Split(uri, slow2_actapi.prefix)
                    Assert_int(len(aa), 2)

                    bb := strings.Split(aa[1], slow2_actapi.middle)
                    Assert_int(len(bb), 2)
                    generic_matched1 := bb[0]
                    generic_matched2 := bb[1]                   
                    
                    actapi1 := slow2_actapi.actapi
                    actapi1.uri = strings.Replace(actapi1.uri, slow2_actapi.pattern1, generic_matched1, 1)
                    actapi1.uri = strings.Replace(actapi1.uri, slow2_actapi.pattern2, generic_matched2, 1)
                    actapi = &actapi1
                    return
                } 
            } else {
                postfix_pos := strings.Index(uri, slow2_actapi.postfix)
                if prefix_pos != -1 && postfix_pos != -1 && middle_pos != -1  && 
                    middle_pos > prefix_end && postfix_pos > middle_end {
                    aa := strings.Split(uri, slow2_actapi.prefix)
                    Assert_int(len(aa), 2)

                    bb := strings.Split(aa[1], slow2_actapi.middle)
                    Assert_int(len(bb), 2)
                    generic_matched1 := bb[0]
                    
                    aa = strings.Split(bb[1], slow2_actapi.postfix)
                    Assert_int(len(bb), 2)
                    generic_matched2 := bb[0]
                    
                    actapi1 := slow2_actapi.actapi
                    actapi1.uri = strings.Replace(actapi1.uri, slow2_actapi.pattern1, generic_matched1, 1)
                    actapi1.uri = strings.Replace(actapi1.uri, slow2_actapi.pattern2, generic_matched2, 1)
                    actapi = &actapi1
                    return
                } 
            }
            
        }
    }
    
    err = ENotFound
    return
}

type DebugResponseWriter struct {
    respwriter http.ResponseWriter
    Body []byte
    Statuscode int
}
func New_DebugResponseWriter(respwriter http.ResponseWriter) *DebugResponseWriter {
    revproxy_respwriter := new(DebugResponseWriter)

    revproxy_respwriter.respwriter = respwriter
    return revproxy_respwriter
}
func (this *DebugResponseWriter)Header() http.Header {
    return this.respwriter.Header()
}
func (this *DebugResponseWriter)Write(body []byte) (int, error) {
    ret, err := this.respwriter.Write(body)

    this.Body = body
    return ret, err
}
func (this *DebugResponseWriter)WriteHeader(code int) { 
    header := this.Header()
    EnforceHeadersPolicy(header)    

    this.respwriter.WriteHeader(code)
    this.Statuscode = code
}

type InnerErrorHidingResponseWriter struct {
    respwriter http.ResponseWriter
    Body []byte
    Statuscode int
}
func New_InnerErrorHidingResponseWriter(respwriter http.ResponseWriter) *InnerErrorHidingResponseWriter {
    revproxy_respwriter := new(InnerErrorHidingResponseWriter)

    revproxy_respwriter.respwriter = respwriter
    return revproxy_respwriter
}
func (this *InnerErrorHidingResponseWriter)Header() http.Header {
    return this.respwriter.Header()
}
func (this *InnerErrorHidingResponseWriter)Write(body []byte) (int, error) {
    ret, err := 0, error(nil)
    if this.Statuscode >= 500 {
        beego.Debug("Clear body for internal error. input body[:32]: ", string(body)[:32]) 
        ret, err = this.respwriter.Write([]byte(""))
    } else {
        ret, err = this.respwriter.Write(body)
        this.Body = body
    }
    return ret, err
}

func EnforceHeadersPolicy(header http.Header) { 
    //beego.Debug("ReverseProxy headers: origin %+v", header)
    header.Del("Server")
    header.Add("Server","Server")
    //"Cache-Control: no-store"
    //"Pragma: no-cache"
    //"Cache-Control: no-cache"
    header.Add("Strict-Transport-Security",
        "max-age=31536000; includeSubDomains")
    header.Add("Cache-control", "no-cache, no-store")
    header.Add("Pragma", "no-cache")
}


func (this *InnerErrorHidingResponseWriter)WriteHeader(code int) {  
    header := this.Header()
    EnforceHeadersPolicy(header)    
    this.Statuscode = code
    if code >= 500 {
        beego.Debug("Hide internal error code: ", code) 
        code = 400
    }
    this.respwriter.WriteHeader(code)
}


func (this *http_redirector_impl)Redirect(req *http.Request, respwriter http.ResponseWriter) (resp string, err error) {
    requrl := req.URL.String()
    beego.Debug("Redirect for original url : ", requrl)
    
    request_uri := ""
    query := ""
    _, _, _, _, request_uri, query, err = this.rest.Parse_URI(requrl)   
    if err != nil {
        beego.Warn("Redirector parse original url fail, err: ", err)
        return 
    }
    method := req.Method

    var actapi *act_api = nil   
    actapi, err = this.get_act_uri(method, request_uri)
    if err != nil {
        beego.Warn("No matched newurl found for : ", method, request_uri, " ,error:", err)
        return 
    }
    newuri := actapi.uri

    revproxy_respwriter := New_DebugResponseWriter(respwriter)
    err = this.rest.ReverseProxy(actapi.host, newuri, query, req, revproxy_respwriter)
    if err != nil {
        beego.Warn("Reverse proxy fail, for : ", actapi.host + newuri, " ,error:", err)
        return 
    }
        
    //beego.Debug("ReverseProxy headers: new %+v", revproxy_respwriter.Header())

    resp = string(revproxy_respwriter.Body)
    beego.Debug("ReverseProxy result, reqUri: ",req.RequestURI, " ,statuscode:", revproxy_respwriter.Statuscode)    
    if revproxy_respwriter.Statuscode /100 == 2 {
        beego.Debug("ReverseProxy succeed: [url] = ", actapi.host + newuri)
    } else {
        beego.Warn("ReverseProxy failed. reqUri: ",req.RequestURI, " ,statuscode:", strconv.Itoa(revproxy_respwriter.Statuscode))
        if revproxy_respwriter.Statuscode == 404 {
            err = ENotFound
        } else {
            err = errors.New("ReverseProxy failed. statuscode:" + strconv.Itoa(revproxy_respwriter.Statuscode))
        }
        return
    }   

    return
}
