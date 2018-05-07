package utils

import (
	"esdk_obs_native_go/com/models"
	"strconv"
)

/**
*函数说明：构造创建桶的xml
*入参说明：input：CreateBucketInput对象实例
*返回值：构造的xm字符串
 */
func CreatBucketXML(input *models.CreateBucketInput) string {
	xml := ""
	if input.Location != "" {
		xml += "<CreateBucketConfiguration xmlns=\"http://s3.amazonaws.com/doc/2006-03-01/\">"
		xml += "<LocationConstraint>" + input.Location + "</LocationConstraint>"
		xml += "</CreateBucketConfiguration>"
	}
	return xml
}

/**
*函数说明：构造设置桶配额的xml
*入参说明：quota：桶配额数
*返回值：构造的xm字符串
 */
func SetBucketQuotaXML(quota int) string {
	xml := ""
	xml += "<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?>"
	xml += "<Quota xmlns=\"http://s3.amazonaws.com/doc/2006-03-01/\">"
	xml += "<StorageQuota>" + strconv.Itoa(quota) + "</StorageQuota>"
	xml += "</Quota>"
	return xml
}

/**
*函数说明：构造删除桶的数据的xml
*入参说明：bucket：桶名
*返回值：构造的xm字符串
 */
func DeleteBucketWithObjectsXML(bucket string) string {
	xml := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>"
	xml += "<DeleteBucket>"
	xml += "<Bucket>" + bucket + "</Bucket>"
	xml += "</DeleteBucket>"
	return xml
}

/**
*函数说明：构造批量删除对象的xml
*入参说明：inputt：DeleteObjectsInput对象实例
*返回值：构造的xm字符串
 */
func DeleteObjectsXML(input *models.DeleteObjectsInput) string {
	xml := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>"
	xml += "<Delete xmlns=\"http://s3.amazonaws.com/doc/2006-03-01/\">"
	if input.Quiet == true {
		xml += "<Quiet>" + "true" + "</Quiet> "
	}
	for _, key := range input.Objects {
		xml += "<Object>"
		xml += "<Key>" + key.Object + "</Key>"
		if key.VersionId != "" {
			xml += "<VersionId>" + key.VersionId + "</VersionId>"
		}
		xml += "</Object>"
	}
	xml += "</Delete>"
	return xml
}
