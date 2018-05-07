package utils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/xml"
	"esdk_obs_native_go/com/models"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Util struct {
	pathMap   map[string]string
	request   *http.Request
	bucket    string
	object    string
	ak        string
	sk        string
	endpoint  string
	pathStyle bool
}

var Response *http.Response

/**
*函数说明：初始化Util实例
*入参：AK：用户的AccessKeyID
*	  SK：用户的SecretAccessKeyID
*	  Endpoint：服务器地址，如（https://129.7.182.2:443）
*     PathStyle：请求方式是否为绝对路径方式，取值True 或 False
*返回值：初始化的Util实例
 */
func NewUtil(AK, SK, Endpoint string, PathStyle bool) *Util {
	util := &Util{ak: AK, sk: SK, endpoint: Endpoint, pathStyle: PathStyle}
	util.pathMap = make(map[string]string)
	return util
}

/**
*函数说明：初始化连接
*入参：mothed：请求方法
*	  bucket：桶名
*	  object：对象名
*     ioread：待发送的数据流
*返回值：执行失败值
 */
func (util *Util) InitConect(mothed, bucket, object string, ioread io.Reader) error {
	util.bucket = bucket
	util.object = object

	path := util.getPath()
	var err error
	util.request, err = http.NewRequest(mothed, path, ioread)
	if err != nil {
		return err
	}

	var host string
	if util.pathStyle == true {
		host = strings.Split(strings.Split(util.endpoint, "//")[1], ":")[0]
	} else {
		if bucket != "" {
			host = bucket + "." + strings.Split(strings.Split(util.endpoint, "//")[1], ":")[0]
		} else {
			host = strings.Split(strings.Split(util.endpoint, "//")[1], ":")[0]
		}
	}
	util.request.Header.Set("Host", host)
	util.request.Header.Set("Date", timestampS3())
	return nil
}

/**
*函数说明：向服务端发送请求
*入参：
*返回值：http.Response实例，Result执行结构信息
 */
func (util *Util) DoExec() (*http.Response, *models.Result) {
	if util.request.Header.Get("Authorization") == "" {
		sign := "AWS " + util.ak + ":"
		sign += signatureS3(stringToSignS3(util), util.sk)
		//fmt.Printf("signStr:<%s>\n", stringToSignS3(util))
		util.request.Header.Set("Authorization", sign)
	}
	var conn *http.Client
	var err error
	if strings.Split(util.endpoint, ":")[0] == "https" {
		tr := &http.Transport{
			TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
			DisableCompression: true,
		}
		conn = &http.Client{Transport: tr}
	} else {
		conn = &http.Client{}
	}

	Result := &models.Result{}
	for i := 0; i < 3; i++ {
		Response, err = conn.Do(util.request)
		if err == nil {
			break
		}
	}
	if err != nil {
		Result.Err = err
		return nil, Result
	}
	//defer Response.Body.Close()
	Result.StatusCode = Response.StatusCode

	if Result.StatusCode >= 300 {
		n, _ := strconv.Atoi(Response.Header.Get("Content-Length"))
		if n > 0 {
			var res models.ErrResponse
			body, err := ioutil.ReadAll(Response.Body)
			if err != nil {
				Result.Err = err
				return nil, Result
			}
			//fmt.Printf("body:%s\n", body)
			err = PareseXML(body, &res)
			if err != nil {
				Result.Err = err
				return nil, Result
			}
			Result.Code = res.Code
			Result.Message = res.Message
			Result.RequestId = res.RequestId
			Result.HostId = res.HostId
		}

	}
	return Response, Result
}

/**
*函数说明：设置请求URI
*入参：key：uri对应key值
*     value：uri的key对应值
*返回值：
 */
func (util *Util) SetPath(key string, value string) {
	util.pathMap[key] = value
}

/**
*函数说明：设置请求头域
*入参：key：head头对应key值
*     value：head头的key对应值
*返回值：
 */
func (util *Util) SetHeader(key string, value string) {
	util.request.Header.Set(key, value)
}

/**
*函数说明：关闭http数据流
*入参：
*返回值：
 */
func (util *Util) Close() {
	if Response != nil {
		defer Response.Body.Close()
	}
}

/**
*函数说明：获取URI
*入参：
*返回值：获取URI值
 */
func (util *Util) getPath() string {
	path := util.endpoint
	if util.bucket != "" {
		path += "/" + util.bucket
	}
	if util.object != "" {
		path += "/" + util.object
	}
	i := 0
	for key, value := range util.pathMap {
		if i == 0 {
			path += "?" + key
			i = 1
		} else {
			path += "&" + key
		}
		if value != "" {
			path += "=" + value
		}
	}
	return path
}

