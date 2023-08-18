package repository

import (
	"errors"
	"minitok/usal"

	"github.com/jinzhu/gorm"
)

type Favorite struct {
	Id      int64 `gorm:"column:favorite_id; primary_key;"` //这是在gorm创建数据表时的主键
	UserId  int64 `gorm:"column:user_id"`
	VideoId int64 `gorm:"column:video_id"`
}

func (Favorite) TableName() string {
	return "favorites"
}

func LikeAction(uid, vid int64) error {
	db := usal.GetDB()
	favorite := Favorite{
		UserId:  uid,
		VideoId: vid,
	}
	err := db.Where("user_id = ? and video_id = ?", uid, vid).Find(&Favorite{}).Error
	if err != gorm.ErrRecordNotFound {
		return errors.New("you have liked this video")
	} //逻辑上不可以重复点赞同一个作品
	err = db.Create(&favorite).Error //未找到点赞记录则在点赞表中创建一条新记录
	if err != nil {
		return err
	}
	authorid, _ := CacheGetAuthor(vid)

	go CacheChangeUserCount(uid, add, "like")
	go CacheChangeUserCount(authorid, add, "liked")
	return nil
}

func UnLikeAction(uid, vid int64) error {
	db := usal.GetDB()
	err := db.Where("user_id = ? and video_id = ?", uid, vid).Delete(&Favorite{}).Error
	if err != nil {
		return err
	}
	authorid, _ := CacheGetAuthor(vid)
	// go func() {
	go CacheChangeUserCount(uid, sub, "like")
	go CacheChangeUserCount(authorid, sub, "liked")
	// }()
	return nil
}

func GetFavoriteList(uid int64) ([]Video, error) {
	var videos []Video
	db := usal.GetDB()
	//使用gorm框架提供的接口查询所有被uid用户所指定的用户点赞过的视频
	err := db.Joins("left join favorites on videos.video_id = favorites.video_id").
		Where("favorites.user_id = ?", uid).Find(&videos).Error
	if err == gorm.ErrRecordNotFound {
		return []Video{}, nil
	} else if err != nil {
		return nil, err
	}
	//填写作者id
	for i, v := range videos {
		author, err := GetUserInfo(v.AuthorId)
		if err != nil {
			return videos, err
		}
		videos[i].Author = author
	}
	return videos, nil
}
