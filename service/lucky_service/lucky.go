package lucky_service

import (
	"xiaoyuzhou/models"
	"xiaoyuzhou/pkg/util"
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
	Id    int
	Lists []string
}

type LuckyInputOne struct {
	Id       int
	Spell    string
	Todo     string
	Song     string
	PageNum  int
	PageSize int
}

func (lk *LuckyInputMany) Add() error {
	return models.AddLucky(lk.Spells, lk.Todos, lk.Songs)
}

func (lk *LuckyInputOne) Edit() error {
	return models.EditLucky(lk.Id, lk.getMapsEdit())
}

func (lk *LuckyInputOne) Delete() error {
	return models.DeleteLucky(lk.Id)
}

func (lk *LuckyInputOne) Get() ([]models.LuckyTodayDto, error) {
	cond, vals, err := util.SqlWhereBuild(lk.getMapsGet(), "and")
	if err != nil {
		return nil, err
	}
	return models.GetLuckys(lk.PageNum, lk.PageSize, cond, vals)
}

func (lk *LuckyInputOne) getMapsGet() map[string]interface{} {
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

func (lk *LuckyInputOne) getMapsEdit() map[string]interface{} {
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
