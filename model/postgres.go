package model

import (
	"codepub-service/crypt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

// Postgres 数据模型
type Postgres struct {
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement;comment:'自增id'"`
	Name     string `json:"name" gorm:"type:varchar(255);unique;not null;comment:'名称'"`
	URL      string `json:"url" gorm:"type:varchar(255);not null;comment:'地址'"`
	Username string `json:"username" gorm:"type:varchar(255);not null;comment:'用户名'"`
	Password string `json:"password" gorm:"type:varchar(255);not null;comment:'密码'"`
}

// TableName 指定表名为 postgres_config
func (Postgres) TableName() string {
	return "postgres_config"
}

// InitPostgresDB 初始化数据库
func InitPostgresDB(db *gorm.DB) {
	_ = db.AutoMigrate(&Postgres{})
}

// CreatePostgresConfig 创建postgres_config
func CreatePostgresConfig(c *gin.Context, db *gorm.DB) {
	var postgres Postgres
	if err := c.ShouldBindJSON(&postgres); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// password进行加密
	encryptionKey := crypt.GetEncryptionKey()
	encryptedPassword, err := crypt.Encrypt(encryptionKey, postgres.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
		return
	}
	postgres.Password = encryptedPassword

	if err := db.Create(&postgres).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "name already exist"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, postgres)
}

// UpdatePostgresConfig 更新postgres_config
func UpdatePostgresConfig(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var postgres Postgres

	// 先获取现有的记录
	if err := db.First(&postgres, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Record not found"})
		return
	}

	// 绑定请求中的新数据
	var updatedData Postgres
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查密码是否为空
	if updatedData.Password == "" {
		// 如果密码字段为空，保持原有密码不变
		updatedData.Password = postgres.Password
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
	postgres.Name = updatedData.Name
	postgres.URL = updatedData.URL
	postgres.Username = updatedData.Username
	postgres.Password = updatedData.Password

	// 保存更新后的数据
	if err := db.Save(&postgres).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, postgres)
}

// DeletePostgresConfig 删除postgres_config
func DeletePostgresConfig(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	result := db.Delete(&Postgres{}, id)
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

// ListPostgresConfig 列出所有postgres_config
func ListPostgresConfig(c *gin.Context, db *gorm.DB) {
	var postgres []Postgres
	db.Find(&postgres)
	c.JSON(http.StatusOK, postgres)
}

// GetPostgresUrlByName 获取单个postgres_config
func GetPostgresUrlByName(db *gorm.DB, name string) (postgres Postgres, err error) {
	if err := db.Where("name = ?", name).First(&postgres).Error; err != nil {
		return postgres, err
	}
	return postgres, nil
}
