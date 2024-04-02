package handler

import (
	deletesession "IHome/DeleteSession/proto"
	getarea "IHome/GetArea/proto"
	getimagecd "IHome/GetImageCd/proto"
	getsession "IHome/GetSession/proto"
	getsmscd "IHome/GetSmsCd/proto"
	getuserhouses "IHome/GetUserHouses/proto"
	getuserinfo "IHome/GetUserInfo/proto"
	models "IHome/IhomeWeb/model"
	"IHome/IhomeWeb/utils"
	postavatar "IHome/PostAvatar/proto"
	posthouses "IHome/PostHouses/proto"
	postlogin "IHome/PostLogin/proto"
	postret "IHome/PostRet/proto"
	postuserauth "IHome/PostUserAuth/proto"
	"context"
	"encoding/json"
	"fmt"
	"github.com/afocus/captcha"
	"github.com/astaxie/beego"
	grpcc "github.com/go-micro/plugins/v4/client/grpc"
	"github.com/go-micro/plugins/v4/registry/consul"
	"github.com/julienschmidt/httprouter"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

var (
	GetAreaServerName       = "go.micro.server.GetArea"
	GetImageCdServerName    = "go.micro.server.GetImageCd"
	GetSmsCdServerName      = "go.micro.server.GetSmsCd"
	PostRetServerName       = "go.micro.server.PostRet"
	GetSessionServerName    = "go.micro.server.GetSession"
	PostLoginServerName     = "go.micro.server.PostLogin"
	DeleteSessionServerName = "go.micro.server.DeleteSession"
	GetUserInfoServerName   = "go.micro.server.GetUserInfo"
	PostAvatarServerName    = "go.micro.server.PostAvatar"
	PostUserAuthServerName  = "go.micro.server.PostUserAuth"
	GetUserHousesServerName = "go.micro.server.GetUserHouses"
	PostHousesServerName    = "go.micro.server.PostHouses"
)

func GetArea(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	beego.Info("请求地区信息 GetArea api/v1.0/areas")
	consulReg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{"127.0.0.1:8500"}
	})
	service := micro.NewService(
		micro.Registry(consulReg),
		micro.Client(grpcc.NewClient()),
	)
	mc := getarea.NewGetAreaService(GetAreaServerName, service.Client())
	rsp, err := mc.GetArea(context.TODO(), &getarea.CallRequest{})
	// 调用服务传回句柄
	if err != nil {
		log.Println(err.Error())
		beego.Info("wrong")
		return
	}
	var areas []models.Area
	for _, value := range rsp.Data {
		temp := models.Area{Id: int(value.Aid), Name: value.Aname}
		areas = append(areas, temp)
	}
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.ErrMsg,
		"data":   areas,
	}

	w.Header().Set("Content-Type", "application/json")
	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

}

func GetSession(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("请求session GetSession api/v1.0/session")

	// cookie里面存储了键值对，这里根据"userlogin"这个键获得键值对，value中保存了sessionid
	// 看源码。当不存在对应的cookie时，返回nil, ErrNoCookie
	cookie, err := r.Cookie("userlogin")
	// golang具有短路性质，这里err判别完成后不考虑cookie.Value == ""。不然肯定会因为野指针报错
	// 所以说这里的cookie.Value == ""是完全应该删除的吗？   也不一定，有可能存在对应的cookie，其value就是“”也说不准
	if err != nil || cookie.Value == "" {
		// 直接返回说明用户未登录
		beego.Info("查看是否进入判断")
		response := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}
		//设置返回数据的格式
		w.Header().Set("Content-Type", "application/json")
		//将map转化为json 返回给前端
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}
	consulReg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{"127.0.0.1:8500"}
	})
	service := micro.NewService(
		micro.Registry(consulReg),
		micro.Client(grpcc.NewClient()),
	)
	mc := getsession.NewGetSessionService(GetSessionServerName, service.Client())
	rsp, err := mc.GetSession(context.TODO(), &getsession.GetSessionRequest{
		SessionId: cookie.Value,
	})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	data := make(map[string]string)
	data["name"] = rsp.UserName

	// 这里是将数据返回给前端，而不是将数据从server返回给client
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.ErrMsg,
		"data":   data,
	}
	//设置返回数据的格式
	w.Header().Set("Content-Type", "application/json")
	//将map转化为json 返回给前端
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

}

