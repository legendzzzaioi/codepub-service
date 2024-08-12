package controllers

import (
	"codepub-service/crypt"
	"codepub-service/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
)

func ExecMysqlSql(c *gin.Context, db *gorm.DB) {
	name := c.Param("name")
	sql := c.PostForm("sql")
	// 获取 Mysql URL、Username、Password
	Mysql, err := model.GetMysqlUrlByName(db, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// password解密
	if Mysql.Password != "" {
		encryptionKey := crypt.GetEncryptionKey()
		decryptedPassword, err := crypt.Decrypt(encryptionKey, Mysql.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
			return
		}
		Mysql.Password = decryptedPassword
	}

	connStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		Mysql.Username, Mysql.Password, Mysql.URL, "information_schema")

	conn, err := gorm.Open(mysql.Open(connStr), &gorm.Config{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	result := conn.Exec(sql)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "true"})
}
