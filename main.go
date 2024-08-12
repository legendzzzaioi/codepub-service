package main

import (
	"codepub-service/middleware"
	"codepub-service/model"
	"codepub-service/routes"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 初始化Viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// 从配置文件读取数据库连接信息
	user := viper.GetString("database.user")
	password := viper.GetString("database.password")
	host := viper.GetString("database.host")
	port := viper.GetInt("database.port")
	name := viper.GetString("database.name")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, name)

	var err error
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 初始化user表
	model.InitUserDB(db)
	// 初始化nacos_config表
	model.InitNacosDB(db)
	// 初始化etcd_config表
	model.InitEtcdDB(db)
	// 初始化mysql_config表
	model.InitMysqlDB(db)
	// 初始化postgres_config表
	model.InitPostgresDB(db)
	// 初始化jenkins_config表
	model.InitJenkinsDB(db)

	// 创建Gin路由
	r := gin.Default()

	// 设置会话存储
	store := cookie.NewStore([]byte("secret"))

	// 配置会话选项
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 24 * 7, // 会话有效期设置为 7 天（秒数）
		HttpOnly: true,          // 仅限 HTTP，防止客户端脚本访问
	})

	r.Use(sessions.Sessions("session", store))

	// 登录路由，不需要会话认证
	r.POST("/api/v1/login", func(c *gin.Context) {
		model.Login(c, db)
	})

	// 其他所有路由都需要会话认证
	r.Use(middleware.AuthMiddleware())

	// 注册各模块的路由
	routes.RegisterUserRoutes(r, db)
	routes.RegisterNacosRoutes(r, db)
	routes.RegisterEtcdRoutes(r, db)
	routes.RegisterMysqlRoutes(r, db)
	routes.RegisterPostgresRoutes(r, db)
	routes.RegisterJenkinsRoutes(r, db)

	// 登出路由
	r.POST("/api/v1/logout", func(c *gin.Context) {
		model.Logout(c)
	})

	//r.GET("/", func(c *gin.Context) {})

	// 运行服务器
	err = r.Run(":8000")
	if err != nil {
		return
	}

}
