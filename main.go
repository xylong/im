package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/user/login", Login)

	http.ListenAndServe(":8080", nil)
}

// Login 登录
func Login(w http.ResponseWriter, r *http.Request) {
	var (
		ok  bool
		err error
	)

	if err = r.ParseForm(); err != nil {
		fmt.Println(err)
	}

	mobile := r.PostForm.Get("mobile")
	password := r.PostForm.Get("passwd")

	if mobile == "19999999999" && password == "123456" {
		ok = true
	}

	if ok {
		data := map[string]interface{}{
			"id":    1,
			"token": "test",
		}
		Resp(w, 0, data, "")
	} else {
		Resp(w, -1, nil, "密码错误")
	}

}

type H struct {
	Code int         `json:"code"`
	Data interface{} `json:"data,omitempty"`
	Msg  string      `json:"msg"`
}

// Resp 响应
func Resp(w http.ResponseWriter, code int, data interface{}, msg string) {
	var (
		jsonStr []byte
		err     error
	)
	h := &H{
		Code: code,
		Data: data,
		Msg:  msg,
	}
	if jsonStr, err = json.Marshal(h); err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-type", "application/json")
	_, _ = w.Write(jsonStr)
}
