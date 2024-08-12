package controllers

import (
	"codepub-service/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// GetNacosAllNamespace 通过name获取对应nacos url下的所有namespace列表
func GetNacosAllNamespace(c *gin.Context, db *gorm.DB) {
	name := c.Param("name")
	// 获取nacos url
	nacosUrl, err := model.GetNacosUrlByName(db, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	// 获取namespace列表
	resp, err := http.Get(nacosUrl + "/v1/console/namespaces")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch namespaces: " + err.Error(),
		})
		return
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	// 处理响应
	body, _ := io.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json; charset=utf-8", body)
}

// GetNacosAllConfig 通过name获取对应nacos url下的所有config列表，query参数：pageSize、tenant
func GetNacosAllConfig(c *gin.Context, db *gorm.DB) {
	name := c.Param("name")
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	tenant := c.Query("tenant")
	// 获取nacos url
	nacosUrl, err := model.GetNacosUrlByName(db, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	nacosUrl = fmt.Sprintf("%s/v1/cs/configs?dataId=&group=&appName=&config_tags=&pageNo=1&pageSize=%d&tenant=%s&search=accurate", nacosUrl, pageSize, tenant)

	// 获取config列表
	resp, err := http.Get(nacosUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch namespaces: " + err.Error(),
		})
		return
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	// 处理响应
	body, _ := io.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json; charset=utf-8", body)
}

// SaveNacosConfig 新增或修改nacos配置，表单参数：tenant、dataId、group、content、type
func SaveNacosConfig(c *gin.Context, db *gorm.DB) {
	name := c.Param("name")
	tenant := c.PostForm("tenant")
	dataId := c.PostForm("dataId")
	group := c.PostForm("group")
	content := c.PostForm("content")
	type_ := c.PostForm("type")

	// 获取 Nacos URL
	nacosUrl, err := model.GetNacosUrlByName(db, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	nacosUrl = fmt.Sprintf("%s/v1/cs/configs", nacosUrl)

	// 替换 content 中的 \n 为换行符
	content = strings.ReplaceAll(content, "\\n", "\n")

	// 发送 POST 请求
	formData := url.Values{
		"dataId":  {dataId},
		"group":   {group},
		"content": {content},
		"type":    {type_},
		"tenant":  {tenant},
	}
	resp, err := http.Post(nacosUrl, "application/x-www-form-urlencoded", strings.NewReader(formData.Encode()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	// 处理响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to read response body",
		})
		return
	}

	if string(body) != "true" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": string(body),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "true"})
}

// DeleteNacosConfig 删除nacos配置，query参数：tenant、dataId、group
func DeleteNacosConfig(c *gin.Context, db *gorm.DB) {
	name := c.Param("name")
	tenant := c.Query("tenant")
	dataId := c.Query("dataId")
	group := c.Query("group")
	// 获取nacos url
	nacosUrl, err := model.GetNacosUrlByName(db, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	nacosUrl = fmt.Sprintf("%s/v1/cs/configs?dataId=%s&group=%s&tenant=%s", nacosUrl, dataId, group, tenant)

	// 创建 DELETE 请求
	req, err := http.NewRequest("DELETE", nacosUrl, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	// 处理响应
	body, _ := io.ReadAll(resp.Body)
	if string(body) != "true" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": string(body),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "true"})
}
