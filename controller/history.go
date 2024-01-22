package controller

import (
	"net/http"
	"strconv"
	"tiktok/logic"
	"tiktok/models"

	"github.com/gin-gonic/gin"
)

// 查看本用户的历史浏览记录
func GetUserHistory(c *gin.Context) {
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

	historyList, err := logic.GetHistoryList(u.ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	//返回历史浏览记录信息
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  historyList,
	})
}

// 删除本用户的历史浏览记录（可批量删除）
func DeleteUserHistory(c *gin.Context) {
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

	//获取待删除历史浏览记录ID
	stringHistoryID := c.Query("history_id")
	//如果找不到，就去Form表单去找
	if stringHistoryID == "" {
		stringHistoryID = c.PostForm("history_id")
	}
	//将历史浏览记录ID转化成int
	historyID, err := strconv.ParseInt(stringHistoryID, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	//调用logic将历史浏览记录从数据库中删除
	getHistory, err := logic.DeleteHistory(u.ID, int(historyID))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"comment": getHistory,
		"msg":     "删除成功",
	})
}
