package common

import (
    _ "fmt"
    "github.com/Unknwon/goconfig"
    "github.com/astaxie/beego"
    "testing"
)

func init() {

}

// t.Error : output some error info and set testcase to failure
// t. Log  : output some error info and set testcase to success
func Test_it_rest_basic(t *testing.T) {
    env, err := goconfig.LoadConfigFile("../../../test/it/env.conf")
    if env == nil || err != nil {
        t.Error("open configuration fail")
        t.FailNow()
    }
    deployer_server, err := env.GetValue("cloudify", "server")

    url := deployer_server

    rest := New_Rest()
    statuscode, res, err := rest.Http_get(url, "")
    beego.Debug(url)

    if err != nil || res == "" {
        t.Error("Test_basic failed")
    } else {
        t.Log("Test_basic succeed, statuscode = ", statuscode)
    }
}

func test_url(t *testing.T, url, user1, pass1, host1, port1, uri1, query1 string) {
    username := ""
    userpass := ""
    host := ""
    port := ""
    uri := ""
    query := ""

    rest := New_Rest()
    username, userpass, host, port, uri, query, err := rest.Parse_URI(url)
    if err != nil {
            t.Log("Parse fail, url: ", url)
    }
    t.Log("user: ", username, " pass: ", userpass, " host: ",host, " port: ",port, " uri:",uri, " query:", query)

    if user1!= username {
        t.Error("wrong username, expect: ", user1, ", acutally: ", username)
    }
    if pass1!= userpass {
        t.Error("wrong pass, expect: ", pass1, ", acutally: ", userpass)
    }
    
    if host1!= host {
        t.Error("wrong host, expect: ", host1, ", acutally: ", host)
    }
    if port1!= port {
        t.Error("wrong port, expect: ", port1, ", acutally: ", port)
    }
    
    if uri1!= uri {
        t.Error("wrong uri, expect: ", uri1, ", acutally: ", uri)
    }
    if query1!= query {
        t.Error("wrong query, expect: ", query1, ", acutally: ", query)
    }
}

func Test_ut_parse_uri(t *testing.T) {
    url1 :="/abc/def?a1"
    url2 :="10.1.1.1/abc/def?a1"
    url3 :="10.1.1.1:200/abc/def?a1"
    url4 :="user:pass@10.1.1.1:200/abc/def?a1"

    test_url(t, url1, "", "", "", "", "/abc/def", "a1")
    test_url(t, url2, "", "", "10.1.1.1", "", "/abc/def", "a1")
    test_url(t, url3, "", "", "10.1.1.1", "200", "/abc/def", "a1")
    test_url(t, url4, "user", "pass", "10.1.1.1", "200", "/abc/def", "a1")
}

