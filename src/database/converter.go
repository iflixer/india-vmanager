package database

type Converter struct {
	ID      int
	Name    string
	CpuQty  int
	Active  bool
	Version string
	PostID  int
	MediaID int
	VideoID int
}

func (c *Converter) TableName() string {
	return "video_converter"
}

func (c *Converter) Register(dbService *Service, name string) (err error) {
	c.Name = name
	err = dbService.DB.Where("name=?", name).FirstOrCreate(c).Error
	return
}

func (c *Converter) Save(dbService *Service) (err error) {
	err = dbService.DB.Save(c).Error
	return
}
