package handler

import (
	pb "IHome/DeleteSession/proto"
	"IHome/IhomeWeb/utils"
	"context"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/gomodule/redigo/redis"
)

type DeleteSession struct{}

var (
	redis_conf = map[string]string{
		"key": utils.G_server_name,
		// 127.0.0.1:6379
		"conn":  utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum": utils.G_redis_dbnum,
	}
)

func (e *DeleteSession) DeleteSession(ctx context.Context, request *pb.DeleteSessionRequest, response *pb.DeleteSessionResponse) error {
	beego.Info("删除session DeleteSession api/v1.0/session")
	// 初始化错误码
	response.Errno = utils.RECODE_OK
	response.ErrMsg = utils.RecodeText(utils.RECODE_OK)
	// prepare to connect to redis
	beego.Info("准备连接redis数据库")
	redis_conf_json, _ := json.Marshal(redis_conf)
	bm, err := cache.NewCache("redis", string(redis_conf_json))
	if err != nil {
		beego.Info("Failed to connect to redis")
		response.Errno = utils.RECODE_DBERR
		response.ErrMsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}

	// delete data from redis using sessionid+name sessionid+mobile   sessionid+user_id
	bm.Delete(request.SessionId + "name")
	bm.Delete(request.SessionId + "mobile")
	bm.Delete(request.SessionId + "user_id")

	return nil
}
