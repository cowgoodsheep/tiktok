package models

import (
	"errors"
	"tiktok/dao"
)

// User Model
type User struct {
	ID          int    `json:"id"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Avatar      string `json:"avatar"`
	Description string `json:"description"`
	Telephone   string `json:"telephone"`
	Is_deleted  int    `json:"is_deleted"`

	//以下字段不需要存储进数据库
	Is_Follow bool `json:"is_follow" gorm:"-"`
}

// 指定User结构体迁移表user
func (u *User) TableName() string {
	return "user"
}

func AddUser(u *User) error {
	if u == nil {
		return errors.New("空指针错误")
	}
	return dao.DB.Create(u).Error
}

// 根据手机号判断该用户是否存在
func IsUserExistByTelephone(telephone string) bool {
	var u User
	dao.DB.Where("telephone=?", telephone).First(&u)
	//如果找不到或者用户已注销
	if u.ID == 0 || u.Is_deleted == 1 {
		return false
	} else {
		return true
	}
}

// 根据ID判断该用户是否存在
func IsUserExistByUserID(userID int) bool {
	var u User
	dao.DB.Where("ID=?", userID).First(&u)
	//如果找不到或者用户已注销
	if u.ID == 0 || u.Is_deleted == 1 {
		return false
	} else {
		return true
	}
}

// 登录时的寻找用户
func LoginFindUser(u *User) error {
	if u == nil {
		return errors.New("空指针错误")
	}
	var u_find User
	dao.DB.Where("telephone=?", u.Telephone).First(&u_find)
	if u_find.ID == 0 || u_find.Is_deleted == 1 {
		return errors.New("该用户不存在")
	} else if u_find.Password != u.Password {
		return errors.New("密码错误")
	}
	u.ID = u_find.ID
	u.Avatar = u_find.Avatar
	u.Description = u_find.Description
	return nil
}

// 根据手机号寻找用户
func FindUserByTelephone(u *User) error {
	if u == nil {
		return errors.New("空指针错误")
	}
	var u_find User
	dao.DB.Where("telephone=?", u.Telephone).First(&u_find)
	if u_find.ID == 0 || u_find.Is_deleted == 1 {
		return errors.New("该用户不存在")
	}
	u.ID = u_find.ID
	u.Username = u_find.Username
	u.Password = u_find.Password
	u.Avatar = u_find.Avatar
	u.Description = u_find.Description
	return nil
}

// 根据手机号查找用户ID
func FindUserIDByTelephone(telephone string) (int, error) {
	var u_find User
	dao.DB.Where("telephone=?", telephone).First(&u_find)
	if u_find.ID == 0 || u_find.Is_deleted == 1 {
		return 0, errors.New("该用户不存在")
	}
	return u_find.ID, nil
}

// 根据用户ID返回用户
func GetUserByID(userID int, user *User) error {
	if user == nil {
		return errors.New("空指针错误")
	}
	err := dao.DB.Where("id = ? AND is_deleted != 1", userID).Find(&user).Error
	return err
}

// 根据关键词查找用户
func GetUserListByKeyWord(userList *[]*User, keyWord string) error {
	if userList == nil {
		return errors.New("空指针错误")
	}
	err := dao.DB.Where("username LIKE ? AND is_deleted != 1", "%"+keyWord+"%").Find(&userList).Error
	return err
}

// 修改用户的所有信息
func UpdateUserInfo(user *User) error {
	if user == nil {
		return errors.New("空指针错误")
	}
	err := dao.DB.Save(&user).Error
	return err
}
