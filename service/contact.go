package service

import (
	"errors"
	"im/model"
	"time"
)

type ContactService struct {
}

func (c *ContactService) AddFriend(userId, dstId int64) error {
	if userId == dstId {
		return errors.New("不能添加自己为好友")
	}
	// 是否已经是好友
	contact := model.Contact{}
	Engine.Where("ownerid=?", userId).And("dstid=?", dstId).And("cate=?", model.CONCAT_CATE_USER).Get(&contact)
	if contact.Id > 0 {
		return errors.New("已经是好友")
	}
	// 事务
	session := Engine.NewSession()
	session.Begin()
	_, e1 := session.InsertOne(model.Contact{
		Ownerid:  userId,
		Dstobj:   dstId,
		Cate:     model.CONCAT_CATE_USER,
		Createat: time.Now(),
	})
	_, e2 := session.InsertOne(model.Contact{
		Ownerid:  dstId,
		Dstobj:   userId,
		Cate:     model.CONCAT_CATE_USER,
		Createat: time.Now(),
	})
	if e1 == nil && e2 == nil {
		session.Commit()
		return nil
	} else {
		session.Rollback()
		if e1 != nil {
			return e1
		} else {
			return e2
		}
	}
}

func (c *ContactService) SearchComunity(userId int64) []model.Community {
	contacts := make([]model.Contact, 0)
	comIds := make([]int64, 0)

	_ = Engine.Where("ownerid = ? and cate = ?", userId, model.CONCAT_CATE_COMUNITY).Find(&contacts)
	for _, v := range contacts {
		comIds = append(comIds, v.Dstobj)
	}
	coms := make([]model.Community, 0)
	if len(comIds) == 0 {
		return coms
	}
	_ = Engine.In("id", comIds).Find(&coms)
	return coms
}

func (c *ContactService) JoinCommunity(userId, comId int64) error {
	cot := model.Contact{
		Ownerid: userId,
		Dstobj:  comId,
		Cate:    model.CONCAT_CATE_COMUNITY,
	}
	Engine.Get(&cot)
	if cot.Id == 0 {
		cot.Createat = time.Now()
		_, err := Engine.InsertOne(cot)
		return err
	}
	return nil
}

func (c *ContactService) CreateCommunity(community model.Community) (comm model.Community, err error) {
	if community.Name == "" {
		err = errors.New("缺少群名称")
		return
	}
	if community.Ownerid == 0 {
		err = errors.New("请先登录")
		return
	}
	comm = model.Community{
		Ownerid: community.Ownerid,
	}
	num, err := Engine.Count(&comm)
	if num > 5 {
		err = errors.New("一个用户最多只能创见5个群")
		return
	}
	community.Createat = time.Now()
	session := Engine.NewSession()
	_ = session.Begin()
	if _, err = session.InsertOne(&community); err != nil {
		return
	}
	if _, err = session.InsertOne(model.Contact{
		Ownerid:  community.Ownerid,
		Dstobj:   community.Id,
		Cate:     model.CONCAT_CATE_USER,
		Createat: time.Now(),
	}); err != nil {
		_ = session.Rollback()
	}
	_ = session.Commit()
	return
}

func (c *ContactService) SearchFriend(userId int64) (users []model.User) {
	contacts := make([]model.Contact, 0)
	objIds := make([]int64, 0)
	_ = Engine.Where("ownerid = ? and cate = ?", userId, model.CONCAT_CATE_USER).Find(&contacts)
	for _, v := range contacts {
		objIds = append(objIds, v.Dstobj)
	}
	if len(objIds) == 0 {
		return
	}
	_ = Engine.In("id", objIds).Find(&users)
	return
}

func (c *ContactService) SearchComunityIds(userId int64) (ids []int64) {
	conacts := make([]model.Contact, 0)
	_ = Engine.Where("ownerid = ? and cate = ?", userId, model.CONCAT_CATE_COMUNITY).Find(&conacts)
	for _, v := range conacts {
		ids = append(ids, v.Dstobj)
	}
	return
}
