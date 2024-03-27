package handler

import (
	models "IHome/IhomeWeb/model"
	"IHome/IhomeWeb/utils"
	pb "IHome/PostRet/proto"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	_ "github.com/gomodule/redigo/redis"
	"time"
)

type PostRet struct{}

var (
	redis_conf = map[string]string{
		"key": utils.G_server_name,
		// 127.0.0.1:6379
		"conn":  utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum": utils.G_redis_dbnum,
	}
)

func Md5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func (e *PostRet) PostRet(ctx context.Context, req *pb.PostRetRequest, rsp *pb.PostRetResponse) error {
	beego.Info("PostRet 注册  /api/v1.0/users")
	rsp.Error = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(utils.RECODE_OK)
	// 验证短信验证码
	// 连接redis
	redis_conf_json, err := json.Marshal(redis_conf)
	//创建redis句柄
	bm, err := cache.NewCache("redis", string(redis_conf_json))
	if err != nil {
		beego.Info("redis连接失败", err)
		rsp.Error = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Error)
		return nil
	}
	// 通过手机号获取短信验证码
	sms_code := bm.Get(req.Mobile)
	if sms_code == nil {
		beego.Info("获取数据失败")
		rsp.Error = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
	}
	//进行短信验证码对比
	sms_code_str, err := redis.String(sms_code, nil)
	SmsCode := req.SmsCode

	if sms_code_str != SmsCode {
		beego.Info("短信验证码错误")
		rsp.Error = utils.RECODE_SMSERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_SMSERR)
		return nil
	}
	// 将数据存入数据库
	o := orm.NewOrm()
	user := models.User{Mobile: req.Mobile, Password_hash: Md5String(req.Password), Name: req.Mobile}
	id, err := o.Insert(&user)
	if err != nil {
		beego.Info("注册失败")
		rsp.Error = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}
	beego.Info("user_id", id)

	//创建sessionid
	session_id := Md5String(req.Mobile + req.Password)
	rsp.SessionId = session_id

	// 以sessionid为key的一部分创建session
	// name
	bm.Put(session_id+"name", user.Mobile, 3600*time.Second)
	// user_id    这里使用创建账户时返回的id
	bm.Put(session_id+"user_id", id, 3600*time.Second)
	// 手机号
	bm.Put(session_id+"mobile", user.Mobile, 3600*time.Second)
	return nil
}