/**
*函数说明：解析xml
*入参：body:待解析的字符串
* 	 obj：解析的格式，传入对象指针
*返回值：
 */
func PareseXML(body []byte, obj interface{}) error {
	err := xml.Unmarshal(body, obj)
	return err
}

/**
*函数说明：构造v2鉴权的signature值
*入参：stringToSign:stringToSign值
* 	 SK：用户的SecretAccessKeyID
*返回值：构造v2鉴权的signature值
 */
func signatureS3(stringToSign string, SK string) string {
	hashed := hmacSHA1([]byte(SK), stringToSign)
	return base64.StdEncoding.EncodeToString(hashed)
}

/**
*函数说明：构造v2鉴权的stringToSign值
*入参：util:Util实例
*返回值：构造v2鉴权的stringToSign值
 */
func stringToSignS3(util *Util) string {
	str := util.request.Method + "\n"

	if util.request.Header.Get("Content-Md5") != "" {
		str += util.request.Header.Get("Content-Md5")
	}
	str += "\n"

	str += util.request.Header.Get("Content-Type") + "\n"

	if util.request.Header.Get("Date") != "" {
		str += util.request.Header.Get("Date")
	} else {
		str += timestampS3()
	}
	str += "\n"

	canonicalHeaders := canonicalAmzHeadersS3(util)
	if canonicalHeaders != "" {
		str += canonicalHeaders
	}

	str += canonicalResourceS3(util)

	return str
}

/**
*函数说明：构造v2鉴权的canonicalAmzHeaders值
*入参：util:Util实例
*返回值：构造v2鉴权的canonicalAmzHeaders值
 */
func canonicalAmzHeadersS3(util *Util) string {
	var headers []string

	for header := range util.request.Header {
		standardized := strings.ToLower(strings.TrimSpace(header))
		if strings.HasPrefix(standardized, "x-amz") {
			headers = append(headers, standardized)
		}
	}

	sort.Strings(headers)

	for i, header := range headers {
		headers[i] = header + ":" + strings.Replace(util.request.Header.Get(header), "\n", " ", -1)
	}

	if len(headers) > 0 {
		return strings.Join(headers, "\n") + "\n"
	} else {
		return ""
	}
}

/**
*函数说明：构造v2鉴权的canonicalResource值
*入参：util:Util实例
*返回值：构造v2鉴权的canonicalResource值
 */
func canonicalResourceS3(util *Util) string {
	res := "/"

	if util.bucket != "" {
		res += util.bucket
	}
	if util.object != "" {
		res += "/" + util.object
	}

	i := 0
	for _, subres := range strings.Split(subresourcesS3, ",") {
		val, ok := util.pathMap[subres]
		if ok {
			if i == 0 {
				res += "?" + subres
				i = 1
			} else {
				res += "&" + subres
			}
			if val != "" {
				res += "=" + val
			}
		}
	}

	return res
}

/**
*函数说明：获取当前RFC时间
*入参：
*返回值：当前RFC时间
 */
func timestampS3() string {
	return now().Format(timeFormatS3)
}

/**
*函数说明：hmacSHA256加密
*入参：key:加密的key值
*     content加密的字符串
*返回值：加密后的值
 */
func hmacSHA256(key []byte, content string) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(content))
	return mac.Sum(nil)
}

/**
*函数说明：hmacSHA1加密
*入参：key:加密的key值
*     content加密的字符串
*返回值：加密后的值
 */
func hmacSHA1(key []byte, content string) []byte {
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(content))
	return mac.Sum(nil)
}

/**
*函数说明：hashSHA256加密
*入参：key:加密的key值
*     content加密的字符串
*返回值：加密后的值
 */
func hashSHA256(content []byte) string {
	h := sha256.New()
	h.Write(content)
	return fmt.Sprintf("%x", h.Sum(nil))
}

/**
*函数说明：md5加密
*入参：content加密的字符串
*返回值：加密后的值
 */
func HashMD5(content []byte) string {
	h := md5.New()
	h.Write(content)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

/**
*函数说明：获取当前UTC时间
*返回值：UTC时间
 */
var now = func() time.Time {
	return time.Now().UTC()
}

const (
	timeFormatS3   = time.RFC1123Z
	subresourcesS3 = "acl,quota,storageinfo,deletebucket,delete,lifecycle,location,logging,notification,partNumber,policy,requestPayment,torrent,uploadId,uploads,versionId,versioning,versions,website"
)
