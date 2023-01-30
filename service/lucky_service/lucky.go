package lucky_service

import (
	"xiaoyuzhou/models"
)

type LuckyInputMany struct {
	Id       int
	Spells   []string
	Todos    []string
	Songs    []string
	PageNum  int
	PageSize int
}

type LuckyInputContent struct {
	PageNum  int
	PageSize int
	Type     string `enums:"spell,todo,song"`
	Lists    []string
}

type LuckyInputOne struct {
	Id       int
	Spell    string
	Todo     string
	Song     string
	PageNum  int
	PageSize int
}

func (lk *LuckyInputContent) Add() error {
	return models.AddLucky(lk.Lists, lk.Type)
}

//func (lk *LuckyInputOne) Delete() error {
//	return models.DeleteLucky(lk.Id)
//}

func (lk *LuckyInputContent) Get() (string, interface{}, int, error) {
	return models.GetLuckys(lk.Type, lk.PageNum, lk.PageSize)
}

//func (lk *LuckyInputOne) getMapsGet() map[string]interface{} {
//	data := make(map[string]interface{})
//	if lk.Spell != "" {
//		data["spell like"] = "%" + lk.Spell + "%"
//	}
//	if lk.Todo != "" {
//		data["todo like"] = "%" + lk.Todo + "%"
//	}
//	if lk.Song != "" {
//		data["song like"] = "%" + lk.Song + "%"
//	}
//	return data
//}
//
//func (lk *LuckyInputOne) getMapsEdit() map[string]interface{} {
//	data := make(map[string]interface{})
//	if lk.Spell != "" {
//		data["spell"] = lk.Spell
//	}
//	if lk.Todo != "" {
//		data["todo"] = lk.Todo
//	}
//	if lk.Song != "" {
//		data["song"] = lk.Song
//	}
//	return data
//}
