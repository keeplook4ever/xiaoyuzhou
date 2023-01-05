package manager

import (
	"gorm.io/gorm"
)

type Author struct {
	Model // gorm.Model 包含了ID，CreatedAt， UpdatedAt， DeletedAt

	Name      string `gorm:"column:name;not null;unique" json:"name"`
	Gender    int    `gorm:"column:gender;not null" json:"gender"`
	Age       int    `gorm:"column:age;not null" json:"age"`
	Desc      string `gorm:"column:desc;not null" json:"desc"` // 简介
	CreatedBy string `gorm:"column:created_by;not null" json:"created_by"`
	UpdatedBy string `gorm:"column:updated_by;not null" json:"updated_by"`

	Articles []Article `json:"articles,omitempty"`
}

type AuthorDto struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Gender    int    `json:"gender"`
	Age       int    `json:"age"`
	Desc      string `json:"desc"`
	CreatedBy string `json:"created_by"`
	UpdatedBy string `json:"updated_by"`
	CreatedAt int    `json:"created_at"`
	UpdatedAt int    `json:"updated_at"`
}

func (a *Author) ToAuthorDto() AuthorDto {
	return AuthorDto{
		ID:        a.ID,
		Name:      a.Name,
		Gender:    a.Gender,
		Age:       a.Age,
		Desc:      a.Desc,
		CreatedBy: a.CreatedBy,
		UpdatedBy: a.UpdatedBy,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

// ExistAuthorByID checks if an author exists based on ID
func ExistAuthorByID(id int) (bool, error) {
	var author Author
	err := db.Model(&Author{}).Select("id").Where("id = ?", id).First(&author).Error
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
	err := db.Model(&Author{}).Select("id").Where("name = ? ", name).First(&author).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if author.ID > 0 {
		return true, nil
	}
	return false, nil
}

func AddAuthor(name string, gender int, age int, desc string, createdBy string, updatedBy string) error {
	author := Author{
		Name:      name,
		Gender:    gender,
		Age:       age,
		Desc:      desc,
		CreatedBy: createdBy,
		UpdatedBy: updatedBy,
	}
	if err := db.Create(&author).Error; err != nil {
		return err
	}
	return nil
}

func EditAuthor(id int, data interface{}) error {
	if err := db.Model(&Author{}).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

func GetAuthorTotal(maps interface{}) (int64, error) {
	var count int64
	if err := db.Model(&Author{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func GetAuthors(pageNum int, pageSize int, maps interface{}) ([]AuthorDto, error) {
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
	resp := make([]AuthorDto, 0)
	for _, a := range authors {
		resp = append(resp, a.ToAuthorDto())
	}
	return resp, nil
}