// 获取首页轮播图
func GetIndex(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	beego.Info("请求首页轮播图 GetIndex api/v1.0/house/index")
	//准备返回给前端的map
	response := map[string]interface{}{
		"errno":  utils.RECODE_OK,
		"errmsg": utils.RecodeText(utils.RECODE_OK),
	}
	//设置返回数据的格式
	w.Header().Set("Content-Type", "application/json")
	//将map转化为json 返回给前端
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func GetImageCd(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	beego.Info("获取图片验证码 url：/api/v1.0/imagecode/:uuid")
	uuid := ps.ByName("uuid")
	consulReg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{"127.0.0.1:8500"}
	})
	service := micro.NewService(
		micro.Registry(consulReg),
		micro.Client(grpcc.NewClient()),
	)
	// call the backend service
	GetImageClient := getimagecd.NewGetImageCdService(GetImageCdServerName, service.Client())
	rsp, err := GetImageClient.GetImageCd(context.TODO(), &getimagecd.Request{
		Uuid: uuid,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	//判断是否返回图片
	if rsp.Errno != "0" {
		w.Header().Set("Content-Type", "application/json")
		response := map[string]interface{}{
			"errno":  rsp.Errno,
			"errmsg": rsp.ErrMsg,
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
	//拼接图片结构体发送给前端
	var img image.RGBA
	for _, value := range rsp.Pix {
		img.Pix = append(img.Pix, uint8(value))
	}
	img.Stride = int(rsp.Stride)
	img.Rect.Min.X = int(rsp.Min.X)
	img.Rect.Min.Y = int(rsp.Min.Y)
	img.Rect.Max.X = int(rsp.Max.X)
	img.Rect.Max.Y = int(rsp.Max.Y)

	var image captcha.Image
	image.RGBA = &img
	fmt.Println(image)
	// 将图片发送给浏览器
	png.Encode(w, image)
}

func Getsmscd(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	beego.Info("获取短信验证码 /api/v1.0/smscode/:mobile")
	//获取手机号
	mobile := ps.ByName("mobile")
	id := r.URL.Query()["id"][0]
	text := r.URL.Query()["text"][0]

	// 判断手机号是否正确
	mobile_reg := regexp.MustCompile(`0?(13|14|15|17|18|19)[0-9]{9}`)
	mobile_bool := mobile_reg.MatchString(mobile)
	// 如果手机号错误直接返回错误，不调用服务
	if mobile_bool == false {
		response := map[string]interface{}{
			"errno":  utils.RECODE_MOBILEERR,
			"errmsg": utils.RecodeText(utils.RECODE_MOBILEERR),
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}
	consulReg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{"127.0.0.1:8500"}
	})
	service := micro.NewService(
		micro.Registry(consulReg),
		micro.Client(grpcc.NewClient()),
	)
	GetSmsCdClient := getsmscd.NewGetSmsCdService(GetSmsCdServerName, service.Client())
	rsp, err := GetSmsCdClient.GetSmsCd(context.TODO(), &getsmscd.SMSRequest{
		Uuid:   id,
		Text:   text,
		Mobile: mobile,
	})
	if err != nil {
		beego.Info("调用GetSmsCd远程服务失败", err)
		return
	}

	response := map[string]string{
		"errno":  rsp.Errno,
		"errmsg": rsp.ErrMsg,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

}

func PostRet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("注册服务 PostRet api/v1.0/users")
	var request map[string]interface{}
	// 将前端 json 数据解析到 map当中
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if request["mobile"] == "" || request["sms_code"] == "" || request["password"] == "" {
		response := map[string]interface{}{
			"errno":  utils.RECODE_DATAERR,
			"errmsg": utils.RecodeText(utils.RECODE_DATAERR),
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

	}
	consulReg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{"127.0.0.1:8500"}
	})
	service := micro.NewService(
		micro.Registry(consulReg),
		micro.Client(grpcc.NewClient()),
	)
	mc := postret.NewPostRetService(PostRetServerName, service.Client())
	rsp, err := mc.PostRet(context.TODO(), &postret.PostRetRequest{
		Mobile:   request["mobile"].(string),
		SmsCode:  request["sms_code"].(string),
		Password: request["password"].(string),
	})
	// 调用服务传回句柄
	if err != nil {
		log.Println(err.Error())
		beego.Info("wrong")
		return
	}
	cookie, err := r.Cookie("userlogin")
	if err != nil || cookie.Value == "" {
		//maxage表示cookie的最大存活时间，单位是秒。MaxAge: 3600意味着这个cookie将在创建或更新后的3600秒（即1小时）内有效。一旦超过这个时间，cookie就会被浏览器自动删除
		// path:"/"表示该cookie对于任何url都适用
		cookie := http.Cookie{Name: "userlogin", Value: rsp.SessionId, MaxAge: 3600, Path: "/"}
		http.SetCookie(w, &cookie)
	}
	response := map[string]interface{}{
		"errno":  rsp.Error,
		"errmsg": rsp.Errmsg,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//fmt.Println(" 注册服务  PostRet  /api/v1.0/users")
	////接受 前端发送过来数据的
	//var request map[string]interface{}
	//// 将前端 json 数据解析到 map当中
	//if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
	//	http.Error(w, err.Error(), 500)
	//	return
	//}
	//if request["mobile"].(string) == "" || request["password"].(string) == "" || request["sms_code"].(string) == "" {
	//	response := map[string]interface{}{
	//		"errno":  utils.RECODE_DATAERR,
	//		"errmsg": utils.RecodeText(utils.RECODE_DATAERR),
	//	}
	//	//设置返回数据的格式
	//	w.Header().Set("Content-Type", "application/json")
	//	//将map转化为json 返回给前端
	//	if err := json.NewEncoder(w).Encode(response); err != nil {
	//		http.Error(w, err.Error(), 500)
	//		return
	//	}
	//}
	//
	////创建 grpc 客户端
	//cli := grpc.NewService()
	////客户端初始化
	//cli.Init()
	//
	////通过protobuf 生成文件 创建 连接服务端 的客户端句柄
	//exampleClient := POSTRET.NewExampleService("go.micro.srv.PostRet", cli.Client())
	////通过句柄调用服务端函数
	//rsp, err := exampleClient.PostRet(context.TODO(), &POSTRET.Request{
	//	Mobile:   request["mobile"].(string),
	//	Password: request["password"].(string),
	//	SmsCode:  request["sms_code"].(string),
	//})
	//
	////判断是否成功
	//if err != nil {
	//	http.Error(w, err.Error(), 500)
	//	return
	//}
	//
	////设置cookie
	//cookie, err := r.Cookie("IHomelogin")
	//if err != nil || cookie.Value == "" {
	//	cookie := http.Cookie{Name: "IHomelogin", Value: rsp.Sessionid, MaxAge: 600, Path: "/"}
	//	http.SetCookie(w, &cookie)
	//}
	////将数据返回前端
	//response := map[string]interface{}{
	//	"errno":  rsp.Errno,
	//	"errmsg": rsp.Errmsg,
	//}
	//w.Header().Set("Content-Type", "application/json")
	//if err := json.NewEncoder(w).Encode(response); err != nil {
	//	http.Error(w, err.Error(), 500)
	//	return
	//}
}

func PostLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("登陆 api/v1.0/sessions")
	//接受 前端发送过来数据的
	var request map[string]interface{}
	// 将前端 json 数据解析到 map当中
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if request["mobile"] == nil || request["password"] == nil {
		response := map[string]interface{}{
			"errno":  utils.RECODE_DATAERR,
			"errmsg": utils.RecodeText(utils.RECODE_DATAERR),
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
	consulReg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{"127.0.0.1:8500"}
	})
	service := micro.NewService(
		micro.Registry(consulReg),
		micro.Client(grpcc.NewClient()),
	)
	mc := postlogin.NewPostLoginService(PostLoginServerName, service.Client())
	rsp, err := mc.PostLogin(context.TODO(), &postlogin.PostLoginRequest{
		Mobile:   request["mobile"].(string),
		Password: request["password"].(string),
	})
	cookie, err := r.Cookie("userlogin")
	beego.Info("返回的sessionid：", rsp.Sessionid)
	if err != nil || cookie.Value == "" {
		cookie := http.Cookie{Name: "userlogin", Value: rsp.Sessionid, MaxAge: 600, Path: "/"}
		http.SetCookie(w, &cookie)
	}

	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.ErrMsg,
	}
	//设置返回数据的格式
	w.Header().Set("Content-Type", "application/json")
	//将map转化为json 返回给前端
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func DeleteSession(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// decode the incoming request as json
	beego.Info("退出登陆 /api/v1.0/session Deletesession()")

	// 获取cookie
	cookie, err := r.Cookie("userlogin")
	if err != nil || cookie.Value == "" {
		beego.Info("退出登录时获取session失败")
		response := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}
	session_id := cookie.Value
	consulReg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{"127.0.0.1:8500"}
	})
	service := micro.NewService(
		micro.Registry(consulReg),
		micro.Client(grpcc.NewClient()),
	)
	mc := deletesession.NewDeleteSessionService(DeleteSessionServerName, service.Client())
	rsp, err := mc.DeleteSession(context.TODO(), &deletesession.DeleteSessionRequest{
		SessionId: session_id,
	})
	if err != nil {
		beego.Info("调用deletesession远程服务失败")
		http.Error(w, err.Error(), 500)
		return

	}
	// 这里为什么要这么写呢？
	if rsp.Errno == "0" {
		_, err := r.Cookie("userlogin")
		if err == nil {
			cookie := http.Cookie{Name: "userlogin", Path: "/", MaxAge: -1, Value: ""}
			http.SetCookie(w, &cookie)
		}
	}
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.ErrMsg,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	//return
	////获取session
	//userlogin, err := r.Cookie("IHomelogin")
	//if err != nil || userlogin.Value == "" {
	//	log.Println("user not login")
	//	response := map[string]interface{}{
	//		"errno":  utils.RECODE_SESSIONERR,
	//		"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
	//	}
	//
	//	w.Header().Set("Content-Type", "application/json")
	//	if err := json.NewEncoder(w).Encode(response); err != nil {
	//		http.Error(w, err.Error(), 500)
	//		return
	//	}
	//	return
	//}
	//rsp, err := exampleClient.DeleteSession(context.TODO(), &DELETESESSION.Request{
	//	Sessionid: userlogin.Value,
	//})
	//if err != nil {
	//	http.Error(w, err.Error(), 500)
	//	return
	//}
	//if rsp.Errno == "0" {
	//	//将cookie中的sessionid设置为空
	//	_, err = r.Cookie("IHomelogin")
	//	if err == nil {
	//		cookie := http.Cookie{Name: "IHomelogin", Path: "/", MaxAge: -1}
	//		http.SetCookie(w, &cookie)
	//	}
	//}
	////返回数据
	//response := map[string]interface{}{
	//	"errno":  rsp.Errno,
	//	"errmsg": rsp.Errmsg,
	//}
	////设置格式
	//w.Header().Set("Content-Type", "application/json")
	//// encode and write the response as json
	//if err := json.NewEncoder(w).Encode(response); err != nil {
	//	http.Error(w, err.Error(), 500)
	//	return
	//}
}

func GetUserInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("获取用户信息 GetUserInfo /api/v1.0/user")
	cookie, err := r.Cookie("userlogin")
	if err != nil || cookie.Value == "" {
		response := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}
		//设置返回数据的格式
		w.Header().Set("Content-Type", "application/json")
		//将map转化为json 返回给前端
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}
	consulReg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{"127.0.0.1:8500"}
	})
	service := micro.NewService(
		micro.Registry(consulReg),
		micro.Client(grpcc.NewClient()),
	)
	mc := getuserinfo.NewGetUserInfoService(GetUserInfoServerName, service.Client())
	rsp, err := mc.GetUserInfo(context.TODO(), &getuserinfo.GetUserInfoRequest{
		SessionId: cookie.Value,
	})
	if err != nil {
		beego.Info("调用deletesession远程服务失败")
		http.Error(w, err.Error(), 500)
		return

	}
	//准备返回数据
	data := make(map[string]interface{})
	data["user_id"] = rsp.UserId
	data["name"] = rsp.Name
	data["mobile"] = rsp.Mobile
	data["real_name"] = rsp.RealName
	data["id_card"] = rsp.IdCard
	data["avatar_url"] = utils.AddDomain2Url(rsp.AvatarUrl)
	//准备返回给前端的map
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.ErrMsg,
		"data":   data,
	}
	// encode and write the response as json
	//设置返回数据的格式
	w.Header().Set("Content-Type", "application/json")
	//将map转化为json 返回给前端
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
func PostAvatar(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	beego.Info("上传用户头像 PostAvatar /api/v1.0/user/avatar")
	//获取sessionid
	cookie, err := r.Cookie("userlogin")
	if err != nil || cookie.Value == "" {
		//准备返回给前端的map
		response := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		//设置返回数据的格式
		w.Header().Set("Content-Type", "application/json")
		//将map转化为json 返回给前端
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
	file, header, err := r.FormFile("avatar")
	if err != nil {
		beego.Info("get file err:", err)
		//准备返回给前端的map
		response := map[string]interface{}{
			"errno":  utils.RECODE_IOERR,
			"errmsg": utils.RecodeText(utils.RECODE_IOERR),
		}

		//设置返回数据的格式
		w.Header().Set("Content-Type", "application/json")
		//将map转化为json 返回给前端
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
	beego.Info("文件大小", header.Size)
	beego.Info("文件名", header.Filename)
	//创建文件大小的切片，这是因为fastdfs中上传文件的操作需要
	filebuffer := make([]byte, header.Size)
	//将file中的数据读入filebuffer
	_, err = file.Read(filebuffer)
	if err != nil {
		beego.Info("get file err:", err)
		//准备返回给前端的map
		response := map[string]interface{}{
			"errno":  utils.RECODE_IOERR,
			"errmsg": utils.RecodeText(utils.RECODE_IOERR),
		}

		//设置返回数据的格式
		w.Header().Set("Content-Type", "application/json")
		//将map转化为json 返回给前端
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	//创建 grpc 客户端
	consulReg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{"127.0.0.1:8500"}
	})
	service := micro.NewService(
		micro.Registry(consulReg),
		micro.Client(grpcc.NewClient()),
	)
	mc := postavatar.NewPostAvatarService(PostAvatarServerName, service.Client())
	rsp, err := mc.PostAvatar(context.TODO(), &postavatar.PostAvatarRequest{
		SessionId: cookie.Value,
		FileExt:   header.Filename,
		FileSize:  header.Size,
		Avatar:    filebuffer,
	})
	// 调用服务传回句柄
	if err != nil {
		log.Println(err.Error())
		beego.Info("wrong")
		return
	}

	//返回数据
	data := make(map[string]interface{})
	//这里是将服务端返回的url进行拼接后返回给前端
	data["avatar_url"] = utils.AddDomain2Url(rsp.AvatarUrl)
	//准备返回给前端的map
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.ErrMsg,
		"data":   data,
	}
	log.Println("data is ", data)
	//设置返回数据的格式
	w.Header().Set("Content-Type", "application/json")
	//将map转化为json 返回给前端
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

// 用户信息检查
func GetUserAuth(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("用户信息检查 GetUserAuth /api/v1.0/user/auth")
	cookie, err := r.Cookie("userlogin")
	if err != nil || cookie.Value == "" {
		response := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}
		//设置返回数据的格式
		w.Header().Set("Content-Type", "application/json")
		//将map转化为json 返回给前端
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}
	consulReg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{"127.0.0.1:8500"}
	})
	service := micro.NewService(
		micro.Registry(consulReg),
		micro.Client(grpcc.NewClient()),
	)
	mc := getuserinfo.NewGetUserInfoService(GetUserInfoServerName, service.Client())
	rsp, err := mc.GetUserInfo(context.TODO(), &getuserinfo.GetUserInfoRequest{
		SessionId: cookie.Value,
	})
	if err != nil {
		beego.Info("调用deletesession远程服务失败")
		http.Error(w, err.Error(), 500)
		return

	}
	//准备返回数据
	data := make(map[string]interface{})
	data["user_id"] = rsp.UserId
	data["name"] = rsp.Name
	data["mobile"] = rsp.Mobile
	data["real_name"] = rsp.RealName
	data["id_card"] = rsp.IdCard
	data["avatar_url"] = utils.AddDomain2Url(rsp.AvatarUrl)
	//准备返回给前端的map
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.ErrMsg,
		"data":   data,
	}
	// encode and write the response as json
	//设置返回数据的格式
	w.Header().Set("Content-Type", "application/json")
	//将map转化为json 返回给前端
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func PostUserAuth(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println("更新实名认证检测  URL: /api/v1.0/user/auth PostUserAuth ")
	//接受 前端发送过来数据的
	var request map[string]interface{}
	// 将前端 json 数据解析到 map当中
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	//数据校验
	if request["real_name"].(string) == "" || request["id_card"].(string) == "" {
		//准备返回给前端的map
		response := map[string]interface{}{
			"errno":  utils.RECODE_NODATA,
			"errmsg": utils.RecodeText(utils.RECODE_NODATA),
		}
		//设置返回数据的格式
		w.Header().Set("Content-Type", "application/json")
		//将map转化为json 返回给前端
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
	//获取sessionid
	cookie, err := r.Cookie("userlogin")
	if err != nil || cookie.Value == "" {
		//准备返回给前端的map
		response := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}
		//设置返回数据的格式
		w.Header().Set("Content-Type", "application/json")
		//将map转化为json 返回给前端
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	//创建 grpc 客户端
	consulReg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{"127.0.0.1:8500"}
	})
	service := micro.NewService(
		micro.Registry(consulReg),
		micro.Client(grpcc.NewClient()),
	)
	mc := postuserauth.NewPostUserAuthService(PostUserAuthServerName, service.Client())
	rsp, err := mc.PostUserAuth(context.TODO(), &postuserauth.PostUserAuthRequest{
		SessionId: cookie.Value,
		RealName:  request["real_name"].(string),
		IdCard:    request["id_card"].(string),
	})
	if err != nil {
		beego.Info("调用postuserauth远程服务失败", err)
		http.Error(w, err.Error(), 500)
		return

	}

	//刷新cookie时间
	cookienew := http.Cookie{Name: "userlogin", Value: cookie.Value, Path: "/", MaxAge: 600}
	http.SetCookie(w, &cookienew)
	//准备返回给前端的map
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.ErrMsg,
	}
	//设置返回数据的格式
	w.Header().Set("Content-Type", "application/json")
	//将map转化为json 返回给前端
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func GetUserHouses(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	beego.Info("获取当前用户所发布的房源 GetUserHouses /api/v1.0/user/houses")

	cookie, err := req.Cookie("userlogin")
	if err != nil || cookie.Value == "" {
		//准备返回给前端的map
		response := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}
		//设置返回数据的格式
		w.Header().Set("Content-Type", "application/json")
		//将map转化为json 返回给前端
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
	//创建 grpc 客户端
	consulReg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{"127.0.0.1:8500"}
	})
	service := micro.NewService(
		micro.Registry(consulReg),
		micro.Client(grpcc.NewClient()),
	)
	mc := getuserhouses.NewGetUserHousesService(GetUserHousesServerName, service.Client())
	rsp, err := mc.GetUserHouses(context.TODO(), &getuserhouses.GetUserHousesRequest{
		SessionId: cookie.Value,
	})
	if err != nil {
		beego.Info("调用getuserhouses远程服务失败", err)
		http.Error(w, err.Error(), 500)
		return
	}
	//房屋切片信息
	house_list := []models.House{}
	json.Unmarshal(rsp.Houses, &house_list)
	//将房屋切片信息转换成map切片返回给前端
	// houses是一个存储了map的切片
	var houses []interface{}
	for _, houseinfo := range house_list {
		houses = append(houses, houseinfo.To_house_info())
	}
	data_map := make(map[string]interface{})
	data_map["houses"] = houses
	//返回数据
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.ErrMsg,
		"data":   data_map,
	}

	w.Header().Set("Content-Type", "application/json")
	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func PostHouses(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// decode the incoming request as json
	beego.Info("PostHouses 发布房源信息 /api/v1.0/houses ")
	// body就是一个json的二进制流
	body, _ := ioutil.ReadAll(r.Body)

	cookie, err := r.Cookie("userlogin")
	if err != nil || cookie.Value == "" {
		//准备返回给前端的map
		response := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}
		//设置返回数据的格式
		w.Header().Set("Content-Type", "application/json")
		//将map转化为json 返回给前端
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
	// 这里我们前端传来的信息转为二进制流，发给服务端
	consulReg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{"127.0.0.1:8500"}
	})
	service := micro.NewService(
		micro.Registry(consulReg),
		micro.Client(grpcc.NewClient()),
	)
	mc := posthouses.NewPostHousesService(PostHousesServerName, service.Client())
	rsp, err := mc.PostHouses(context.TODO(), &posthouses.PostHousesRequest{
		SessionId: cookie.Value,
		Houses:    body,
	})
	if err != nil {
		beego.Info("调用posthouses远程服务失败", err)
		http.Error(w, err.Error(), 500)
		return
	}
	data := make(map[string]interface{})
	data["house_id"] = rsp.HouseId
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.ErrMsg,
		"data":   data,
	}
	w.Header().Set("Content-Type", "application/json")
	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

