package models

import (
	"gorm.io/gorm"
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
	Spell string `json:"spell"` //今日好运咒语
	Todo  string `json:"todo"`  //今日适宜
	Song  string `json:"song"`  //今日好运歌曲
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

func AddLucky(spell, todo, song string) error {
	lkToAdd := LuckyToday{
		Spell: spell,
		Todo:  todo,
		Song:  song,
	}
	if err := Db.Create(&lkToAdd).Error; err != nil {
		return err
	}
	return nil
}

func EditLucky(id int, data map[string]interface{}) error {
	if err := Db.Model(&LuckyToday{}).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

func DeleteLucky(id int) error {
	if err := Db.Where("id = ?", id).Delete(&LuckyToday{}).Error; err != nil {
		return err
	}
	return nil
}

func GetLuckys(pageNum int, pageSize int, cond string, vals []interface{}) ([]LuckyTodayDto, error) {
	var lucks []LuckyToday

	err := Db.Where(cond, vals...).Offset(pageNum).Limit(pageSize).Find(&lucks).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	resp := make([]LuckyTodayDto, len(lucks))

	for i, aa := range lucks {
		resp[i] = aa.ToLuckyTodayDto()
	}
	return resp, nil
}
