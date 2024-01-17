package middleware

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// 定义secret
var MySecret = []byte("cowgoodsheep")

type Claims struct {
	Telephone string `json:"telephone"`
	jwt.StandardClaims
}

// 生成token
func MakeToken(telephone string) (tokenString string, err error) {
	claim := Claims{
		Telephone: telephone,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(), //过期时间
			IssuedAt:  time.Now().Unix(),                     //签发时间
			NotBefore: time.Now().Unix(),                     //生效时间
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err = token.SignedString(MySecret)
	return tokenString, err
}

// 解析token
func ParseToken(tokenString string) (*Claims, bool) {
	//调用jwt.ParseWithClaims函数来解析JWT Token
	//该函数接受三个参数：待解析的Token字符串、声明信息的结构体指针（&Claims{}）、以及一个回调函数
	token, _ := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return MySecret, nil
	})
	if token != nil {
		if key, ok := token.Claims.(*Claims); ok == true {
			//如果token无效，则返回错误
			if token.Valid == false {
				return key, false
			} else {
				return key, true
			}
		}
	}
	//如果token为空，则解析过程发生错误
	return nil, false
}

// jwt鉴权中间件,设置telephone
func JWTMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		//从请求的查询参数中获取token
		tokenString := c.Query("token")
		//如果找不到，就去Form表单去找
		if tokenString == "" {
			tokenString = c.PostForm("token")
		}
		//如果用户不存在
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "用户不存在",
			})
			c.Abort() //验证失败，跳过后续操作
			return
		}
		//验证token是否正确
		claims, ok := ParseToken(tokenString)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "token错误",
			})
			c.Abort() //验证失败，跳过后续操作
			return
		}
		//token过期了
		if time.Now().Unix() > claims.ExpiresAt {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "token过期",
			})
			c.Abort() //验证失败，跳过后续操作
			return
		}
		//token验证成功，返回用户手机号
		c.Set("telephone", claims.Telephone)
		c.Next()
	}
}
