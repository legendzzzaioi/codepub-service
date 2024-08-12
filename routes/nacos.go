package routes

import (
	"codepub-service/controllers"
	"codepub-service/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterNacosRoutes(r *gin.Engine, db *gorm.DB) {
	// --------------------------------nacos api-------------------------------------
	// 通过name获取对应nacos地址的所有namespace
	r.GET("/api/v1/nacos/namespace/:name", func(c *gin.Context) {
		controllers.GetNacosAllNamespace(c, db)
	})
	// 通过name获取对应nacos地址的所有config，其中query参数：pageSize、tenant
	r.GET("/api/v1/nacos/config/:name", func(c *gin.Context) {
		controllers.GetNacosAllConfig(c, db)
	})
	// 通过name创建或修改对应nacos地址的config，其中表单参数：：tenant、dataId、group、content、type
	r.POST("/api/v1/nacos/config/:name", func(c *gin.Context) {
		controllers.SaveNacosConfig(c, db)
	})
	// 通过name删除对应nacos地址的config，其中query参数：：tenant、dataId、group
	r.DELETE("/api/v1/nacos/config/:name", func(c *gin.Context) {
		controllers.DeleteNacosConfig(c, db)
	})

	// --------------------------------nacos表-------------------------------------
	// 获取nacos_config表中的配置列表
	r.GET("/api/v1/nacos_config/list", func(c *gin.Context) {
		model.ListNacosConfig(c, db)
	})
	// 新增nacos_config表中的配置，提交字段name、url
	r.POST("/api/v1/nacos_config/list", func(c *gin.Context) {
		model.CreateNacosConfig(c, db)
	})
	// 通过id更新nacos_config表中的配置，提交字段name、url
	r.PUT("/api/v1/nacos_config/:id", func(c *gin.Context) {
		model.UpdateNacosConfig(c, db)
	})
	// 通过id删除nacos_config表中的配置
	r.DELETE("/api/v1/nacos_config/:id", func(c *gin.Context) {
		model.DeleteNacosConfig(c, db)
	})
}
