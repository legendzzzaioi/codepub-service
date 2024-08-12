package model

import (
	"codepub-service/crypt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

// Mysql 数据模型
type Mysql struct {
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement;comment:'自增id'"`
	Name     string `json:"name" gorm:"type:varchar(255);unique;not null;comment:'名称'"`
	URL      string `json:"url" gorm:"type:varchar(255);not null;comment:'地址'"`
	Username string `json:"username" gorm:"type:varchar(255);not null;comment:'用户名'"`
	Password string `json:"password" gorm:"type:varchar(255);not null;comment:'密码'"`
}

// TableName 指定表名为 mysql_config
func (Mysql) TableName() string {
	return "mysql_config"
}

// InitMysqlDB 初始化数据库
func InitMysqlDB(db *gorm.DB) {
	_ = db.AutoMigrate(&Mysql{})
}

// CreateMysqlConfig 创建mysql_config
func CreateMysqlConfig(c *gin.Context, db *gorm.DB) {
	var mysql Mysql
	if err := c.ShouldBindJSON(&mysql); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// password进行加密
	encryptionKey := crypt.GetEncryptionKey()
	encryptedPassword, err := crypt.Encrypt(encryptionKey, mysql.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
		return
	}
	mysql.Password = encryptedPassword

	if err := db.Create(&mysql).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "name already exist"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, mysql)
}

// UpdateMysqlConfig 更新mysql_config
func UpdateMysqlConfig(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var mysql Mysql

	// 先获取现有的记录
	if err := db.First(&mysql, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Record not found"})
		return
	}

	// 绑定请求中的新数据
	var updatedData Mysql
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查密码是否为空
	if updatedData.Password == "" {
		// 如果密码字段为空，保持原有密码不变
		updatedData.Password = mysql.Password
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
	mysql.Name = updatedData.Name
	mysql.URL = updatedData.URL
	mysql.Username = updatedData.Username
	mysql.Password = updatedData.Password

	// 保存更新后的数据
	if err := db.Save(&mysql).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mysql)
}

// DeleteMysqlConfig 删除mysql_config
func DeleteMysqlConfig(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	result := db.Delete(&Mysql{}, id)
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

// ListMysqlConfig 列出所有mysql_config
func ListMysqlConfig(c *gin.Context, db *gorm.DB) {
	var mysql []Mysql
	db.Find(&mysql)
	c.JSON(http.StatusOK, mysql)
}

// GetMysqlUrlByName 获取单个mysql_config
func GetMysqlUrlByName(db *gorm.DB, name string) (mysql Mysql, err error) {
	if err := db.Where("name = ?", name).First(&mysql).Error; err != nil {
		return mysql, err
	}
	return mysql, nil
}
