package manager

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model        // gorm.Model 包含了ID，CreatedAt， UpdatedAt， DeletedAt
	Name       string `gorm:"column:name;unique;not null" json:"name" ` // 用户名唯一
	Passwd     string `gorm:"column:passwd;not null" json:"passwd"`
	CreatedBy  string `gorm:"column:created_by;not null" json:"created_by"`
	UpdatedBy  string `gorm:"column:updated_by;not null" json:"updated_by"`
}

type UserDto struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedBy string    `json:"updated_by"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) ToUserDto() UserDto {
	return UserDto{
		ID:        u.ID,
		Name:      u.Name,
		CreatedBy: u.CreatedBy,
		CreatedAt: u.CreatedAt,
		UpdatedBy: u.UpdatedBy,
		UpdatedAt: u.UpdatedAt,
	}
}

func ExistUserByName(name string) (bool, error) {
	var user User
	err := db.Model(&User{}).Select("id").Where("name = ? ", name).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if user.ID > 0 {
		return true, nil
	}
	return false, nil
}

func AddUser(name, passwd, createdBy, updatedBy string) error {
	user := User{
		Name:      name,
		Passwd:    passwd,
		CreatedBy: createdBy,
		UpdatedBy: updatedBy,
	}

	if err := db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func GetUser(maps map[string]interface{}) ([]User, error) {
	var user []User

	err := db.Where(maps).Find(&user).Error
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
