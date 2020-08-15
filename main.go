package main

import (
	"html/template"
	"im/ctrl"
	"net/http"
)

func main() {
	http.Handle("/asset/", http.FileServer(http.Dir(".")))
	http.Handle("/mnt/", http.FileServer(http.Dir(".")))
	RegisterView()

	http.HandleFunc("/user/login", ctrl.Login)
	http.HandleFunc("/user/register", ctrl.Register)
	http.HandleFunc("/contact/loadcommunity", ctrl.LoadCommunity)
	http.HandleFunc("/contact/loadfriend", ctrl.LoadFriend)
	http.HandleFunc("/contact/joincommunity", ctrl.JoinCommunity)
	http.HandleFunc("/contact/addfriend", ctrl.Addfriend)
	http.HandleFunc("/chat", ctrl.Chat)
	http.HandleFunc("/attach/upload", ctrl.Upload)

	http.HandleFunc("/contact/createcommunity", ctrl.CreateCommunity)

	http.ListenAndServe(":8080", nil)
}

func RegisterView() {
	var (
		err error
		tpl *template.Template
	)
	if tpl, err = template.ParseGlob("view/**/*"); err != nil {
		panic(err)
	}
	for _, v := range tpl.Templates() {
		name := v.Name()
		http.HandleFunc(name, func(writer http.ResponseWriter, request *http.Request) {
			_ = tpl.ExecuteTemplate(writer, name, nil)
		})
	}
}

func Register(w http.ResponseWriter, r *http.Request) {

}

type H struct {
	Code int         `json:"code"`
	Data interface{} `json:"data,omitempty"`
	Msg  string      `json:"msg"`
}
