package models

import (
	"errors"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

type LuckySpell struct {
	Model
	Spell string `gorm:"column:spell;not null;type:varchar(191)" json:"spell"` //今日好运咒语
}

type LuckyTodo struct {
	Model
	Todo string `gorm:"column:todo;not null;type:varchar(191)" json:"todo"` //今日适宜
}

type LuckySong struct {
	Model
	Song string `gorm:"column:song;not null;type:varchar(191)" json:"song"` //今日幸运之歌
}

type LuckyTodayDto struct {
	Spell string `json:"spell"` //今日好运咒语
	Todo  string `json:"todo"`  //今日适宜
	Song  string `json:"song"`  //今日好运歌曲
}

//func (l *LuckyToday) ToLuckyTodayDto() LuckyTodayDto {
//	return LuckyTodayDto{
//		Spell: l.Spell,
//		Todo:  l.Todo,
//		Song:  l.Song,
//	}
//}

// GetOneRandomLuckyToday 为用户产生一个随机的今日好运内容
func GetOneRandomLuckyToday() (*LuckyTodayDto, error) {
	var luckToday LuckyTodayDto
	var luckSpells []LuckySpell
	var luckTodos []LuckyTodo
	var luckSongs []LuckySong
	if err := Db.Model(&LuckyTodo{}).Find(&luckTodos).Error; err != nil {
		return nil, err
	}

	if err := Db.Model(&LuckySong{}).Find(&luckSongs).Error; err != nil {
		return nil, err
	}

	if err := Db.Model(&LuckySpell{}).Find(&luckSpells).Error; err != nil {
		return nil, err
	}
	if len(luckSpells) == 0 || len(luckSongs) == 0 || len(luckTodos) == 0 {
		return nil, errors.New("没有内容了")
	}
	rand.Seed(time.Now().Unix())
	luckySpellChose := luckSpells[rand.Intn(len(luckSpells))]
	rand.Seed(time.Now().Unix())
	luckSongChose := luckSongs[rand.Intn(len(luckSongs))]
	rand.Seed(time.Now().Unix())
	luckTodoChose := luckTodos[rand.Intn(len(luckTodos))]
	luckToday.Todo = luckTodoChose.Todo
	luckToday.Spell = luckySpellChose.Spell
	luckToday.Song = luckSongChose.Song

	return &luckToday, nil
}

func AddLucky(data []string, _type string) error {
	if _type == "spell" {
		toAdd := make([]LuckySpell, 0)
		for _, value := range data {
			toAdd = append(toAdd, LuckySpell{Spell: value})
		}
		if err := Db.Create(&toAdd).Error; err != nil {
			return err
		}
	}

	if _type == "song" {
		toAdd := make([]LuckySong, 0)
		for _, value := range data {
			toAdd = append(toAdd, LuckySong{Song: value})
		}
		if err := Db.Create(&toAdd).Error; err != nil {
			return err
		}
	}

	if _type == "todo" {
		toAdd := make([]LuckyTodo, 0)
		for _, value := range data {
			toAdd = append(toAdd, LuckyTodo{Todo: value})
		}
		if err := Db.Create(&toAdd).Error; err != nil {
			return err
		}
	}

	return nil
}

func EditLuckySpell(id int, data map[string]interface{}) error {
	if err := Db.Model(&LuckySpell{}).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

//func DeleteLucky(id int) error {
//	if err := Db.Where("id = ?", id).Delete(&LuckyToday{}).Error; err != nil {
//		return err
//	}
//	return nil
//}

func GetLuckys(_type string, pageNum int, pageSize int) (string, interface{}, int, error) {
	if _type == "spell" {
		var lucks []LuckySpell
		if err := Db.Offset(pageNum).Limit(pageSize).Find(&lucks).Error; err != nil && err != gorm.ErrRecordNotFound {
			return _type, nil, 0, err
		}
		return _type, lucks, len(lucks), nil
	} else if _type == "song" {
		var lucks []LuckySong
		if err := Db.Offset(pageNum).Limit(pageSize).Find(&lucks).Error; err != nil && err != gorm.ErrRecordNotFound {
			return _type, nil, 0, err
		}
		return _type, lucks, len(lucks), nil
	} else if _type == "todo" {
		var lucks []LuckyTodo
		if err := Db.Offset(pageNum).Limit(pageSize).Find(&lucks).Error; err != nil && err != gorm.ErrRecordNotFound {
			return _type, nil, 0, err
		}
		return _type, lucks, len(lucks), nil
	}
	return "", nil, 0, nil
}
