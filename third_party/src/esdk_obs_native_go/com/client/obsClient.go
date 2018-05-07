package client

import (
	"errors"
	"esdk_obs_native_go/com/models"
	"esdk_obs_native_go/com/utils"
	//"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Client struct {
	AK        string
	SK        string
	Endpoint  string
	PathStyle bool
}

/**
*函数原型：func Factory(AK, SK, Endpoint string, PathStyle bool) (client *Client)
*函数功能：初始化Client类实例
*参数说明：AK：用户的AccessKeyID
*		 SK：用户的SecretAccessKeyID
*		 Endpoint：服务器地址，如（https://129.7.182.2:443）
*		 PathStyle：请求方式是否为绝对路径方式，取值True 或 False
*返回值：Client实例化对象
 */
func Factory(AK, SK, Endpoint string, PathStyle bool) (client *Client) {
	return &Client{AK, SK, Endpoint, PathStyle}
}

/**
*函数原型：func (client *Client) CreateBucket(input *models.CreateBucketInput) (requst *models.Result)
*函数功能：创建桶
*参数说明：input: CreateBucketInput对象实例
*返回值：requst: Resul对象实例
 */
func (client *Client) CreateBucket(input *models.CreateBucketInput) (requst *models.Result) {
	requst = &models.Result{}
	if input.Bucket == "" {
		requst.Err = errors.New(models.BUCKET_NIL)
		return requst
	}
	util := utils.NewUtil(client.AK, client.SK, client.Endpoint, client.PathStyle)
	if input.ACL != "" {
		util.SetPath("x-amz-acl", input.ACL)
	}
	xml := utils.CreatBucketXML(input)
	ioRead := strings.NewReader(xml)

	err := util.InitConect("PUT", input.Bucket, "", ioRead)
	if err != nil {
		requst.Err = err
		return requst
	}
	util.SetHeader("Content-Length", strconv.Itoa(len(xml)))
	_, requst = util.DoExec()
	util.Close()
	return requst
}

/**
*函数原型：func (client *Client) HeadBucket(bucketName string) (requst *models.Result)
*函数功能：查询桶是否存在
*参数说明：bucketName: 桶名
*返回值：requst: Resul对象实例
 */
func (client *Client) HeadBucket(bucketName string) (requst *models.Result) {
	requst = &models.Result{}
	if bucketName == "" {
		requst.Err = errors.New(models.BUCKET_NIL)
		return requst
	}
	util := utils.NewUtil(client.AK, client.SK, client.Endpoint, client.PathStyle)
	err := util.InitConect("HEAD", bucketName, "", nil)
	if err != nil {
		requst.Err = err
		return requst
	}
	_, requst = util.DoExec()
	util.Close()
	return requst
}

/**
*函数原型：func (client *Client) ListBuckets() (requst *models.Result, listBucketOutPut *models.ListBucketsOutput)
*函数功能：获取桶列表
*参数说明：
*返回值：requst: Resul对象实例
*	   listBucketOutPut：ListBucketsOutput对象实例
 */
func (client *Client) ListBuckets() (requst *models.Result, listBucketOutPut *models.ListBucketsOutput) {
	requst = &models.Result{}
	util := utils.NewUtil(client.AK, client.SK, client.Endpoint, client.PathStyle)
	err := util.InitConect("GET", "", "", nil)
	if err != nil {
		requst.Err = err
		return requst, nil
	}
	Response, requst := util.DoExec()
	if requst.Err != nil {
		return requst, nil
	}
	xml, err := ioutil.ReadAll(Response.Body)
	if err != nil {
		requst.Err = err
		return requst, nil
	}
	util.Close()
	listBucketOutPut = &models.ListBucketsOutput{}
	err = utils.PareseXML(xml, listBucketOutPut)
	if err != nil {
		requst.Err = err
		return requst, nil
	}
	return requst, listBucketOutPut
}

/**
*函数原型：func (client *Client) DeleteBucket(bucketName string) (requst *models.Result)
*函数功能：删除桶
*参数说明：bucketName：桶名
*返回值：requst: Resul对象实例
 */
