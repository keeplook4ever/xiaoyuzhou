package user_service

import (
	"xiaoyuzhou/models"
	"xiaoyuzhou/pkg/util"
)

type UserInput struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Passwd    string `json:"passwd"`
	CreatedBy string `json:"created_by"`
	UpdatedBy string `json:"updated_by"`
	Role      string `json:"role"`
}

func (u *UserInput) ExistByName() (bool, error) {
	return models.ExistUserByName(u.Name)
}

func (u *UserInput) GetUser() ([]models.UserDto, int64, error) {
	cond, vals, err := util.SqlWhereBuild(u.getMaps(), "and")
	if err != nil {
		return nil, 0, err
	}

	users, count, err := models.GetUser(cond, vals)
	if err != nil {
		return nil, 0, err
	}
	resp := make([]models.UserDto, 0)
	for _, u := range users {
		resp = append(resp, u.ToUserDto())
	}
	return resp, count, nil
}

func (u *UserInput) Add() error {
	return models.AddUser(u.Name, u.Passwd, u.CreatedBy, u.UpdatedBy, u.Role)
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

	if u.Role != "" {
		maps["role"] = u.Role
	}
	return maps
}
