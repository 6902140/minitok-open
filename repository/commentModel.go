package repository

import (
	"encoding/json"
	"minitok/log"
	"minitok/usal"
	"strconv"
	"time"
)

type Comment struct {
	CommentId int64  `gorm:"column:comment_id; primary_key;"`
	UserId    int64  `gorm:"column:user_id"`
	VideoId   int64  `gorm:"column:video_id"`
	Comment   string `gorm:"column:comment"`
	Time      string `gorm:"column:time"`
}

func CommentAdd(userId, videoId int64, comment_text string) (*Comment, error) {
	db := usal.GetDB()

	nowtime := time.Now().Format("01-02")
	comment := Comment{
		UserId:  userId,
		VideoId: videoId,
		Comment: comment_text,
		Time:    nowtime,
	}
	result := db.Create(&comment)

	if result.Error != nil {
		return nil, result.Error
	}
	log.Infof("comment:%+v", comment)

	CacheDelCommentAll(videoId)

	return &comment, nil
}

func CommentDelete(videoId, comment_id int64) error {
	db := usal.GetDB()
	commentTemp := Comment{}

	err := db.Where("comment_id = ?", comment_id).Take(&commentTemp).Error
	if err != nil {
		return err
	}

	CacheDelCommentAll(videoId)

	db.Delete(&commentTemp)
	return nil

}

func CommentList(videoId int64) ([]Comment, error) {
	var comments []Comment
	db := usal.GetDB()
	var err error

	comments, _ = CacheGetComment(videoId)
	log.Infof("comments-------------------------:%+v\n", comments)

	if comments != nil {
		return comments, nil
	}

	err = db.Where("video_id = ?", videoId).Order("comment_id DESC").Find(&comments).Error
	//查询所有评论进行返回
	CacheSetComment(videoId, comments)
	log.Infof("comments:%+v", comments)

	if err != nil {
		return nil, err
	}
	return comments, nil

}

func CacheSetComment(videoId int64, c []Comment) {

	vid := strconv.FormatInt(videoId, 10)
	err := usal.CacheHSet("comment", vid, c)
	if err != nil {
		log.Errorf("set cache error:%+v", err)
	}
}

//}

func CacheGetComment(vid int64) ([]Comment, error) {
	key := strconv.FormatInt(vid, 10)
	data, err := usal.CacheHGet("comment", key)

	//var comments = make([]map[string]interface{},2)

	var comments []Comment
	//comment := Comment{}
	if err != nil {
		return comments, err
	}
	err = json.Unmarshal(data, &comments)
	if err != nil {
		return comments, err
	}
	return comments, nil
}

func CacheDelCommentAll(videoId int64) {

	vid := strconv.FormatInt(videoId, 10)
	err := usal.CacheDelHash("comment", vid)
	if err != nil {
		log.Errorf("set cache error:%+v", err)
	}
}

/* func CacheDelCommentOne(videoId, comment_id int64) {

	vid := strconv.FormatInt(videoId, 10)
	cid := strconv.FormatInt(comment_id, 10)
	err := usal.CacheDelHash2("comment", vid, cid)
	if err != nil {
		log.Errorf("set cache error:%+v", err)
	}
} */
