package ctrl

import (
	"im/args"
	service2 "im/service"
	"im/util"
	"net/http"
)

var (
	service service2.ContactService
)

func LoadFriend(w http.ResponseWriter, req *http.Request) {
	var arg args.ContactArg
	//如果这个用的上,那么可以直接
	_ = util.Bind(req, &arg)

	users := service.SearchFriend(arg.Userid)
	util.RespOkList(w, users, len(users))
}

func LoadCommunity(w http.ResponseWriter, req *http.Request) {
	var arg args.ContactArg
	//如果这个用的上,那么可以直接
	util.Bind(req, &arg)
	comunitys := service.SearchComunity(arg.Userid)
	util.RespOkList(w, comunitys, len(comunitys))
}

func JoinCommunity(w http.ResponseWriter, req *http.Request) {
	var arg args.ContactArg

	//如果这个用的上,那么可以直接
	util.Bind(req, &arg)
	err := service.JoinCommunity(arg.Userid, arg.Dstid)
	if err != nil {
		util.Fail(w, err.Error())
	} else {
		util.Success(w, nil, "")
	}
}

func Addfriend(w http.ResponseWriter, req *http.Request) {
	var arg args.ContactArg
	util.Bind(req, &arg)
	//调用service
	err := service.AddFriend(arg.Userid, arg.Dstid)
	//
	if err != nil {
		util.Fail(w, err.Error())
	} else {
		util.Success(w, nil, "好友添加成功")
	}
}
