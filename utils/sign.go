package utils

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func NewSign() *Sign {
	return &Sign{
		signHeaders: map[string]string{},
		ctx:         nil,
		secret:		 "",
		expired:	 10,
		openExpired: false,
	}
}

type Sign struct {
	//签名的头信息
	signHeaders map[string]string
	// gin Context
	ctx *gin.Context
	//应用key
	apiKey string
	//秘钥
	secret string
	//签名数据
	str string
	//过期时间(秒)
	expired int
	//开启验证签名有效期
	openExpired bool
}

//设置 gin Context
func (s *Sign) SetCtx(ctx *gin.Context) *Sign{
	s.ctx = ctx
	return s
}

//设置秘钥
func (s *Sign) SetSecret(secret string) *Sign{
	s.secret = secret
	return s
}

//设置应用key
func (s *Sign) SetApiKey(apiKey string) *Sign{
	s.apiKey = apiKey
	return s
}

//设置签名有效期
//expired 秒
func (s *Sign) SetExpired(expired int) *Sign{
	s.expired = expired
	return s
}

//设置签名有效期
func (s *Sign) SetOpenExpired(open bool) *Sign{
	s.openExpired = open
	return s
}

//获取签名头信息
func (s *Sign) getSignHeaders() (error) {
	key := s.ctx.GetHeader("HTTP_API_KEY")
	if key == "" {
		return errors.New("The HTTP_API_KEY cannot be empty")
	}
	s.signHeaders["key"] = key

	nonce := s.ctx.GetHeader("HTTP_API_NONCE")
	if nonce == "" {
		return errors.New("The HTTP_API_NONCE cannot be empty")
	}
	s.signHeaders["nonce"] = nonce

	timestamp := s.ctx.GetHeader("HTTP_API_TIMESTAMP")
	if timestamp == "" {
		return errors.New("The HTTP_API_TIMESTAMP cannot be empty")
	}
	s.signHeaders["timestamp"] = timestamp

	signature := s.ctx.GetHeader("HTTP_API_SIGNATURE")
	if signature == ""{
		return errors.New("The HTTP_API_SIGNATURE cannot be empty")
	}
	s.signHeaders["signature"] = signature

	return nil
}


//验证签名是否过期
func (s *Sign) verifyExpired() error {
	if s.openExpired == false {
		return nil
	}
	timestamp, err := strconv.ParseInt(s.signHeaders["timestamp"], 10, 64)
	if err != nil {
		return err
	}
	if time.Now().Unix()  - timestamp > int64(s.expired) {
		return errors.New("Signature expired")
	}
	return nil
}

//验证签名
func (s *Sign) Verify(str string) error {
	if s.secret == "" {
		return errors.New("The secret key cannot be empty")
	}
	if err := s.getSignHeaders(); err != nil {
		return err
	}
	if s.apiKey != s.signHeaders["key"] {
		fmt.Println(s.apiKey, s.signHeaders["key"])
		return errors.New("The application doesn't exist")
	}
	if err := s.verifyExpired(); err != nil {
		return err
	}
	if s.signHeaders["signature"] !=
		s.GenerateSignature(s.signHeaders["key"],
			s.signHeaders["nonce"],
			s.signHeaders["timestamp"],
			str) {
		return errors.New("signature error")
	}
	return nil
}

//生成签名
func (s *Sign) GenerateSignature(apikey, nonce, timestamp, str string) string {
	signStr := s.secret +
		apikey +
		nonce +
		timestamp +
		str
	return fmt.Sprintf("%x", sha1.Sum([]byte(signStr)))
}

