package models

import (
	"errors"
	"tiktok/dao"
	"time"
)

// History Model
type History struct {
	ID           int       `json:"id"`
	User_ID      int       `json:"user_id"`
	Video_ID     int       `json:"video_id"`
	History_Time time.Time `json:"history_time"`

	//以下字段不需要存储进数据库
	Video Video `json:"video" gorm:"-"`
}

// 指定History结构体迁移表history
func (h *History) TableName() string {
	return "history"
}

// 添加历史浏览记录
func AddHistory(h *History) error {
	if h == nil {
		return errors.New("空指针错误")
	}
	return dao.DB.Create(h).Error
}

// 根据用户ID和视频ID获取历史浏览记录信息
func GetHistoryByUserIDAndVideoID(userID, videoID int, history *History) error {
	if history == nil {
		return errors.New("空指针错误")
	}
	err := dao.DB.Where("user_id = ? AND video_id = ?", userID, videoID).Find(&history).Error
	return err
}

// 根据历史浏览记录ID获取历史浏览记录信息
func GetHistoryByID(historyID int, history *History) error {
	if history == nil {
		return errors.New("空指针错误")
	}
	err := dao.DB.Where("id = ?", historyID).Find(&history).Error
	return err
}

// 根据历史浏览记录ID删除历史浏览记录
func DeleteHistoryByID(historyID int) error {
	err := dao.DB.Delete(&History{}, historyID).Error
	return err
}

// 根据用户ID查询用户历史浏览记录列表
func GetHistoryListByUserID(userID int, historyList *[]*History) error {
	if historyList == nil {
		return errors.New("空指针错误")
	}
	var err error
	if err = dao.DB.Raw("SELECT h.* FROM history h WHERE h.user_id = ?", userID).Scan(historyList).Error; err != nil {
		return err
	}
	if len(*historyList) == 0 || (*historyList)[0].ID == 0 {
		return errors.New("历史浏览记录列表为空")
	}
	return nil
}
