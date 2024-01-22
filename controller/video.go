package controller

import (
	"net/http"
	"path/filepath"
	"strconv"
	"tiktok/cache"
	"tiktok/logic"
	"tiktok/middleware"
	"tiktok/models"
	"tiktok/util"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	videoExtMap = map[string]struct{}{
		".avi": {},

		".wmv": {},

		".mpg":  {},
		".mpeg": {},
		".vob":  {},
		".dat":  {},
		".3gp":  {},
		".mp4":  {},

		".mkv": {},

		".rm":   {},
		".rmvb": {},

		".mov": {},

		".flv": {},
	}

	pictureExtMap = map[string]struct{}{
		".xbm":   {},
		".tif":   {},
		".pjp":   {},
		".svgz":  {},
		".jpg":   {},
		".jpeg":  {},
		".ico":   {},
		".tiff":  {},
		".gif":   {},
		".svg":   {},
		".jfif":  {},
		".webp":  {},
		".png":   {},
		".bmp":   {},
		".pjpeg": {},
		".avif":  {},
	}
)

// 用户上传视频
func UploadVideo(c *gin.Context) {
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
	var u models.User
	u.Telephone = telephone.(string)
	if err := models.FindUserByTelephone(&u); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	//获取视频标题和视频描述
	title := c.PostForm("title")
	description := c.PostForm("description")
	//获取视频和封面数据
	video, err := c.FormFile("video")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
			"msg":   "视频获取出错",
		})
		return
	}
	cover, err := c.FormFile("cover")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
			"msg":   "封面获取出错",
		})
		return
	}

	//获取文件扩展名
	videoExt := filepath.Ext(video.Filename)
	coverExt := filepath.Ext(cover.Filename)

	//判断视频扩展名是否正确
	if _, ok := videoExtMap[videoExt]; ok == false {
		c.JSON(http.StatusOK, gin.H{
			"file": video,
			"msg":  "视频文件格式不正确",
		})
		return
	}

	//判断封面扩展名是否正确
	if _, ok := pictureExtMap[coverExt]; ok == false {
		c.JSON(http.StatusOK, gin.H{
			"file": cover,
			"msg":  "封面文件格式不正确",
		})
		return
	}

	//将文件名唯一化
	videofileName := util.GenerateUniqueFileName(video.Filename)
	coverfileName := util.GenerateUniqueFileName(cover.Filename)

	//保存视频文件
	if err := c.SaveUploadedFile(video, "./static/videos/"+videofileName); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
			"msg":   "视频文件保存失败",
		})
		return
	}
	//保存封面文件
	if err := c.SaveUploadedFile(cover, "./static/covers/"+coverfileName); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
			"msg":   "封面文件保存失败",
		})
		return
	}

	//将视频存入数据库
	if err := logic.PostVideo(u.ID, videofileName, coverfileName, title, description); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  title + "上传成功",
	})
}

// 查看本用户已上传的视频
func GetUserVideos(c *gin.Context) {
	telephone, ok := c.Get("telephone")
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "用户手机号码解析错误",
		})
		return
	}
	var u models.User
	u.Telephone = telephone.(string)

	//在数据库中寻找该用户信息
	if err := models.FindUserByTelephone(&u); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	//在数据库中寻找该用户ID上传的视频
	var videoList []*models.Video
	if err := models.GetVideoListByUserID(&videoList, u.ID); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	//返回视频信息
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  videoList,
	})
}

// 用户删除自己上传的视频
func DeleteUserVideos(c *gin.Context) {
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
	var u models.User
	u.Telephone = telephone.(string)
	if err := models.FindUserByTelephone(&u); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	//获取待删除视频的ID
	stringVideoID := c.Query("video_id")
	//如果找不到，就去Form表单去找
	if stringVideoID == "" {
		stringVideoID = c.PostForm("video_id")
	}
	//将视频ID转化成int
	videoID, err := strconv.ParseInt(stringVideoID, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	//调用logic将视频从本地和数据库中删除
	getVideo, err := logic.DeleteVideo(int(videoID), u.ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":  200,
		"video": getVideo,
		"msg":   "删除成功",
	})
}

// 视频详情页
func VideoDetail(c *gin.Context) {
	//获取视频的ID
	stringVideoID := c.Query("video_id")
	//如果找不到，就去Form表单去找
	if stringVideoID == "" {
		stringVideoID = c.PostForm("video_id")
	}
	//将视频ID转化成int
	videoID, err := strconv.ParseInt(stringVideoID, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	//增加视频播放量
	if err := models.AddVideoViewByVideoID(int(videoID)); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	//获取token
	token := c.Query("token")
	//如果找不到，就去Form表单去找
	if token == "" {
		token = c.PostForm("token")
	}

	//从数据库中寻找该视频的信息
	var video models.Video
	if err := models.GetVideoByID(int(videoID), &video); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	//获得视频作者信息
	if err := models.GetUserByID(video.User_ID, &video.Author); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	//如果是登录状态
	if len(token) > 0 {
		claim, ok := middleware.ParseToken(token)
		if ok == false {
			c.JSON(http.StatusOK, gin.H{"error": "token解析失败"})
			return
		}

		//token过期
		if time.Now().Unix() > claim.ExpiresAt {
			c.JSON(http.StatusOK, gin.H{"error": "token已失效"})
			return
		}

		//获取用户ID
		userID, err := models.FindUserIDByTelephone(claim.Telephone)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}

		//检查该作者是否已被用户关注
		video.Author.Is_Follow = cache.GetUserFollow(userID, video.Author.ID)
		//检查该视频是否已被点赞
		video.Is_Like = cache.GetUserLike(userID, video.ID)
		//检查该视频是否已被收藏
		video.Is_Favorite = cache.GetUserFavorite(userID, video.ID)

		//调用logic，添加历史浏览记录
		err = logic.PostHistory(userID, int(videoID))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "添加历史浏览记录成功",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":  200,
		"video": video,
		"msg":   "获取视频详细信息成功",
	})
}
