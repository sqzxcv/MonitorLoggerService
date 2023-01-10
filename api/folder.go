package api

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/sqzxcv/glog"
    model "github.com/sqzxcv/monitorloggerservice/model/gen"
    "github.com/sqzxcv/monitorloggerservice/model/modeltools"
    "io/ioutil"
    "net/http"
    "os"
)

func GetFolderInfo(c *gin.Context) {

    var req model.ReqFolderInfo
    _, err := modeltools.FetchBodyFromGinContext(c, &req)
    if err != nil {
        //return result, err
        modeltools.ResponseFolderInfoResult(c, nil, err)
    }

    currInfo, err2 := os.Stat(req.Path)
    if err2 != nil {
        glog.Error("get folder info failed:", err2.Error())
        err = &model.Error{
            Code: int32(model.ErrNum_ParamErrCode),
            Msg:  err2.Error(),
        }
        modeltools.ResponseFolderInfoResult(c, nil, err)
    }
    if currInfo.IsDir() == false {
        glog.Error("get folder info failed, the folder does not exist:", req.Path)
        err = &model.Error{
            Code: int32(model.ErrNum_ParamErrCode),
            Msg:  "get folder info failed, the folder does not exist",
        }
        modeltools.ResponseFolderInfoResult(c, nil, err)
    }

    children, totalSize, err := getFolderInfo(req.Path)
    if err != nil {
        modeltools.ResponseFolderInfoResult(c, nil, err)
    }
    fileInfo := &model.FileInfo{
        Name:         currInfo.Name(),
        ModifyTime:   currInfo.ModTime().Unix(),
        Size:         totalSize,
        Children:     children,
        IsDir:        true,
        AbsolutePath: req.Path,
        Formatter:    nil,
        Exist:        false,
    }

    modeltools.ResponseFolderInfoResult(c, fileInfo, nil)

}

func getFolderInfo(path string) (children []*model.FileInfo, totalSize int64, err *model.Error) {
    files, err1 := ioutil.ReadDir(path)
    if err1 != nil {
        glog.Error("get folder info failed:", err1.Error())
        err = &model.Error{
            Code: int32(model.ErrNum_ParamErrCode),
            Msg:  err1.Error(),
        }
        return nil, 0, err
    }

    for _, file := range files {

        absPath := fmt.Sprintf("%s/%s", path, file.Name())
        var cchildren []*model.FileInfo
        var currentSize = int64(0)
        if file.IsDir() {

            cchildren, currentSize, err = getFolderInfo(absPath)
            if err != nil {
                return nil, totalSize, err
            }

        } else if file.Mode().IsRegular() == false {
            // 忽略特殊文件
            continue

        } else {
            currentSize = file.Size()
            cchildren = nil
        }
        totalSize += currentSize

        info := &model.FileInfo{
            Name:         file.Name(),
            ModifyTime:   file.ModTime().Unix(),
            Size:         currentSize,
            Children:     cchildren,
            IsDir:        file.IsDir(),
            AbsolutePath: absPath,
            Formatter:    nil,
            Exist:        false,
        }
        children = append(children, info)
    }
    return children, totalSize, nil
}

func DownloadFileService(c *gin.Context) {
    fileDir := c.Query("path")
    //打开文件
    currInfo, err2 := os.Stat(fileDir)
    if err2 != nil {
        glog.Error("download file failed with err:", err2.Error())
        c.Redirect(http.StatusNotFound, "/404")
        return
    }
    if currInfo.IsDir() {
        c.Redirect(http.StatusFound, "/404")
    }
    c.Header("Content-Type", "application/octet-stream")
    c.Header("Content-Disposition", "attachment; filename="+currInfo.Name())
    c.Header("Content-Transfer-Encoding", "binary")
    c.File(fileDir)
    return
}
