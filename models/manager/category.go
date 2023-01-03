package manager

import (
	"github.com/jinzhu/gorm"
)

type Category struct {
	Model

	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"` //0表示禁用，1表示启用

	Articles []Article `json:"articles,omitempty"`
}

// ExistCategoryByName checks if there is a Category with the same name
func ExistCategoryByName(name string) (bool, error) {
	var tag Category
	err := db.Select("id").Where("name = ? AND deleted_on = ? ", name, 0).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if tag.ID > 0 {
		return true, nil
	}

	return false, nil
}

// AddCategory Add a Category
func AddCategory(name string, state int, createdBy string) error {
	tag := Category{
		Name:      name,
		State:     state,
		CreatedBy: createdBy,
	}
	if err := db.Create(&tag).Error; err != nil {
		return err
	}

	return nil
}

// GetCategory gets a list of tags based on paging and constraints
func GetCategory(pageNum int, pageSize int, maps interface{}) ([]Category, error) {
	var (
		tags []Category
		err  error
	)

	if pageSize > 0 && pageNum > 0 {
		//err = db.Set("gorm:auto_preload", true).Where(maps).Find(&tags).Offset(pageNum).Limit(pageSize).Error
		err = db.Where(maps).Find(&tags).Offset(pageNum).Limit(pageSize).Error
	} else {
		//err = db.Set("gorm:auto_preload", true).Where(maps).Find(&tags).Error
		err = db.Where(maps).Find(&tags).Error
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return tags, nil
}

// GetCategoryTotal counts the total number of tags based on the constraint
func GetCategoryTotal(maps interface{}) (int, error) {
	var count int
	if err := db.Model(&Category{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// ExistCategoryByID determines whether a Category exists based on the ID
func ExistCategoryByID(id int) (bool, error) {
	var tag Category
	err := db.Select("id").Where("id = ? AND deleted_on = ? ", id, 0).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if tag.ID > 0 {
		return true, nil
	}

	return false, nil
}

// DeleteCategory delete a Category
func DeleteCategory(id int) error {
	if err := db.Where("id = ?", id).Delete(&Category{}).Error; err != nil {
		return err
	}

	return nil
}

// EditCategory modify a single Category
func EditCategory(id int, data interface{}) error {
	if err := db.Model(&Category{}).Where("id = ? AND deleted_on = ? ", id, 0).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

// CleanAllCategory clear all Category
func CleanAllCategory() (bool, error) {
	if err := db.Unscoped().Where("deleted_on != ? ", 0).Delete(&Category{}).Error; err != nil {
		return false, err
	}

	return true, nil
}

func GetCategoryByID(id int) (tag Category, err error) {
	err = db.Select("id").Where("id = ? AND deleted_on = ? ", id, 0).First(&tag).Error
	//if err != nil && err != gorm.ErrRecordNotFound {
	//	return nil, err
	//}
	return tag, err
}
