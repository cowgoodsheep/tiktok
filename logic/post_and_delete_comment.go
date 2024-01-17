package logic

import (
	"fmt"
	"tiktok/models"
	"time"
)

type CommentFlow struct {
	userID  int
	videoID int
	content string
}

// 上传评论
func PostComment(userID, videoID int, content string) (*models.Comment, error) {
	return NewPostCommentFlow(userID, videoID, content).Do()
}

func NewPostCommentFlow(userID, videoID int, content string) *CommentFlow {
	return &CommentFlow{userID: userID, videoID: videoID, content: content}
}

func (p *CommentFlow) Do() (*models.Comment, error) {
	var err error

	//检查ID是否正确
	if !models.IsUserExistByUserID(p.userID) {
		return nil, fmt.Errorf("用户%d不存在", p.userID)
	}
	if !models.IsVideoExistByVideoID(p.videoID) {
		return nil, fmt.Errorf("视频%d不存在", p.videoID)
	}

	//整理数据，上传数据库
	comment := models.Comment{User_ID: p.userID, Video_ID: p.videoID, Content: p.content, Comment_Time: time.Now()}
	err = models.AddCommentAndUpdateVideoCommentCount(&comment)
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

// 删除评论
func DeleteComment(commentID, videoID, userID int) (*models.Comment, error) {
	//先获取comment
	var comment models.Comment
	err := models.GetCommentById(commentID, &comment)
	if err != nil {
		return nil, err
	}

	//检查这条评论的合法性
	//如果评论ID和用户ID不对应
	if comment.User_ID != userID {
		return nil, fmt.Errorf("评论ID与用户ID不匹配")
	}
	//如果评论ID和视频ID不对应
	if comment.Video_ID != videoID {
		return nil, fmt.Errorf("评论ID与视频ID不匹配")
	}

	//删除comment
	err = models.DeleteCommentAndUpdateVideoCommentCount(commentID, videoID)
	if err != nil {
		return nil, err
	}

	//删除该评论的回复
	if err := models.DelectReplyByCommentID(commentID); err != nil {
		return nil, err
	}

	return &comment, nil
}
