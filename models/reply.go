package models

import (
	"errors"
	"tiktok/dao"
	"time"

	"github.com/jinzhu/gorm"
)

// Reply Model
type Reply struct {
	ID          int       `json:"id"`
	User_ID     int       `json:"user_id"`
	Comment_ID  int       `json:"comment_id"`
	Reply_ID    int       `json:"reply_id"`
	Content     string    `json:"content"`
	Replyt_Time time.Time `json:"comment_time"`
}

// 指定Reply结构体迁移表reply
func (R *Reply) TableName() string {
	return "reply"
}

// 存入一条回复，并且给该视频评论回复数加一
func AddReplyAndUpdateCommentReplyCount(reply *Reply) error {
	if reply == nil {
		return errors.New("空指针错误")
	}

	//执行事务
	return dao.DB.Transaction(func(tx *gorm.DB) error {
		//添加评论数据
		if err := tx.Create(reply).Error; err != nil {
			// 返回任何错误都会回滚事务
			return err
		}

		//增加视频的评论数
		if err := tx.Exec("UPDATE comment c SET c.reply_count = c.reply_count+1 WHERE c.id=?", reply.Comment_ID).Error; err != nil {
			return err
		}

		// 返回 nil 提交事务
		return nil
	})
}

// 根据回复ID获取回复
func GetReplyById(replyID int, reply *Reply) error {
	if reply == nil {
		return errors.New("空指针错误")
	}
	return dao.DB.Where("id=?", replyID).First(reply).Error
}

// 删除回复，并且给视频评论回复数减一
func DeleteReplyAndUpdateCommentReplyCount(replyID, commentID int) error {
	//执行事务
	return dao.DB.Transaction(func(tx *gorm.DB) error {
		//删除回复数据
		if err := tx.Exec("DELETE FROM reply WHERE id = ?", replyID).Error; err != nil {
			// 返回任何错误都会回滚事务
			return err
		}

		//减少评论的回复数
		if err := tx.Exec("UPDATE comment c SET c.reply_count = c.reply_count-1 WHERE c.id = ? AND c.reply_count > 0", commentID).Error; err != nil {
			return err
		}

		// 返回 nil 提交事务
		return nil
	})
}

// 根据评论ID获取回复列表
func GetReplyListByVideoID(commentID int, reply *[]*Reply) error {
	if reply == nil {
		return errors.New("空指针错误")
	}
	if err := dao.DB.Model(&Reply{}).Where("comment_id=?", commentID).Find(reply).Error; err != nil {
		return err
	}
	return nil
}

// 根据评论ID删除回复
func DelectReplyByCommentID(commentID int) error {
	err := dao.DB.Delete(&Reply{}, "comment_id = ?", commentID).Error
	return err
}
