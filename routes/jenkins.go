package routes

import (
	"codepub-service/controllers"
	"codepub-service/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterJenkinsRoutes(r *gin.Engine, db *gorm.DB) {
	// --------------------------------jenkins api-------------------------------------
	// 通过name获取对应jenkins地址的所有视图
	r.GET("/api/v1/jenkins/view/:name", func(c *gin.Context) {
		controllers.GetJenkinsAllView(c, db)
	})
	// 通过name获取对应jenkins地址下指定视图下的所有jobs，其中query参数：viewName
	r.GET("/api/v1/jenkins/jobs/:name", func(c *gin.Context) {
		controllers.GetJenkinsJobsByView(c, db)
	})
	// 通过name获取对应jenkins地址，获取job的构建参数，其中Query参数：jobName
	r.GET("/api/v1/jenkins/job_param/:name", func(c *gin.Context) {
		controllers.GetJenkinsJobBuildParam(c, db)
	})
	// 通过name获取对应jenkins地址，参数化构建job，其中表单参数：jobName,params
	r.PUT("/api/v1/jenkins/job/:name", func(c *gin.Context) {
		controllers.BuildJenkinsJob(c, db)
	})

	// --------------------------------jenkins表-------------------------------------
	// 获取jenkins_config表中的配置列表
	r.GET("/api/v1/jenkins_config/list", func(c *gin.Context) {
		model.ListJenkinsConfig(c, db)
	})
	// 新增jenkins_config表中的配置，提交字段name、url、username、password
	r.POST("/api/v1/jenkins_config/list", func(c *gin.Context) {
		model.CreateJenkinsConfig(c, db)
	})
	// 通过id更新jenkins_config表中的配置，提交字段name、url、username、password
	r.PUT("/api/v1/jenkins_config/:id", func(c *gin.Context) {
		model.UpdateJenkinsConfig(c, db)
	})
	// 通过id删除jenkins_config表中的配置
	r.DELETE("/api/v1/jenkins_config/:id", func(c *gin.Context) {
		model.DeleteJenkinsConfig(c, db)
	})
}
