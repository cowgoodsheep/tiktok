package controller

import (
	"net/http"
	"strconv"
	"tiktok/logic"
	"tiktok/models"

	"github.com/gin-gonic/gin"
)

// 用户给视频发表评论
func PostVideoComment(c *gin.Context) {
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

	//获取被评论视频的ID
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

	//获取评论内容
	content := c.Query("content")
	//如果找不到，就去Form表单去找
	if content == "" {
		content = c.PostForm("content")
	}

	//调用logic存入数据库
	getComment, err := logic.PostComment(u.ID, int(videoID), content)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"comment": getComment,
		"msg":     "上传成功",
	})
}

// 用户删除视频评论
func DeleteVideoComment(c *gin.Context) {
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

	//获取被评论视频的ID
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

	//获取评论的ID
	stringCommentID := c.Query("comment_id")
	//如果找不到，就去Form表单去找
	if stringCommentID == "" {
		stringCommentID = c.PostForm("comment_id")
	}
	//将评论ID转化成int
	commentID, err := strconv.ParseInt(stringCommentID, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	//调用logic将评论从数据库中删除
	getComment, err := logic.DeleteComment(int(commentID), int(videoID), u.ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"comment": getComment,
		"msg":     "删除成功",
	})
}

// 获取视频评论列表
func GetCommentList(c *gin.Context) {
	//获取被评论视频的ID
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

	//调用logic从数据库中寻找该视频的评论列表
	commentList, err := logic.GetCommentList(int(videoID))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"comment": commentList,
		"msg":     "获取成功",
	})
}
