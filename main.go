package main

import (
    "flag"
    "github.com/sqzxcv/golang_public_util/exception"
    "github.com/sqzxcv/monitorloggerservice/global"
    "github.com/sqzxcv/monitorloggerservice/httpsrv"
    "os"
    "path/filepath"
    "strings"

    "github.com/sqzxcv/glog"
)

func main() {

    defer exception.CatchException("主线程出现异常")
    logsDir := getCurrentDirectory()
    glog.SetConsole(true)
    glog.Info("日志目录:", logsDir)
    glog.SetRollingDaily(logsDir, "ss_backend.log", false)
    glog.Info(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> App 启动成功 <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
    //var configPath = flag.String("c", "", "配置文件地址")
    var port = flag.Int("port", 9080, "服务监听端口")
    var versionFlag = flag.Bool("v", false, "输出版本号")
    flag.Parse()
    if *versionFlag {
        println(global.VersionCode)
        return
    }
    glog.Info("当前版本:", global.VersionCode)

    restartService(*port)
    done := make(chan bool)
    <-done
    glog.Info("程序退出")
}

func restartService(port int) {
    go func() {
        httpsrv.SetupHTTPSrv(port)
    }()
}

func getCurrentDirectory() string {
    dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
    if err != nil {
        glog.Error(err)
    }
    return strings.Replace(dir, "\\", "/", -1) //将\替换成/
}
