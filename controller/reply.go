package controller

import (
	"net/http"
	"strconv"
	"tiktok/logic"
	"tiktok/models"

	"github.com/gin-gonic/gin"
)

// 用户给评论发表回复
func PostCommentReply(c *gin.Context) {
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

	//获取被回复评论的ID
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

	//如果回复的评论也是一条回复
	stringReplyID := c.Query("reply_id")
	//如果找不到，就去Form表单去找
	if stringReplyID == "" {
		stringReplyID = c.PostForm("reply_id")
	}
	//如果获取不到回复ID
	if stringReplyID == "" {
		stringReplyID = "0"
	}
	//将回复ID转化成int
	replyID, err := strconv.ParseInt(stringReplyID, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	//获取回复内容
	content := c.Query("content")
	//如果找不到，就去Form表单去找
	if content == "" {
		content = c.PostForm("content")
	}

	//调用logic存入数据库
	getReply, err := logic.PostReply(u.ID, int(commentID), int(replyID), content)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"comment": getReply,
		"msg":     "上传成功",
	})
}

// 用户删除视频回复
func DeleteCommentReply(c *gin.Context) {
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

	//获取被删除回复评论的ID
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

	//获取被删除回复评论的ID
	stringReplyID := c.Query("reply_id")
	//如果找不到，就去Form表单去找
	if stringReplyID == "" {
		stringReplyID = c.PostForm("reply_id")
	}
	//将回复ID转化成int
	replyID, err := strconv.ParseInt(stringReplyID, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	//调用logic将回复从数据库中删除
	getReply, err := logic.DeleteReply(int(replyID), int(commentID), u.ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"comment": getReply,
		"msg":     "删除成功",
	})
}

// 获取视频回复列表
func GetReplyList(c *gin.Context) {
	//获取被回复评论的ID
	stringCommentID := c.Query("comment_id")
	//如果找不到，就去Form表单去找
	if stringCommentID == "" {
		stringCommentID = c.PostForm("comment_id")
	}
	//将视频ID转化成int
	commentID, err := strconv.ParseInt(stringCommentID, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	//调用logic从数据库中寻找该评论的回复列表
	replyList, err := logic.GetReplyList(int(commentID))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"comment": replyList,
		"msg":     "获取成功",
	})
}
