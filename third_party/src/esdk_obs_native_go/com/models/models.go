package models

//接口执行信息集合
type Result struct {
	Err        error  //执行的错误信息，当Err为nil是，其他值有效
	StatusCode int    //http 返回的状态码
	Code       string //访问服务的错误码
	Message    string //访问服务的错误描述
	RequestId  string //本次错误请求的请求ID
	HostId     string //返回该消息的服务端ID
}

//创建桶入参信息集合
type CreateBucketInput struct {
	Bucket   string //桶名，必选
	ACL      string //x-amz-acl 对应值，可选
	Location string //服务器所在区域，可选
}

//上传对象入参信息集合
type PutObjectInput struct {
	Bucket                  string            //桶名，必选
	Object                  string            //对象名，必选
	ACL                     string            //x-amz-acl 对应值，可选
	Metadata                map[string]string //自定义的元数据，可选
	WebsiteRedirectLocation string            //当桶设置了Website配置，可以将获取这个对象的请求重定向到桶内另一个对象或一个外部的URL，可选
	Body                    string            //待上传的字符串，和SourceFile不能同时存在，可选
	SourceFile              string            //待上传的文件，和Body不能同时存在，可选
}

//上传对象返回信息集合
type PutObjectOutPut struct {
	VersionId string //对象的版本号
	ETag      string //对象的Etag值
}

//获取桶内对象列表入参集合信息
type ListObjectsInput struct {
	Bucket    string //桶名
	Prefix    string //列举以指定的字符串prefix开头的对象,可选
	Marker    string //返回的对象列表将是按照字典顺序排序后这个标识符以后的所有对象，可选
	MaxKeys   int    //指定返回的最大对象数，范围是[1，1000]，可选
	Delimiter string //用来分组桶内对象的字符串，可选
}

//获取桶内对象列表返回值集合信息
type ListObjectsOutput struct {
	Delimiter      string          `xml:"Delimiter"`      //请求中携带的delimiter参数
	IsTruncated    string          `xml:"IsTruncated"`    //表明是否本次返回的ListBucketResult结果列表被截断
	Marker         string          `xml:"Marker"`         //列举对象时的起始位置
	NextMarker     string          `xml:"NextMarker"`     //标明本次请求列举到的最后一个对象
	MaxKeys        int             `xml:"MaxKeys"`        //列举时最多返回的对象个数
	Name           string          `xml:"Name"`           //本次请求的桶名
	Prefix         string          `xml:"Prefix"`         //本次请求只列举对象名能匹配该前缀的所有对象
	Contents       []Content       `xml:"Contents"`       //对象的元数据信息
	CommonPrefixes []CommonPrefixe `xml:"CommonPrefixes"` //请求中带delimiter参数时，返回消息带CommonPrefixes分组信息
}

//获取对象内容入参信息集合
type GetObjectInput struct {
	Bucket                     string //桶名，必选
	Object                     string //对象名，必选
	IfMatch                    string //如果对象的ETag和请求中指定的ETag相同，则返回对象内容，可选
	IfNoneMatch                string //如果对象的ETag和请求中指定的ETag不相同，则返回对象内容，可选
	IfModifiedSince            string //如果对象在请求中指定的时间之后有修改，则返回对象内容，值为HTTP时间字符串，可选
	IfUnmodifiedSince          string //如果对象在请求中指定的时间之后没有修改，则返回对象内容，值为HTTP时间字符串，可选
	Range                      string //获取对象时获取在Range范围内的对象内容，格式如（bytes=0-4），可选
	ResponseCacheControl       string //重写响应中的Cache-Control头，可选
	ResponseContentDisposition string //重写响应中的Content-Disposition头，可选
	ResponseContentEncoding    string //重写响应中的Content-Encoding头，可选
	ResponseContentLanguage    string //重写响应中的Content-Language头，可选
	ResponseContentType        string //重写响应中的Content-Type头，可选
	ResponseExpires            string //重写响应中的Expires头，可选
	VersionId                  string //指定获取对象的版本号，可选
	SaveAsFilePath             string //下载对象到本地的路径，必选
}

//获取对象内容返回值信息集合
type GetObjectOutput struct {
	DeleteMarker            string //标识对象是否是删除标记
	Expiration              string //如果对象配置了过期时间，那么响应的消息中应该包含此消息头
	LastModified            string //对象的最近一次修改时间
	ETag                    string //对象的ETag值
	WebsiteRedirectLocation string //请求重定向到该属性指定的桶内的另一个对象或外部的URL
	VersionId               string //对象版本号
	ContentLength           int64  //对象的字节数
}

