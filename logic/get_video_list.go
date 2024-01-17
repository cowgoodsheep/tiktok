package logic

import (
	"tiktok/models"
	"tiktok/util"
	"time"
)

const (
	MaxVideoNum = 10 //每次返回的视频数量
)

type VideoList struct {
	Videos   []*models.Video `json:"video_list"`
	NextTime int             `json:"next_time"` //下一个视频列表时间戳
}

type VideoListFlow struct {
	userID     int
	nextTime   int       //下一个视频列表时间戳
	latestTime time.Time //最新时间
	videos     []*models.Video

	Video_List *VideoList
}

func GetVideoList(userID int, latestTime time.Time) (*VideoList, error) {
	return NewVideoListFlow(userID, latestTime).Do()
}

func NewVideoListFlow(userID int, latestTime time.Time) *VideoListFlow {
	return &VideoListFlow{userID: userID, latestTime: latestTime}
}

func (n *VideoListFlow) Do() (*VideoList, error) {
	//上层通过把userId置零，表示userId不存在或不需要
	if n.userID > 0 {
		//这里说明userId是有效的，可以定制性的做一些登录用户的专属视频推荐

	}

	//如果时间出错，则设置时间为现在
	if n.latestTime.IsZero() {
		n.latestTime = time.Now()
	}

	if err := n.getVideos(); err != nil {
		return nil, err
	}

	//整理发送
	n.Video_List = &VideoList{
		Videos:   n.videos,
		NextTime: n.nextTime,
	}
	return n.Video_List, nil
}

func (n *VideoListFlow) getVideos() error {
	//获取视频列表
	if err := models.GetVideoListByLimitAndLastestTime(&n.videos, MaxVideoNum, n.latestTime); err != nil {
		return err
	}

	//检查该视频是否被已登录用户点赞和收藏
	if err := util.CheckVideoState(n.userID, &n.videos); err != nil {
		return err
	}

	//准备下一个视频列表的时间戳
	if len(n.videos) == 0 {
		n.nextTime = int(time.Now().Unix() / 1e6)
		return nil
	}
	lastTime := n.videos[len(n.videos)-1].Upload_time
	n.nextTime = int(lastTime.UnixNano() / 1e6)
	return nil
}
