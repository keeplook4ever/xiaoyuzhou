package user_service

import (
	"xiaoyuzhou/models/manager"
)

type UserInput struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Passwd    string `json:"passwd"`
	CreatedBy string `json:"created_by"`
	UpdatedBy string `json:"updated_by"`
}

func (u *UserInput) ExistByName() (bool, error) {
	return manager.ExistUserByName(u.Name)
}

func (u *UserInput) GetUser() ([]manager.UserDto, error) {
	users, err := manager.GetUser(u.getMaps())
	if err != nil {
		return nil, err
	}
	resp := make([]manager.UserDto, 0)
	for _, u := range users {
		resp = append(resp, u.ToUserDto())
	}
	return resp, nil
}

func (u *UserInput) Add() error {
	return manager.AddUser(u.Name, u.Passwd, u.CreatedBy, u.UpdatedBy)
}

func (u *UserInput) Check() (bool, error) {
	return manager.CheckUser(u.Name, u.Passwd)
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
