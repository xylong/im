package service

import (
	"errors"
	"fmt"
	"im/model"
	"im/util"
	"math/rand"
	"time"
)

type UserService struct {
}

//注册函数
func (s *UserService) Register(mobile, plainpwd, nickname, avatar, sex string) (user model.User, err error) {
	//检测手机号码是否存在,
	tmp := model.User{}
	_, err = Engine.Where("mobile=? ", mobile).Get(&tmp)
	if err != nil {
		return tmp, err
	}
	//如果存在则返回提示已经注册
	if tmp.Id > 0 {
		return tmp, errors.New("该手机号已经注册")
	}
	//否则拼接插入数据
	tmp.Mobile = mobile
	tmp.Avatar = avatar
	tmp.Nickname = nickname
	tmp.Sex = sex
	tmp.Salt = fmt.Sprintf("%06d", rand.Int31n(10000))
	tmp.Passwd = util.MakePasswd(plainpwd, tmp.Salt)
	tmp.Createat = time.Now()
	//token 可以是一个随机数
	tmp.Token = fmt.Sprintf("%08d", rand.Int31())
	//passwd =
	//md5 加密
	//返回新用户信息

	//插入 InserOne
	_, err = Engine.InsertOne(&tmp)
	//前端恶意插入特殊字符
	//数据库连接操作失败
	return tmp, err
}

//登录函数
func (s *UserService) Login(mobile, plainpwd string) (user *model.User, err error) {
	user = &model.User{}
	Engine.Where("mobile = ?", mobile).Get(user)
	if user.Id == 0 {
		err = errors.New("该用户不存在")
	}
	if !util.ValidatePasswd(plainpwd, user.Salt, user.Passwd) {
		err = errors.New("密码不正确")
	}
	str := fmt.Sprintf("%d", time.Now().Unix())
	user.Token = util.MD5Encode(str)
	Engine.ID(user.Id).Cols("token").Update(user)
	return user, nil
}
