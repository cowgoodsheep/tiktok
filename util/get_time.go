package util

import (
	"strconv"
	"time"
)

func GetDate() string {
	return strconv.Itoa(time.Now().Year()) + "-" + strconv.Itoa(int(time.Now().Month())) + "-" + strconv.Itoa(time.Now().Day())
}
