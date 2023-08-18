package storage

import (
	"minitok/config"
	"minitok/log"
	"minitok/util"
	"strconv"
	"strings"

	"github.com/minio/minio-go/v6"
)

//封装了与 MinIO 相关的配置和客户端对象，方便在代码中使用和管理。

type Minio struct {
	MinioClient  *minio.Client
	endpoint     string
	port         string
	VideoBuckets string
	PicBuckets   string
}

var client Minio

func GetMinio() Minio {
	return client
}

func InitMinio() {
	conf := config.GetConfig()
	endpoint := conf.Minio.Host
	port := conf.Minio.Port
	endpoint = endpoint + ":" + port
	accessKeyID := conf.Minio.AccessKeyID
	secretAccessKey := conf.Minio.SecretAccessKey
	videoBucket := conf.Minio.Videobuckets
	picBucket := conf.Minio.Picbuckets
	//useSSL := false
	//读取配置信息

	// 初使化 minio client对象。
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, false) //默认`不使用ssl，如果需要使用，请将openSSL改为true
	if err != nil {
		log.Error(err)
	}

	//创建存储桶 一个用来存储上传视频 另外一个用于存储图片文件
	creatBucket(minioClient, videoBucket)
	creatBucket(minioClient, picBucket)
	client = Minio{minioClient, endpoint, port, videoBucket, picBucket}
}

// 制作存储桶的函数 m是Client的指针 bucket是桶的名字
func creatBucket(m *minio.Client, bucket string) {
	// log.Debug("bucketname", bucket)
	found, err := m.BucketExists(bucket)
	if err != nil {
		log.Errorf("check %s bucketExists err:%s", bucket, err.Error())
	}
	if !found {
		err := m.MakeBucket(bucket, "us-east-1")
		if err != nil {
			return
		}
	}
	//设置桶策略
	policy := `{"Version": "2012-10-17",
				"Statement": 
					[{
						"Action":["s3:GetObject"],
						"Effect": "Allow",
						"Principal": {"AWS": ["*"]},
						"Resource": ["arn:aws:s3:::` + bucket + `/*"],
						"Sid": ""
					}]
				}`
	err = m.SetBucketPolicy(bucket, policy)
	if err != nil {
		log.Errorf("SetBucketPolicy %s  err:%s", bucket, err.Error())
	}
}

func (m *Minio) UploadFile(filetype, file, userID string) (string, error) {
	var fileName strings.Builder
	var contentType, Suffix, bucket string
	if filetype == "video" {
		contentType = "video/mp4"
		Suffix = ".mp4"
		bucket = m.VideoBuckets
	} else {
		contentType = "image/jpeg"
		Suffix = ".jpg"
		bucket = m.PicBuckets
	}
	fileName.WriteString(userID)
	fileName.WriteString("_")
	fileName.WriteString(strconv.FormatInt(util.GetCurrentTime(), 10))
	fileName.WriteString(Suffix)
	n, err := m.MinioClient.FPutObject(bucket, fileName.String(), file, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		log.Errorf("upload file error:%s", err.Error())
		return "", err
	}
	log.Infof("upload file %dbyte success,fileName:%s", n, fileName)
	url := "http://" + m.endpoint + "/" + bucket + "/" + fileName.String()
	return url, nil
}
