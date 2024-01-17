package logic

import (
	"fmt"
	"tiktok/cache"
	"tiktok/models"
)

type FollowFlow struct {
	userID       int
	followUserID int
}

// 添加关注
func PostFollow(userID, followUserID int) (*models.Follow, error) {
	return NewPostFollowFlow(userID, followUserID).Do()
}

func NewPostFollowFlow(userID, followUserID int) *FollowFlow {
	return &FollowFlow{userID: userID, followUserID: followUserID}
}

func (p *FollowFlow) Do() (*models.Follow, error) {
	var err error

	//检查所关注用户ID是否正确
	if !models.IsUserExistByUserID(p.followUserID) {
		return nil, fmt.Errorf("所关注用户%d不存在", p.userID)
	}
	//不允许关注自己
	if p.userID == p.followUserID {
		return nil, fmt.Errorf("不能关注自己")
	}

	//检查该用户是否已被关注
	if models.IsFollowExistByUserIDAndFollowUserID(p.userID, p.followUserID) {
		return nil, fmt.Errorf("请勿重复关注")
	}

	//整理数据，上传数据库
	follow := models.Follow{User_ID: p.userID, Follow_User_ID: p.followUserID}
	err = models.AddFollow(&follow)
	if err != nil {
		return nil, err
	}

	//更新redis的关注信息
	cache.UpdateUserFollow(p.userID, p.followUserID, true)

	return &follow, nil
}

// 取消关注
func DeleteFollow(userID, followUserID int) (*models.Follow, error) {
	//先获取follow
	var follow models.Follow
	err := models.GetFollowByUserIDAndFollowUserID(userID, followUserID, &follow)
	if err != nil {
		return nil, err
	}

	//检查该用户是否已被关注
	if models.IsFollowExistByUserIDAndFollowUserID(userID, followUserID) == false {
		return nil, fmt.Errorf("不能取关未关注用户")
	}

	//删除follow
	err = models.DeletefollowByID(follow.ID)
	if err != nil {
		return nil, err
	}

	//更新redis的关注信息
	cache.UpdateUserFollow(userID, followUserID, false)

	return &follow, nil
}
