package database

type Media struct {
	ID            int
	PostID        int
	Orig          string
	LengthSeconds int
	Status        int
	LangId        string
}

func (s *Media) TableName() string {
	return "media"
}

/*func (s *Media) Save(dbService *Service) (err error) {
	err = dbService.DB.Save(s).Error
	return
}*/

/*func MediaGetByVideoId(dbService *Service, videoId int) (m *Media, err error) {
	err = dbService.DB.Where("").Find(m).Error
	return
}*/

// MediaGetReadyToConvert returns Media list
func MediaSearchReadyToConvert(dbService *Service) (m []*Media, err error) {
	err = dbService.DB.Order("id desc").Where("status=0 AND orig!='' AND isnull(deleted_at)").Find(&m).Error
	return
}

// MediaGet returns Media by ID
func MediaGet(dbService *Service, id int) (m *Media, err error) {
	err = dbService.DB.Where("id=? AND isnull(deleted_at)", id).Find(m).Error
	return
}

// MediaReadyToPlay should be called when we have converter at least 1 video for this media (ready to play)
func MediaReadyToPlay(dbService *Service, mediaId int) (err error) {
	m := &Media{}
	err = dbService.DB.Where("id=?", mediaId).Find(m).Error
	if err != nil {
		return
	}
	m.Status = 1
	err = dbService.DB.Save(m).Error
	return
}
