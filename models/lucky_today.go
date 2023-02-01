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
	switch _type {
	case "spell":
		toAdd := make([]LuckySpell, 0)
		for _, value := range data {
			toAdd = append(toAdd, LuckySpell{Spell: value})
		}
		if err := Db.Create(&toAdd).Error; err != nil {
			return err
		}
	case "song":
		toAdd := make([]LuckySong, 0)
		for _, value := range data {
			toAdd = append(toAdd, LuckySong{Song: value})
		}
		if err := Db.Create(&toAdd).Error; err != nil {
			return err
		}
	case "todo":
		toAdd := make([]LuckyTodo, 0)
		for _, value := range data {
			toAdd = append(toAdd, LuckyTodo{Todo: value})
		}
		if err := Db.Create(&toAdd).Error; err != nil {
			return err
		}
	default:
		return errors.New("type not supported")
	}

	return nil
}

func EditLucky(xtype string, id int, data string) error {
	switch xtype {
	case "spell":
		if err := Db.Model(&LuckySpell{}).Where("id = ?", id).Update("spell", data).Error; err != nil {
			return err
		}
	case "todo":
		if err := Db.Model(&LuckyTodo{}).Where("id = ?", id).Update("todo", data).Error; err != nil {
			return err
		}
	case "song":
		if err := Db.Model(&LuckySong{}).Where("id = ?", id).Update("song", data).Error; err != nil {
			return err
		}
	default:
		return errors.New("type not supported")
	}
	return nil
}

func DeleteLucky(xtype string, idSlice []int) error {
	switch xtype {
	case "spell":
		if err := Db.Delete(&LuckySpell{}, idSlice).Error; err != nil {
			return err
		}
	case "todo":
		if err := Db.Delete(&LuckyTodo{}, idSlice).Error; err != nil {
			return err
		}
	case "song":
		if err := Db.Delete(&LuckySong{}, idSlice).Error; err != nil {
			return err
		}
	default:
		return errors.New("type not supported")
	}
	return nil
}

func GetLuckys(_type string, pageNum int, pageSize int) (string, interface{}, int64, error) {
	switch _type {
	case "spell":
		var lucks []LuckySpell
		if err := Db.Offset(pageNum).Limit(pageSize).Find(&lucks).Error; err != nil && err != gorm.ErrRecordNotFound {
			return _type, nil, 0, err
		}
		//获取总数
		var count int64
		Db.Model(&LuckySpell{}).Count(&count)
		return _type, lucks, count, nil

	case "song":
		var lucks []LuckySong
		if err := Db.Offset(pageNum).Limit(pageSize).Find(&lucks).Error; err != nil && err != gorm.ErrRecordNotFound {
			return _type, nil, 0, err
		}
		//获取总数
		var count int64
		Db.Model(&LuckySong{}).Count(&count)
		return _type, lucks, count, nil
	case "todo":
		var lucks []LuckyTodo
		if err := Db.Offset(pageNum).Limit(pageSize).Find(&lucks).Error; err != nil && err != gorm.ErrRecordNotFound {
			return _type, nil, 0, err
		}
		//获取总数
		var count int64
		Db.Model(&LuckyTodo{}).Count(&count)
		return _type, lucks, count, nil
	default:
		return _type, nil, 0, errors.New("type not supported")
	}
}
