package controllers

import (
	"codepub-service/crypt"
	"codepub-service/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

func ExecPostgresSql(c *gin.Context, db *gorm.DB) {
	name := c.Param("name")
	sql := c.PostForm("sql")
	// 获取 Postgres URL、Username、Password
	Postgres, err := model.GetPostgresUrlByName(db, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// password解密
	if Postgres.Password != "" {
		encryptionKey := crypt.GetEncryptionKey()
		decryptedPassword, err := crypt.Decrypt(encryptionKey, Postgres.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
			return
		}
		Postgres.Password = decryptedPassword
	}

	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable",
		strings.Split(Postgres.URL, ":")[0], Postgres.Username, Postgres.Password, strings.Split(Postgres.URL, ":")[1])

	conn, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
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
