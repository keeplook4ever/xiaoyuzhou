package models

type LuckyToday struct {
	Model
	Spell string `gorm:"column:spell,not null" json:"spell"` //今日好运咒语
	Todo  string `gorm:"column:todo,not null" json:"todo"`   //今日适宜
	Song  string `gorm:"column:song,not null" json:"song"`   //今日幸运之歌
}

type LuckyTodayDto struct {
	Spell string `json:"spell"`
	Todo  string `json:"todo"`
	Song  string `json:"song"`
}

func (l *LuckyToday) ToLuckyTodayDto() LuckyTodayDto {
	return LuckyTodayDto{
		Spell: l.Spell,
		Todo:  l.Todo,
		Song:  l.Song,
	}
}
