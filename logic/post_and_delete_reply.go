package logic

import (
	"fmt"
	"tiktok/models"
	"time"
)

type ReplyFlow struct {
	userID    int
	commentID int
	replyID   int
	content   string
}

// 上传评论
func PostReply(userID, commentID, replyID int, content string) (*models.Reply, error) {
	return NewPostReplyFlow(userID, commentID, replyID, content).Do()
}

func NewPostReplyFlow(userID, commentID, replyID int, content string) *ReplyFlow {
	return &ReplyFlow{userID: userID, commentID: commentID, replyID: replyID, content: content}
}

func (p *ReplyFlow) Do() (*models.Reply, error) {
	var err error

	//检查ID是否正确
	if !models.IsUserExistByUserID(p.userID) {
		return nil, fmt.Errorf("用户%d不存在", p.userID)
	}
	//检查评论ID是否正确
	if !models.IsCommentExistByCommentID(p.commentID) {
		return nil, fmt.Errorf("评论%d不存在", p.commentID)
	}

	//检查一下（记得删）
	fmt.Println(p.replyID)

	//整理数据，上传数据库
	reply := models.Reply{User_ID: p.userID, Comment_ID: p.commentID, Reply_ID: p.replyID, Content: p.content, Replyt_Time: time.Now()}
	err = models.AddReplyAndUpdateCommentReplyCount(&reply)
	if err != nil {
		return nil, err
	}

	return &reply, nil
}

// 删除回复
func DeleteReply(replyID, commentID, userID int) (*models.Reply, error) {
	//先获取reply
	var reply models.Reply
	err := models.GetReplyById(replyID, &reply)
	if err != nil {
		return nil, err
	}

	//检查这条回复的合法性
	//如果回复ID和用户ID不对应
	if reply.User_ID != userID {
		return nil, fmt.Errorf("回复ID与用户ID不匹配")
	}
	//如果回复ID和评论ID不对应
	if reply.Comment_ID != commentID {
		return nil, fmt.Errorf("回复ID与评论ID不匹配")
	}

	//删除reply
	err = models.DeleteReplyAndUpdateCommentReplyCount(replyID, commentID)
	if err != nil {
		return nil, err
	}

	return &reply, nil
}
