package player

type Tarot struct {
	Id       int    `json:"id,required"`
	ImgUrl   string `json:"imgUrl,required"`
	Name     string `json:"name,required"`
	Type_    string `json:"type,required"`
	CardRead string `json:"cardRead,required"`
	TypeRead string `json:"typeRead,required"`
	Answer   string `json:"answer,required"`
}
