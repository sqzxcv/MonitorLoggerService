package modeltools

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"

	"github.com/sqzxcv/glog"
	model "github.com/sqzxcv/monitorloggerservice/model/gen"
	"net/http"
	"time"
)

func ByteToProto(buf []byte, pb proto.Message) (err *model.Error) {

	err2 := proto.Unmarshal(buf, pb)
	if err2 != nil {
		glog.Error("解析pb失败" + err2.Error())
		err = &model.Error{
			Code: int32(model.ErrNum_PbDecode),
			Msg:  "解析pb失败" + err2.Error(),
		}
	}
	return err
}

func FetchBodyFromReqBase(base *model.ReqBase, pb proto.Message) (err *model.Error) {

	return ByteToProto(base.PbBody, pb)
}

func FetchBodyFromGinContext(c *gin.Context, pb proto.Message) (base1 *model.ReqBase, err *model.Error) {

	base, exist := c.Get("reqbase")
	if exist == false {
		glog.Error("当前请求环境不存在 reqbase")
		err = &model.Error{
			Code: int32(model.ErrNum_PbDecode),
			Msg:  "当前请求环境不存在 reqbase",
		}
		return nil, err

	}
	return base.(*model.ReqBase), FetchBodyFromReqBase(base.(*model.ReqBase), pb)
}

func ResponseFolderInfoResult(c *gin.Context, result interface{}, pberr *model.Error) {

	var msg string
	var retbase *model.RetBase
	if pberr == nil || pberr.Code == int32(model.ErrNum_OK) {

		data, err2 := proto.Marshal(result.(proto.Message))
		if err2 != nil {
			msg = "service err:" + err2.Error()
			retbase = &model.RetBase{
				MsgCtx:     msg,
				Code:       int32(model.ErrNum_PbEncode),
				ServerTime: time.Now().Unix(),
			}
		} else {

			// 客户端 的protobuf对msg的大小限制(小于64M=67108864)
			if len(data) < 67108800 {
				retbase = &model.RetBase{
					Code:       int32(model.ErrNum_OK),
					PbBody:     data,
					ServerTime: time.Now().Unix(),
				}
			} else {
				glog.Warn("Too many files, parsing failed:", len(data))
				retbase = &model.RetBase{
					MsgCtx:     "Too many files, parsing failed",
					Code:       int32(model.ErrNum_PbEncode),
					ServerTime: time.Now().Unix(),
				}
			}
		}
	} else {
		retbase = &model.RetBase{
			MsgCtx:     pberr.Msg,
			Code:       pberr.Code,
			ServerTime: time.Now().Unix(),
		}
	}

	c.ProtoBuf(http.StatusOK, retbase)
}

func ResponseResult(c *gin.Context, result interface{}, pberr *model.Error) {

	var msg string
	var retbase *model.RetBase
	if pberr == nil || pberr.Code == int32(model.ErrNum_OK) {

		data, err2 := proto.Marshal(result.(proto.Message))
		if err2 != nil {
			msg = "service err:" + err2.Error()
			retbase = &model.RetBase{
				MsgCtx:     msg,
				Code:       int32(model.ErrNum_PbEncode),
				ServerTime: time.Now().Unix(),
			}
		} else {

			// flutter 的protobuf对msg的大小限制(小于64M=67108864)
			if len(data) < 67108800 {
				retbase = &model.RetBase{
					Code:       int32(model.ErrNum_OK),
					PbBody:     data,
					ServerTime: time.Now().Unix(),
				}
			} else {
				retbase = &model.RetBase{
					MsgCtx:     "data too big",
					Code:       int32(model.ErrNum_PbEncode),
					ServerTime: time.Now().Unix(),
				}
			}
		}
	} else {
		retbase = &model.RetBase{
			MsgCtx:     pberr.Msg,
			Code:       pberr.Code,
			ServerTime: time.Now().Unix(),
		}
	}

	c.ProtoBuf(http.StatusOK, retbase)
}
