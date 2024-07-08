package database

type Converter struct {
	ID      int
	Name    string
	CpuQty  int
	Active  bool
	Version string
}

func (c *Converter) TableName() string {
	return "video_converter"
}

func (c *Converter) Register(dbService *Service) (err error) {
	err = dbService.DB.Where(Converter{Name: c.Name}).FirstOrCreate(c).Error
	return
}

func (c *Converter) Load(dbService *Service) (err error) {
	err = dbService.DB.Where(`id=?`, c.ID).Find(c).Error
	return
}

func (c *Converter) Save(dbService *Service) (err error) {
	err = dbService.DB.Save(c).Error
	return
}

func (c *Converter) UpdateVersion(dbService *Service) (err error) {
	err = dbService.DB.Save(c).Error
	return
}
