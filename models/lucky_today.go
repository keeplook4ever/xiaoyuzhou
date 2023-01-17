package models

import (
	"math/rand"
	"time"
)

type LuckyToday struct {
	Model
	Spell string `gorm:"column:spell;not null;type:varchar(191)" json:"spell"` //今日好运咒语
	Todo  string `gorm:"column:todo;not null;type:varchar(191)" json:"todo"`   //今日适宜
	Song  string `gorm:"column:song;not null;type:varchar(191)" json:"song"`   //今日幸运之歌
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

func GetOneRandomLuckyToday() (*LuckyTodayDto, error) {
	var luckList []LuckyToday
	if err := Db.Model(&LuckyToday{}).Find(&luckList).Error; err != nil {
		return nil, err
	}
	idList := make([]uint, 0)
	for _, v := range luckList {
		idList = append(idList, v.ID)
	}

	rand.Seed(time.Now().Unix())
	luckyChose := luckList[rand.Intn(len(idList))].ToLuckyTodayDto()

	return &luckyChose, nil
}
