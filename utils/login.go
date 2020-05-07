package utils

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
)

type Login struct {
	//验证有效期
	verifyExpire bool
	//cookie 有效期(时间戳)
	expire int64
	//cookie 有效期（秒）
	maxAge int
	//cookie 路径。
	path string
	//cookie域名。
	domain string;
	// cookie 名字
	name string
	//cookie 信息
	value string;
	// gin Context
	ctx *gin.Context
}

func NewLogin(ctx *gin.Context) *Login {
	return &Login{
		verifyExpire: true,
		expire: 0,
		maxAge: 0,
		path:   "/",
		domain: "",
		name:   "login",
		value:  "",
		ctx: ctx,
	}
}

func (l *Login) SetCtx(ctx *gin.Context) *Login{
	l.ctx = ctx
	return l
}

//设置cookie 的有效期。
//verify true:验证/false:不验证
//默认验证
func (l *Login) SetVerifyExpire(verify bool) *Login{
	l.verifyExpire = verify
	return l
}

//设置cookie 的有时间。
//expire 秒
func (l *Login) SetExpire(maxAge int) *Login{
	l.expire = time.Now().Unix() + int64(maxAge)
	l.maxAge = maxAge
	return l
}

//设置cookie路径。
func (l *Login) SetPath(path string) *Login{
	if path == "" {
		path = "/"
	}
	l.path = path
	return l
}

//设置cookie域名
func (l *Login) SetDomain(domain string) *Login{
	l.domain = domain
	return l
}

//设置cookie名字
func (l *Login) SetName(name string) *Login{
	l.name = name
	return l
}

//获取客户端IP
func (l *Login) getClientIP() string{
	return l.ctx.ClientIP()
}

//获取浏览器信息
func (l *Login) getUserAgent() string{
	return l.ctx.GetHeader("HTTP_USER_AGENT")
}

//token
func (l *Login) token() string{
	str := l.value +
		l.getUserAgent() +
		l.getClientIP() +
		strconv.FormatInt(l.expire, 10)
	return l.md5(str)
}

func (l *Login) md5(str string) string{
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

//删除cookie
func (l *Login) Delete() *Login{
	l.ctx.SetCookie(l.name, "", 0, l.path, l.domain, false, true)
	return l
}

//设置cookie值
//使用此方法会最终生成cookie，放在最后调用
func (l *Login) SetValue(value string) *Login{
	l.value = value
	val := l.token() + "-" + l.value + "-" + strconv.FormatInt(l.expire, 10)
	l.ctx.SetCookie(l.name, val, int(l.maxAge), l.path, l.domain, false, true)
	return l
}

//获取cookie
func  (l *Login) GetValue() (string, error){
	if err := l.verify(); err != nil {
		return "", err
	}
	return l.value, nil
}

//设置cookie
func  (l *Login) verify() error {
	val, err := l.ctx.Cookie(l.name)
	if err != nil {
		return err
	}
	vals := strings.Split(val, "-")
	if len(vals) != 3 {
		return errors.New("Parameter error")
	}
	expire, err := strconv.ParseInt(vals[2], 10, 64)
	if err != nil {
		return err
	}
	l.value = vals[1]
	l.expire = expire
	if l.verifyExpire && l.expire < time.Now().Unix() {
		return errors.New("cookie expired")
	}
	if vals[0] != l.token() {
		return errors.New("verify failure")
	}
	return nil
}


