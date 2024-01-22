package routers

import (
	"tiktok/config"
	"tiktok/controller"
	"tiktok/middleware"

	"github.com/gin-gonic/gin"
)

// 跨域访问处理
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // 允许所有来源访问，也可设置特定域名
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	//设置静态路径
	r.Static("static", config.Conf.StaticSourcePath)

	// 设置跨域访问处理中间件
	r.Use(corsMiddleware())

	/*
		所有功能

		显示主页√

		搜索：
		视频√
		用户√
		直播间

		用户：
		登陆√
		注册√
		查看用户信息√
		修改个人信息√
		修改密码√
		注销账号（软删除）√
		上传视频√
		查看已上传视频√
		删除上传视频√
		查看收藏视频√
		删除收藏视频√
		关注用户√
		查看已关注用户√
		查看粉丝√
		删除已关注用户√
		查看历史浏览记录√
		删除历史浏览记录√

		视频：
		视频详情（若已登录则添加历史浏览记录）√
		查看视频评论√
		用户发表视频评论√
		用户删除视频评论√
		查看视频评论回复√
		用户回复视频评论√
		用户删除回复视频评论√
		用户点赞视频√
		用户取消点赞视频√
		获取该视频的弹幕×
		用户发送视频弹幕×
		用户删除视频弹幕×

		直播：
		用户创建直播间
		用户开启直播间
		用户关闭直播间
		用户观看直播间
		用户送礼物
		用户给主播发弹幕
		用户关注主播用户
		主播开播提醒已关注用户
	*/

	//Home主页
	r.GET("/home", controller.Home)

	//搜索路由组
	searchGroup := r.Group("/search")
	{
		//搜索视频
		searchGroup.GET("/videos", controller.SearchVideos)
		//搜索用户
		searchGroup.GET("/users", controller.SearchUsers)
		//搜索直播间
	}

	//用户路由组
	userGroup := r.Group("/user")
	{
		//用户注册
		userGroup.POST("/register", middleware.SHAMiddleWare(), controller.UserRegister)
		//用户登录
		userGroup.POST("/login", middleware.SHAMiddleWare(), controller.UserLogin)

		//查看用户信息
		userGroup.GET("/", controller.GetUserInfo)

		//用户注销账号（软删除）
		userGroup.DELETE("/delete", middleware.JWTMiddleWare(), controller.DeleteUser)

		//用户修改个人信息
		userGroup.PUT("/update/info", middleware.JWTMiddleWare(), controller.UpdateUserInfo)
		//用户修改密码
		userGroup.PUT("/update/password", middleware.JWTMiddleWare(), controller.UpdateUserPassword)

		//查看本用户已上传的视频
		userGroup.GET("/videos", middleware.JWTMiddleWare(), controller.GetUserVideos)
		//用户上传视频
		userGroup.POST("/videos", middleware.JWTMiddleWare(), controller.UploadVideo)
		//用户删除自己上传的视频
		userGroup.DELETE("/videos", middleware.JWTMiddleWare(), controller.DeleteUserVideos)

		//查看本用户收藏的视频
		userGroup.GET("/favorite", middleware.JWTMiddleWare(), controller.GetUserFavorite)
		//删除本用户收藏的视频
		userGroup.DELETE("/favorite", middleware.JWTMiddleWare(), controller.DeleteUserFavorite)

		//用户关注其他用户
		userGroup.POST("/follow", middleware.JWTMiddleWare(), controller.PostUserFollow)
		//用户查看已关注用户列表
		userGroup.GET("/follow", middleware.JWTMiddleWare(), controller.GetUserFollow)
		//用户取关已关注用户
		userGroup.DELETE("/follow", middleware.JWTMiddleWare(), controller.DeleteUserFollow)

		//用户查看粉丝
		userGroup.GET("/fans", middleware.JWTMiddleWare(), controller.GetUserFans)

		//用户查看历史浏览记录
		userGroup.GET("/history", middleware.JWTMiddleWare(), controller.GetUserHistory)
		//用户删除历史浏览记录
		userGroup.DELETE("/history", middleware.JWTMiddleWare(), controller.DeleteUserHistory)
	}

	//视频路由组
	videoGroup := r.Group("/video")
	{
		//视频详情页（若已登录则添加历史浏览记录）
		videoGroup.GET("/detail", controller.VideoDetail)

		//查看该视频的评论列表
		videoGroup.GET("/comment", controller.GetCommentList)
		//用户给视频发表评论
		videoGroup.POST("/comment", middleware.JWTMiddleWare(), controller.PostVideoComment)
		//用户删除视频评论
		videoGroup.DELETE("/comment", middleware.JWTMiddleWare(), controller.DeleteVideoComment)

		//查看该视频评论的回复列表
		videoGroup.GET("/comment/reply", controller.GetReplyList)
		//用户给视频评论回复
		videoGroup.POST("/comment/reply", middleware.JWTMiddleWare(), controller.PostCommentReply)
		//用户删除视频评论回复
		videoGroup.DELETE("/comment/reply", middleware.JWTMiddleWare(), controller.DeleteCommentReply)

		//用户收藏该视频
		videoGroup.POST("/favorite", middleware.JWTMiddleWare(), controller.PostUserFavorite)

		//用户点赞该视频
		videoGroup.POST("/like", middleware.JWTMiddleWare(), controller.PostUserLike)
		//用户取消点赞视频
		videoGroup.DELETE("/like", middleware.JWTMiddleWare(), controller.DeleteUserLike)
	}

	return r
}
