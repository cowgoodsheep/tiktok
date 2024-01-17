package models

import (
	"errors"
	"tiktok/dao"
	"time"

	"github.com/jinzhu/gorm"
)

// Comment Model
type Comment struct {
	ID             int       `json:"id"`
	User_ID        int       `json:"user_id"`
	Video_ID       int       `json:"video_id"`
	Content        string    `json:"content"`
	Comment_Time   time.Time `json:"comment_time"`
	Like_counts    int       `json:"like_counts"`
	Disgust_counts int       `json:"disgust_counts"`
	Reply_count    int       `json:"reply_count"`
}

// 指定Comment结构体迁移表comment
func (C *Comment) TableName() string {
	return "comment"
}

// 存入一条评论，并且给视频评论数加一
func AddCommentAndUpdateVideoCommentCount(comment *Comment) error {
	if comment == nil {
		return errors.New("空指针错误")
	}

	//执行事务
	return dao.DB.Transaction(func(tx *gorm.DB) error {
		//添加评论数据
		if err := tx.Create(comment).Error; err != nil {
			// 返回任何错误都会回滚事务
			return err
		}

		//增加视频的评论数
		if err := tx.Exec("UPDATE video v SET v.comment_count = v.comment_count+1 WHERE v.id=?", comment.Video_ID).Error; err != nil {
			return err
		}

		// 返回 nil 提交事务
		return nil
	})
}

// 根据评论ID获取评论
func GetCommentById(commentID int, comment *Comment) error {
	if comment == nil {
		return errors.New("空指针错误")
	}
	return dao.DB.Where("id=?", commentID).First(comment).Error
}

// 删除评论，并且给视频评论数减一
func DeleteCommentAndUpdateVideoCommentCount(commentID, videoID int) error {
	//执行事务
	return dao.DB.Transaction(func(tx *gorm.DB) error {
		//删除评论数据
		if err := tx.Exec("DELETE FROM comment WHERE id = ?", commentID).Error; err != nil {
			// 返回任何错误都会回滚事务
			return err
		}

		//减少视频的评论数
		if err := tx.Exec("UPDATE video v SET v.comment_count = v.comment_count-1 WHERE v.id = ? AND v.comment_count > 0", videoID).Error; err != nil {
			return err
		}

		// 返回 nil 提交事务
		return nil
	})
}

// 根据视频ID获取评论列表
func GetCommentListByVideoID(videoID int, comments *[]*Comment) error {
	if comments == nil {
		return errors.New("空指针错误")
	}
	if err := dao.DB.Model(&Comment{}).Where("video_id=?", videoID).Find(comments).Error; err != nil {
		return err
	}
	return nil
}

// 根据视频ID删除评论
func DelectCommentByVideoID(videoID int) error {
	err := dao.DB.Delete(&Comment{}, "video_id = ?", videoID).Error
	return err
}

// 根据ID判断该评论是否存在
func IsCommentExistByCommentID(commentID int) bool {
	var c Comment
	dao.DB.Where("ID=?", commentID).First(&c)
	if c.ID == 0 {
		return false
	} else {
		return true
	}
}
