package database

import (
	"errors"
	"fmt"
	"log"
)

type Video struct {
	ID          int
	MediaId     int
	FormatId    int
	Status      int
	FileSize    int64
	Progress    int
	ConverterId int
}

func (v *Video) TableName() string {
	return "media_video"
}

// VideoFindJobForConverter searches the media without video and returns newly created video
func VideoFindJobForConverter(dbService *Service, converterId int) (media *Media, video *Video, format *Format, err error) {

	// get all media with status=0 (convert not done)
	medias, err := MediaSearchReadyToConvert(dbService)
	if err != nil || len(medias) == 0 {
		return
	}

	// get all formats
	formats, err := FormatGetAuto(dbService)
	if err != nil {
		return
	}

	// find media which has no video in some format
	media, format = searchVideo(dbService, medias, formats)
	if media == nil || format == nil {
		return
	}

	video = &Video{
		MediaId:     media.ID,
		FormatId:    format.ID,
		ConverterId: converterId,
	}
	// here we have an error "duplicate key" in case race condition. its ok!
	err = dbService.DB.Create(video).Error
	return
}

func searchVideo(dbService *Service, medias []*Media, formats []*Format) (media *Media, format *Format) {
	video := &Video{}
	var Found bool
	for _, m := range medias {
		for _, f := range formats {
			dbService.DB.Raw("SELECT EXISTS("+
				"SELECT 1 FROM "+video.TableName()+" WHERE media_id=? and format_id=?"+
				") AS found",
				m.ID, f.ID).Scan(&Found)
			if !Found {
				_, err := PostGet(dbService, m.PostID)
				if err != nil {
					continue
				}
				format = f
				media = m
				return
			}
		}
	}
	return
}

// VideoUpdateProgress should be called when converter is working
func VideoUpdateProgress(dbService *Service, videoId, progress int) (err error) {
	v := &Video{}
	err = dbService.DB.Where("id=?", videoId).Find(v).Error
	if err != nil {
		return
	}
	v.Progress = 1
	err = dbService.DB.Save(v).Error
	return
}

// VideoDone should be called when converter finished the work
func VideoDone(dbService *Service, videoId int, fileSize int64) (v *Video, err error) {
	v = &Video{}
	err = dbService.DB.Where("id=?", videoId).Find(v).Error
	if err != nil {
		return
	}
	if v.ID == 0 {
		err = errors.New(fmt.Sprintf("Video not found with id %d", videoId))
		return
	}

	v.FileSize = fileSize
	v.Status = 1
	v.Progress = 100

	err = dbService.DB.Save(v).Error
	if err != nil {
		log.Println(err)
		return
	}

	return
}
