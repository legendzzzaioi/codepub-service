package routes

import (
	"codepub-service/controllers"
	"codepub-service/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterMysqlRoutes(r *gin.Engine, db *gorm.DB) {
	// --------------------------------mysql api-------------------------------------
	// 通过name获取对应的mysql地址，并执行sql，其中表单参数：sql
	r.POST("/api/v1/mysql/sql/:name", func(c *gin.Context) {
		controllers.ExecMysqlSql(c, db)
	})

	// --------------------------------mysql表-------------------------------------
	// 获取mysql_config表中的配置列表
	r.GET("/api/v1/mysql_config/list", func(c *gin.Context) {
		model.ListMysqlConfig(c, db)
	})
	// 新增mysql_config表中的配置，提交字段name、url、username、password
	r.POST("/api/v1/mysql_config/list", func(c *gin.Context) {
		model.CreateMysqlConfig(c, db)
	})
	// 通过id更新mysql_config表中的配置，提交字段name、url、username、password
	r.PUT("/api/v1/mysql_config/:id", func(c *gin.Context) {
		model.UpdateMysqlConfig(c, db)
	})
	// 通过id删除mysql_config表中的配置
	r.DELETE("/api/v1/mysql_config/:id", func(c *gin.Context) {
		model.DeleteMysqlConfig(c, db)
	})
}
