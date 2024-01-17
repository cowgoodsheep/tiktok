package logic

import (
	"fmt"
	"tiktok/cache"
	"tiktok/models"
	"time"
)

type FavoriteFlow struct {
	userID  int
	videoID int
}

// 添加收藏
func PostFavorite(userID, videoID int) (*models.Favorite, error) {
	return NewPostFavoriteFlow(userID, videoID).Do()
}

func NewPostFavoriteFlow(userID, videoID int) *FavoriteFlow {
	return &FavoriteFlow{userID: userID, videoID: videoID}
}

func (p *FavoriteFlow) Do() (*models.Favorite, error) {
	var err error

	//检查ID是否正确
	if !models.IsUserExistByUserID(p.userID) {
		return nil, fmt.Errorf("用户%d不存在", p.userID)
	}
	if !models.IsVideoExistByVideoID(p.videoID) {
		return nil, fmt.Errorf("视频%d不存在", p.videoID)
	}

	//检查该视频是否已被该用户收藏
	if models.IsFavoriteExistByUserIDAndVideoID(p.userID, p.videoID) {
		return nil, fmt.Errorf("请勿重复收藏同一个视频")
	}

	//整理数据，上传数据库
	favorite := models.Favorite{User_ID: p.userID, Video_ID: p.videoID, Favorite_Time: time.Now()}
	err = models.AddFavorite(&favorite)
	if err != nil {
		return nil, err
	}

	//更新redis的点赞信息
	cache.UpdateUserFavorite(p.userID, p.videoID, true)

	return &favorite, nil
}

// 删除收藏
func DeleteFavorite(userID, videoID int) (*models.Favorite, error) {
	//先获取favorite
	var favorite models.Favorite
	err := models.GetFavoriteByUserIDAndVideoID(userID, videoID, &favorite)
	if err != nil {
		return nil, err
	}

	//检查该收藏的合法性
	//如果收藏ID和用户ID不对应
	if favorite.User_ID != userID {
		return nil, fmt.Errorf("收藏ID与用户ID不匹配")
	}
	//如果收藏ID和视频ID不对应
	if favorite.Video_ID != videoID {
		return nil, fmt.Errorf("收藏ID与视频ID不匹配")
	}

	//删除favorite
	err = models.DeleteFavoriteByID(favorite.ID)
	if err != nil {
		return nil, err
	}

	//更新redis的点赞信息
	cache.UpdateUserFavorite(userID, videoID, false)

	return &favorite, nil
}
