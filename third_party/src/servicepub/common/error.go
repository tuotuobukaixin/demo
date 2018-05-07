package common
import (
    "fmt"
    "net/http"
)
type Error interface {
    Error()
}

type error_impl struct {
    Msg   string
    Errno int
}

//not used actually, because err is nil when succeed
var ESucceed = New_Error("Succeed", 0)

// this is very general error, should be used carefully, because it does not make sense actually
var EFailure = New_Error("Failure", 1)

//var EOngoing = New_Error("Ongoing", 10)
var EAgain = New_Error("Try again later", 100)
var ENotFound = New_Error("Not found", 101)
var EHttpError = New_Error("Http internal error", 102)
var EEmptyString = New_Error("string is empty", 103)
var EObjectIsNull = New_Error("Object is null", 104)
var EServiceNotReady = New_Error("service is not ready", 105)
var EFileNotExist = New_Error("file is not existed", 106)
var EAssert = New_Error("Assert not satisfied", 107)
var EParseFail = New_Error("Parse fail", 108)
var EReqInvaild = New_Error("Request invalid", 109)
var EObjectExist = New_Error("Object already exist", 110)
var EJsonParseFail = New_Error("Json parse fail", 111)
var EInvalidParam = New_Error("Invalid Param", 112)
var ENoMoreData = New_Error("No more data", 113)
var ENotAllowed = New_Error("Not Allowed", 114)
var ENotImplemented = New_Error("Not Implemented", 115)
var EOperationTimeout = New_Error("Operation timeout", 116)
var EAlreadyExist = New_Error("Object already exist", 117)
var EAlreadyBound = New_Error("Object already bound", 118)
var EObjectInUse = New_Error("Object is in use", 119)
var EInvalidBillingData = New_Error("Invalid billing data", 120)
var EHasGenBilling = New_Error("has been gen billing", 121)
var EReachUpLimit = New_Error("Object has reached the upper limit", 122)
var ENotExist = New_Error("Object is not exist", 123)
var ENotAuth = New_Error("Authorization failed", 401)
var ELocked = New_Error("the user has been locked", 124)

var ENeedNoBill = New_Error("Need no bill", 200)
var ENotBillingTime = New_Error("Not billing time", 201)

var ECreateDeployFailed = New_Error("Creat deployment failed", 300)
var EDeployFailed = New_Error("Deploy failed", 301)
var EDeleteDeployFailed = New_Error("Delete deploy failed", 302)

var ETooLong = New_Error("Too long", 400)

var EBindingConflict = New_Error("Binding is conflict", 500)
var EBindingGone = New_Error("Binding is gone", 501)
var EInstanceGone = New_Error("Service instance is gone", 502)

func New_Error(msg string, errno int) error {
    errimpl := new(error_impl)
    errimpl.Msg = msg
    errimpl.Errno = errno

    return errimpl
}

func (this *error_impl) Error() string {
    errinfo := fmt.Sprintf("%s. errno:%d", this.Msg, this.Errno)
    return errinfo
}


func Error_to_statuscode(err error) int {
    statuscode := 0
    
    if err == ESucceed {
        statuscode = http.StatusOK
    } else if err == ENotFound {
        statuscode = http.StatusNotFound
    } else if err == ENotAllowed || err == ELocked {
        statuscode = http.StatusForbidden
    } else if err == EBindingConflict {
        statuscode = http.StatusConflict
    } else if err == EBindingGone || err == EInstanceGone {
        statuscode = http.StatusGone
    } else if err == EEmptyString  || err == EReqInvaild || err == EObjectExist || 
         err == EJsonParseFail || err == EInvalidParam || err == EParseFail ||
         err == EAlreadyExist  || err == EAlreadyBound ||  err == ETooLong {
        statuscode = http.StatusBadRequest
    } else if err == ENotImplemented {
        statuscode = http.StatusNotImplemented
    } else if err == EObjectIsNull || err == EOperationTimeout ||   
        err == ECreateDeployFailed || err == EDeployFailed || err == EDeleteDeployFailed ||
        err == EServiceNotReady || err == EHttpError || err == EObjectInUse {
        statuscode = http.StatusInternalServerError
    } else if err == EFailure{ 
        statuscode = http.StatusInternalServerError
    } else{
        statuscode = http.StatusInternalServerError
    }
    
    return statuscode
}

