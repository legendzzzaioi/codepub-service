package routes

import (
	"codepub-service/controllers"
	"codepub-service/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterEtcdRoutes(r *gin.Engine, db *gorm.DB) {
	// --------------------------------etcd api-------------------------------------
	// 通过name获取对应etcd地址的所有keys
	r.GET("/api/v1/etcd/key/:name", func(c *gin.Context) {
		controllers.GetEtcdAllKeys(c, db)
	})
	// 通过name获取对应etcd地址的value，其中query参数：key
	r.GET("/api/v1/etcd/value/:name", func(c *gin.Context) {
		controllers.GetEtcdValueByKey(c, db)
	})
	// 通过name创建或修改对应etcd地址的config，其中表单参数：：key、value
	r.POST("/api/v1/etcd/config/:name", func(c *gin.Context) {
		controllers.SaveEtcdValueByKey(c, db)
	})
	// 通过name删除对应etcd地址的config，其中query参数：：key
	r.DELETE("/api/v1/etcd/key/:name", func(c *gin.Context) {
		controllers.DeleteEtcdValueByKey(c, db)
	})

	// --------------------------------etcd表-------------------------------------
	// 获取etcd_config表中的配置列表
	r.GET("/api/v1/etcd_config/list", func(c *gin.Context) {
		model.ListEtcdConfig(c, db)
	})
	// 新增etcd_config表中的配置，提交字段name、url
	r.POST("/api/v1/etcd_config/list", func(c *gin.Context) {
		model.CreateEtcdConfig(c, db)
	})
	// 通过id更新etcd_config表中的配置，提交字段name、url
	r.PUT("/api/v1/etcd_config/:id", func(c *gin.Context) {
		model.UpdateEtcdConfig(c, db)
	})
	// 通过id删除etcd_config表中的配置
	r.DELETE("/api/v1/etcd_config/:id", func(c *gin.Context) {
		model.DeleteEtcdConfig(c, db)
	})
}
