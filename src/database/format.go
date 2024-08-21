package database

type Format struct {
	ID          int
	Format      string
	Bitrate     int
	Active      bool
	Height      int
	Quality     string
	Priority    int
	AutoConvert bool
	Name        string
	Description string
	Pass1       string
	Pass2       string
}

func (f *Format) TableName() string {
	return "video_format"
}

func (f *Format) Load(dbService *Service) (err error) {
	err = dbService.DB.Limit(1).Find(f).Error
	return
}

// FormatGetAuto returns formats which should be converter automatically
func FormatGetAuto(dbService *Service) (formats []*Format, err error) {
	err = dbService.DB.Where("auto_convert=1").Order("priority").Find(&formats).Error
	return
}