func (client *Client) DeleteBucket(bucketName string) (requst *models.Result) {
	requst = &models.Result{}
	if bucketName == "" {
		requst.Err = errors.New(models.BUCKET_NIL)
		return requst
	}
	util := utils.NewUtil(client.AK, client.SK, client.Endpoint, client.PathStyle)
	err := util.InitConect("DELETE", bucketName, "", nil)
	if err != nil {
		requst.Err = err
		return requst
	}
	_, requst = util.DoExec()
	util.Close()
	return requst
}

/**
*函数原型：func (client *Client) SetBucketQuota(bucketName string, storageQuota int) (requst *models.Result)
*函数功能：设置桶的配额
*参数说明：bucketName：桶名
*		 storageQuota：桶空间的配额，单位为B字节
*返回值：requst: Resul对象实例
 */
func (client *Client) SetBucketQuota(bucketName string, storageQuota int) (requst *models.Result) {
	requst = &models.Result{}
	if bucketName == "" {
		requst.Err = errors.New(models.BUCKET_NIL)
		return requst
	}
	util := utils.NewUtil(client.AK, client.SK, client.Endpoint, client.PathStyle)
	util.SetPath("quota", "")

	xml := utils.SetBucketQuotaXML(storageQuota)
	ioRead := strings.NewReader(xml)

	err := util.InitConect("PUT", bucketName, "", ioRead)
	if err != nil {
		requst.Err = err
		return requst
	}
	util.SetHeader("Content-Length", strconv.Itoa(len(xml)))
	_, requst = util.DoExec()
	util.Close()
	return requst
}

/**
*函数原型：func (client *Client) GetBucketQuota(bucketName string) (requst *models.Result, getBucketQuotaOutput *models.GetBucketQuotaOutput)
*函数功能：获取桶的空间配额
*参数说明：bucketName：桶名
*返回值：requst: Resul对象实例
*	   getBucketQuotaOutput：GetBucketQuotaOutput对象实例
 */
func (client *Client) GetBucketQuota(bucketName string) (requst *models.Result, getBucketQuotaOutput *models.GetBucketQuotaOutput) {
	requst = &models.Result{}
	if bucketName == "" {
		requst.Err = errors.New(models.BUCKET_NIL)
		return requst, nil
	}
	util := utils.NewUtil(client.AK, client.SK, client.Endpoint, client.PathStyle)
	util.SetPath("quota", "")

	err := util.InitConect("GET", bucketName, "", nil)
	if err != nil {
		requst.Err = err
		return requst, nil
	}
	Response, requst := util.DoExec()
	if requst.Err != nil {
		return requst, nil
	}
	xml, err := ioutil.ReadAll(Response.Body)
	if err != nil {
		requst.Err = err
		return requst, nil
	}

	util.Close()
	getBucketQuotaOutput = &models.GetBucketQuotaOutput{}
	err = utils.PareseXML(xml, getBucketQuotaOutput)
	if err != nil {
		requst.Err = err
		return requst, nil
	}
	return requst, getBucketQuotaOutput
}

/**
*函数原型：func (client *Client) GetBucketStorageInfo(bucketName string) (requst *models.Result, getBucketStorageInfoOutput *models.GetBucketStorageInfoOutput)
*函数功能：获取桶的存储信息
*参数说明：bucketName：桶名
*返回值：requst: Resul对象实例
*	   getBucketStorageInfoOutput：GetBucketStorageInfoOutput对象实例
 */
