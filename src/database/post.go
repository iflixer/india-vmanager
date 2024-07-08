package database

type Post struct {
	ID int
}

func (f *Post) TableName() string {
	return "post"
}

func PostGet(dbService *Service, id int) (m *Post, err error) {
	m = &Post{}
	err = dbService.DB.Where("id=? AND isnull(deleted_at)", id).Find(m).Error
	return
}