//删除对象入参信息集合
type DeleteObjectInput struct {
	Bucket    string //桶名，必选
	Object    string //对象名，必选
	VersionId string //对象版本号，可选
}

//删除对象返回值信息集合
type DeleteObjectOutput struct {
	DeleteMarker string //标识对象是否标记删除
	VersionId    string //对象版本号
}

//批量删除对象入参信息集合
type DeleteObjectsInput struct {
	Bucket  string   //桶名，必选
	Quiet   bool     //用于指定使用quiet模式，只返回删除失败的对象结果，可选
	Objects []Object //待删除的对象列表，必选
}

//批量删除对象返回信息集合
type DeleteObjectsOutput struct {
	Deleteds []Deleted `xml:"Deleted"` //删除成功结果列表
	Errors   []Error   `xml:"Error"`   //删除失败结果的列表
}

//获取对象元数据入参信息集合
type GetObjectMetadataInput struct {
	Bucket    string //桶名，必选
	Object    string //对象名，必选
	VersionId string //对象版本号，可选
}

//获取对象元数据返回值信息集合
type GetObjectMetadataOutput struct {
	Expiration              string //对象的过期时间
	LastModified            string //对象的最后修改时间
	ETag                    string //对象的ETag值
	VersionId               string //对象版本号
	WebsiteRedirectLocation string //请求重定向到该属性指定的桶内的另一个对象或外部的URL
}

//对象信息
type Object struct {
	Object    string //对象名，必选
	VersionId string //对象版本号，可选
}

//CommonPrefixe信息
type CommonPrefixe struct {
	Prefix string `xml:"Prefix"` //对象名的前缀，表示本次请求只列举对象名能匹配该前缀的所有对象
}

//对象的元数据信息
type Content struct {
	Owner        Owner  `xml:"Owner"`        //用户信息
	ETag         string `xml:"ETag"`         //对象的MD5值
	Key          string `xml:"Key"`          //对象名
	LastModified string `xml:"LastModified"` //对象最近一次被修改的时间
	Size         int    `xml:"Size"`         //对象的字节数
	StorageClass string `xml:"StorageClass"` //对象的存储类型
}

//列举桶列表返回信息
type ListBucketsOutput struct {
	Owner   Owner    `xml:"Owner"`          //用户信息
	Buckets []Bucket `xml:"Buckets>Bucket"` //桶信息列表
}

//获取桶配额返回信息
type GetBucketQuotaOutput struct {
	Quota int `xml:"StorageQuota"` //桶的配额值
}

//获取桶的存储信息返回值信息
type GetBucketStorageInfoOutput struct {
	Size         int `xml:"Size"`         //存储大小
	ObjectNumber int `xml:"ObjectNumber"` //桶内的对对象数
}

//用户信息
type Owner struct {
	ID          string `xml:"ID"`          //用户ID
	DisplayName string `xml:"DisplayName"` //用户名
}

//桶信息
type Bucket struct {
	Name         string `xml:"Name"`         //桶名
	CreationDate string `xml:"CreationDate"` //桶的创建时间
}

//删除对象信息
type Deleted struct {
	Key                   string `xml:"Key"`                   //对象名
	VersionId             string `xml:"VersionId"`             //对象版本号
	DeleteMarker          bool   `xml:"DeleteMarker"`          //当批量删除请求访问的桶是多版本桶时，如果创建或删除一个删除标记，返回消息中该元素的值为true
	DeleteMarkerVersionId string `xml:"DeleteMarkerVersionId"` //删除的删除标记版本号
}

//删除对象错误信息
type Error struct {
	Key       string `xml:"Key"`       //对象名
	Code      string `xml:"Code"`      //删除失败结果的错误码
	Message   string `xml:"Message"`   //删除失败结果的错误消息
	VersionId string `xml:"VersionId"` //对象版本号
}

//http 返回错误信息
type ErrResponse struct {
	Code      string `xml:"Code"`      //错误信息码
	Message   string `xml:"Message"`   //错误信息
	RequestId string `xml:"RequestId"` //本次错误请求的请求ID
	HostId    string `xml:HostId`      //返回该消息的服务端ID
}
