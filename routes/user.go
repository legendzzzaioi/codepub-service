package routes

import (
	"codepub-service/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUserRoutes(r *gin.Engine, db *gorm.DB) {
	// 获取用户列表
	r.GET("/api/v1/user", func(c *gin.Context) {
		model.GetUser(c, db)
	})
	// 创建
	r.POST("/api/v1/user", func(c *gin.Context) {
		model.CreateUser(c, db)
	})
	// 更新
	r.PUT("/api/v1/user/:id", func(c *gin.Context) {
		model.UpdateUser(c, db)
	})
	// 删除
	r.DELETE("/api/v1/user/:id", func(c *gin.Context) {
		model.DeleteUser(c, db)
	})
	// 获取当前用户信息
	r.GET("/api/v1/user_info", func(c *gin.Context) {
		model.GetUserInfo(c, db)
	})
}
