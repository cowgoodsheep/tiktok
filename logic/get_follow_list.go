package logic

import (
	"errors"
	"tiktok/models"
)

type FollowList struct {
	Follow_User_List []*models.User `json:"follow_user_list"`
}

type FollowListFlow struct {
	userID   int
	userList []*models.User

	*FollowList
}

func GetFollowList(userID int) (*FollowList, error) {
	return NewFollowListFlow(userID).Do()
}

func NewFollowListFlow(userID int) *FollowListFlow {
	return &FollowListFlow{userID: userID}
}

func (q *FollowListFlow) Do() (*FollowList, error) {
	var err error

	//检查该用户ID的合法性
	if !models.IsUserExistByUserID(q.userID) {
		return nil, errors.New("用户不存在或已注销")
	}

	//从数据库中获取关注列表
	var userList []*models.User
	err = models.GetFollowListByUserID(q.userID, &userList)
	if err != nil {
		return nil, err
	}

	for _, v := range userList {
		//将isFollow标记为已关注
		v.Is_Follow = true
	}

	q.userList = userList
	q.FollowList = &FollowList{Follow_User_List: q.userList}

	return q.FollowList, nil
}
