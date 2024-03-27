package handler

import (
	models "IHome/IhomeWeb/model"
	"IHome/IhomeWeb/utils"
	pb "IHome/PostLogin/proto"
	"context"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	_ "github.com/gomodule/redigo/redis"
	"time"
)

type PostLogin struct{}

var (
	redis_conf = map[string]string{
		"key": utils.G_server_name,
		// 127.0.0.1:6379
		"conn":  utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum": utils.G_redis_dbnum,
	}
)

func (e *PostLogin) PostLogin(_ context.Context, req *pb.PostLoginRequest, rsp *pb.PostLoginResponse) error {
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

	o := orm.NewOrm()
	var user models.User
	qs := o.QueryTable("user")

	// 这里的one尝试从数据库中查询满足条件的单条记录。  注意再进行查询的时候，Filter("mobile", req.Mobile)后面一定要加All或者是One，并且All返回的是一个切片
	err = qs.Filter("mobile", req.Mobile).One(&user)
	if err != nil {
		beego.Info("尝试登陆时，从数据库中查询用户信息失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.ErrMsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}
	beego.Info(utils.Getmd5string(req.Password))
	beego.Info(user.Password_hash)
	// 数据库中存的是密文，这里对request中的password加密后，与数据库中密文保存的密码相比较
	if utils.Getmd5string(req.Password) != user.Password_hash {
		beego.Info("用户密码不匹配", err)
		rsp.Errno = utils.RECODE_PWDERR
		rsp.ErrMsg = utils.RecodeText(utils.RECODE_PWDERR)
		return nil
	}

	// 创建sessionid=hash(mobile+password)
	sessionid := utils.Getmd5string(req.Mobile + req.Password)
	rsp.Sessionid = sessionid

	// 拼接key
	// user_id
	session_userid := sessionid + "user_id"
	// name
	session_name := sessionid + "name"
	// mobile
	session_mobile := sessionid + "mobile"

	// 将登陆信息进行缓存

	_ = bm.Put(session_userid, user.Id, time.Second*600)
	_ = bm.Put(session_name, user.Name, time.Second*600)
	_ = bm.Put(session_mobile, user.Mobile, time.Second*600)

	return nil
}
