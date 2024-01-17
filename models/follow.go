package models

import (
	"errors"
	"tiktok/dao"
)

// Follow Model
type Follow struct {
	ID             int `json:"id"`
	User_ID        int `json:"user_id"`
	Follow_User_ID int `json:"follow_user_id"`
}

// 指定Follow结构体迁移表follow
func (f *Follow) TableName() string {
	return "follow"
}

// 添加关注
func AddFollow(f *Follow) error {
	if f == nil {
		return errors.New("空指针错误")
	}
	return dao.DB.Create(f).Error
}

// 根据用户ID和被关注用户ID检查该用户是否已被关注
func IsFollowExistByUserIDAndFollowUserID(userID, followUserID int) bool {
	var f Follow
	dao.DB.Where("user_id = ? AND follow_user_id = ?", userID, followUserID).First(&f)
	if f.ID == 0 {
		return false
	} else {
		return true
	}
}

// 根据用户ID和被关注用户ID获取关注信息
func GetFollowByUserIDAndFollowUserID(userID, followUserID int, follow *Follow) error {
	if follow == nil {
		return errors.New("空指针错误")
	}
	err := dao.DB.Where("user_id = ? AND follow_user_id = ?", userID, followUserID).Find(&follow).Error
	return err
}

// 根据关注ID取消关注
func DeletefollowByID(followID int) error {
	err := dao.DB.Delete(&Follow{}, followID).Error
	return err
}

// 根据用户ID查询用户关注列表
func GetFollowListByUserID(userID int, userList *[]*User) error {
	if userList == nil {
		return errors.New("空指针错误")
	}
	var err error
	if err = dao.DB.Raw("SELECT u.* FROM follow f, user u WHERE f.user_id = ? AND f.follow_user_id = u.id", userID).Scan(userList).Error; err != nil {
		return err
	}
	if len(*userList) == 0 || (*userList)[0].ID == 0 {
		return errors.New("关注列表为空")
	}
	return nil
}

// 根据用户ID查询用户粉丝列表
func GetFansListByUserID(userID int, userList *[]*User) error {
	if userList == nil {
		return errors.New("空指针错误")
	}
	var err error
	if err = dao.DB.Raw("SELECT u.* FROM follow f, user u WHERE f.follow_user_id = ? AND f.user_id = u.id", userID).Scan(userList).Error; err != nil {
		return err
	}
	if len(*userList) == 0 || (*userList)[0].ID == 0 {
		return errors.New("粉丝列表为空")
	}
	return nil
}