func (client *Client) GetBucketStorageInfo(bucketName string) (requst *models.Result, getBucketStorageInfoOutput *models.GetBucketStorageInfoOutput) {
	requst = &models.Result{}
	if bucketName == "" {
		requst.Err = errors.New(models.BUCKET_NIL)
		return requst, nil
	}
	util := utils.NewUtil(client.AK, client.SK, client.Endpoint, client.PathStyle)
	util.SetPath("storageinfo", "")

	err := util.InitConect("GET", bucketName, "", nil)
	if err != nil {
		requst.Err = err
		return requst, nil
	}
	Response, requst := util.DoExec()
	if requst.Err != nil {
		return requst, nil
	}
	xml, err := ioutil.ReadAll(Response.Body)
	if err != nil {
		requst.Err = err
		return requst, nil
	}
	util.Close()
	getBucketStorageInfoOutput = &models.GetBucketStorageInfoOutput{}
	err = utils.PareseXML(xml, getBucketStorageInfoOutput)
	if err != nil {
		requst.Err = err
		return requst, nil
	}
	return requst, getBucketStorageInfoOutput
}

/**
*函数原型：func (client *Client) DeleteBucketWithObjects(bucketName string) (requst *models.Result)
*函数功能：删除桶数据
*参数说明：bucketName：桶名
*返回值：requst: Resul对象实例
 */
func (client *Client) DeleteBucketWithObjects(bucketName string) (requst *models.Result) {
	requst = &models.Result{}
	if bucketName == "" {
		requst.Err = errors.New(models.BUCKET_NIL)
		return requst
	}
	util := utils.NewUtil(client.AK, client.SK, client.Endpoint, client.PathStyle)
	util.SetPath("deletebucket", "")

	xml := utils.DeleteBucketWithObjectsXML(bucketName)
	ioRead := strings.NewReader(xml)

	err := util.InitConect("POST", bucketName, "", ioRead)
	if err != nil {
		requst.Err = err
		return requst
	}
	util.SetHeader("Content-Length", strconv.Itoa(len(xml)))
	_, requst = util.DoExec()
	util.Close()
	return requst
}

/**
*函数原型：func (client *Client) ListObjects(input *models.ListObjectsInput) (requst *models.Result, listObjectsOutput *models.ListObjectsOutput)
*函数功能：获取桶内对象
*参数说明：bucketName：桶名
*返回值：requst: Resul对象实例
*		listObjectsOutput：ListObjectsOutput对象实例
 */
func (client *Client) ListObjects(input *models.ListObjectsInput) (requst *models.Result, listObjectsOutput *models.ListObjectsOutput) {
	requst = &models.Result{}
	if input.Bucket == "" {
		requst.Err = errors.New(models.BUCKET_NIL)
		return requst, nil
	}
	util := utils.NewUtil(client.AK, client.SK, client.Endpoint, client.PathStyle)
	if input.Prefix != "" {
		util.SetPath("prefix", input.Prefix)
	}
	if input.Marker != "" {
		util.SetPath("marker", input.Marker)
	}
	if input.MaxKeys != 0 {
		util.SetPath("max-keys", strconv.Itoa(input.MaxKeys))
	}
	if input.Delimiter != "" {
		util.SetPath("delimiter", input.Delimiter)
	}
	err := util.InitConect("GET", input.Bucket, "", nil)
	if err != nil {
		requst.Err = err
		return requst, nil
	}
	Response, requst := util.DoExec()
	if requst.Err != nil {
		return requst, nil
	}
	xml, err := ioutil.ReadAll(Response.Body)
	if err != nil {
		requst.Err = err
		return requst, nil
	}
	util.Close()
	listObjectsOutput = &models.ListObjectsOutput{}
	err = utils.PareseXML(xml, listObjectsOutput)
	if err != nil {
		requst.Err = err
		return requst, nil
	}
	return requst, listObjectsOutput
}

/**
*函数原型：func (client *Client) PutObject(input *models.PutObjectInput) (requst *models.Result, putObjectOutPut *models.PutObjectOutPut)
*函数功能：上传对象
*参数说明：input：PutObjectInput对象实例
*返回值：requst: Resul对象实例
*		putObjectOutPut：PutObjectOutPut对象实例
 */
