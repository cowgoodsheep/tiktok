package util

import (
	"tiktok/cache"
	"tiktok/models"
)

func CheckVideoState(userID int, videos *[]*models.Video) error {
	for i := 0; i < len(*videos); i++ {
		//查询视频作者信息
		var author models.User
		if err := models.GetUserByID((*videos)[i].User_ID, &author); err != nil {
			continue
		}
		(*videos)[i].Author = author

		if userID > 0 {
			//检查该作者是否已被用户关注
			(*videos)[i].Author.Is_Follow = cache.GetUserFollow(userID, author.ID)
			//检查该视频是否已被点赞
			(*videos)[i].Is_Like = cache.GetUserLike(userID, (*videos)[i].ID)
			//检查该视频是否已被收藏
			(*videos)[i].Is_Favorite = cache.GetUserFavorite(userID, (*videos)[i].ID)
		}
	}
	return nil
}
