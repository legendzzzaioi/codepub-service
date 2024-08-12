package model

import (
	"codepub-service/crypt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

// Jenkins 数据模型
type Jenkins struct {
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement;comment:'自增id'"`
	Name     string `json:"name" gorm:"type:varchar(255);unique;not null;comment:'名称'"`
	URL      string `json:"url" gorm:"type:varchar(255);not null;comment:'地址'"`
	Username string `json:"username" gorm:"type:varchar(255);not null;comment:'用户名'"`
	Password string `json:"password" gorm:"type:varchar(255);not null;comment:'密码'"`
}

// TableName 指定表名为 jenkins_config
func (Jenkins) TableName() string {
	return "jenkins_config"
}

// InitJenkinsDB 初始化数据库
func InitJenkinsDB(db *gorm.DB) {
	_ = db.AutoMigrate(&Jenkins{})
}

// CreateJenkinsConfig 创建jenkins_config
func CreateJenkinsConfig(c *gin.Context, db *gorm.DB) {
	var jenkins Jenkins
	if err := c.ShouldBindJSON(&jenkins); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// password进行加密
	encryptionKey := crypt.GetEncryptionKey()
	encryptedPassword, err := crypt.Encrypt(encryptionKey, jenkins.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
		return
	}
	jenkins.Password = encryptedPassword

	if err := db.Create(&jenkins).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "name already exist"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, jenkins)
}

// UpdateJenkinsConfig 更新jenkins_config
func UpdateJenkinsConfig(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var jenkins Jenkins

	// 先获取现有的记录
	if err := db.First(&jenkins, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Record not found"})
		return
	}

	// 绑定请求中的新数据
	var updatedData Jenkins
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查密码是否为空
	if updatedData.Password == "" {
		// 如果密码字段为空，保持原有密码不变
		updatedData.Password = jenkins.Password
	} else {
		// 如果密码字段非空，则进行加密
		encryptionKey := crypt.GetEncryptionKey()
		encryptedPassword, err := crypt.Encrypt(encryptionKey, updatedData.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
			return
		}
		updatedData.Password = encryptedPassword
	}

	// 更新其他字段
	jenkins.Name = updatedData.Name
	jenkins.URL = updatedData.URL
	jenkins.Username = updatedData.Username
	jenkins.Password = updatedData.Password

	// 保存更新后的数据
	if err := db.Save(&jenkins).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, jenkins)
}

// DeleteJenkinsConfig 删除jenkins_config
func DeleteJenkinsConfig(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	result := db.Delete(&Jenkins{}, id)
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

// ListJenkinsConfig 列出所有jenkins_config
func ListJenkinsConfig(c *gin.Context, db *gorm.DB) {
	var jenkins []Jenkins
	db.Find(&jenkins)
	c.JSON(http.StatusOK, jenkins)
}

// GetJenkinsUrlByName 获取单个jenkins_config
func GetJenkinsUrlByName(db *gorm.DB, name string) (jenkins Jenkins, err error) {
	if err := db.Where("name = ?", name).First(&jenkins).Error; err != nil {
		return jenkins, err
	}
	return jenkins, nil
}
