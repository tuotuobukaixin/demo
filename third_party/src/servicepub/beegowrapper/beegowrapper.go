// const
package beegowrapper
import (
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/context"
)

const (
    HTTP_CONENT_TYPE_NULL = ""
    HTTP_BODY_NULL = "" 
    HTTP_CONTENT_TYPE_JSON = "application/json;charset=UTF-8"
)

func BeegoSetFailOutput(ctx *context.Context, statuscode int) {
    BeegoSetOutput(ctx, statuscode, HTTP_CONENT_TYPE_NULL, HTTP_BODY_NULL)
}

func BeegoSetFailOutputWithJsonBody(ctx *context.Context, statuscode int, body string) {
    BeegoSetOutput(ctx, statuscode, HTTP_CONTENT_TYPE_JSON, body)
}

func BeegoSetSucceedOutput(ctx *context.Context, statuscode int, body string) {
    content_type := ""
    if body != "" {
        content_type = HTTP_CONTENT_TYPE_JSON
    }
    BeegoSetOutput(ctx, statuscode, content_type, body)
}

func BeegoEnforceHeadersPolicy(out *context.BeegoOutput) {  
    //"Cache-Control: no-store"
    //"Pragma: no-cache"
    //"Cache-Control: no-cache"
    out.Header("Strict-Transport-Security",
        "max-age=31536000; includeSubDomains")
    out.Header("Cache-control", "no-cache, no-store")
    out.Header("Pragma", "no-cache")
}

func BeegoSetOutput(ctx *context.Context, statuscode int, content_type string, body string) {
    if content_type != "" {
        ctx.Output.Header("Content-Type", content_type)
    }
    //no-cache
    BeegoEnforceHeadersPolicy(ctx.Output)
    
    ctx.Output.SetStatus(statuscode)
    ctx.Output.Body([]byte(body))

    beego.Debug("write to resp: statuscode:", statuscode)
}