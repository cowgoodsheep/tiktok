package controller

import (
	"net/http"
	"strconv"
	"tiktok/logic"
	"tiktok/models"

	"github.com/gin-gonic/gin"
)

// 用户收藏视频
func PostUserFavorite(c *gin.Context) {
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

	//获取被收藏视频的ID
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
	getFavorite, err := logic.PostFavorite(u.ID, int(videoID))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"comment": getFavorite,
		"msg":     "收藏成功",
	})
}

// 查看本用户收藏的视频
func GetUserFavorite(c *gin.Context) {
	telephone, ok := c.Get("telephone")
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "用户手机号码解析错误",
		})
		return
	}
	var u models.User
	u.Telephone = telephone.(string)

	//在数据库中寻找该用户信息
	if err := models.FindUserByTelephone(&u); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	//在数据库中寻找该用户ID收藏的视频
	var favList []*models.Video
	if err := models.GetFavListByUserID(&favList, u.ID); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	//返回视频信息
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  favList,
	})
}

// 删除本用户收藏的视频
func DeleteUserFavorite(c *gin.Context) {
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

	//获取待删除收藏视频的ID
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

	//调用logic将收藏从数据库中删除
	getFavorite, err := logic.DeleteFavorite(u.ID, int(videoID))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"comment": getFavorite,
		"msg":     "删除成功",
	})
}
