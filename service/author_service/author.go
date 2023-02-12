package author_service

import (
	"time"
	"xiaoyuzhou/models"
	"xiaoyuzhou/pkg/util"
)

type AuthorInput struct {
	ID        int
	Name      string
	Gender    int
	Age       int
	Desc      string // 简介
	CreatedBy string // 创建者
	CreatedAt time.Time
	UpdatedBy string
	PageNum   int
	PageSize  int
	AvatarUrl string //头像URL
}

func (a *AuthorInput) ExistByID() (bool, error) {
	return models.ExistAuthorByID(a.ID)
}

func (a *AuthorInput) ExistByName() (bool, error) {
	return models.ExistAuthorByName(a.Name)
}

func (a *AuthorInput) Add() error {
	return models.AddAuthor(a.Name, a.Gender, a.Age, a.Desc, a.CreatedBy, a.UpdatedBy, a.AvatarUrl)
}

func (a *AuthorInput) Edit() error {

	data := make(map[string]interface{})
	data["updated_by"] = a.UpdatedBy
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

	if a.AvatarUrl != "" {
		data["avatar_url"] = a.AvatarUrl
	}
	return models.EditAuthor(a.ID, data)
}

func (a *AuthorInput) GetAll() ([]models.AuthorDto, int64, error) {
	var (
		authors []models.AuthorDto
	)

	cond, vals, err := util.SqlWhereBuild(a.getMaps(), "and")
	if err != nil {
		return nil, 0, err
	}
	authors, count, err := models.GetAuthors(a.PageNum, a.PageSize, cond, vals)

	return authors, count, nil
}

func (a *AuthorInput) Count() (int64, error) {
	cond, vals, err := util.SqlWhereBuild(a.getMaps(), "and")
	if err != nil {
		return 0, err
	}

	return models.GetAuthorTotal(cond, vals)
}

func (a *AuthorInput) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})

	if a.Name != "" {
		maps["name like"] = "%" + a.Name + "%"
	}

	if a.ID != 0 {
		maps["id ="] = a.ID
	}
	return maps
}
