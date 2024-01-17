package models

import (
	"errors"
	"tiktok/dao"
	"time"
)

// Video Model
type Video struct {
	ID            int       `json:"id"`
	User_ID       int       `json:"user_id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Upload_time   time.Time `json:"upload_time"`
	Likes_count   int       `json:"likes_count"`
	Comment_count int       `json:"comment_counts"`
	Video_url     string    `json:"video_url"`
	Cover_url     string    `json:"cover_url"`

	//以下字段不需要存储进数据库
	Is_Like     bool `json:"is_like" gorm:"-"`
	Is_Favorite bool `json:"is_favorite" gorm:"-"`
	Author      User `json:"author" gorm:"-"`
}

// 指定Video结构体迁移表video
func (v *Video) TableName() string {
	return "video"
}

// 存入一个视频
func AddVideo(v *Video) error {
	if v == nil {
		return errors.New("空指针错误")
	}
	return dao.DB.Create(v).Error
}

// 根据ID判断该视频是否存在
func IsVideoExistByVideoID(videoID int) bool {
	var v Video
	dao.DB.Where("ID=?", videoID).First(&v)
	if v.ID == 0 {
		return false
	} else {
		return true
	}
}

// 按照投稿时间倒叙返回视频列表，最多返回limit个视频
func GetVideoListByLimitAndLastestTime(videoList *[]*Video, limit int, lastestTime time.Time) error {
	if videoList == nil {
		return errors.New("空指针错误")
	}
	err := dao.DB.Model(&Video{}).Where("upload_time<?", lastestTime).
		Order("upload_time DESC").
		Limit(limit).
		Select([]string{"id", "user_id", "title", "description", "upload_time", "likes_count", "comment_count", "video_url", "cover_url"}).
		Find(&videoList).
		Error
	return err
}

// 根据用户ID按照时间倒叙返回视频列表
func GetVideoListByUserID(videoList *[]*Video, userID int) error {
	if videoList == nil {
		return errors.New("空指针错误")
	}
	err := dao.DB.Where("user_id=?", userID).Order("upload_time DESC").Find(&videoList).Error
	return err
}

// 根据关键词按照时间倒叙返回视频列表
func GetVideoListByKeyWord(videoList *[]*Video, keyWord string) error {
	if videoList == nil {
		return errors.New("空指针错误")
	}
	err := dao.DB.Where("title LIKE ?", "%"+keyWord+"%").Order("upload_time DESC").Find(&videoList).Error
	return err
}

// 根据视频ID返回视频
func GetVideoByID(videoID int, video *Video) error {
	if video == nil {
		return errors.New("空指针错误")
	}
	err := dao.DB.Where("id=?", videoID).Order("upload_time DESC").Find(&video).Error
	return err
}

// 根据视频ID删除视频
func DeleteVideoByID(videoID int) error {
	err := dao.DB.Delete(&Video{}, videoID).Error
	return err
}

// 根据视频ID给点赞数加一
func AddVideoLikeByVideoID(videoID int) error {
	var video Video
	if err := GetVideoByID(videoID, &video); err != nil {
		return err
	}

	if err := dao.DB.Model(&video).Update("likes_count", video.Likes_count+1).Error; err != nil {
		return err
	}
	return nil
}

// 根据视频ID给点赞数减一
func SubVideoLikeByVideoID(videoID int) error {
	var video Video
	if err := GetVideoByID(videoID, &video); err != nil {
		return err
	}

	if err := dao.DB.Model(&video).Update("likes_count", video.Likes_count-1).Error; err != nil {
		return err
	}
	return nil
}