/*

func PutUserInfo(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println("更新用户名   PutUserInfo   /api/v1.0/user/name")
	//接受 前端发送过来数据的
	var request map[string]interface{}
	// 将前端 json 数据解析到 map当中
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	//数据校验
	username := request["name"].(string)
	if username == "" {
		response := map[string]interface{}{
			"errno":  utils.RECODE_NODATA,
			"errmsg": utils.RecodeText(utils.RECODE_NODATA),
		}
		// encode and write the response as json
		//设置返回数据的格式
		w.Header().Set("Content-Type", "application/json")
		//将map转化为json 返回给前端
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}
	//获取sessionid
	cookie, err := r.Cookie("IHomelogin")
	if err != nil {
		log.Println("获取cookie失败")
		//准备返回给前端的map
		response := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}
		// encode and write the response as json
		//设置返回数据的格式
		w.Header().Set("Content-Type", "application/json")
		//将map转化为json 返回给前端
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
	//创建 grpc 客户端
	cli := grpc.NewService()
	//客户端初始化
	cli.Init()

	//通过protobuf 生成文件 创建 连接服务端 的客户端句柄
	exampleClient := PUTUSERINFO.NewExampleService("go.micro.srv.PutUserInfo", cli.Client())
	//通过句柄调用服务端函数
	rsp, err := exampleClient.PutUserInfo(context.TODO(), &PUTUSERINFO.Request{
		Sessionid: cookie.Value,
		Username:  request["name"].(string),
	})
	//判断是否成功
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	//刷新cookie时间
	cookienew := http.Cookie{Name: "IHomelogin", Value: cookie.Value, Path: "/", MaxAge: 600}
	http.SetCookie(w, &cookienew)
	//返回数据
	data := make(map[string]interface{})
	data["name"] = rsp.Username
	//准备返回给前端的map
	response := map[string]interface{}{
		"errno":  utils.RECODE_MOBILEERR,
		"errmsg": utils.RecodeText(utils.RECODE_MOBILEERR),
		"data":   data,
	}
	//设置返回数据的格式
	w.Header().Set("Content-Type", "application/json")
	//将map转化为json 返回给前端
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

}


func PostHouseImage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	beego.Info("发送房屋图片PostHousesImage  /api/v1.0/houses/:id/images")
	//获取houseid
	houseid := params.ByName("id")
	//获取sessionid
	userlogin, err := r.Cookie("ihomelogin")
	if err != nil {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}

		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			beego.Info(err)
			return
		}
		return
	}
	file, header, err := r.FormFile("house_image")
	if err != nil {
		beego.Info("Postupavatar   c.GetFile(avatar) err", err)

		resp := map[string]interface{}{
			"errno":  utils.RECODE_IOERR,
			"errmsg": utils.RecodeText(utils.RECODE_IOERR),
		}
		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			beego.Info(err)
			return
		}
		return
	}

	filebuffer := make([]byte, header.Size)
	_, err = file.Read(filebuffer)
	if err != nil {
		beego.Info("Postupavatar   file.Read(filebuffer) err", err)
		resp := map[string]interface{}{
			"errno":  utils.RECODE_IOERR,
			"errmsg": utils.RecodeText(utils.RECODE_IOERR),
		}
		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			beego.Info(err)
			return
		}
		return
	}

	cli := grpc.NewService()
	cli.Init()
	// call the backend service
	exampleClient := POSTHOUSESIMAGE.NewExampleService("go.micro.srv.PostHouseImage", cli.Client())
	rsp, err := exampleClient.PostHousesImage(context.TODO(), &POSTHOUSESIMAGE.Request{
		Sessionid: userlogin.Value,
		Id:        houseid,
		Image:     filebuffer,
		Filesize:  header.Size,
		Filename:  header.Filename,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//接收数据
	data := make(map[string]interface{})
	data["url"] = utils.AddDomain2Url(rsp.Url)
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   data,
	}

	w.Header().Set("Content-Type", "application/json")
	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
*/
