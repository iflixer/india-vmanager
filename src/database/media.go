package database

type Media struct {
	ID       int
	PostID   int
	Orig     string
	Duration int
	Status   int
	LangId   int
}

func (s *Media) TableName() string {
	return "media"
}

func (v *Media) Load(dbService *Service, id int) (err error) {
	return dbService.DB.Where("id=?", id).Find(v).Error
}

func (v *Media) Save(dbService *Service) (err error) {
	return dbService.DB.Save(v).Error
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
func MediaSearchReadyToConvert(dbService *Service, formatsQty int) (m []*Media, err error) {
	// disable telegram videos
	//err = dbService.DB.Order("id desc").Where("orig!='' AND isnull(deleted_at)").Find(&m).Error
	// err = dbService.DB.Order("id desc").Where("status=0 AND orig!='' AND orig like 'inbox/%' AND isnull(deleted_at)").Find(&m).Error
	err = dbService.DB.Raw(`SELECT m.*, count(mv.id) as mv_qty FROM media m 
		LEFT JOIN media_video mv ON (mv.media_id=m.id)
		WHERE m.orig like 'inbox/%' AND isnull(m.deleted_at)
		GROUP BY m.id
		HAVING mv_qty < ? 
		ORDER BY post_id DESC`, formatsQty).Find(&m).Error
	return
}
