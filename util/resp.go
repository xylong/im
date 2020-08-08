package util

import (
	"encoding/json"
	"net/http"
)

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
	if jsonStr, err = json.Marshal(&H{
		Code: code,
		Data: data,
		Msg:  msg,
	}); err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-type", "application/json")
	_, _ = w.Write(jsonStr)
}

// Fail 失败响应
func Fail(w http.ResponseWriter, msg string) {
	Resp(w, -1, nil, msg)
}

// Success 成功响应
func Success(w http.ResponseWriter, data interface{}, msg string) {
	Resp(w, 0, data, msg)
}
