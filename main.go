package main

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/lyun0ne/webook/internal/repository"
	"github.com/lyun0ne/webook/internal/repository/dao"
	"github.com/lyun0ne/webook/internal/service"
	"github.com/lyun0ne/webook/internal/web"
	"github.com/lyun0ne/webook/internal/web/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db := initDB()
	u := initUser(db)

	server := initWebserver()
	u.RegisterRoutes(server)

	server.Run() // 监听并在 0.0.0.0:8080上启动服务
}

func initWebserver() *gin.Engine {
	server := gin.Default()
	server.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		//不写AllowMethod 就是所有方法都允许
		AllowMethods:  []string{"PUT", "PATCH", "POST"},
		AllowHeaders:  []string{"Origin", "content-type"},
		ExposeHeaders: []string{"Content-Length"},
		//是否允许带cookie一类的东西
		AllowCredentials: true,
		//通过一个方法允许origin是否允许
		AllowOriginFunc: func(origin string) bool {
			if strings.Contains(origin, "http://localhost") {
				return true
			}
			if strings.Contains(origin, "实际环境") {
				return true
			}
			return false
		},
		MaxAge: 12 * time.Hour,
	}))

	store := cookie.NewStore([]byte("secret"))
	server.Use(sessions.Sessions("webookId", store))

	server.Use(middleware.
		NewLoginMiddlewareBuilder().
		IngorePaths("/users/login").
		IngorePaths("/users/signup").
		Build())
	return server
}

func initUser(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDao(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	return u
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	if err != nil {
		panic(err)
	}

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}

	return db
}
