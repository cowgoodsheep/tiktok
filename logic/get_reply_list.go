package logic

import (
	"errors"
	"fmt"
	"tiktok/models"
)

type ReplyList struct {
	Replys []*models.Reply `json:"reply_list"`
}

type ReplyListFlow struct {
	commentID int
	replys    []*models.Reply

	ReplyList *ReplyList
}

func GetReplyList(commentID int) (*ReplyList, error) {
	return NewReplyListFlow(commentID).Do()
}

func NewReplyListFlow(commentID int) *ReplyListFlow {
	return &ReplyListFlow{commentID: commentID}
}

func (n *ReplyListFlow) Do() (*ReplyList, error) {
	//检查评论ID是否有效
	if !models.IsCommentExistByCommentID(n.commentID) {
		return nil, fmt.Errorf("评论%d不存在或已经被删除", n.commentID)
	}

	//从数据库中获取回复列表
	err := models.GetReplyListByVideoID(n.commentID, &n.replys)
	if err != nil {
		return nil, err
	}

	//如果该评论没有回复
	if &n.replys == nil || len(n.replys) == 0 {
		return nil, errors.New("该评论暂时还没有人回复")
	}

	//整理发送返回
	n.ReplyList = &ReplyList{Replys: n.replys}
	return n.ReplyList, nil
}
