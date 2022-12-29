package manager

import "github.com/jinzhu/gorm"

type Author struct {
	Model

	Name       string `json:"name"`
	Gender     int    `json:"gender"`
	Age        int    `json:"age"`
	Desc       string `json:"desc"` // 简介
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
}

// ExistAuthorByID checks if an author exists based on ID
func ExistAuthorByID(id int) (bool, error) {
	var author Author
	err := db.Model(&Author{}).Select("id").Where("id = ? AND deleted_on = ? ", id, 0).First(&author).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if author.ID > 0 {
		return true, nil
	}
	return false, nil
}

func ExistAuthorByName(name string) (bool, error) {
	var author Author
	err := db.Model(&Author{}).Select("id").Where("name = ? AND deleted_on = ?", name, 0).First(&author).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if author.ID > 0 {
		return true, nil
	}
	return false, nil
}

func AddAuthor(name string, gender int, age int, desc string, createdBy string) error {
	author := Author{
		Name:      name,
		Gender:    gender,
		Age:       age,
		Desc:      desc,
		CreatedBy: createdBy,
	}
	if err := db.Create(&author).Error; err != nil {
		return err
	}
	return nil
}

func EditAuthor(id int, data interface{}) error {
	if err := db.Model(&Author{}).Where("id = ? AND deleted_on = ? ", id, 0).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

func GetAuthorTotal(maps interface{}) (int, error) {
	var count int
	if err := db.Model(&Author{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func GetAuthors(pageNum int, pageSize int, maps interface{}) ([]Author, error) {
	var (
		authors []Author
		err     error
	)

	if pageSize > 0 && pageNum > 0 {
		err = db.Where(maps).Find(&authors).Offset(pageNum).Limit(pageSize).Error
	} else {
		err = db.Where(maps).Find(&authors).Error
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return authors, nil
}
