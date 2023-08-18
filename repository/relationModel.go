package repository

import (
	"errors"
	"minitok/log"
	"minitok/usal"
	"strconv"

	"github.com/jinzhu/gorm"
)

const (
	add = int64(1)
	sub = int64(-1)
)

type Relation struct {
	// gorm.Model
	Id       int64 `gorm:"column:relation_id; primary_key;"`
	Follow   int64 `gorm:"column:follow_id"`
	Follower int64 `gorm:"column:follower_id"`
}

func (Relation) TableName() string {
	return "relations"
}

func FollowAction(userId, toUserId int64) error {
	db := usal.GetDB()
	relation := Relation{
		Follow:   userId,
		Follower: toUserId,
	}
	err := db.Where("follow_id = ? and follower_id = ?", userId, toUserId).Find(&Relation{}).Error
	if err != gorm.ErrRecordNotFound {
		return errors.New("you have followed this user")
	}
	err = db.Create(&relation).Error
	if err != nil {
		return err
	}
	go CacheChangeUserCount(userId, add, "follow")
	go CacheChangeUserCount(toUserId, add, "follower")
	return nil
}

func UnFollowAction(userId, toUserId int64) error {
	db := usal.GetDB()
	err := db.Where("follow_id = ? and follower_id = ?", userId, toUserId).Delete(&Relation{}).Error
	if err != nil {
		return err
	}
	log.Debug("unfollow update user cache")
	go CacheChangeUserCount(userId, sub, "follow")
	go CacheChangeUserCount(toUserId, sub, "follower")
	return nil
}

func GetFollowList(userId int64, usertype string) ([]User, error) {
	log.Debug("try to get FollowerList!")

	db := usal.GetDB()        //获得唯一数据库
	relations := []Relation{} //用于存储用户关注的列表

	err := db.Where("relations."+usertype+"_id = ?", userId).Find(&relations).Error //用于填充relations列表
	if err == gorm.ErrRecordNotFound {
		return []User{}, nil
	} else if err != nil {
		return nil, err
	}

	numOfRelations := len(relations)
	list := make([]User, numOfRelations)

	for i, r := range relations {
		uid := int64(0)
		if usertype == "follow" {
			uid = r.Follower
		} else {
			uid = r.Follow
		}
		list[i], err = GetUserInfo(uid)
		if err != nil {
			log.Debug("get user info err AT GetFollowerlist !")
		}
	}
	return list, nil
}

func CacheChangeUserCount(userid, op int64, ftype string) {
	uid := strconv.FormatInt(userid, 10)
	mutex, _ := usal.GetLock("user_" + uid)
	defer usal.UnLock(mutex)
	user, err := CacheGetUser(userid)
	if err != nil {
		log.Infof("user:%v miss cache", userid)
		return
	}
	switch ftype {
	case "follow":
		user.Follow += op
	case "follower":
		user.Follower += op
	case "like":
		user.FavCount += op
	case "liked":
		user.TotalFav += op
	}
	CacheSetUser(user)
}
