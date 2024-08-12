package model

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

// Nacos 数据模型
type Nacos struct {
	ID   uint   `json:"id" gorm:"primaryKey;autoIncrement;comment:'自增id'"`
	Name string `json:"name" gorm:"type:varchar(255);unique;not null;comment:'名称'"`
	URL  string `json:"url" gorm:"type:varchar(255);not null;comment:'地址'"`
}

// TableName 指定表名为 nacos_config
func (Nacos) TableName() string {
	return "nacos_config"
}

// InitNacosDB 初始化数据库
func InitNacosDB(db *gorm.DB) {
	_ = db.AutoMigrate(&Nacos{})
}

// CreateNacosConfig 创建Nacos_config
func CreateNacosConfig(c *gin.Context, db *gorm.DB) {
	var nacos Nacos
	if err := c.ShouldBindJSON(&nacos); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Create(&nacos).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "name already exist"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nacos)
}

// UpdateNacosConfig 更新Nacos_config
func UpdateNacosConfig(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var nacos Nacos
	if err := db.First(&nacos, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Record not found"})
		return
	}
	if err := c.ShouldBindJSON(&nacos); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Save(&nacos)
	c.JSON(http.StatusOK, nacos)
}

// DeleteNacosConfig 删除Nacos_config
func DeleteNacosConfig(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	result := db.Delete(&Nacos{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Record not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Record deleted"})
}

// ListNacosConfig 列出所有Nacos_config
func ListNacosConfig(c *gin.Context, db *gorm.DB) {
	var nacos []Nacos
	db.Find(&nacos)
	c.JSON(http.StatusOK, nacos)
}

// GetNacosUrlByName 获取单个Nacos_config
func GetNacosUrlByName(db *gorm.DB, name string) (url string, err error) {
	var nacos Nacos
	if err := db.Where("name = ?", name).First(&nacos).Error; err != nil {
		return "", err
	}
	return nacos.URL, nil
}
