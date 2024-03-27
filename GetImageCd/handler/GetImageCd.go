package handler

import (
	pb "IHome/GetImageCd/proto"
	"IHome/IhomeWeb/utils"
	"context"
	"encoding/json"
	"github.com/afocus/captcha"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/gomodule/redigo/redis"
	"image/color"
	"time"
)

type GetImageCd struct{}

var (
	redis_conf = map[string]string{
		"key": utils.G_server_name,
		// 127.0.0.1:6379
		"conn":  utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum": utils.G_redis_dbnum,
	}
)

func (e *GetImageCd) GetImageCd(ctx context.Context, request *pb.Request, response *pb.Response) error {
	//TODO implement me
	beego.Info("获取图片验证码 url：/api/v1.0/imagecode/:uuid")
	cap := captcha.New()

	if err := cap.SetFont("comic.ttf"); err != nil {
		panic(err.Error())
	}

	// 大小根据前端要素决定
	cap.SetSize(90, 41)
	cap.SetDisturbance(captcha.NORMAL)
	cap.SetFrontColor(color.RGBA{255, 255, 255, 255})
	cap.SetBkgColor(color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 153, 0, 255})

	img, str := cap.Create(6, captcha.ALL)
	// 准备连接redis
	beego.Info(redis_conf)
	// 将map转化为json
	redis_conf_json, err := json.Marshal(redis_conf)
	//创建redis句柄
	bm, err := cache.NewCache("redis", string(redis_conf_json))

	if err != nil {
		beego.Info("redis连接失败", err)
		response.Errno = utils.RECODE_DBERR
		response.ErrMsg = utils.RecodeText(response.Errno)
		return nil
	}

	str_json, _ := json.Marshal(str)
	err = bm.Put(request.Uuid, str_json, 300*time.Second)
	if err != nil {
		beego.Info("redis缓存失败", err)
		response.Errno = utils.RECODE_DBERR
		response.ErrMsg = utils.RecodeText(response.Errno)
		return nil
	}
	response.Pix = []byte(img.Pix)
	response.Min = &pb.Response_Point{X: int64(img.Rect.Min.X), Y: int64(img.Rect.Min.Y)}
	response.Max = &pb.Response_Point{X: int64(img.Rect.Max.X), Y: int64(img.Rect.Max.Y)}
	response.Stride = int64(img.Stride)
	response.Errno = utils.RECODE_OK
	response.ErrMsg = utils.RecodeText(response.Errno)
	return nil
}
