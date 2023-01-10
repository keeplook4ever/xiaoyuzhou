package models

import (
	"gorm.io/gorm"
)

type Category struct {
	Model
	Name      string    `gorm:"column:name;not null;unique" json:"name"`
	CreatedBy string    `gorm:"column:created_by;not null" json:"created_by"`
	UpdatedBy string    `gorm:"column:updated_by;not null" json:"updated_by"`
	State     int       `gorm:"column:state;not null;default:1" json:"state"` //0表示禁用，1表示启用
	Articles  []Article `json:"articles,omitempty"`
}

type CategoryDto struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	CreatedBy string `json:"created_by"`
	UpdatedBy string `json:"updated_by"`
	State     int    `json:"state"`
	CreatedAt int    `json:"created_at"`
	UpdatedAt int    `json:"updated_at"`
}

func (c *Category) ToCategoryDto() CategoryDto {
	return CategoryDto{
		ID:        c.ID,
		Name:      c.Name,
		CreatedBy: c.CreatedBy,
		UpdatedBy: c.UpdatedBy,
		State:     c.State,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

// ExistCategoryByName checks if there is a Category with the same name
func ExistCategoryByName(name string) (bool, error) {
	var tag Category
	err := Db.Select("id").Where("name = ?", name).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if tag.ID > 0 {
		return true, nil
	}

	return false, nil
}

// AddCategory Add a Category
func AddCategory(name string, state int, createdBy string, updatedBy string) error {
	tag := Category{
		Name:      name,
		State:     state,
		CreatedBy: createdBy,
		UpdatedBy: updatedBy,
	}
	if err := Db.Create(&tag).Error; err != nil {
		return err
	}

	return nil
}

// GetCategory gets a list of tags based on paging and constraints
func GetCategory(pageNum int, pageSize int, maps interface{}) ([]CategoryDto, error) {
	var (
		tags []Category
		err  error
	)

	if pageSize > 0 && pageNum > 0 {
		//err = db.Set("gorm:auto_preload", true).Where(maps).Find(&tags).Offset(pageNum).Limit(pageSize).Error
		err = Db.Where(maps).Find(&tags).Offset(pageNum).Limit(pageSize).Error
	} else {
		//err = db.Set("gorm:auto_preload", true).Where(maps).Find(&tags).Error
		err = Db.Where(maps).Find(&tags).Error
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	resp := make([]CategoryDto, 0)
	for _, c := range tags {
		resp = append(resp, c.ToCategoryDto())
	}
	return resp, nil
}

// GetCategoryTotal counts the total number of tags based on the constraint
func GetCategoryTotal(maps interface{}) (int64, error) {
	var count int64
	if err := Db.Model(&Category{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// ExistCategoryByID determines whether a Category exists based on the ID
func ExistCategoryByID(id int) (bool, error) {
	var tag Category
	err := Db.Select("id").Where("id = ? ", id).First(&tag).Error
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
	if err := Db.Where("id = ?", id).Delete(&Category{}).Error; err != nil {
		return err
	}

	return nil
}

// EditCategory modify a single Category
func EditCategory(id int, data interface{}) error {
	if err := Db.Model(&Category{}).Where("id = ? ", id).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

// CleanAllCategory clear all Category
func CleanAllCategory() (bool, error) {
	if err := Db.Unscoped().Delete(&Category{}).Error; err != nil {
		return false, err
	}

	return true, nil
}

func GetCategoryByID(id int) (tag *Category, err error) {
	err = Db.Select("id").Where("id = ? ", id).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return tag, err
}
