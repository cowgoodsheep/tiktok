package dao

import (
	"fmt"
	"tiktok/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func InitMySQL() {
	var err error
	DB, err = gorm.Open("mysql", config.DBConnectString())
	if err != nil {
		fmt.Println("mysql connect fault")
		return
	}
}
