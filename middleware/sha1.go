package middleware

import (
	"crypto/sha1"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	MinPasswordSize = 6
	MaxPasswordSize = 20
)

// SHA1加密
func SHA1(s string) string {
	//创建一个SHA1哈希对象
	o := sha1.New()
	//将输入的字符串转化成字节数组，写入o
	o.Write([]byte(s))
	//计算并返回SHA1哈希值
	return hex.EncodeToString(o.Sum(nil))
}

// SHA1加密用户密码中间件
func SHAMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		//从请求的查询参数中获取password
		password := c.Query("password")
		//如果找不到，就去Form表单去找
		if password == "" {
			password = c.PostForm("password")
		}
		if len(password) < MinPasswordSize || len(password) > MaxPasswordSize {
			c.JSON(http.StatusOK, gin.H{"err": "密码长度小于或大于限制"})
			//密码输入出错，跳过后续操作
			c.Abort()
			return
		}
		//对传入的password进行SHA1加密，并存入下文中
		c.Set("password", SHA1(password))
		c.Next()
	}
}
