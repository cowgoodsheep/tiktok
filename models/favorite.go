package models

import (
	"errors"
	"tiktok/dao"
	"time"
)

// Favorite Model
type Favorite struct {
	ID            int       `json:"id"`
	User_ID       int       `json:"user_id"`
	Video_ID      int       `json:"video_id"`
	Favorite_Time time.Time `json:"favorite_time"`
}

// 指定Favorite结构体迁移表favorite
func (F *Favorite) TableName() string {
	return "favorite"
}

// 添加收藏
func AddFavorite(c *Favorite) error {
	if c == nil {
		return errors.New("空指针错误")
	}
	return dao.DB.Create(c).Error
}

// 根据用户ID获得收藏视频列表
func GetFavListByUserID(favList *[]*Video, userId int) error {
	if favList == nil {
		return errors.New("空指针错误")
	}

	//多表查询，左连接得到结果，再映射到数据
	if err := dao.DB.Raw("SELECT v.* FROM favorite f, video v WHERE f.user_id = ? AND f.video_id = v.id", userId).Scan(favList).Error; err != nil {
		return err
	}

	//如果id为0，则说明没有查到收藏数据
	if len(*favList) == 0 || (*favList)[0].ID == 0 {
		return errors.New("收藏列表为空")
	}
	return nil
}

// 根据收藏ID返回收藏
func GetFavoriteByID(favoriteID int, favorite *Favorite) error {
	if favorite == nil {
		return errors.New("空指针错误")
	}
	err := dao.DB.Where("id=?", favoriteID).Order("favorite_time DESC").Find(&favorite).Error
	return err
}

// 根据用户ID和视频ID获取收藏
func GetFavoriteByUserIDAndVideoID(userID, videoID int, favorite *Favorite) error {
	if favorite == nil {
		return errors.New("空指针错误")
	}
	err := dao.DB.Where("user_id = ? AND video_id = ?", userID, videoID).Order("favorite_time DESC").Find(&favorite).Error
	return err
}

// 根据收藏ID删除收藏
func DeleteFavoriteByID(favoriteID int) error {
	err := dao.DB.Delete(&Favorite{}, favoriteID).Error
	return err
}

// 根据用户ID和视频ID检查该视频是否已被收藏
func IsFavoriteExistByUserIDAndVideoID(userID, videoID int) bool {
	var f Favorite
	dao.DB.Where("user_id = ? AND video_id = ?", userID, videoID).First(&f)
	if f.ID == 0 {
		return false
	} else {
		return true
	}
}
