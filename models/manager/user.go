package manager

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	Model

	Name       string `json:"name" gorm:"unique"` // 用户名唯一
	Passwd     string `json:"passwd"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
}

func ExistUserByName(name string) (bool, error) {
	var user User
	err := db.Model(&User{}).Select("id").Where("name = ? AND deleted_on = ?", name, 0).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if user.ID > 0 {
		return true, nil
	}
	return false, nil
}

func AddUser(name, passwd string) error {
	var user User
	user.Name = name
	user.Passwd = passwd
	if err := db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func GetUser(name string) ([]User, error) {
	var user []User

	err := db.Where(&User{Name: name}).Find(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return user, nil
}

// CheckUser checks if user exists
func CheckUser(username, password string) (bool, error) {
	var user User
	err := db.Select("id").Where(User{Name: username, Passwd: password}).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if user.ID > 0 {
		return true, nil
	}

	return false, nil
}
