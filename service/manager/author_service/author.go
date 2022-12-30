package author_service

import (
	"encoding/json"
	"xiaoyuzhou/models/manager"
	"xiaoyuzhou/pkg/gredis"
	"xiaoyuzhou/pkg/logging"
	"xiaoyuzhou/service/manager/cache_service"
)

type Author struct {
	ID         int
	Name       string
	Gender     int
	Age        int
	Desc       string // 简介
	CreatedBy  string // 创建者
	ModifiedBy string
	PageNum    int
	PageSize   int
}

func (a *Author) ExistByID() (bool, error) {
	return manager.ExistAuthorByID(a.ID)
}

func (a *Author) ExistByName() (bool, error) {
	return manager.ExistAuthorByName(a.Name)
}

func (a *Author) Add() error {
	return manager.AddAuthor(a.Name, a.Gender, a.Age, a.Desc, a.CreatedBy)
}

func (a *Author) Edit() error {

	data := make(map[string]interface{})
	data["modified_by"] = a.ModifiedBy
	if a.Name != "" {
		data["name"] = a.Name
	}
	if a.Age > 20 && a.Age < 60 {
		data["age"] = a.Age
	}
	if a.Desc != "" {
		data["desc"] = a.Desc
	}
	if a.Gender > 0 && a.Gender < 3 {
		data["gender"] = a.Gender
	}

	return manager.EditAuthor(a.ID, data)
}

func (a *Author) GetAll() ([]manager.Author, error) {
	var (
		authors, cacheTags []manager.Author
	)

	cache := cache_service.Tag{
		PageNum:  a.PageNum,
		PageSize: a.PageSize,
	}
	key := cache.GetAuthorsKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheTags)
			return cacheTags, nil
		}
	}

	authors, err := manager.GetAuthors(a.PageNum, a.PageSize, a.getMaps())
	if err != nil {
		return nil, err
	}

	gredis.Set(key, authors, 3600)
	return authors, nil
}

func (a *Author) Count() (int, error) {
	return manager.GetAuthorTotal(a.getMaps())
}

func (a *Author) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0

	if a.Name != "" {
		maps["name"] = a.Name
	}

	if a.ID != 0 {
		maps["id"] = a.ID
	}
	return maps
}
