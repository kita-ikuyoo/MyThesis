package handler

import (
	models "IHome/IhomeWeb/model"
	"IHome/IhomeWeb/utils"
	pb "IHome/PostAvatar/proto"
	"context"
	"encoding/json"
	"path"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	_ "github.com/garyburd/redigo/redis"
	_ "github.com/gomodule/redigo/redis"
)

type PostAvatar struct{}

var (
	redis_conf = map[string]string{
		"key": utils.G_server_name,
		// 192.168.87.198
		"conn":  utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum": utils.G_redis_dbnum,
	}
)

func (e *PostAvatar) PostAvatar(ctx context.Context, req *pb.PostAvatarRequest, rsp *pb.PostAvatarResponse) error {
	beego.Info("获取用户信息 GetUserInfo /api/v1.0/user")
	//Initialize the return value
	rsp.Errno = utils.RECODE_OK
	rsp.ErrMsg = utils.RecodeText(rsp.Errno)
	// check validity of the avatar
	size := len(req.Avatar)
	if req.FileSize != int64(size) {
		beego.Info("用户头像传输数据丢失")
		rsp.Errno = utils.RECODE_DATAERR
		rsp.ErrMsg = utils.RecodeText(rsp.Errno)
	}

	//get fileext  ,example: (dot).jpg   pay attention to the dot
	file_ext := path.Ext(req.FileExt)

	//call fdfs func to upload avatar to server
	file_id, err := utils.UploadByBuffer(req.Avatar, file_ext[1:])
	if err != nil {
		beego.Info("向fastdfs服务器中上传图片失败", err)
		rsp.Errno = utils.RECODE_DATAERR
		rsp.ErrMsg = utils.RecodeText(rsp.Errno)
	}

	// show fileid(url)
	beego.Info(file_id)
	// get sessionid
	session_id := req.SessionId
	// connect redis
	redis_conf_json, err := json.Marshal(redis_conf)
	if err != nil {
		beego.Info("Failed to connect to redis")
		rsp.Errno = utils.RECODE_DBERR
		rsp.ErrMsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}
	bm, _ := cache.NewCache("redis", string(redis_conf_json))
	// get user_id
	user_id := bm.Get(session_id + "user_id")
	if user_id == nil {
		beego.Info("根据session从redis中获取用户id失败")
		rsp.Errno = utils.RECODE_DBERR
		rsp.ErrMsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}
	user_id_str, _ := redis.String(user_id, nil)
	//为什么这里还要转？不是已经用redis.String转过了吗？这是因为我们model.User中的id是int类型
	id, _ := strconv.Atoi(user_id_str)

	// write fileid into the mysql corresponding with user_id

	user := models.User{
		Id:         id,
		Avatar_url: file_id,
	}
	o := orm.NewOrm()
	_, err = o.Update(&user, "avatar_url")
	if err != nil {
		beego.Info("在数据库中更新用户头像链接失败")
		rsp.Errno = utils.RECODE_DBERR
		rsp.ErrMsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}
	// return fileid to client

	rsp.AvatarUrl = file_id

	return nil
}
