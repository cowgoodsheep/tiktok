package logic

import (
	"errors"
	"fmt"
	"tiktok/cache"
	"tiktok/models"
)

type LikeFlow struct {
	userID  int
	videoID int
}

// 添加点赞
func PostLike(userID, videoID int) error {
	return NewPostLikeFlow(userID, videoID).Do()
}

func NewPostLikeFlow(userID, videoID int) *LikeFlow {
	return &LikeFlow{userID: userID, videoID: videoID}
}

func (p *LikeFlow) Do() error {
	var err error

	//检查用户ID是否正确
	if !models.IsUserExistByUserID(p.userID) {
		return fmt.Errorf("用户%d不存在", p.userID)
	}

	//检查视频ID是否正确
	if !models.IsVideoExistByVideoID(p.videoID) {
		return fmt.Errorf("视频%d不存在", p.videoID)
	}

	//检查该用户是否已点赞该视频
	if cache.GetUserLike(p.userID, p.videoID) {
		return fmt.Errorf("请勿重复点赞")
	}

	//使视频点赞数加一
	err = models.AddVideoLikeByVideoID(p.videoID)
	if err != nil {
		return errors.New("不要重复点赞")
	}

	//更新redis的点赞信息
	cache.UpdateUserLike(p.userID, p.videoID, true)

	return nil
}

// 取消点赞
func DeleteLike(userID, videoID int) error {
	//检查用户ID是否正确
	if !models.IsUserExistByUserID(userID) {
		return fmt.Errorf("用户%d不存在", userID)
	}

	//检查视频ID是否正确
	if !models.IsVideoExistByVideoID(videoID) {
		return fmt.Errorf("视频%d不存在", videoID)
	}

	//检查该用户是否已点赞该视频
	if !cache.GetUserLike(userID, videoID) {
		return fmt.Errorf("请勿重复取消点赞")
	}

	//删除Like
	err := models.SubVideoLikeByVideoID(videoID)
	if err != nil {
		return err
	}

	//更新redis的关注信息
	cache.UpdateUserLike(userID, videoID, false)

	return nil
}
