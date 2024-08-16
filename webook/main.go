package main

import (
	"basic-go/webook/config"
	"basic-go/webook/internal/service/sms/tencent"
	"github.com/gin-gonic/gin"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"net/http"
)

func main() {
	server := InitWebServer()
	////
	////db := initDB()
	//u := initUser(db)

	//web.RegisterRoutes(server, u)
	//server := gin.Default()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello world")
	})
	err := server.Run(":8081")
	if err != nil {
		return
	}
}

//func initWebServer() *gin.Engine {
//	server := gin.Default()
//	server.Use(cors.New(cors.Config{
//		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
//		AllowHeaders:     []string{"Content-Type", "Authorization"},
//		AllowCredentials: true,
//		ExposeHeaders:    []string{"Token"},
//		//第一种方式
//		//AllowOrigins: []string{"http://localhost:8081"},
//		//第二种方式
//
//		AllowOriginFunc: func(origin string) bool {
//			if strings.HasPrefix(origin, "http://localhost") ||
//				strings.HasPrefix(origin, "https://live.webook.com") {
//				return true
//			}
//			fmt.Println("origin:", origin)
//			return false
//		},
//		MaxAge: 12 * time.Hour,
//	}))
//	server.Use(func(context *gin.Context) {
//		fmt.Println("这是第一个路由")
//	})
//
//	//redisClient := redis.NewClient(&redis.Options{
//	//	Addr:     config.Config.Redisaddr.Addr,
//	//	Password: "root",
//	//	DB:       0,
//	//})
//	//server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())
//	//步骤1
//	//store1 := &mysqlmdb.Store{}
//	//store := memstore.NewStore([]byte("pI4(fR0}mB4]fS2*eR3:lL9[iG1*qH8#"), []byte("oN1)tV1{xA6#xM2/nR5/hU1#fH2$bU0$"))
//	//store, _ := redis.NewStore(16, "tcp", "localhost:16379", "",
//	//	[]byte("pI4(fR0}mB4]fS2*eR3:lL9[iG1*qH8#"), []byte("oN1)tV1{xA6#xM2/nR5/hU1#fH2$bU0$"))
//	//store := cookie.NewStore([]byte("secret"))
//	//server.Use(sessions.Sessions("mysession", store))
//
//	server.Use(middleware.NewLoginJwtMiddlewareBuilder().IgnorePath("/user/login").
//		IgnorePath("/user/signup").IgnorePath("/user/login/jwt").
//		IgnorePath("/user/login_sms/code/send").
//		IgnorePath("/user/login_sms/code/verify").Build())
//
//	return server
//}

//func initUser(db *gorm.DB) *web.UserHandler {
//	ud := dao.NewUserDao(db)
//	redisClient := redis.NewClient(&redis.Options{
//		Addr:     config.Config.Redisaddr.Addr,
//		Password: "",
//		DB:       0,
//	})
//
//	ch := cache.NewUserCache(redisClient)
//	repo := repository.NewUserRepository(ud, ch)
//	svc := service.NewUserService(repo)
//	codeCache := cache.NewCodeCache(redisClient)
//	codeRepository := repository.NewCodeRepository(codeCache)
//	smsService := memory.NewService()
//	//smsService := initTenCentClient()
//	codeService := service.NewCodeService(codeRepository, smsService)
//	u := web.NewUserHandler(svc, codeService)
//
//	return u
//}

//func initDB() *gorm.DB {
//	dsn := config.Config.DB.DSN
//	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
//	if err != nil {
//		panic(err)
//	}
//	err = dao.InitTable(db)
//	if err != nil {
//		panic(err)
//	}
//	return db
//}

func initTenCentClient() *tencent.Service {
	credential := common.NewCredential(config.SecretId, config.SecretKey)
	smsTxService, err := sms.NewClient(credential, "ap-nanjing", profile.NewClientProfile())
	if err != nil {
		panic(err)
	}
	return tencent.NewService(config.SdkAPPId, config.SigName, smsTxService)
}
