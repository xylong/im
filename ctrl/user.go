package ctrl

import (
	"fmt"
	"im/model"
	service2 "im/service"
	"im/util"
	"math/rand"
	"net/http"
)

var (
	userService service2.UserService
)

// Login 登录
func Login(w http.ResponseWriter, r *http.Request) {
	var (
		user *model.User
		err  error
	)

	if err = r.ParseForm(); err != nil {
		fmt.Println(err)
	}

	mobile := r.PostForm.Get("mobile")
	password := r.PostForm.Get("passwd")

	if user, err = userService.Login(mobile, password); err != nil {
		util.Fail(w, err.Error())
	} else {
		util.Success(w, user, "")
	}
}

// Register 注册
func Register(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		user model.User
	)
	if err = r.ParseForm(); err != nil {
		fmt.Println(err)
	}
	mobile := r.PostForm.Get("mobile")
	plainpwd := r.PostForm.Get("passwd")
	nickname := fmt.Sprintf("user%06d", rand.Int31())
	avatar := ""
	sex := model.SEX_UNKNOW

	if user, err = userService.Register(mobile, plainpwd, nickname, avatar, sex); err != nil {
		util.Fail(w, err.Error())
	} else {
		util.Success(w, user, "")
	}
}
