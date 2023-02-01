package models

import (
	"gorm.io/gorm"
)

type User struct {
	Model            // gorm.Model 包含了ID，CreatedAt， UpdatedAt， DeletedAt
	Name      string `gorm:"column:name;unique;not null;type:varchar(50)" json:"name" ` // 用户名唯一
	Passwd    string `gorm:"column:passwd;not null;type:varchar(50)" json:"passwd"`
	CreatedBy string `gorm:"column:created_by;not null;type:varchar(50)" json:"created_by"`
	UpdatedBy string `gorm:"column:updated_by;not null;type:varchar(50)" json:"updated_by"`
	Role      string `gorm:"column:role;not null;default:user;type:varchar(50)" json:"role"`
}

type UserDto struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	CreatedBy string `json:"created_by"`
	CreatedAt int    `json:"created_at"`
	UpdatedBy string `json:"updated_by"`
	UpdatedAt int    `json:"updated_at"`
	Role      string `json:"role"`
}

func (u *User) ToUserDto() UserDto {
	return UserDto{
		ID:        u.ID,
		Name:      u.Name,
		CreatedBy: u.CreatedBy,
		CreatedAt: u.CreatedAt,
		UpdatedBy: u.UpdatedBy,
		UpdatedAt: u.UpdatedAt,
		Role:      u.Role,
	}
}

func ExistUserByName(name string) (bool, error) {
	var user User
	err := Db.Model(&User{}).Select("id").Where("name = ? ", name).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if user.ID > 0 {
		return true, nil
	}
	return false, nil
}

func AddUser(name, passwd, createdBy, updatedBy, role string) error {
	user := User{
		Name:      name,
		Passwd:    passwd,
		CreatedBy: createdBy,
		UpdatedBy: updatedBy,
		Role:      role,
	}

	if err := Db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func GetUser(cond string, vals []interface{}) ([]User, int64, error) {
	var user []User
	var count int64
	Db.Model(&User{}).Where(cond, vals...).Count(&count)
	err := Db.Where(cond, vals...).Find(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, count, err
	}
	return user, count, nil
}

// CheckUser checks if user exists
func CheckUser(username, password string) (bool, error) {
	var user User
	err := Db.Select("id").Where(User{Name: username, Passwd: password}).Find(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if user.ID > 0 {
		return true, nil
	}

	return false, nil
}
