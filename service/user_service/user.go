package user_service

import (
	"xiaoyuzhou/models"
)

type UserInput struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Passwd    string `json:"passwd"`
	CreatedBy string `json:"created_by"`
	UpdatedBy string `json:"updated_by"`
}

func (u *UserInput) ExistByName() (bool, error) {
	return models.ExistUserByName(u.Name)
}

func (u *UserInput) GetUser() ([]models.UserDto, error) {
	users, err := models.GetUser(u.getMaps())
	if err != nil {
		return nil, err
	}
	resp := make([]models.UserDto, 0)
	for _, u := range users {
		resp = append(resp, u.ToUserDto())
	}
	return resp, nil
}

func (u *UserInput) Add() error {
	return models.AddUser(u.Name, u.Passwd, u.CreatedBy, u.UpdatedBy)
}

func (u *UserInput) Check() (bool, error) {
	return models.CheckUser(u.Name, u.Passwd)
}

func (u *UserInput) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})

	if u.Name != "" {
		maps["name"] = u.Name
	}

	if u.ID > 0 {
		maps["id"] = u.ID
	}

	return maps
}
