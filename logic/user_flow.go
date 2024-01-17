package logic

import (
	"errors"
	"tiktok/middleware"
	"tiktok/models"
)

const (
	MaxUsernameSize    = 100
	MaxDescriptionSize = 255
)

type UserFlow struct {
	username  string
	telephone string
	password  string
	ID        int    `json:"id"`
	Token     string `json:"token"`
}

// 注册用户并得到token和id
func PostUser(telephone, username, password string) (*UserFlow, error) {
	return NewPostUserFlow(telephone, username, password).PostDo()
}

func NewPostUserFlow(telephone, username, password string) *UserFlow {
	return &UserFlow{telephone: telephone, username: username, password: password}
}

func (u *UserFlow) PostDo() (*UserFlow, error) {
	//对注册信息进行合法性验证
	if err := u.postCheck(); err != nil {
		return nil, err
	}
	//将注册信息上传至数据库
	if err := u.postUpload(); err != nil {
		return nil, err
	}
	return u, nil
}

func (u *UserFlow) postCheck() error {
	if u.username == "" {
		return errors.New("用户名为空")
	}
	if len(u.username) > MaxUsernameSize {
		return errors.New("用户名长度超出限制")
	}
	if len(u.telephone) != 11 {
		return errors.New("手机号码长度不为11位")
	}
	for _, v := range u.telephone {
		if v < '0' || v > '9' {
			return errors.New("手机号码存在非数字字符")
		}
	}
	if u.password == "" {
		return errors.New("密码为空")
	}
	return nil
}

func (u *UserFlow) postUpload() error {
	userinfo := models.User{
		Username:  u.username,
		Telephone: u.telephone,
		Password:  u.password,
	}

	//判断手机号是否已被注册
	if models.IsUserExistByTelephone(u.telephone) {
		return errors.New("该手机号已被注册")
	}

	//上传数据库
	if err := models.AddUser(&userinfo); err != nil {
		return err
	}

	//根据用户手机号生成token
	token, err := middleware.MakeToken(u.telephone)
	if err != nil {
		return err
	}
	u.Token = token
	u.ID = userinfo.ID
	return nil
}

// 查询登录用户并返回token和id
func GetUser(telephone, password string) (*UserFlow, error) {
	return NewGetUserFlow(telephone, password).GetDo()
}

func NewGetUserFlow(telephone, password string) *UserFlow {
	return &UserFlow{telephone: telephone, password: password}
}

func (u *UserFlow) GetDo() (*UserFlow, error) {
	//对登录信息进行合法性验证
	if err := u.getCheck(); err != nil {
		return nil, err
	}
	//从数据库中得到登录用户的信息
	if err := u.getUpload(); err != nil {
		return nil, err
	}
	return u, nil
}

func (u *UserFlow) getCheck() error {
	if u.telephone == "" {
		return errors.New("手机号为空")
	}
	if len(u.telephone) != 11 {
		return errors.New("手机号码长度不为11位")
	}
	for _, v := range u.telephone {
		if v < '0' || v > '9' {
			return errors.New("手机号码存在非数字字符")
		}
	}
	if u.password == "" {
		return errors.New("密码为空")
	}
	return nil
}

func (u *UserFlow) getUpload() error {
	userinfo := models.User{
		Telephone: u.telephone,
		Password:  u.password,
	}

	//从数据库中寻找用户
	if err := models.LoginFindUser(&userinfo); err != nil {
		return err
	}
	//根据用户手机号生成token
	token, err := middleware.MakeToken(u.telephone)
	if err != nil {
		return err
	}
	u.username = userinfo.Username
	u.ID = userinfo.ID
	u.Token = token
	return nil
}
