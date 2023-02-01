package category_service

import (
	"xiaoyuzhou/models"
	"xiaoyuzhou/pkg/util"
)

type CategoryInput struct {
	ID        int
	Name      string
	CreatedBy string
	UpdatedBy string
	State     int

	PageNum  int
	PageSize int
}

func (t *CategoryInput) ExistByName() (bool, error) {
	return models.ExistCategoryByName(t.Name)
}

func (t *CategoryInput) ExistByID() (bool, error) {
	return models.ExistCategoryByID(t.ID)
}

func (t *CategoryInput) GetByID() (*models.Category, error) {
	return models.GetCategoryByID(t.ID)
}

func (t *CategoryInput) Add() error {
	return models.AddCategory(t.Name, t.State, t.CreatedBy, t.UpdatedBy)
}

func (t *CategoryInput) Edit() error {
	data := make(map[string]interface{})
	data["updated_by"] = t.UpdatedBy
	if t.Name != "" {
		data["name"] = t.Name
	}
	if t.State >= 0 {
		data["state"] = t.State
	}

	return models.EditCategory(t.ID, data)
}

func (t *CategoryInput) Delete() error {
	return models.DeleteCategory(t.ID)
}

func (t *CategoryInput) Count() (int64, error) {
	return models.GetCategoryTotal(t.getMaps())
}

func (t *CategoryInput) GetAll() ([]models.CategoryDto, int64, error) {
	var (
		categories []models.CategoryDto
	)
	cond, vals, err := util.SqlWhereBuild(t.getMaps(), "and")
	if err != nil {
		return nil, 0, err
	}
	categories, count, err := models.GetCategory(t.PageNum, t.PageSize, cond, vals)
	if err != nil {
		return nil, count, err
	}

	return categories, count, nil
}

//
//func (t *CategoryInput) Export() (string, error) {
//	categories, err := t.GetAll()
//	if err != nil {
//		return "", err
//	}
//
//	xlsFile := xlsx.NewFile()
//	sheet, err := xlsFile.AddSheet("标签信息")
//	if err != nil {
//		return "", err
//	}
//
//	titles := []string{"ID", "名称", "创建人", "创建时间", "修改人", "修改时间"}
//	row := sheet.AddRow()
//
//	var cell *xlsx.Cell
//	for _, title := range titles {
//		cell = row.AddCell()
//		cell.Value = title
//	}
//
//	for _, v := range categories {
//		values := []string{
//			strconv.Itoa(v.ID),
//			v.Name,
//			v.CreatedBy,
//			strconv.Itoa(v.CreatedAt),
//			v.UpdatedBy,
//			strconv.Itoa(v.ModifiedAt),
//		}
//
//		row = sheet.AddRow()
//		for _, value := range values {
//			cell = row.AddCell()
//			cell.Value = value
//		}
//	}
//
//	time := strconv.Itoa(int(time.Now().Unix()))
//	filename := "category-" + time + export.EXT
//
//	dirFullPath := export.GetExcelFullPath()
//	err = file.IsNotExistMkDir(dirFullPath)
//	if err != nil {
//		return "", err
//	}
//
//	err = xlsFile.Save(dirFullPath + filename)
//	if err != nil {
//		return "", err
//	}
//
//	return filename, nil
//}

func (t *CategoryInput) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})

	if t.Name != "" {
		maps["name like"] = "%" + t.Name + "%"
	}
	if t.State >= 0 {
		maps["state"] = t.State
	}
	if t.ID > 0 {
		maps["id"] = t.ID
	}

	return maps
}
