package controllers

import (
	"codepub-service/model"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	clientv3 "go.etcd.io/etcd/client/v3"
	"gorm.io/gorm"
)

// GetEtcdAllKeys 获取所有的keys
func GetEtcdAllKeys(c *gin.Context, db *gorm.DB) {
	name := c.Param("name")
	// 获取etcd url
	etcdUrl, err := model.GetEtcdUrlByName(db, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	// 连接到etcd
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdUrl},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	defer func(cli *clientv3.Client) {
		_ = cli.Close()
	}(cli)
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 获取所有keys
	resp, err := cli.Get(ctx, "", clientv3.WithPrefix(), clientv3.WithKeysOnly())
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	// 遍历所有keys
	var keys []string
	for _, kv := range resp.Kvs {
		keys = append(keys, string(kv.Key))
	}
	body, _ := json.Marshal(keys)
	c.Data(http.StatusOK, "application/json; charset=utf-8", body)
}

// GetEtcdValueByKey 通过key获取指定的value
func GetEtcdValueByKey(c *gin.Context, db *gorm.DB) {
	name := c.Param("name")
	key := c.Query("key")
	// 获取etcd url
	etcdUrl, err := model.GetEtcdUrlByName(db, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	// 连接到etcd
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdUrl},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	defer func(cli *clientv3.Client) {
		_ = cli.Close()
	}(cli)
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 获取指定key的value
	resp, err := cli.Get(ctx, key)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	// 没有找到指定key
	if len(resp.Kvs) == 0 {
		c.JSON(500, gin.H{"message": "key not found"})
		return
	}
	// 返回指定key的value
	c.JSON(200, gin.H{
		"value": string(resp.Kvs[0].Value),
	})
}

// SaveEtcdValueByKey 新增或更新key、value
func SaveEtcdValueByKey(c *gin.Context, db *gorm.DB) {
	name := c.Param("name")
	key := c.PostForm("key")
	value := c.PostForm("value")
	// 替换 value 中的 \n 为换行符
	value = strings.ReplaceAll(value, "\\n", "\n")
	// 获取etcd url
	etcdUrl, err := model.GetEtcdUrlByName(db, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	// 连接到etcd
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdUrl},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	defer func(cli *clientv3.Client) {
		_ = cli.Close()
	}(cli)
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 提交key、value，新增或更新
	_, err = cli.Put(ctx, key, value)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	// 返回
	c.JSON(http.StatusOK, gin.H{"message": "true"})
}

// DeleteEtcdValueByKey 删除指定的key
func DeleteEtcdValueByKey(c *gin.Context, db *gorm.DB) {
	name := c.Param("name")
	key := c.Query("key")
	// 获取etcd url
	etcdUrl, err := model.GetEtcdUrlByName(db, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	// 连接到etcd
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdUrl},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	defer func(cli *clientv3.Client) {
		_ = cli.Close()
	}(cli)
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 删除key
	_, err = cli.Delete(ctx, key)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	// 返回
	c.JSON(http.StatusOK, gin.H{"message": "true"})
}
