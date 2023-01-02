package user_service

import "xiaoyuzhou/models/manager"

type User struct {
	Name   string `json:"name"`
	Passwd string `json:"passwd"`
}

func (u *User) ExistByName() (bool, error) {
	return manager.ExistUserByName(u.Name)
}

func (u *User) GetUserByName() ([]manager.User, error) {
	return manager.GetUser(u.Name)
}

func (u *User) Add() error {
	return manager.AddUser(u.Name, u.Passwd)
}

func (u *User) Check() (bool, error) {
	return manager.CheckUser(u.Name, u.Passwd)
}
