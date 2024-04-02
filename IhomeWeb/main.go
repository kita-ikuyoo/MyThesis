package main

import (
	Ihandler "IHome/IhomeWeb/handler"
	_ "IHome/IhomeWeb/model"
	"github.com/go-micro/plugins/v4/registry/consul"
	"github.com/julienschmidt/httprouter"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/web"
	"log"
	"net/http"
)

func main() {
	consulRegistry := consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{
			"127.0.0.1:8500", // 这里假设您的Consul服务运行在本机的8500端口
		}
	})
	// 构造web服务
	service := web.NewService(
		web.Name("go.micro.web.IhomeWeb"),
		web.Version("latest"),
		web.Address(":22333"),
		web.Registry(consulRegistry),
	)
	// 服务初始化
	if err := service.Init(); err != nil {
		log.Fatal(err)
	}
	//构建路由
	rou := httprouter.New()
	//将路由注册到服务
	// http.Dir("html")创建了一个http.FileSystem的实现，它以"html"目录作为根目录。这意味着服务器将在这个目录下寻找请求的文件。
	rou.NotFound = http.FileServer(http.Dir("html"))
	rou.GET("/api/v1.0/areas", Ihandler.GetArea)
	//欺骗浏览器  session index
	rou.GET("/api/v1.0/session", Ihandler.GetSession)
	//session
	rou.GET("/api/v1.0/house/index", Ihandler.GetIndex)
	//获取图片验证码
	rou.GET("/api/v1.0/imagecode/:uuid", Ihandler.GetImageCd)
	//获取短信验证码
	rou.GET("/api/v1.0/smscode/:mobile", Ihandler.Getsmscd)
	//用户注册
	rou.POST("/api/v1.0/users", Ihandler.PostRet)
	//用户登陆
	rou.POST("/api/v1.0/sessions", Ihandler.PostLogin)
	//退出登陆   注意退出登录是session而用户登录是sessions
	rou.DELETE("/api/v1.0/session", Ihandler.DeleteSession)
	//获取用户详细信息
	rou.GET("/api/v1.0/user", Ihandler.GetUserInfo)
	//用户上传图片
	rou.POST("/api/v1.0/user/avatar", Ihandler.PostAvatar)
	////请求更新用户名
	//rou.PUT("/api/v1.0/user/name", Ihandler.PutUserInfo)
	//身份认证检查 同  获取用户信息   所调用的服务是 GetUserInfo
	rou.GET("/api/v1.0/user/auth", Ihandler.GetUserAuth)
	//实名认证服务   这里一个是posti一个是get所以可以写一样的url
	rou.POST("/api/v1.0/user/auth", Ihandler.PostUserAuth)
	//获取用户已发布房源信息服务
	rou.GET("/api/v1.0/user/houses", Ihandler.GetUserHouses)
	//发送（发布）房源信息服务
	rou.POST("/api/v1.0/houses", Ihandler.PostHouses)
	////发送（上传）房屋图片服务
	//rou.POST("/api/v1.0/houses/:id/images", handler.PostHouseImage)
	// 根据Handle源码来看 我需要一个http.Handler
	// 再进一步查看，发现http.Handler就是实现了type Handler interface {
	//	ServeHTTP(ResponseWriter, *Request)
	//}的东西。
	service.Handle("/", rou)

	// 服务运行
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
