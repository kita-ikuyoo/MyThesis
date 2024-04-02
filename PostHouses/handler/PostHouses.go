package handler

import (
	models "IHome/IhomeWeb/model"
	"IHome/IhomeWeb/utils"
	pb "IHome/PostHouses/proto"
	"context"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	"reflect"
	"strconv"

	"encoding/json"
)

type PostHouses struct{}

var (
	redis_conf = map[string]string{
		"key": utils.G_server_name,
		// 127.0.0.1:6379
		"conn":  utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum": utils.G_redis_dbnum,
	}
)

func (e *PostHouses) PostHouses(ctx context.Context, req *pb.PostHousesRequest, rsp *pb.PostHousesResponse) error {
	beego.Info("PostHouses 发布房源信息 /api/v1.0/houses ")
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

	id, err := redis.String(user_id, nil)
	beego.Info(reflect.TypeOf(id), id)
	new_id, _ := strconv.Atoi(id)

	// 解析客户端发送的houses信息
	var Requestmap = make(map[string]interface{})
	json.Unmarshal(req.Houses, &Requestmap)

	// 准备插入数据库
	house := models.House{}
	house.Title = Requestmap["title"].(string)
	price, _ := strconv.Atoi(Requestmap["price"].(string))
	house.Price = price * 100
	house.Address = Requestmap["address"].(string)
	house.Room_count, _ = strconv.Atoi(Requestmap["room_count"].(string))
	house.Acreage, _ = strconv.Atoi(Requestmap["acreage"].(string))
	house.Unit = Requestmap["unit"].(string)
	house.Capacity, _ = strconv.Atoi(Requestmap["capacity"].(string))
	house.Beds, _ = Requestmap["beds"].(string)
	depsoit, _ := strconv.Atoi(Requestmap["deposit"].(string))
	house.Deposit = depsoit * 100
	house.Min_days, _ = strconv.Atoi(Requestmap["min_days"].(string))
	house.Max_days, _ = strconv.Atoi(Requestmap["max_days"].(string))
	area_id, _ := strconv.Atoi(Requestmap["area_id"].(string))
	area := models.Area{Id: area_id}
	house.Area = &area
	user := models.User{Id: new_id}
	house.User = &user

	facility := []*models.Facility{}
	for _, f_id := range Requestmap["facility"].([]interface{}) {
		fid, _ := strconv.Atoi(f_id.(string))
		fac := &models.Facility{Id: fid}
		facility = append(facility, fac)
	}

	o := orm.NewOrm()
	house_id, err := o.Insert(&house)
	if err != nil {
		beego.Info("数据插入失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.ErrMsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	m2m := o.QueryM2M(&house, "facility")
	_, err = m2m.Add(facility)
	if err != nil {
		beego.Info("房屋设施多对多插入失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.ErrMsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	rsp.HouseId = strconv.Itoa(int(house_id))

	return nil
}
