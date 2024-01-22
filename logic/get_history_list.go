package logic

import (
	"errors"
	"tiktok/models"
)

type HistoryList struct {
	History_List []*models.History `json:"history_list"`
}

type HistoryListFlow struct {
	userID      int
	historyList []*models.History

	*HistoryList
}

func GetHistoryList(userID int) (*HistoryList, error) {
	return NewHistoryListFlow(userID).Do()
}

func NewHistoryListFlow(userID int) *HistoryListFlow {
	return &HistoryListFlow{userID: userID}
}

func (q *HistoryListFlow) Do() (*HistoryList, error) {
	var err error

	//检查该用户ID的合法性
	if !models.IsUserExistByUserID(q.userID) {
		return nil, errors.New("用户不存在或已注销")
	}

	//从数据库中获取历史浏览记录列表
	var historyList []*models.History
	err = models.GetHistoryListByUserID(q.userID, &historyList)
	if err != nil {
		return nil, err
	}

	//获取视频信息
	for _, v := range historyList {
		if err := models.GetVideoByID(v.Video_ID, &v.Video); err != nil {
			continue
		}
	}

	q.historyList = historyList
	q.HistoryList = &HistoryList{History_List: q.historyList}

	return q.HistoryList, nil
}
