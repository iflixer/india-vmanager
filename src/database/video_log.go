package database

import (
	"log"
)

type VideoLog struct {
	ID          int
	Status      int
	ConverterId int
	PostId      int
	MediaId     int
	VideoId     int
	Message     string
}

func (f *VideoLog) TableName() string {
	return "video_log"
}

func VideoLogAdd(dbService *Service, status, converterId, postId, mediaId, videoId int, message string) {
	m := &VideoLog{
		Status:      status,
		ConverterId: converterId,
		PostId:      postId,
		MediaId:     mediaId,
		VideoId:     videoId,
		Message:     message,
	}
	err := dbService.DB.Save(m).Error
	if err != nil {
		log.Println(err)
	}
	return
}
