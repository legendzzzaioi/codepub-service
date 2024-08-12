package model

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

// User 数据模型
type User struct {
	Id       uint   `json:"id" gorm:"primaryKey;autoIncrement;comment:'自增id'"`
	Username string `json:"username" gorm:"type:varchar(255);unique;not null;comment:'用户名'"`
	Password string `json:"password" gorm:"type:varchar(255);not null;comment:'密码'"`
}

// UserDTO 返回给前端的数据
type UserDTO struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
}

// InitUserDB 初始化数据库
func InitUserDB(db *gorm.DB) {
	// 自动迁移 User 模型
	_ = db.AutoMigrate(&User{})

	// 检查是否存在用户名为 admin 的用户
	var user User
	result := db.Where("username = ?", "admin").First(&user)

	// 如果没有找到用户名为 admin 的用户，则创建一个默认用户
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		password, err := hashPassword("admin")
		if err != nil {
			// 如果加密失败，直接返回
			return
		}
		adminUser := User{
			Username: "admin",
			Password: password,
		}
		db.Create(&adminUser)
	}
}

// hashPassword 密码hash加密
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// checkPasswordHash 密码hash校验
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// CreateUser 创建user
func CreateUser(c *gin.Context, db *gorm.DB) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// password进行加密
	password, err := hashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
		return
	}
	user.Password = password

	if err := db.Create(&user).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "name already exist"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// UpdateUser 更新user
func UpdateUser(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var user User

	// 先获取现有的记录
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Record not found"})
		return
	}

	// 绑定请求中的新数据
	var updatedData User
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查密码是否为空
	if updatedData.Password == "" {
		// 如果密码字段为空，保持原有密码不变
		updatedData.Password = user.Password
	} else {
		// 如果密码字段非空，则进行加密
		password, err := hashPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
			return
		}
		updatedData.Password = password
	}

	// 更新其他字段
	user.Username = updatedData.Username
	user.Password = updatedData.Password

	// 保存更新后的数据
	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser 删除user
func DeleteUser(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	result := db.Delete(&User{}, id)
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

// GetUser 列出所有user
func GetUser(c *gin.Context, db *gorm.DB) {
	var users []User
	db.Find(&users)

	var userDTOs []UserDTO
	for _, user := range users {
		userDTOs = append(userDTOs, UserDTO{
			Id:       user.Id,
			Username: user.Username,
		})
	}

	c.JSON(http.StatusOK, userDTOs)
}

// Login 登陆
func Login(c *gin.Context, db *gorm.DB) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	password := user.Password

	// 查询数据库
	result := db.Where("username = ?", user.Username).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username not found"})
		return
	}

	// 校验密码
	if !checkPasswordHash(password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}
	// 登录成功，设置会话
	session := sessions.Default(c)
	session.Set("user_id", user.Id)
	_ = session.Save()

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

// Logout 登出
func Logout(c *gin.Context) {
	// 清除会话
	session := sessions.Default(c)
	session.Clear()
	_ = session.Save()

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

// GetUserInfo 获取当前用户信息
func GetUserInfo(c *gin.Context, db *gorm.DB) {
	session := sessions.Default(c)
	userID := session.Get("user_id")

	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var user User
	result := db.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username not found"})
		return
	}

	userDTO := UserDTO{
		Id:       user.Id,
		Username: user.Username,
	}

	c.JSON(http.StatusOK, userDTO)
}
