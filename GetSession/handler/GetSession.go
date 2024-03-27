package handler

import (
	pb "IHome/GetSession/proto"
	"IHome/IhomeWeb/utils"
	"context"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/gomodule/redigo/redis"
	_ "github.com/gomodule/redigo/redis"
)

type GetSession struct{}

var (
	redis_conf = map[string]string{
		"key": utils.G_server_name,
		// 127.0.0.1:6379
		"conn":  utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum": utils.G_redis_dbnum,
	}
)

func (e *GetSession) GetSession(ctx context.Context, req *pb.GetSessionRequest, rsp *pb.GetSessionResponse) error {
	beego.Info("获取session信息GetSession url：/api/v1.0/session")
	// 初始化返回值
	rsp.Errno = utils.RECODE_OK
	rsp.ErrMsg = utils.RecodeText(utils.RECODE_OK)

	// 获取username
	// 准备连接redis
	redis_conf_json, err := json.Marshal(redis_conf)
	//创建redis句柄
	//
	bm, err := cache.NewCache("redis", string(redis_conf_json))
	if err != nil {
		beego.Info("redis连接失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.ErrMsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	// 注意不要忘了查找name的话要在session后面加上“name”标识符
	username := bm.Get(req.SessionId + "name")
	// 无则返回成功
	if username == nil {
		beego.Info("根据session从redis中获取username失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.ErrMsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	// 有则返回成功
	rsp.UserName, _ = redis.String(username, nil)

	return nil
}
