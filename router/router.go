package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sqzxcv/glog"
	"github.com/sqzxcv/monitorloggerservice/global"
	model "github.com/sqzxcv/monitorloggerservice/model/gen"
	"github.com/sqzxcv/monitorloggerservice/websocket"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"net/http"
	"strconv"
)

//SetupHTTPRouter 设置http路由
func SetupHTTPRouter(r *gin.Engine) *gin.Engine {

    // r.UseRawPath = true

    api := r.Group("/api", ProcessRequest())
    {
        setupApi(api)
    }

    r.GET("/version", APIVersion)

    r.GET("/ws", websocket.WsPage)
    return r
}

// APIVersion apiversion
func APIVersion(c *gin.Context) {

    c.JSON(http.StatusOK, gin.H{"version": global.VersionCode})
}

func ProcessRequest() gin.HandlerFunc {
    return func(c *gin.Context) {
        var reqbase model.ReqBase
        data, err1 := ioutil.ReadAll(c.Request.Body)
        if err1 != nil {
            glog.Error("get request, read request body failed: " + err1.Error())
        }
        err := proto.Unmarshal(data, &reqbase)
        if err != nil {
            glog.Info("收到请求, 解析pb失败" + err.Error())
        }

        glog.Info("收到请求, cmd:", "0x"+strconv.FormatInt(int64(reqbase.Cmd), 16))
        c.Set("reqbase", &reqbase)
        c.Next()
    }
}
