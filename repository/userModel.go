package repository

import (
	"encoding/json"
	"minitok/log"
	"minitok/usal"
	"strconv"

	"errors"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	// gorm.Model
	Id              int64  `gorm:"column:user_id; primary_key;"`
	Name            string `gorm:"column:user_name"`
	Password        string `gorm:"column:password"`
	Follow          int64  `gorm:"column:follow_count"`
	Follower        int64  `gorm:"column:follower_count"`
	Avatar          string `gorm:"column:avatar"`
	BackgroundImage string `gorm:"column:background_image"`
	Signature       string `gorm:"column:signature"`
	TotalFav        int64  `gorm:"column:total_favorited"`
	FavCount        int64  `gorm:"column:favorite_count"`
}

func (User) TableName() string {
	return "users"
}

// 检查该用户名是否已经存在
func UserNameIsExist(userName string) error {
	db := usal.GetDB()
	user := User{}
	err := db.Where("user_name = ?", userName).Find(&user).Error
	if err == nil {
		return errors.New("username exist")
	} else if err != gorm.ErrRecordNotFound {
		return err
	}
	return nil
}

// 创建用户
func InsertUser(userName, password string) (*User, error) {
	db := usal.GetDB()
	hasedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := User{
		Name:            userName,
		Password:        string(hasedPassword),
		Follow:          0,
		Follower:        0,
		TotalFav:        0,
		FavCount:        0,
		Avatar:          "https://pica.zhimg.com/80/v2-cb7fa2dd512cf9be98403c7ea85733f2_720w.webp?source=1940ef5c",
		BackgroundImage: "https://pic1.zhimg.com/80/v2-6e1a2821499d14f25a0844815015da10_720w.webp",
		Signature:       "test sign",
	}
	result := db.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	log.Infof("regist user:%+v", user)
	go CacheSetUser(user)
	return &user, nil
}

// 获取用户信息
func GetUserInfo(u interface{}) (User, error) {
	db := usal.GetDB()
	user := User{}
	var err error
	switch u := u.(type) { // 根据传入的参数类型进行不同的查询
	case int64:
		user, err = CacheGetUser(u)
		if err == nil {
			return user, nil
		}
		err = db.Where("user_id = ?", u).Find(&user).Error

	case string:
		err = db.Where("user_name = ?", u).Find(&user).Error
	default:
		err = errors.New("")
	}
	if err != nil {
		return user, errors.New("user error")
	}
	go CacheSetUser(user)
	log.Infof("%+v", user)
	return user, nil //返回查询得到的用户信息
}

func CacheSetUser(u User) {
	uid := strconv.FormatInt(u.Id, 10)
	err := usal.CacheSet("user_"+uid, u)
	if err != nil {
		log.Errorf("set cache error:%+v", err)
	}
}

func CacheGetUser(uid int64) (User, error) {
	key := strconv.FormatInt(uid, 10)
	data, err := usal.CacheGet("user_" + key)
	user := User{}
	if err != nil {
		return user, err
	}
	err = json.Unmarshal(data, &user)
	if err != nil {
		return user, err
	}
	return user, nil
}
