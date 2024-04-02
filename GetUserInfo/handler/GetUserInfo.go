package handler

import (
	pb "IHome/GetUserInfo/proto"
	models "IHome/IhomeWeb/model"
	"IHome/IhomeWeb/utils"
	"context"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	"reflect"
	"strconv"
)

type GetUserInfo struct{}

var (
	redis_conf = map[string]string{
		"key": utils.G_server_name,
		// 127.0.0.1:6379
		"conn":  utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum": utils.G_redis_dbnum,
	}
)

func (e *GetUserInfo) GetUserInfo(_ context.Context, req *pb.GetUserInfoRequest, rsp *pb.GetUserInfoResponse) error {
	beego.Info("获取用户信息 GetUserInfo /api/v1.0/user")

	// 初始化错误码
	rsp.Errno = utils.RECODE_OK
	rsp.ErrMsg = utils.RecodeText(utils.RECODE_OK)

	// connect to redis
	redis_conf_json, _ := json.Marshal(redis_conf)

	bm, err := cache.NewCache("redis", string(redis_conf_json))
	if err != nil {
		beego.Info("Failed to connect to redis")
		rsp.Errno = utils.RECODE_DBERR
		rsp.ErrMsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}
	user_id := bm.Get(req.SessionId + "user_id")
	if user_id == nil {
		beego.Info("根据session从redis中获取用户id失败")
		rsp.Errno = utils.RECODE_DBERR
		rsp.ErrMsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}
	beego.Info(reflect.TypeOf(user_id), user_id)
	// 这里为什么要user_id.([]uint8)呢？这是为了进行类型断言
	id, err := redis.String(user_id, nil)
	beego.Info(reflect.TypeOf(id), id)
	new_id, _ := strconv.Atoi(id)
	//创建user对象
	user := models.User{Id: new_id}

	o := orm.NewOrm()
	err = o.Read(&user)
	if err != nil {
		beego.Info("从mysql中获取数据失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.ErrMsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}

	rsp.UserId = strconv.Itoa(user.Id)
	rsp.Name = user.Name
	rsp.Mobile = user.Mobile
	rsp.AvatarUrl = user.Avatar_url
	rsp.RealName = user.Real_name
	rsp.IdCard = user.Id_card

	return nil
}
