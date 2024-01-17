package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// MySQL配置信息
type MySQL struct {
	Host      string //MySQL服务器主机名或IP地址
	Port      int    //MySQL服务器端口号
	Database  string //要连接的数据库名称
	Username  string //登录MySQL服务器的用户名
	Password  string //登录MySQL服务器的密码
	Charset   string //连接使用的字符集
	ParseTime bool   `toml:"parse_time"` //是否将MySQL返回的时间类型解析为Go的本地时间类型
	Loc       string //指定本地时区
}

// Redis配置信息
type Redis struct {
	IP       string
	Port     int
	Database int
}

// 服务器配置信息
type Server struct {
	IP   string //服务器IP地址
	Port int    //服务器端口号
}

// 路径配置信息
type Path struct {
	StaticSourcePath string `toml:"static_source_path"` //静态文件路径
}

type Config struct {
	DB     MySQL `toml:"mysql"`
	RDB    Redis `toml:"redis"`
	Server `toml:"server"`
	Path   `toml:"path"`
}

var Conf Config

func ensurePathValid() {
	var err error

	//使用os.Stat()函数判断StaticSourcePath对应的路径是否存在
	//如果路径不存在，则使用os.Mkdir()函数创建该路径，并设置权限为0755
	//如果创建路径失败，则使用log.Fatalf()函数输出错误信息并终止程序执行
	if _, err = os.Stat(Conf.StaticSourcePath); os.IsNotExist(err) {
		if err = os.Mkdir(Conf.StaticSourcePath, 0755); err != nil {
			log.Fatalf("mkdir error:path %s", Conf.StaticSourcePath)
		}
	}

	//使用filepath.Abs()函数将StaticSourcePath路径转换为绝对路径
	//并将结果赋值给StaticSourcePath变量
	//如果转换失败，则使用log.Fatalln()函数输出错误信息并终止程序执行
	Conf.StaticSourcePath, err = filepath.Abs(Conf.StaticSourcePath)
	if err != nil {
		log.Fatalln("get abs path failed:", Conf.StaticSourcePath)
	}
}

// 初始化函数
// 程序启动时会自动读取配置文件，并确保配置信息的有效性。这样可以避免在代码中硬编码配置信息，使得程序更具可维护性和可扩展性
func init() {
	//使用toml包的DecodeFile()函数解析配置文件config.toml
	if _, err := toml.DecodeFile("./config/config.toml", &Conf); err != nil {
		panic(err)
	}

	//去除左右的空格
	strings.Trim(Conf.Server.IP, " ")
	// strings.Trim(Conf.RDB.IP, " ") Redis 还没学
	strings.Trim(Conf.DB.Host, " ")

	//确保路径的有效性和转换为绝对路径
	ensurePathValid()
}

// 填充得到数据库连接字符串
func DBConnectString() string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%v&loc=%s",
		Conf.DB.Username, Conf.DB.Password, Conf.DB.Host, Conf.DB.Port, Conf.DB.Database,
		Conf.DB.Charset, Conf.DB.ParseTime, Conf.DB.Loc)
	log.Println(dsn) //输出一下
	return dsn
}
