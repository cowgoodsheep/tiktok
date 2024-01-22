package logic

import (
	"fmt"
	"tiktok/models"
	"time"
)

type HistoryFlow struct {
	userID  int
	videoID int
}

// 添加历史浏览记录
func PostHistory(userID, videoID int) error {
	return NewPostHistoryFlow(userID, videoID).Do()
}

func NewPostHistoryFlow(userID, videoID int) *HistoryFlow {
	return &HistoryFlow{userID: userID, videoID: videoID}
}

func (p *HistoryFlow) Do() error {
	var err error

	//检查用户ID是否正确
	if !models.IsUserExistByUserID(p.userID) {
		return fmt.Errorf("所关注用户%d不存在", p.userID)
	}

	//整理数据，上传数据库
	history := models.History{User_ID: p.userID, Video_ID: p.videoID, History_Time: time.Now()}
	err = models.AddHistory(&history)
	if err != nil {
		return err
	}

	return nil
}

// 删除历史浏览记录
func DeleteHistory(userID, historyID int) (*models.History, error) {
	//检查用户ID是否正确
	if !models.IsUserExistByUserID(userID) {
		return nil, fmt.Errorf("用户%d不存在", userID)
	}

	//先获取history
	var history models.History
	err := models.GetHistoryByID(historyID, &history)
	if err != nil {
		return nil, err
	}

	//删除history
	err = models.DeleteHistoryByID(history.ID)
	if err != nil {
		return nil, err
	}

	return &history, nil
}
