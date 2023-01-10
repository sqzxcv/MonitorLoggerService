package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sqzxcv/monitorloggerservice/api"
)

func setupApi(route *gin.RouterGroup) {

	// 获取文件夹信息
	route.POST("folder", api.GetFolderInfo)
	// 下载文件
	route.GET("downloadfile", api.DownloadFileService)
}
