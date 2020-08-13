package service

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"im/model"
)

var (
	Engine *xorm.Engine
)

func init() {
	var err error
	if Engine, err = xorm.NewEngine("mysql", "root:root@(127.0.0.1:3306)/im?charset=utf8"); err != nil {
		panic(err)
	}
	Engine.ShowSQL(true)
	Engine.SetMaxOpenConns(10)
	Engine.Sync2(new(model.User), new(model.Community), new(model.Contact))
}
