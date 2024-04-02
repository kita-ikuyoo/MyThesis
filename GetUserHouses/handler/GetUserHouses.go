package handler

import (
	pb "IHome/GetUserHouses/proto"
	models "IHome/IhomeWeb/model"
	"IHome/IhomeWeb/utils"
	"context"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	"strconv"
)

type GetUserHouses struct{}

var (
	redis_conf = map[string]string{
		"key": utils.G_server_name,
		// 127.0.0.1:6379
		"conn":  utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum": utils.G_redis_dbnum,
	}
)

func (e *GetUserHouses) GetUserHouses(_ context.Context, req *pb.GetUserHousesRequest, rsp *pb.GetUserHousesResponse) error {
	beego.Info("获取当前用户所发布的房源 GetUserHouses /api/v1.0/user/houses")

	// 初始化错误码
	rsp.Errno = utils.RECODE_OK
	rsp.ErrMsg = utils.RecodeText(utils.RECODE_OK)

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
	id, err := redis.String(user_id, nil)
	new_id, _ := strconv.Atoi(id)
	//创建user对象
	o := orm.NewOrm()
	qs := o.QueryTable("house")

	house_list := []models.House{}
	_, err = qs.Filter("user_id", new_id).All(&house_list)
	if err != nil {
		beego.Info("查询房屋数据失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.ErrMsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}
	house, _ := json.Marshal(house_list)
	rsp.Houses = house

	return nil
}
