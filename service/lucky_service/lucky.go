package lucky_service

import (
	"xiaoyuzhou/models"
	"xiaoyuzhou/pkg/util"
)

type LuckyInput struct {
	Id       int
	Spell    string
	Todo     string
	Song     string
	PageNum  int
	PageSize int
}

func (lk *LuckyInput) Add() error {
	return models.AddLucky(lk.Spell, lk.Todo, lk.Song)
}

func (lk *LuckyInput) Edit() error {
	return models.EditLucky(lk.Id, lk.getMapsEdit())
}

func (lk *LuckyInput) Delete() error {
	return models.DeleteLucky(lk.Id)
}

func (lk *LuckyInput) Get() ([]models.LuckyTodayDto, error) {
	cond, vals, err := util.SqlWhereBuild(lk.getMapsGet(), "and")
	if err != nil {
		return nil, err
	}
	return models.GetLuckys(lk.PageNum, lk.PageSize, cond, vals)
}

func (lk *LuckyInput) getMapsGet() map[string]interface{} {
	data := make(map[string]interface{})
	if lk.Spell != "" {
		data["spell like"] = "%" + lk.Spell + "%"
	}
	if lk.Todo != "" {
		data["todo like"] = "%" + lk.Todo + "%"
	}
	if lk.Song != "" {
		data["song like"] = "%" + lk.Song + "%"
	}
	return data
}

func (lk *LuckyInput) getMapsEdit() map[string]interface{} {
	data := make(map[string]interface{})
	if lk.Spell != "" {
		data["spell"] = lk.Spell
	}
	if lk.Todo != "" {
		data["todo"] = lk.Todo
	}
	if lk.Song != "" {
		data["song"] = lk.Song
	}
	return data
}
