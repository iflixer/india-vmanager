package database

import "encoding/json"

type Post struct {
	ID   int
	Lang string
}

func (f *Post) TableName() string {
	return "post"
}

func PostGet(dbService *Service, id int) (m *Post, err error) {
	m = &Post{}
	err = dbService.DB.Where("id=? AND isnull(deleted_at)", id).Find(m).Error
	return
}

// PostAddLang adds new lang to lang field for post table
func PostAddLang(dbService *Service, postId, langId int) (err error) {
	p := &Post{}

	err = dbService.DB.Where("id=?", postId).Find(p).Error
	if err != nil {
		return
	}

	langs := []int{}
	err = json.Unmarshal([]byte(p.Lang), &langs)
	for _, l := range langs {
		if l == langId {
			return
		}
	}

	langs = append(langs, langId)

	lang, err := json.Marshal(langs)
	if err != nil {
		return
	}
	dbService.DB.Model(&Post{}).Where("id = ?", postId).Update("lang", string(lang))

	return
}
