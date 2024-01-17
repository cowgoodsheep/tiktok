package logic

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"tiktok/config"
	"tiktok/models"
	"tiktok/util"
	"time"
)

type PostVideoFlow struct {
	videoUrl    string //视频URL
	coverUrl    string //封面URL
	title       string //视频标题
	description string //视频描述
	userID      int    //上传用户ID
}

// 上传视频
func PostVideo(userID int, videoUrl, coverUrl, title, description string) error {
	return NewPostVideoFlow(userID, videoUrl, coverUrl, title, description).Do()
}

func NewPostVideoFlow(userID int, videoUrl, coverUrl, title, description string) *PostVideoFlow {
	return &PostVideoFlow{videoUrl: videoUrl, coverUrl: coverUrl, userID: userID, title: title, description: description}
}

func (p *PostVideoFlow) Do() error {

	// 获取Url
	p.videoUrl = util.GetVideoFileUrl(p.videoUrl)
	p.coverUrl = util.GetCoverFileUrl(p.coverUrl)

	if err := p.uploadVideo(); err != nil {
		return err
	}
	return nil
}

// 组合并添加到数据库
func (p *PostVideoFlow) uploadVideo() error {
	video := &models.Video{
		User_ID:     p.userID,
		Title:       p.title,
		Description: p.description,
		Video_url:   p.videoUrl,
		Cover_url:   p.coverUrl,
		Upload_time: time.Now(),
	}
	return models.AddVideo(video)
}

// 删除视频
func DeleteVideo(videoID, userID int) (*models.Video, error) {
	//先获取video
	var video models.Video
	if err := models.GetVideoByID(videoID, &video); err != nil {
		return nil, err
	}

	//检查这个视频的合法性
	//如果视频ID和用户ID不对应
	if video.User_ID != userID {
		return nil, fmt.Errorf("评论ID与用户ID不匹配")
	}

	//删除video
	if err := models.DeleteVideoByID(videoID); err != nil {
		return nil, err
	}

	//删除该视频的评论
	if err := models.DelectCommentByVideoID(videoID); err != nil {
		return nil, err
	}

	//删除本地的视频文件和封面文件
	//先拼出url前半部分
	ip := config.Conf.Server.IP + ":" + strconv.Itoa(config.Conf.Server.Port)
	//接着找出地址下标开始的地方
	videoIndex := strings.Index(video.Video_url, ip) + len(ip)
	coverIndex := strings.Index(video.Cover_url, ip) + len(ip)
	//拼接出视频和封面的本地地址
	videoPath := "." + video.Video_url[videoIndex:]
	coverPath := "." + video.Cover_url[coverIndex:]
	if err := os.Remove(videoPath); err != nil {
		return nil, err
	}
	if err := os.Remove(coverPath); err != nil {
		return nil, err
	}

	return &video, nil
}
