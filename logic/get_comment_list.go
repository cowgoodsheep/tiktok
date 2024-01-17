package logic

import (
	"errors"
	"fmt"
	"tiktok/models"
)

type CommentList struct {
	Comments []*models.Comment `json:"comment_list"`
}

type CommentListFlow struct {
	videoID  int
	comments []*models.Comment

	CommentList *CommentList
}

func GetCommentList(videoID int) (*CommentList, error) {
	return NewCommentListFlow(videoID).Do()
}

func NewCommentListFlow(videoID int) *CommentListFlow {
	return &CommentListFlow{videoID: videoID}
}

func (n *CommentListFlow) Do() (*CommentList, error) {
	//检查视频ID是否有效
	if !models.IsVideoExistByVideoID(n.videoID) {
		return nil, fmt.Errorf("视频%d不存在或已经被删除", n.videoID)
	}

	//从数据库中获取评论列表
	err := models.GetCommentListByVideoID(n.videoID, &n.comments)
	if err != nil {
		return nil, err
	}

	//如果该视频没有评论
	if &n.comments == nil || len(n.comments) == 0 {
		return nil, errors.New("该视频暂时还没有人评论")
	}

	//整理发送返回
	n.CommentList = &CommentList{Comments: n.comments}
	return n.CommentList, nil
}
