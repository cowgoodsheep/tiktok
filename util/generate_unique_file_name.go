package util

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// 获取唯一文件名
func GenerateUniqueFileName(originalName string) string {

	//创建一个随机数生成器randomGenerator，使用当前时间的纳秒级Unix时间戳作为种子。
	randomGenerator := rand.New(rand.NewSource(time.Now().UnixNano()))
	//生成一个范围在0到999之间的随机整数randomInt。
	randomInt := randomGenerator.Intn(1000)
	//创建了一个MD5哈希实例hash
	hash := md5.New()
	//将原始文件名和随机整数转换为字节数组后写入哈希实例。
	hash.Write([]byte(originalName + strconv.Itoa(randomInt)))
	//通过hash.Sum(nil)计算出MD5哈希值，并将其转换为十六进制字符串表示，存储在变量md5Hash中。
	md5Hash := hex.EncodeToString(hash.Sum(nil))

	//根据原始文件名中最后一个.的位置，将原始文件名拆分为多个部分
	//如果原始文件名包含扩展名，则将MD5哈希值与原始文件名的扩展名拼接起来作为最终的唯一文件名返回
	//否则，只返回MD5哈希值作为唯一文件名。
	parts := strings.Split(originalName, ".")
	if len(parts) > 1 {
		return md5Hash + "." + parts[len(parts)-1]
	}

	return md5Hash
}