func (client *Client) PutObject(input *models.PutObjectInput) (requst *models.Result, putObjectOutPut *models.PutObjectOutPut) {
	requst = &models.Result{}
	if input.Bucket == "" {
		requst.Err = errors.New(models.BUCKET_NIL)
		return requst, nil
	}
	if input.Object == "" {
		requst.Err = errors.New(models.OBJECT_NIL)
		return requst, nil
	}
	if input.Body != "" && input.SourceFile != "" {
		requst.Err = errors.New(models.BODY_FILE_EXIST)
		return requst, nil
	}
	util := utils.NewUtil(client.AK, client.SK, client.Endpoint, client.PathStyle)
	var ioRead io.Reader = nil
	var length string = "0"
	var err error = nil
	if input.Body != "" {
		ioRead = strings.NewReader(input.Body)
		length = strconv.Itoa(len(input.Body))
	} else if input.SourceFile != "" {
		ioRead, err = os.Open(input.SourceFile)
		if err != nil {
			requst.Err = err
			return requst, nil
		}
		var fi os.FileInfo
		fi, err = os.Stat(input.SourceFile)
		if err != nil {
			requst.Err = err
			return requst, nil
		}
		length = strconv.FormatInt(int64(fi.Size()), 10)
	}

	err = util.InitConect("PUT", input.Bucket, input.Object, ioRead)
	if err != nil {
		requst.Err = err
		return requst, nil
	}
	util.SetHeader("Content-Length", length)
	if input.ACL != "" {
		util.SetHeader("x-amz-acl", input.ACL)
	}
	if input.WebsiteRedirectLocation != "" {
		util.SetHeader("x-amz-website-redirect-location", input.WebsiteRedirectLocation)
	}

	for key, val := range input.Metadata {
		key = "x-amz-meta-" + key
		util.SetHeader(key, val)
	}
	Response, requst := util.DoExec()
	if requst.Err != nil {
		return requst, nil
	}

	putObjectOutPut = &models.PutObjectOutPut{}
	putObjectOutPut.VersionId = Response.Header.Get("X-Amz-Version-Id")
	putObjectOutPut.ETag = Response.Header.Get("Etag")
	util.Close()
	return requst, putObjectOutPut
}

/**
*函数原型：func (client *Client) GetObject(input *models.GetObjectInput) (requst *models.Result, getObjectOutput *models.GetObjectOutput)
*函数功能：获取对象内容
*参数说明：input：GetObjectInput对象实例
*返回值：requst: Resul对象实例
*		getObjectOutput：GetObjectOutput对象实例
 */
