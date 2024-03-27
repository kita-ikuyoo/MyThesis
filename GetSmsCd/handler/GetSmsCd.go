package handler

import (
	pb "IHome/GetSmsCd/proto"
	models "IHome/IhomeWeb/model"
	"IHome/IhomeWeb/utils"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	_ "github.com/gomodule/redigo/redis"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	sms_account  = "C01095301"
	sms_password = "34230e77fa9a60709b5473628f9fefb4"
	sms_tmplate  = "您的验证码是：。请不要把验证码泄露给其他人。"
)

type GetSmsCd struct{}

func (e *GetSmsCd) ClientStream(ctx context.Context, stream pb.GetSmsCd_ClientStreamStream) error {
	//TODO implement me
	panic("implement me")
}

func (e *GetSmsCd) ServerStream(ctx context.Context, request *pb.ServerStreamRequest, stream pb.GetSmsCd_ServerStreamStream) error {
	//TODO implement me
	panic("implement me")
}

func (e *GetSmsCd) BidiStream(ctx context.Context, stream pb.GetSmsCd_BidiStreamStream) error {
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
)

func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
func ihuyi(cd int64, mobile string) {
	v := url.Values{}
	//fmt.Printf(_now)
	account := sms_account   //用户名是登录用户中心->验证码短信->产品总览->APIID
	password := sms_password //查看密码请登录用户中心->验证码短信->产品总览->APIKEY

	SMSUrl := fmt.Sprintf("http://106.ihuyi.com/webservice/sms.php?method=Submit&account=%s&password=%s&mobile=%s&content=您的验证码是：%v。请不要把验证码泄露给其他人。", account, password, mobile, cd)
	//body := ioutil.NopCloser(strings.NewReader(v.Encode())) //把form数据编下码
	body := strings.NewReader(v.Encode()) //把form数据编下码
	client := &http.Client{}
	req, _ := http.NewRequest("POST", SMSUrl, body)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	//fmt.Printf("%+v\n", req) //看下发送的结构

	resp, err := client.Do(req) //发送
	defer resp.Body.Close()     //一定要关闭resp.Body
	data, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(data), err)
}

func (e *GetSmsCd) GetSmsCd(ctx context.Context, request *pb.SMSRequest, response *pb.SMSResponse) error {
	beego.Info("获取短信验证码 /api/v1.0/smscode/:mobile")
	response.Errno = utils.RECODE_OK
	response.ErrMsg = utils.RecodeText(utils.RECODE_OK)

	// 验证手机号是否存在
	// 创建数据库句柄
	o := orm.NewOrm()
	// 使用手机号作为查询条件
	user := models.User{Mobile: request.Mobile}

	err := o.Read(&user)
	// 如果查询到，err就为nil，也就代表着该手机号存在，不得再次创建
	if err == nil {
		beego.Info("该手机号已存在", err)
		response.Errno = utils.RECODE_MOBILEERR
		response.ErrMsg = utils.RecodeText(utils.RECODE_MOBILEERR)
		return nil
	}
	// 验证图片验证码是否正确
	// 连接redis

	redis_conf_json, err := json.Marshal(redis_conf)
	//创建redis句柄
	bm, err := cache.NewCache("redis", string(redis_conf_json))
	if err != nil {
		beego.Info("redis连接失败", err)
		response.Errno = utils.RECODE_DBERR
		response.ErrMsg = utils.RecodeText(response.Errno)
		return nil
	}
	// value 是uint8类型
	value := bm.Get(request.Uuid)
	if value == nil {
		beego.Info("redis获取失败", err)
		response.Errno = utils.RECODE_DBERR
		response.ErrMsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}
	value_str, _ := redis.String(value, nil)
	Text := `"` + request.Text + `"`
	value_str_upper := strings.ToUpper(value_str)
	Text_upper := strings.ToUpper(Text)
	if value_str_upper != Text_upper {
		beego.Info("图片验证码错误", err)
		response.Errno = utils.RECODE_DATAERR
		response.ErrMsg = utils.RecodeText(utils.RECODE_DATAERR)
		return nil
	}
	cd := int64(rand.Intn(9000) + 1000)
	// 发送短信验证码

	//ihuyi(cd, request.Mobile)
	err = bm.Put(request.Mobile, cd, 300*time.Second)
	beego.Info(cd)
	if err != nil {
		beego.Info("redis插入失败", err)
		response.Errno = utils.RECODE_DBERR
		response.ErrMsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}
	return nil

	// 通过uuid查找图片验证码的值
	//调用短信接口发送短信

	//将短信验证码存入缓存库

}
