package controller

import (
	"net/http"
	"tiktok/models"

	"github.com/gin-gonic/gin"
)

// 搜索视频
func SearchVideos(c *gin.Context) {

	//在数据库中模糊寻找title包含keyword的视频
	keyWord := c.PostForm("keyWord")
	var videoList []*models.Video
	if err := models.GetVideoListByKeyWord(&videoList, keyWord); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	//返回视频列表信息
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  videoList,
	})
}

// 搜索用户
func SearchUsers(c *gin.Context) {

	//在数据库中模糊寻找username包含keyword的用户
	keyWord := c.PostForm("keyWord")
	var userList []*models.User
	if err := models.GetUserListByKeyWord(&userList, keyWord); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	//返回用户列表信息
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  userList,
	})
}
