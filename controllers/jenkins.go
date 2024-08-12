package controllers

import (
	"codepub-service/crypt"
	"codepub-service/model"
	"context"
	"encoding/json"
	"net/http"

	"github.com/bndr/gojenkins"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetJenkinsAllView 获取jenkins所有视图
func GetJenkinsAllView(c *gin.Context, db *gorm.DB) {
	name := c.Param("name")
	// 获取jenkins url
	Jenkins, err := model.GetJenkinsUrlByName(db, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	// password解密
	if Jenkins.Password != "" {
		encryptionKey := crypt.GetEncryptionKey()
		decryptedPassword, err := crypt.Decrypt(encryptionKey, Jenkins.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
			return
		}
		Jenkins.Password = decryptedPassword
	}

	// 初始化 Jenkins 客户端
	jenkins := gojenkins.CreateJenkins(nil, Jenkins.URL, Jenkins.Username, Jenkins.Password)
	_, err = jenkins.Init(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// 获取所有视图
	views, err := jenkins.GetAllViews(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// 创建一个视图名称列表
	viewNames := make([]string, len(views))
	for i, view := range views {
		viewNames[i] = view.GetName()
	}

	// 将视图名称列表转换为 JSON
	viewNamesJSON, err := json.Marshal(viewNames)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// 处理响应
	c.Data(http.StatusOK, "application/json; charset=utf-8", viewNamesJSON)
}

// GetJenkinsJobsByView 获取jenkins指定视图下的所有jobs，其中Query参数：viewName
func GetJenkinsJobsByView(c *gin.Context, db *gorm.DB) {
	name := c.Param("name")
	viewName := c.Query("viewName")

	// 获取jenkins url
	Jenkins, err := model.GetJenkinsUrlByName(db, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	// password解密
	if Jenkins.Password != "" {
		encryptionKey := crypt.GetEncryptionKey()
		decryptedPassword, err := crypt.Decrypt(encryptionKey, Jenkins.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
			return
		}
		Jenkins.Password = decryptedPassword
	}

	// 初始化 Jenkins 客户端
	jenkins := gojenkins.CreateJenkins(nil, Jenkins.URL, Jenkins.Username, Jenkins.Password)
	_, err = jenkins.Init(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// 获取指定视图
	view, err := jenkins.GetView(context.Background(), viewName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// 获取视图下的所有 jobs
	jobs := view.GetJobs()

	// 将视图名称列表转换为 JSON
	jobNamesJSON, err := json.Marshal(jobs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// 处理响应
	c.Data(http.StatusOK, "application/json; charset=utf-8", jobNamesJSON)
}

// GetJenkinsJobBuildParam 获取job的构建参数，其中Query参数：jobName
func GetJenkinsJobBuildParam(c *gin.Context, db *gorm.DB) {
	name := c.Param("name")
	jobName := c.Query("jobName")

	// 获取jenkins url
	Jenkins, err := model.GetJenkinsUrlByName(db, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	// password解密
	if Jenkins.Password != "" {
		encryptionKey := crypt.GetEncryptionKey()
		decryptedPassword, err := crypt.Decrypt(encryptionKey, Jenkins.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
			return
		}
		Jenkins.Password = decryptedPassword
	}

	// 初始化 Jenkins 客户端
	jenkins := gojenkins.CreateJenkins(nil, Jenkins.URL, Jenkins.Username, Jenkins.Password)
	_, err = jenkins.Init(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// 获取job对象
	job, err := jenkins.GetJob(context.Background(), jobName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// 获取 Job 的参数列表
	parameters, err := job.GetParameters(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// 过滤掉 type 为 "WHideParameterDefinition" 的隐藏参数
	var filteredParameters []gojenkins.ParameterDefinition
	for _, param := range parameters {
		if param.Type != "WHideParameterDefinition" {
			filteredParameters = append(filteredParameters, param)
		}
	}

	// 将过滤后的参数列表转换为 JSON
	paramListJSON, err := json.Marshal(filteredParameters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// 处理响应
	c.Data(http.StatusOK, "application/json; charset=utf-8", paramListJSON)
}

// BuildJenkinsJob 参数化构建job，其中json数据：jobName,params{"param1": "value1",.....}
func BuildJenkinsJob(c *gin.Context, db *gorm.DB) {
	name := c.Param("name")
	var request struct {
		JobName string            `json:"jobName"`
		Params  map[string]string `json:"params"`
	}

	// 解析 JSON 请求体
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request payload"})
		return
	}

	jobName := request.JobName
	params := request.Params

	// 获取jenkins url
	Jenkins, err := model.GetJenkinsUrlByName(db, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	// password解密
	if Jenkins.Password != "" {
		encryptionKey := crypt.GetEncryptionKey()
		decryptedPassword, err := crypt.Decrypt(encryptionKey, Jenkins.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
			return
		}
		Jenkins.Password = decryptedPassword
	}

	// 初始化 Jenkins 客户端
	jenkins := gojenkins.CreateJenkins(nil, Jenkins.URL, Jenkins.Username, Jenkins.Password)
	_, err = jenkins.Init(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// 获取job对象
	job, err := jenkins.GetJob(context.Background(), jobName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// 参数化构建job
	_, err = job.InvokeSimple(context.Background(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// 处理响应
	c.JSON(http.StatusOK, gin.H{"message": "true"})
}
