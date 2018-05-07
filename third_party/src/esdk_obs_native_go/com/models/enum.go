package models

//enum x-amz-acl
const (
	PRIVATE                   = "private"
	PUBLIC_READ               = "public-read"
	PUBLIC_READ_WRITE         = "public-read-write"
	AUTHENTICATED_READ        = "authenticated-read"
	BUCKET_OWNER_READ         = "bucket-owner-read"
	BUCKET_OWNER_FULL_CONTROL = "bucket-owner-full-control"
	LOG_DELIVERY_WRITE        = "log-delivery-write"
)

//自定义错误信息
const (
	BUCKET_NIL      = "the input bucket name is nil,please enter the correct value!"
	OBJECT_NIL      = "the input object name is nil,please enter the correct value!"
	BODY_FILE_EXIST = "the input body and sourcefile exist at same time,please specify one of eigther a string or file to be send!"
)
