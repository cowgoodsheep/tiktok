package controller

import (
	"net/http"
	"strconv"
	"tiktok/logic"
	"tiktok/models"

	"github.com/gin-gonic/gin"
)

// 用户点赞视频
func PostUserLike(c *gin.Context) {
	//解析token 得到用户手机号
	telephone, ok := c.Get("telephone")
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "用户手机号码解析错误",
		})
		return
	}

	//根据用户手机号得到用户信息
	var u models.User
	u.Telephone = telephone.(string)
	if err := models.FindUserByTelephone(&u); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	//获取被点赞视频的ID
	stringVideoID := c.Query("video_id")
	//如果找不到，就去Form表单去找
	if stringVideoID == "" {
		stringVideoID = c.PostForm("video_id")
	}
	//将视频ID转化成int
	videoID, err := strconv.ParseInt(stringVideoID, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	//调用logic存入数据库
	err = logic.PostLike(u.ID, int(videoID))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":     200,
		"user_id":  u.ID,
		"video_id": videoID,
		"msg":      "点赞成功",
	})
}

// 用户取消点赞视频
func DeleteUserLike(c *gin.Context) {
	//解析token 得到用户手机号
	telephone, ok := c.Get("telephone")
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "用户手机号码解析错误",
		})
		return
	}

	//根据用户手机号得到用户信息
	var u models.User
	u.Telephone = telephone.(string)
	if err := models.FindUserByTelephone(&u); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	//获取取消点赞视频的ID
	stringVideoID := c.Query("video_id")
	//如果找不到，就去Form表单去找
	if stringVideoID == "" {
		stringVideoID = c.PostForm("video_id")
	}
	//将视频ID转化成int
	videoID, err := strconv.ParseInt(stringVideoID, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	//调用logic存入数据库
	err = logic.DeleteLike(u.ID, int(videoID))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":     200,
		"user_id":  u.ID,
		"video_id": videoID,
		"msg":      "取消点赞成功",
	})
}
