package category_service

import (
	"encoding/json"
	"strconv"
	"time"
	"xiaoyuzhou/models/manager"
	"xiaoyuzhou/service/manager/cache_service"

	"github.com/tealeg/xlsx"

	"xiaoyuzhou/pkg/export"
	"xiaoyuzhou/pkg/file"
	"xiaoyuzhou/pkg/gredis"
	"xiaoyuzhou/pkg/logging"
)

type Category struct {
	ID         int
	Name       string
	CreatedBy  string
	ModifiedBy string
	State      int

	PageNum  int
	PageSize int
}

func (t *Category) ExistByName() (bool, error) {
	return manager.ExistCategoryByName(t.Name)
}

func (t *Category) ExistByID() (bool, error) {
	return manager.ExistCategoryByID(t.ID)
}

func (t *Category) GetByID() (manager.Category, error) {
	return manager.GetCategoryByID(t.ID)
}

func (t *Category) Add() error {
	return manager.AddCategory(t.Name, t.State, t.CreatedBy)
}

func (t *Category) Edit() error {
	data := make(map[string]interface{})
	data["modified_by"] = t.ModifiedBy
	if t.Name != "" {
		data["name"] = t.Name
	}
	if t.State >= 0 {
		data["state"] = t.State
	}

	return manager.EditCategory(t.ID, data)
}

func (t *Category) Delete() error {
	return manager.DeleteCategory(t.ID)
}

func (t *Category) Count() (int, error) {
	return manager.GetCategoryTotal(t.getMaps())
}

func (t *Category) GetAll() ([]manager.Category, error) {
	var (
		categories, cacheTags []manager.Category
	)

	cache := cache_service.Category{
		State:    t.State,
		PageNum:  t.PageNum,
		PageSize: t.PageSize,
	}
	key := cache.GetCategoryKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheTags)
			return cacheTags, nil
		}
	}

	categories, err := manager.GetCategory(t.PageNum, t.PageSize, t.getMaps())
	if err != nil {
		return nil, err
	}

	gredis.Set(key, categories, 3600)
	return categories, nil
}

func (t *Category) Export() (string, error) {
	categories, err := t.GetAll()
	if err != nil {
		return "", err
	}

	xlsFile := xlsx.NewFile()
	sheet, err := xlsFile.AddSheet("标签信息")
	if err != nil {
		return "", err
	}

	titles := []string{"ID", "名称", "创建人", "创建时间", "修改人", "修改时间"}
	row := sheet.AddRow()

	var cell *xlsx.Cell
	for _, title := range titles {
		cell = row.AddCell()
		cell.Value = title
	}

	for _, v := range categories {
		values := []string{
			strconv.Itoa(v.ID),
			v.Name,
			v.CreatedBy,
			strconv.Itoa(v.CreatedOn),
			v.ModifiedBy,
			strconv.Itoa(v.ModifiedOn),
		}

		row = sheet.AddRow()
		for _, value := range values {
			cell = row.AddCell()
			cell.Value = value
		}
	}

	time := strconv.Itoa(int(time.Now().Unix()))
	filename := "category-" + time + export.EXT

	dirFullPath := export.GetExcelFullPath()
	err = file.IsNotExistMkDir(dirFullPath)
	if err != nil {
		return "", err
	}

	err = xlsFile.Save(dirFullPath + filename)
	if err != nil {
		return "", err
	}

	return filename, nil
}

func (t *Category) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0

	if t.Name != "" {
		maps["name"] = t.Name
	}
	if t.State >= 0 {
		maps["state"] = t.State
	}
	if t.ID > 0 {
		maps["id"] = t.ID
	}

	return maps
}
