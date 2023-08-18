package services

import (
	"minitok/config"
	"minitok/log"
	message "minitok/proto/pkg"
	"minitok/repository"
	"minitok/storage"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func PublishVideo(userId int64, saveFile, title string) (*message.DouyinPublishActionResponse, error) {
	client := storage.GetMinio()
	videourl, err := client.UploadFile("video", saveFile, strconv.FormatInt(userId, 10))
	//调用minIO服务器的upload方法将视频文件上传
	//video是桶的名字，saveFile是路径，第三个参数是唯一标识符
	//并且还会返回视屏的url
	if err != nil {
		return nil, err
	}
	imageFile, err := GetImageFile(saveFile) //截取一帧作为标题图片

	if err != nil {
		return nil, err
	}

	log.Debugf("imageFile %v\n", imageFile)

	picurl, err := client.UploadFile("pic", imageFile, strconv.FormatInt(userId, 10))
	//会返回图片的url地址

	if err != nil {
		//截取失败就随便找一张图片代替
		picurl = "https://p6-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/7909abe413ec4a1e82032d2beb810157~tplv-k3u1fbpfcp-zoom-in-crop-mark:1304:0:0:0.awebp?"
	}

	err = repository.InsertVideo(userId, videourl, picurl, title)
	//这是更新数据库的调用,将新增的视频条目插入到video表单中

	if err != nil {
		return nil, err
	}
	return &message.DouyinPublishActionResponse{}, nil //返回调用成功的消息
}

func PublishList(tokenUserId, userId int64) (*message.DouyinPublishListResponse, error) {
	videos, err := repository.GetVideoList(userId)
	if err != nil {
		return nil, err
	}
	list := &message.DouyinPublishListResponse{
		VideoList: VideoList(videos, tokenUserId),
	}

	return list, nil
}

func GetImageFile(videoPath string) (string, error) {
	temp := strings.Split(videoPath, "/")
	videoName := temp[len(temp)-1]
	b := []byte(videoName)
	videoName = string(b[:len(b)-3]) + "jpg"
	picpath := config.GetConfig().Path.Picfile
	picName := filepath.Join(picpath, videoName)
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-ss", "1", "-f", "image2", "-t", "0.01", "-y", picName)
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	log.Debugf(picName)
	return picName, nil
}
