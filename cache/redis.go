package cache

import (
	"context"
	"fmt"
	"tiktok/config"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client
var c = context.Background()

const (
	like     = "like"
	follow   = "follow"
	favorite = "favorite"
)

// 初始化redis
func init() {
	rdb = redis.NewClient(
		&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", config.Conf.RDB.IP, config.Conf.RDB.Port),
			Password: "", //没有设置密码
			DB:       config.Conf.RDB.Database,
		})
}

// 使用redis中名为follow的集合set来存储关注和被关注的关系
// 更新关注集合，state为true是添加关注，为false是取消关注
func UpdateUserFollow(userID, followUserID int, state bool) {
	key := fmt.Sprintf("%s:%d", follow, userID)
	value := followUserID
	if state == true {
		//往set中添加（key，value）
		rdb.SAdd(c, key, value)
	} else {
		//从set中删除（key，value）
		rdb.SRem(c, key, value)
	}
}

// 获取关注状态，true为已关注，false为未关注
func GetUserFollow(userID, followUserID int) bool {
	key := fmt.Sprintf("%s:%d", follow, userID)
	//检查set中是否有（key，value）
	return rdb.SIsMember(c, key, followUserID).Val()
}

// 使用redis中名为like的集合set来存储点赞用户和点赞视频的关系
// 更新关注集合，state为true是点赞，为false是取消点赞
func UpdateUserLike(userID, videoID int, state bool) {
	key := fmt.Sprintf("%s:%d", like, userID)
	value := videoID
	if state == true {
		//往set中添加（key，value）
		rdb.SAdd(c, key, value)
	} else {
		//从set中删除（key，value）
		rdb.SRem(c, key, value)
	}
}

// 获取点赞状态，true为已点赞，false为未点赞
func GetUserLike(userID, videoID int) bool {
	key := fmt.Sprintf("%s:%d", like, userID)
	//检查set中是否有（key，value）
	return rdb.SIsMember(c, key, videoID).Val()
}

// 使用redis中名为favorite的集合set来存储点赞用户和点赞视频的关系
// 更新收藏集合，state为true是收藏，为false是取消收藏
func UpdateUserFavorite(userID, videoID int, state bool) {
	key := fmt.Sprintf("%s:%d", favorite, userID)
	value := videoID
	if state == true {
		//往set中添加（key，value）
		rdb.SAdd(c, key, value)
	} else {
		//从set中删除（key，value）
		rdb.SRem(c, key, value)
	}
}

// 获取收藏状态，true为已收藏，false为未收藏
func GetUserFavorite(userID, videoID int) bool {
	key := fmt.Sprintf("%s:%d", favorite, userID)
	//检查set中是否有（key，value）
	return rdb.SIsMember(c, key, videoID).Val()
}
