package database

import (
	"errors"
)

type Video struct {
	ID          int
	MediaId     int
	FormatId    int
	Status      int
	FileSize    int
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
	if err != nil {
		return
	}

	// get all formats
	formats, err := FormatGetAuto(dbService)
	if err != nil {
		return
	}

	// find media which has no video in some format
	media, format = searchMedia(dbService, medias, formats)
	if media == nil || format == nil {
		err = errors.New("media or format is nil")
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

func searchMedia(dbService *Service, medias []*Media, formats []*Format) (media *Media, format *Format) {
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
	v := &Video{
		ID:       videoId,
		Progress: progress,
	}
	err = dbService.DB.Save(v).Error
	return
}

// VideoDone should be called when converter finished the work
func VideoDone(dbService *Service, videoId, fileSize int) (v *Video, err error) {
	err = dbService.DB.Where("id=?", videoId).Find(v).Error
	if err != nil {
		return
	}

	v.FileSize = fileSize
	v.Status = 1
	v.Progress = 100

	err = dbService.DB.Save(v).Error
	if err != nil {
		return
	}

	return
}
