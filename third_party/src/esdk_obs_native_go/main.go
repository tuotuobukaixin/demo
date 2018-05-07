// esdk_obs_native_go project main.go
package main

import (
	"esdk_obs_native_go/com/client"
	"esdk_obs_native_go/com/models"
	"fmt"
)

func main() {

	//初始化服务类工厂
	obs := client.Factory("C44F0CE76614BE1A79B5", "XLZxUt6l+ygUcgn6h87kERS8n+wAAAFOZhS+Gptt", "https://129.7.182.2:443", true)

	//创建桶
	fmt.Println("Create bucket start...")
	input := &models.CreateBucketInput{Bucket: "123456-3", ACL: models.LOG_DELIVERY_WRITE, Location: "DC1"}
	requst := obs.CreateBucket(input)
	fmt.Printf("err:%s,statusCode:%d,code:%s,message:%s\n", requst.Err, requst.StatusCode, requst.Code, requst.Message)
	fmt.Println("Create bucket end")

	// Head桶
	fmt.Println("Head bucket start...")
	requst := obs.HeadBucket("123456-3")
	fmt.Printf("err:%s,statusCode:%d,code:%s,message:%s\n", requst.Err, requst.StatusCode, requst.Code, requst.Message)
	fmt.Println("Head bucket end")

	//查询桶列表
	fmt.Println("List bukect start...")
	requst, output := obs.ListBuckets()
	fmt.Printf("err:%s,statusCode:%d,code:%s,message:%s\n", requst.Err, requst.StatusCode, requst.Code, requst.Message)
	if output != nil {
		fmt.Printf("owner_name:%s,owner_id:%s\n", output.Owner.DisplayName, output.Owner.ID)
		for _, val := range output.Buckets {
			fmt.Printf("bucket:%s,create_date:%s\n", val.Name, val.CreationDate)
		}
	}
	fmt.Println("List bukect end")

	//删除桶
	fmt.Println("Delete bucket start...")
	requst := obs.DeleteBucket("123456-3")
	fmt.Printf("err:%s,statusCode:%d,code:%s,message:%s\n", requst.Err, requst.StatusCode, requst.Code, requst.Message)
	fmt.Println("Delete bucket end")

	//设置桶的配额信息
	fmt.Println("Set bucket quota start...")
	requst := obs.SetBucketQuota("123456-10", 10485760)
	fmt.Printf("err:%s,statusCode:%d,code:%s,message:%s\n", requst.Err, requst.StatusCode, requst.Code, requst.Message)
	fmt.Println("Set bucket quota end")

	//获取桶的配额信息
	fmt.Println("Get bucket quota start...")
	requst, output := obs.GetBucketQuota("123456-3")
	fmt.Printf("err:%s,statusCode:%d,code:%s,message:%s\n", requst.Err, requst.StatusCode, requst.Code, requst.Message)
	if output != nil {
		fmt.Printf("quota:%d\n", output.Quota)
	}
	fmt.Println("Get bucket quota end")

	//获取桶的存储信息
	fmt.Println("Get bucket storage info start...")
	requst, output := obs.GetBucketStorageInfo("123456-1")
	fmt.Printf("err:%s,statusCode:%d,code:%s,message:%s\n", requst.Err, requst.StatusCode, requst.Code, requst.Message)
	if output != nil {
		fmt.Printf("size:%d,object_number:%d\n", output.Size, output.ObjectNumber)
	}
	fmt.Println("Get bucket storage info end")

	//删除桶数据
	fmt.Println("Delete bucket with objects start...")
	requst := obs.DeleteBucketWithObjects("123456-1")
	fmt.Printf("err:%s,statusCode:%d,code:%s,message:%s\n", requst.Err, requst.StatusCode, requst.Code, requst.Message)
	fmt.Println("Delete bucket with objects end")

	//获取桶内对象
	fmt.Println("List objects start...")
	input := new(models.ListObjectsInput)
	input.Bucket = "123456-1"
	input.MaxKeys = 10
	input.Prefix = "T"
	requst, output := obs.ListObjects(input)
	fmt.Printf("err:%s,statusCode:%d,code:%s,message:%s\n", requst.Err, requst.StatusCode, requst.Code, requst.Message)
	if output != nil {
		fmt.Printf("Delimiter:%s,IsTruncated:%s,Marker%s,NextMarker:%s,MaxKeys:%d,Name:%s,Prefix:%s\n",
			output.Delimiter, output.IsTruncated, output.Marker, output.NextMarker, output.MaxKeys, output.Name, output.Prefix)
		for i, val := range output.Contents {
			fmt.Printf("conten[%d]:owner_id:%s,owner_name:%s,ETag:%s,Key:%s,LastModified:%s,Size:%s,StorageClass:%s\n",
				i, val.Owner.ID, val.Owner.DisplayName, val.ETag, val.Key, val.LastModified, val.Size, val.StorageClass)
		}
		for i, val := range output.CommonPrefixes {
			fmt.Printf("CommonPrefixe[%d]:Prefix:%s\n", i, val.Prefix)
		}
	}
	fmt.Println("List objects end")

	//上传对象
	fmt.Println("Put object start...")
	input := new(models.PutObjectInput)
	input.Bucket = "123456-10"
	input.Object = "BB.log"
	input.ACL = models.PUBLIC_READ_WRITE
	input.Metadata = make(map[string]string)
	input.Metadata["test1"] = "test-1"
	input.Metadata["test2"] = "test-2"
	input.Body = "123456789987654321"
	//input.SourceFile = "C:\\BB.rar"
	requst, output := obs.PutObject(input)
	fmt.Printf("err:%s,statusCode:%d,code:%s,message:%s\n", requst.Err, requst.StatusCode, requst.Code, requst.Message)
	if output != nil {
		fmt.Printf("ETag:%s,VersionId:%s\n", output.ETag, output.VersionId)
	}
	fmt.Println("Put object end")

	//下载对象
	fmt.Println("Get object start...")
	input := new(models.GetObjectInput)
	input.Bucket = "123456-1"
	input.Object = "Test.log"
	input.Range = "bytes=0-1024"
	input.SaveAsFilePath = "C:\\ooo.log"
	requst, output := obs.GetObject(input)
	fmt.Printf("err:%s,statusCode:%d,code:%s,message:%s\n", requst.Err, requst.StatusCode, requst.Code, requst.Message)
	if output != nil {
		fmt.Printf("DeleteMarker:%s,ETag:%s,Expiration:%s,LastModified:%s,VersionId:%s,WebsiteRedirectLocation:%s,ContentLength:%d\n",
			output.DeleteMarker, output.ETag, output.Expiration, output.LastModified, output.VersionId, output.WebsiteRedirectLocation, output.ContentLength)
	}
	fmt.Println("Get object end")

	//删除对象
	fmt.Println("Delete object start...")
	input := &models.DeleteObjectInput{Bucket: "123456-1", Object: "Test.log", VersionId: ""}
	requst, output := obs.DeleteObject(input)
	fmt.Printf("err:%s,statusCode:%d,code:%s,message:%s\n", requst.Err, requst.StatusCode, requst.Code, requst.Message)
	if output != nil {
		fmt.Printf("DeleteMarker:%s,VersionId:%s\n", output.DeleteMarker, output.VersionId)
	}
	fmt.Println("Delete object end")

	//批量删除对象
	fmt.Println("Delete objects start...")
	input := new(models.DeleteObjectsInput)
	input.Bucket = "123456-1"
	var key [2]models.Object
	key[0].Object = "test"
	key[1].Object = "users.dat"
	input.Objects = append(input.Objects, key[0], key[1])
	requst, output := obs.DeleteObjects(input)
	fmt.Printf("err:%s,statusCode:%d,code:%s,message:%s\n", requst.Err, requst.StatusCode, requst.Code, requst.Message)
	if output != nil {
		for i, val := range output.Deleteds {
			fmt.Printf("Deleteds[%d]:DeleteMarker:%s,DeleteMarkerVersionId:%s,Key:%s,VersionId:%s\n",
				i, val.DeleteMarker, val.DeleteMarkerVersionId, val.Key, val.VersionId)
		}
		for i, val := range output.Errors {
			fmt.Printf("Errors[%d]:Code:%s,Key:%s,Message:%s,VersionId:%s\n",
				i, val.Code, val.Key, val.Message, val.VersionId)
		}
	}
	fmt.Println("Delete objects end")

	//获取对象元数据
	fmt.Println("Get object metadata start...")
	input := &models.GetObjectMetadataInput{Bucket: "123456-1", Object: "AA.rar"}
	requst, output := obs.GetObjectMetadata(input)
	fmt.Printf("err:%s,statusCode:%d,code:%s,message:%s\n", requst.Err, requst.StatusCode, requst.Code, requst.Message)
	if output != nil {
		fmt.Printf("ETag:%s,Expiration:%s,LastModified:%s,VersionId:%s,WebsiteRedirectLocation:%s\n",
			output.ETag, output.Expiration, output.LastModified, output.VersionId, output.WebsiteRedirectLocation)
	}
	fmt.Println("Get object metadata end")

	//检测对象是否存在
	fmt.Println("Check object exist start...")
	requst, output := obs.IsObjectExist("123456-1", "AAA.rar")
	fmt.Printf("err:%s,statusCode:%d,code:%s,message:%s\n", requst.Err, requst.StatusCode, requst.Code, requst.Message)
	if requst.Err == nil {
		if output {
			fmt.Print("IsExist:true\n")
		} else {
			fmt.Print("IsExist:false\n")
		}
	}
	fmt.Println("Check object exist end")
}
