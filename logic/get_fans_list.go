package logic

import (
	"errors"
	"tiktok/cache"
	"tiktok/models"
)

type FansList struct {
	Fans_User_List []*models.User `json:"fans_user_list"`
}

type FansListFlow struct {
	userID   int
	userList []*models.User

	*FansList
}

func GetFansList(userID int) (*FansList, error) {
	return NewFansListFlow(userID).Do()
}

func NewFansListFlow(userID int) *FansListFlow {
	return &FansListFlow{userID: userID}
}

func (q *FansListFlow) Do() (*FansList, error) {
	var err error

	//检查该用户ID的合法性
	if !models.IsUserExistByUserID(q.userID) {
		return nil, errors.New("用户不存在或已注销")
	}

	//从数据库中获取粉丝列表
	err = models.GetFansListByUserID(q.userID, &q.userList)
	if err != nil {
		return nil, err
	}

	//从redis中获取is_follow
	for _, v := range q.userList {
		v.Is_Follow = cache.GetUserFollow(q.userID, v.ID)
	}

	q.FansList = &FansList{Fans_User_List: q.userList}

	return q.FansList, nil
}
