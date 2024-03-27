package handler

import (
	models "IHome/IhomeWeb/model"
	"IHome/IhomeWeb/utils"
	"context"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	_ "github.com/gomodule/redigo/redis"
	"time"

	pb "IHome/GetArea/proto"
)

type GetArea struct{}

func (e *GetArea) ClientStream(ctx context.Context, stream pb.GetArea_ClientStreamStream) error {
	//TODO implement me
	panic("implement me")
}

func (e *GetArea) ServerStream(ctx context.Context, request *pb.ServerStreamRequest, stream pb.GetArea_ServerStreamStream) error {
	//TODO implement me
	panic("implement me")
}

func (e *GetArea) BidiStream(ctx context.Context, stream pb.GetArea_BidiStreamStream) error {
	//TODO implement me
	panic("implement me")
}

var (
	redis_conf = map[string]string{
		"key": utils.G_server_name,
		// 127.0.0.1:6379
		"conn":  utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum": utils.G_redis_dbnum,
	}
	redis_key = "area_info"
)

func (e *GetArea) GetArea(ctx context.Context, req *pb.CallRequest, rsp *pb.CallResponse) error {
	beego.Info("请求地区信息 GetArea api/v1.0/areas")
	// 初始化 错误码。之后可能会修改
	rsp.Errno = utils.RECODE_OK
	rsp.ErrMsg = utils.RecodeText(rsp.Errno)
	// 1 从缓存中获取数据  如果有数据就发给前端
	// 准备连接redis
	beego.Info(redis_conf)
	// 将map转化为json
	redis_conf_json, err := json.Marshal(redis_conf)
	//创建redis句柄
	bm, err := cache.NewCache("redis", string(redis_conf_json))
	if err != nil {
		beego.Info("redis连接失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.ErrMsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	// 获取数据  需要定义一个key来用作area查询  area_info
	area_value := bm.Get(redis_key)
	// 读出json字符串
	if area_value != nil {
		beego.Info("从redis中获取到地域信息")
		area_map := []map[string]interface{}{}
		// 将获取到的数据json反序列化
		json.Unmarshal(area_value.([]byte), &area_map)
		beego.Info("从缓存中得到area数据", area_map)
		for _, value := range area_map {
			//beego.Info(key, value)
			// CallResponse_Areas是嵌套的那个message类型
			tmp := pb.CallResponse_Areas{Aid: int32(value["aid"].(float64)), Aname: value["aname"].(string)}
			rsp.Data = append(rsp.Data, &tmp)
		}
		// 已经从缓存中获取数据，后面无需执行
		return nil
	}

	// 2 没有数据就从mysql中查找数据
	// 创建orm句柄
	o := orm.NewOrm()
	// 查询语句
	qs := o.QueryTable("area")
	// 用什么来接收数据
	area := []models.Area{}
	num, err := qs.All(&area)
	// 如果查询失败，修改以上定义的错误码
	if err != nil {
		beego.Info("数据库查询失败")
		rsp.Errno = utils.RECODE_DATAERR
		rsp.ErrMsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	if num == 0 {
		beego.Info("数据库未查询到数据")
		rsp.Errno = utils.RECODE_NODATA
		rsp.ErrMsg = utils.RecodeText(rsp.Errno)

	}
	// 3 将查找到的数据存入缓存
	// 需要将获取到的数据转为json
	area_json, _ := json.Marshal(area)
	// 操作redis存入数据
	err = bm.Put(redis_key, area_json, 3600*time.Second)
	if err != nil {
		beego.Info("地域信息数据存入缓存失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.ErrMsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	// 4 将查到的数据发送给前端
	// 将查询到的数据按照proto格式发送给web服务
	for _, value := range area {
		//beego.Info(key, value)
		// CallResponse_Areas是嵌套的那个message类型
		tmp := pb.CallResponse_Areas{Aid: int32(value.Id), Aname: value.Name}
		rsp.Data = append(rsp.Data, &tmp)
	}

	return nil
}
