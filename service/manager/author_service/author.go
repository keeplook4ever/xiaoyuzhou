package author_service

import "xiaoyuzhou/models/manager"

type Author struct {
	ID        int
	Name      string
	Gender    int
	Age       int
	Desc      string // 简介
	CreatedBy string // 创建者
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
