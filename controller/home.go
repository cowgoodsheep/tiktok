package controller

import (
	"net/http"
	"strconv"
	"tiktok/logic"
	"tiktok/middleware"
	"tiktok/models"
	"time"

	"github.com/gin-gonic/gin"
)

// 主页显示
func Home(c *gin.Context) {
	//获取token
	token := c.Query("token")
	//如果找不到，就去Form表单去找
	if token == "" {
		token = c.PostForm("token")
	}

	//检查是否登录
	if len(token) == 0 {
		//未登录的视频推流
		intTime, err := strconv.ParseInt(c.PostForm("lastest_time"), 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}
		lastestTime := time.Unix(0, intTime*1e6) //将时间转化成纳秒

		videoList, err := logic.GetVideoList(0, lastestTime)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  videoList,
		})
	} else {
		//已登录的视频推流
		claim, ok := middleware.ParseToken(token)
		if ok == false {
			c.JSON(http.StatusOK, gin.H{"error": "token解析失败"})
			return
		}

		//token过期
		if time.Now().Unix() > claim.ExpiresAt {
			c.JSON(http.StatusOK, gin.H{"error": "token已失效"})
			return
		}

		intTime, err := strconv.ParseInt(c.PostForm("lastest_time"), 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}
		lastestTime := time.Unix(0, intTime*1e6) //将时间转化成纳秒

		//获取用户ID
		userID, err := models.FindUserIDByTelephone(claim.Telephone)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}

		//获取视频推流
		videoList, err := logic.GetVideoList(userID, lastestTime)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  videoList,
		})
	}
}
