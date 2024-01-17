package controller

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"tiktok/config"
	"tiktok/logic"
	"tiktok/middleware"
	"tiktok/models"
	"tiktok/util"

	"github.com/gin-gonic/gin"
)

// 用户注册
func UserRegister(c *gin.Context) {
	//获取用户注册信息
	var u models.User
	u.Username = c.PostForm("username")
	u.Telephone = c.PostForm("telephone")
	passwordTemp, ok := c.Get("password")
	u.Password = passwordTemp.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "密码解析错误",
		})
		return
	}

	//上传用户注册信息至用户登录服务，进行用户注册
	registerMsg, err := logic.PostUser(u.Telephone, u.Username, u.Password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"msg": err.Error()})
		return
	}

	//注册成功
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  registerMsg,
	})
}

// 用户登录
func UserLogin(c *gin.Context) {
	//获取用户登录信息
	var u models.User
	u.Telephone = c.PostForm("telephone")
	passwordTemp, ok := c.Get("password")
	u.Password = passwordTemp.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "密码解析错误",
		})
		return
	}

	//获取用户流信息
	loginMsg, err := logic.GetUser(u.Telephone, u.Password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"msg": err.Error()})
		return
	}

	//用户登录成功
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  loginMsg,
	})
}

// 查看用户信息
func GetUserInfo(c *gin.Context) {
	//从请求的查询参数中获取用户ID
	stringUserID := c.Query("user_id")
	//如果找不到，就去Form表单去找
	if stringUserID == "" {
		stringUserID = c.PostForm("user_id")
	}

	//如果用户不存在
	if stringUserID == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "用户不存在",
		})
		return
	}

	//将用户ID转化成int
	userID, err := strconv.ParseInt(stringUserID, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	//在数据库中寻找该用户信息
	if err := models.GetUserByID(int(userID), &user); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	//查询成功
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"user": user,
		"msg":  "查询成功",
	})
}

// 用户修改个人信息
func UpdateUserInfo(c *gin.Context) {
	//解析token 得到用户手机号
	telephone, ok := c.Get("telephone")
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "用户手机号码解析错误",
		})
		return
	}

	//根据用户手机号得到用户信息
	var user models.User
	user.Telephone = telephone.(string)
	if err := models.FindUserByTelephone(&user); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	//获取修改后的用户昵称
	username := c.Query("username")
	//如果找不到，就去Form表单去找
	if username == "" {
		username = c.PostForm("username")
	}
	//如果用户昵称不为空，则修改
	if username != "" {
		//检查用户名合法性
		if len(username) > logic.MaxUsernameSize {
			c.JSON(http.StatusOK, gin.H{"error": "用户名长度超出限制"})
			return
		}
		user.Username = username
	}

	//获取修改后的个人描述
	description := c.Query("description")
	//如果找不到，就去Form表单去找
	if description == "" {
		description = c.PostForm("description")
	}
	//如果个人描述不为空，则修改
	if description != "" {
		//检查个人描述合法性
		if len(description) > logic.MaxDescriptionSize {
			c.JSON(http.StatusOK, gin.H{"error": "个人描述长度超出限制"})
			return
		}
		user.Description = description
	}

	//获取修改后的头像文件
	avatarFile, err := c.FormFile("avatar")
	if err != nil && err != http.ErrMissingFile {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	//如果头像文件不为空，则修改
	if err != http.ErrMissingFile {
		//如果用户存在头像，则先删除本地的头像文件
		if user.Avatar != "" {
			//先拼出url前半部分
			ip := config.Conf.Server.IP + ":" + strconv.Itoa(config.Conf.Server.Port)
			//接着找出地址下标开始的地方
			avatarIndex := strings.Index(user.Avatar, ip) + len(ip)
			//拼接出头像文件的本地地址
			avatarPath := "." + user.Avatar[avatarIndex:]
			if err := os.Remove(avatarPath); err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
				return
			}
		}

		//获取头像文件扩展名
		avatarExt := filepath.Ext(avatarFile.Filename)

		//判断头像文件扩展名是否正确
		if _, ok := pictureExtMap[avatarExt]; ok == false {
			c.JSON(http.StatusOK, gin.H{
				"file": avatarFile,
				"msg":  "头像文件格式不正确",
			})
			return
		}

		//将文件名唯一化
		avatarFileName := util.GenerateUniqueFileName(avatarFile.Filename)

		//更新用户的头像URL
		user.Avatar = util.GetAvatarFileUrl(avatarFileName)

		//保存头像文件至服务器本地
		if err := c.SaveUploadedFile(avatarFile, "./static/avatars/"+avatarFileName); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error": err.Error(),
				"msg":   "头像文件保存失败",
			})
			return
		}
	}

	//更新数据库数据
	if err := models.UpdateUserInfo(&user); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	//修改成功
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"user": user,
		"msg":  "修改成功",
	})
}

// 用户修改密码
func UpdateUserPassword(c *gin.Context) {
	//解析token 得到用户手机号
	telephone, ok := c.Get("telephone")
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "用户手机号码解析错误",
		})
		return
	}

	//根据用户手机号得到用户信息
	var user models.User
	user.Telephone = telephone.(string)
	if err := models.FindUserByTelephone(&user); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	//这里可以进行各种身份验证，但为了方便就不验证了

	//获取修改后的密码
	password := c.Query("password")
	//如果找不到，就去Form表单去找
	if password == "" {
		password = c.PostForm("password")
	}
	//如果密码不为空，则修改
	if password != "" {
		//检查密码合法性
		if len(password) < middleware.MinPasswordSize || len(password) > middleware.MaxPasswordSize {
			c.JSON(http.StatusOK, gin.H{"err": "密码长度小于或大于限制"})
			return
		}
		//加密修改后的密码
		user.Password = middleware.SHA1(password)
	}

	//更新数据库数据
	if err := models.UpdateUserInfo(&user); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	//修改成功
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"user": user,
		"msg":  "密码修改成功",
	})
}

// 用户注销账号
func DeleteUser(c *gin.Context) {
	//解析token 得到用户手机号
	telephone, ok := c.Get("telephone")
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "用户手机号码解析错误",
		})
		return
	}

	//根据用户手机号得到用户信息
	var user models.User
	user.Telephone = telephone.(string)
	if err := models.FindUserByTelephone(&user); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	//软删除
	user.Is_deleted = 1

	//更新数据库数据
	if err := models.UpdateUserInfo(&user); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	//注销成功
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"user": user,
		"msg":  "注销成功",
	})
}
