package author_service

import (
	"time"
	"xiaoyuzhou/models"
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
}

func (a *AuthorInput) ExistByID() (bool, error) {
	return models.ExistAuthorByID(a.ID)
}

func (a *AuthorInput) ExistByName() (bool, error) {
	return models.ExistAuthorByName(a.Name)
}

func (a *AuthorInput) Add() error {
	return models.AddAuthor(a.Name, a.Gender, a.Age, a.Desc, a.CreatedBy, a.UpdatedBy)
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

	return models.EditAuthor(a.ID, data)
}

func (a *AuthorInput) GetAll() ([]models.AuthorDto, error) {
	var (
		authors []models.AuthorDto
	)
	authors, err := models.GetAuthors(a.PageNum, a.PageSize, a.getMaps())
	if err != nil {
		return nil, err
	}
	return authors, nil
}

func (a *AuthorInput) Count() (int64, error) {
	return models.GetAuthorTotal(a.getMaps())
}

func (a *AuthorInput) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})

	if a.Name != "" {
		maps["name"] = a.Name
	}

	if a.ID != 0 {
		maps["id"] = a.ID
	}
	return maps
}
