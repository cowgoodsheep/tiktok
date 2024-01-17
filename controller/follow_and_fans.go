package controller

import (
	"net/http"
	"strconv"
	"tiktok/logic"
	"tiktok/models"

	"github.com/gin-gonic/gin"
)

// 用户关注其他用户
func PostUserFollow(c *gin.Context) {
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

	//获取关注的用户ID
	stringFollowUserID := c.Query("follow_user_id")
	//如果找不到，就去Form表单去找
	if stringFollowUserID == "" {
		stringFollowUserID = c.PostForm("follow_user_id")
	}
	//将关注用户ID转化成int
	videoFollowUserIDID, err := strconv.ParseInt(stringFollowUserID, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	//调用logic存入数据库,并写入缓存，加速读过程
	getFollow, err := logic.PostFollow(u.ID, int(videoFollowUserIDID))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"comment": getFollow,
		"msg":     "关注成功",
	})
}

// 用户查看已关注用户
func GetUserFollow(c *gin.Context) {
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

	//调用logic得到关注列表
	followList, err := logic.GetFollowList(u.ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":        200,
		"follow_list": followList,
		"msg":         "查询成功",
	})
}

// 用户取关已关注用户
func DeleteUserFollow(c *gin.Context) {
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

	//获取待取消关注的用户ID
	stringFollowUserID := c.Query("follow_user_id")
	//如果找不到，就去Form表单去找
	if stringFollowUserID == "" {
		stringFollowUserID = c.PostForm("follow_user_id")
	}
	//将待取消关注用户ID转化成int
	videoFollowUserIDID, err := strconv.ParseInt(stringFollowUserID, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	//调用logic存入数据库,并写入缓存，加速读过程
	getFollow, err := logic.DeleteFollow(u.ID, int(videoFollowUserIDID))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"comment": getFollow,
		"msg":     "取消关注成功",
	})
}

// 用户查看粉丝
func GetUserFans(c *gin.Context) {
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

	//调用logic得到粉丝列表
	fansList, err := logic.GetFansList(u.ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  fansList,
	})
}
