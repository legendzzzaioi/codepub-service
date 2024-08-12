package routes

import (
	"codepub-service/controllers"
	"codepub-service/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterPostgresRoutes(r *gin.Engine, db *gorm.DB) {
	// --------------------------------postgres api-------------------------------------
	// 通过name获取对应的postgres地址，并执行sql，其中表单参数：sql
	r.POST("/api/v1/postgres/sql/:name", func(c *gin.Context) {
		controllers.ExecPostgresSql(c, db)
	})

	// --------------------------------postgres表-------------------------------------
	// 获取postgres_config表中的配置列表
	r.GET("/api/v1/postgres_config/list", func(c *gin.Context) {
		model.ListPostgresConfig(c, db)
	})
	// 新增postgres_config表中的配置，提交字段name、url、username、password
	r.POST("/api/v1/postgres_config/list", func(c *gin.Context) {
		model.CreatePostgresConfig(c, db)
	})
	// 通过id更新postgres_config表中的配置，提交字段name、url、username、password
	r.PUT("/api/v1/postgres_config/:id", func(c *gin.Context) {
		model.UpdatePostgresConfig(c, db)
	})
	// 通过id删除postgres_config表中的配置
	r.DELETE("/api/v1/postgres_config/:id", func(c *gin.Context) {
		model.DeletePostgresConfig(c, db)
	})
}
