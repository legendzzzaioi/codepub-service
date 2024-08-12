package model

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

// Etcd 数据模型
type Etcd struct {
	ID   uint   `json:"id" gorm:"primaryKey;autoIncrement;comment:'自增id'"`
	Name string `json:"name" gorm:"type:varchar(255);unique;not null;comment:'名称'"`
	URL  string `json:"url" gorm:"type:varchar(255);not null;comment:'地址'"`
}

// TableName 指定表名为 etcd_config
func (Etcd) TableName() string {
	return "etcd_config"
}

// InitEtcdDB 初始化数据库
func InitEtcdDB(db *gorm.DB) {
	_ = db.AutoMigrate(&Etcd{})
}

// CreateEtcdConfig 创建etcd_config
func CreateEtcdConfig(c *gin.Context, db *gorm.DB) {
	var etcd Etcd
	if err := c.ShouldBindJSON(&etcd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Create(&etcd).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "name already exist"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, etcd)
}

// UpdateEtcdConfig 更新etcd_config
func UpdateEtcdConfig(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var etcd Etcd
	if err := db.First(&etcd, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Record not found"})
		return
	}
	if err := c.ShouldBindJSON(&etcd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Save(&etcd)
	c.JSON(http.StatusOK, etcd)
}

// DeleteEtcdConfig 删除etcd_config
func DeleteEtcdConfig(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	result := db.Delete(&Etcd{}, id)
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

// ListEtcdConfig 列出所有etcd_config
func ListEtcdConfig(c *gin.Context, db *gorm.DB) {
	var etcd []Etcd
	db.Find(&etcd)
	c.JSON(http.StatusOK, etcd)
}

// GetEtcdUrlByName 获取单个etcd_config
func GetEtcdUrlByName(db *gorm.DB, name string) (url string, err error) {
	var etcd Etcd
	if err := db.Where("name = ?", name).First(&etcd).Error; err != nil {
		return "", err
	}
	return etcd.URL, nil
}
