package httpsrv

import (
    "context"
    "fmt"
    "github.com/sqzxcv/monitorloggerservice/router"
    "net/http"
    "os"
    "time"

    "github.com/sqzxcv/glog"

    "github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
    // Disable Console Color
    // gin.DisableConsoleColor()
    r := gin.Default()
    // r.Use(cors.Default())
    r.Use(crosHandler())
    // r.POST("/om/", func(c *gin.Context) {
    // 	go func() {
    // 		routeraction.RouterHandle(c)
    // 	}()
    // })

    return router.SetupHTTPRouter(r)
}

//var quit chan os.Signal
var httpSrv *http.Server
var isRestart bool

// SetupHTTPSrv 启动http srv
func SetupHTTPSrv(httPort int) {

    isRestart = false
    if httpSrv != nil {
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        glog.Info("等待5秒 http srv 重启中")
        defer cancel()
        isRestart = true
        if err := httpSrv.Shutdown(ctx); err != nil {
            glog.Error("http srv Shutdown failed with error:", err.Error())
        }
    }
    //signal.Notify(quit, os.Interrupt)
    r := setupRouter()
    // Listen and Server in 0.0.0.0:8080
    add := fmt.Sprintf(":%d", httPort)
    httpSrv = &http.Server{
        Addr:    add,
        Handler: r,
    }
    //err := r.Run(add)

    glog.Info("HTTP listen port :", add)
    err := httpSrv.ListenAndServe()
    if err != nil && isRestart == false {
        glog.Error("start http service failed with error:", err.Error())
        os.Exit(0)
        return
    }
}

//跨域访问：cross  origin resource share
func crosHandler() gin.HandlerFunc {
    return func(context *gin.Context) {
        method := context.Request.Method
        context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        context.Header("Access-Control-Allow-Origin", "*") // 设置允许访问所有域
        context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
        context.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma,token,openid,opentoken")
        context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
        context.Header("Access-Control-Max-Age", "172800")
        context.Header("Access-Control-Allow-Credentials", "false")
        // context.Set("content-type", "application/json")  //设置返回格式是json

        // if method == "OPTIONS" {
        // 	context.JSON(http.StatusOK, result.Result{Code: result.OK, Data: "Options Request!"})
        // }

        //处理请求
        if method == "OPTIONS" {
            context.JSON(http.StatusOK, "Options Request!")
        }
        context.Next()
    }

}

func init() {
    //quit = make(chan os.Signal)
}
