package util

import (
	"fmt"
	"tiktok/config"
)

func GetVideoFileUrl(fileName string) string {
	base := fmt.Sprintf("http://%s:%d/static/videos/%s", config.Conf.IP, config.Conf.Port, fileName)
	return base
}

func GetCoverFileUrl(fileName string) string {
	base := fmt.Sprintf("http://%s:%d/static/covers/%s", config.Conf.IP, config.Conf.Port, fileName)
	return base
}

func GetAvatarFileUrl(fileName string) string {
	base := fmt.Sprintf("http://%s:%d/static/avatars/%s", config.Conf.IP, config.Conf.Port, fileName)
	return base
}
