package manager

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID         int       `gorm:"column:id;primary_key;autoIncrement;" json:"id"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime;not null;type:datetime;default:CURRENT_TIMESTAMP;" json:"created_at"`
	ModifiedAt time.Time `gorm:"column:modified_at;autoUpdateTime;not null;type:datetime;default:CURRENT_TIMESTAMP;" json:"modified_at"`

	Name       string `gorm:"unique;not null" json:"name" ` // 用户名唯一
	Passwd     string `gorm:"not null;not null" json:"passwd"`
	CreatedBy  string `gorm:"not null;not null" json:"created_by"`
	ModifiedBy string `gorm:"not null;not null" json:"modified_by"`
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
