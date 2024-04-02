package handler

import (
	models "IHome/IhomeWeb/model"
	"IHome/IhomeWeb/utils"
	pb "IHome/PostUserAuth/proto"
	"context"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	_ "github.com/gomodule/redigo/redis"
	"strconv"
	"time"
)

type PostUserAuth struct{}

var (
	redis_conf = map[string]string{
		"key": utils.G_server_name,
		// 127.0.0.1:6379
		"conn":  utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum": utils.G_redis_dbnum,
	}
)

func (e *PostUserAuth) PostUserAuth(_ context.Context, req *pb.PostUserAuthRequest, rsp *pb.PostUserAuthResponse) error {
	beego.Info("PostUserAuth 实名认证  /api/v1.0/user/auth")
	// 初始化返回值
	rsp.Errno = utils.RECODE_OK
	rsp.ErrMsg = utils.RecodeText(utils.RECODE_OK)
	session_id := req.SessionId

	redis_conf_json, _ := json.Marshal(redis_conf)

	bm, err := cache.NewCache("redis", string(redis_conf_json))
	if err != nil {
		beego.Info("Failed to connect to redis")
		rsp.Errno = utils.RECODE_DBERR
		rsp.ErrMsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}

	user_id_raw := bm.Get(session_id + "user_id")

	user_id_string, _ := redis.String(user_id_raw, nil)
	id, _ := strconv.Atoi(user_id_string)

	user := models.User{Id: id, Id_card: req.IdCard, Real_name: req.RealName}

	o := orm.NewOrm()

	_, err = o.Update(&user, "real_name", "id_card")
	if err != nil {
		beego.Info("更新数据库失败")
		rsp.Errno = utils.RECODE_DBERR
		rsp.ErrMsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}

	_ = bm.Put(session_id+"user_id", user_id_string, time.Second*600)

	return nil
}