func (client *Client) GetObject(input *models.GetObjectInput) (requst *models.Result, getObjectOutput *models.GetObjectOutput) {
	requst = &models.Result{}
	if input.Bucket == "" {
		requst.Err = errors.New(models.BUCKET_NIL)
		return requst, nil
	}
	if input.Object == "" {
		requst.Err = errors.New(models.OBJECT_NIL)
		return requst, nil
	}
	util := utils.NewUtil(client.AK, client.SK, client.Endpoint, client.PathStyle)
	if input.ResponseCacheControl != "" {
		util.SetPath("response-cache-control", input.ResponseCacheControl)
	}
	if input.ResponseContentDisposition != "" {
		util.SetPath("response-content-disposition", input.ResponseContentDisposition)
	}
	if input.ResponseContentEncoding != "" {
		util.SetPath("response-content-encoding", input.ResponseContentEncoding)
	}
	if input.ResponseContentLanguage != "" {
		util.SetPath("response-content-language", input.ResponseContentLanguage)
	}
	if input.ResponseContentType != "" {
		util.SetPath("response-content-type", input.ResponseContentType)
	}
	if input.ResponseExpires != "" {
		util.SetPath("response-expires", input.ResponseExpires)
	}
	if input.VersionId != "" {
		util.SetPath("versionId", input.VersionId)
	}

	err := util.InitConect("GET", input.Bucket, input.Object, nil)
	if err != nil {
		requst.Err = err
		return requst, nil
	}

	if input.IfMatch != "" {
		util.SetHeader("If-Match", input.IfMatch)
	}
	if input.IfModifiedSince != "" {
		util.SetHeader("If-Modified-Since", input.IfModifiedSince)
	}
	if input.IfNoneMatch != "" {
		util.SetHeader("If-None-Match", input.IfNoneMatch)
	}
	if input.IfUnmodifiedSince != "" {
		util.SetHeader("If-Unmodified-Since", input.IfUnmodifiedSince)
	}
	if input.Range != "" {
		util.SetHeader("Range", input.Range)
	}
	Response, requst := util.DoExec()
	if requst.Err != nil {
		return requst, nil
	}
	getObjectOutput = &models.GetObjectOutput{}
	getObjectOutput.DeleteMarker = Response.Header.Get("X-Amz-Delete-Marker")
	getObjectOutput.Expiration = Response.Header.Get("X-Amz-Expiration")
	getObjectOutput.LastModified = Response.Header.Get("Last-Modified")
	getObjectOutput.WebsiteRedirectLocation = Response.Header.Get("x-amz-website-redirect-location")
	getObjectOutput.ContentLength, _ = strconv.ParseInt(Response.Header.Get("Content-Length"), 10, 64)
	getObjectOutput.VersionId = Response.Header.Get("X-Amz-Version-Id")
	getObjectOutput.ETag = Response.Header.Get("Etag")
	if input.SaveAsFilePath != "" && requst.StatusCode < 300 {
		var writeLen int64 = 0
		//file, err := os.Create(input.SaveAsFilePath)
		file, err := os.OpenFile(input.SaveAsFilePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC|os.O_APPEND, 0644)
		if err != nil {
			requst.Err = err
			return requst, getObjectOutput
		}
		defer file.Close()
		buf := make([]byte, 1048576)
		for {
			n, _ := Response.Body.Read(buf)
			if n == 0 || writeLen >= getObjectOutput.ContentLength {
				break
			}
			writeLen += int64(n)
			file.Write(buf[:n])
		}
	}
	util.Close()
	return requst, getObjectOutput
}

/**
*函数原型：func (client *Client) DeleteObject(input *models.DeleteObjectInput) (requst *models.Result, deleteObjectOutput *models.DeleteObjectOutput)
*函数功能：删除对象
*参数说明：input：DeleteObjectInput对象实例
*返回值：requst: Resul对象实例
*		deleteObjectOutput：DeleteObjectOutput对象实例
 */
func (client *Client) DeleteObject(input *models.DeleteObjectInput) (requst *models.Result, deleteObjectOutput *models.DeleteObjectOutput) {
	requst = &models.Result{}
	if input.Bucket == "" {
		requst.Err = errors.New(models.BUCKET_NIL)
		return requst, nil
	}
	if input.Object == "" {
		requst.Err = errors.New(models.OBJECT_NIL)
		return requst, nil
	}
	util := utils.NewUtil(client.AK, client.SK, client.Endpoint, client.PathStyle)
	if input.VersionId != "" {
		util.SetPath("versionId", input.VersionId)
	}
	err := util.InitConect("DELETE", input.Bucket, input.Object, nil)
	if err != nil {
		requst.Err = err
		return requst, nil
	}
	Response, requst := util.DoExec()
	if requst.Err != nil {
		return requst, nil
	}
	util.Close()
	deleteObjectOutput = &models.DeleteObjectOutput{}
	deleteObjectOutput.DeleteMarker = Response.Header.Get("X-Amz-Delete-Marker")
	deleteObjectOutput.VersionId = Response.Header.Get("X-Amz-Version-Id")

	return requst, deleteObjectOutput
}

/**
*函数原型：func (client *Client) DeleteObjects(input *models.DeleteObjectsInput) (requst *models.Result, deleteObjectsOutput *models.DeleteObjectsOutput)
*函数功能：删除对象列表
*参数说明：input：DeleteObjectsInput对象实例
*返回值：requst: Resul对象实例
*		deleteObjectsOutput：DeleteObjectsOutput对象实例
 */
func (client *Client) DeleteObjects(input *models.DeleteObjectsInput) (requst *models.Result, deleteObjectsOutput *models.DeleteObjectsOutput) {
	requst = &models.Result{}
	if input.Bucket == "" {
		requst.Err = errors.New(models.BUCKET_NIL)
		return requst, nil
	}
	util := utils.NewUtil(client.AK, client.SK, client.Endpoint, client.PathStyle)
	xml := utils.DeleteObjectsXML(input)
	ioRead := strings.NewReader(xml)

	util.SetPath("delete", "")
	err := util.InitConect("POST", input.Bucket, "", ioRead)
	if err != nil {
		requst.Err = err
		return requst, nil
	}
	util.SetHeader("Content-Length", strconv.Itoa(len(xml)))
	util.SetHeader("Content-MD5", utils.HashMD5([]byte(xml)))
	Response, requst := util.DoExec()
	if requst.Err != nil {
		return requst, nil
	}

	deleteObjectsOutput = &models.DeleteObjectsOutput{}
	RspXML, err := ioutil.ReadAll(Response.Body)
	util.Close()
	if err != nil {
		requst.Err = err
		return requst, nil
	}
	err = utils.PareseXML(RspXML, deleteObjectsOutput)
	if err != nil {
		requst.Err = err
		return requst, nil
	}
	return requst, deleteObjectsOutput
}

/**
*函数原型：func (client *Client) GetObjectMetadata(input *models.GetObjectMetadataInput) (requst *models.Result, getObjectMetadataOutput *models.GetObjectMetadataOutput)
*函数功能：获取对象元数据
*参数说明：input：GetObjectMetadataInput对象实例
*返回值：requst: Resul对象实例
*		getObjectMetadataOutput：GetObjectMetadataOutput对象实例
 */
func (client *Client) GetObjectMetadata(input *models.GetObjectMetadataInput) (requst *models.Result, getObjectMetadataOutput *models.GetObjectMetadataOutput) {
	requst = &models.Result{}
	if input.Bucket == "" {
		requst.Err = errors.New(models.BUCKET_NIL)
		return requst, nil
	}
	if input.Object == "" {
		requst.Err = errors.New(models.OBJECT_NIL)
		return requst, nil
	}
	util := utils.NewUtil(client.AK, client.SK, client.Endpoint, client.PathStyle)
	if input.VersionId != "" {
		util.SetPath("versionId", input.VersionId)
	}

	err := util.InitConect("HEAD", input.Bucket, input.Object, nil)
	if err != nil {
		requst.Err = err
		return requst, nil
	}
	Response, requst := util.DoExec()
	if requst.Err != nil {
		return requst, nil
	}
	util.Close()
	getObjectMetadataOutput = &models.GetObjectMetadataOutput{}
	getObjectMetadataOutput.Expiration = Response.Header.Get("X-Amz-Expiration")
	getObjectMetadataOutput.LastModified = Response.Header.Get("Last-Modified")
	getObjectMetadataOutput.ETag = Response.Header.Get("Etag")
	getObjectMetadataOutput.VersionId = Response.Header.Get("X-Amz-Version-Id")
	getObjectMetadataOutput.WebsiteRedirectLocation = Response.Header.Get("X-Amz-Website-Redirect-Location")
	return requst, getObjectMetadataOutput
}

/**
*函数原型：func (client *Client) IsObjectExist(bucket, object string) (requst *models.Result, isExist bool)
*函数功能：查询对象是否存在
*参数说明：bucket：桶名
*		 object：对象名
*返回值：requst: Resul对象实例
*		isExist：对象是否存在，是：true，否：false
 */
func (client *Client) IsObjectExist(bucket, object string) (requst *models.Result, isExist bool) {
	requst = &models.Result{}
	isExist = false
	if bucket == "" {
		requst.Err = errors.New(models.BUCKET_NIL)
		return requst, isExist
	}
	if object == "" {
		requst.Err = errors.New(models.OBJECT_NIL)
		return requst, isExist
	}
	input := &models.ListObjectsInput{Bucket: bucket}
	requst, output := client.ListObjects(input)
	if requst.Err != nil {
		return requst, isExist
	}
	for _, val := range output.Contents {
		if val.Key == object {
			isExist = true
			break
		}
	}
	return requst, isExist
}
