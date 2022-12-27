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

type Tag struct {
	ID         int
	Name       string
	CreatedBy  string
	ModifiedBy string
	State      int

	PageNum  int
	PageSize int
}

func (t *Tag) ExistByName() (bool, error) {
	return manager.ExistTagByName(t.Name)
}

func (t *Tag) ExistByID() (bool, error) {
	return manager.ExistTagByID(t.ID)
}

func (t *Tag) Add() error {
	return manager.AddTag(t.Name, t.State, t.CreatedBy)
}

func (t *Tag) Edit() error {
	data := make(map[string]interface{})
	data["modified_by"] = t.ModifiedBy
	data["name"] = t.Name
	if t.State >= 0 {
		data["state"] = t.State
	}

	return manager.EditTag(t.ID, data)
}

func (t *Tag) Delete() error {
	return manager.DeleteTag(t.ID)
}

func (t *Tag) Count() (int, error) {
	return manager.GetTagTotal(t.getMaps())
}

func (t *Tag) GetAll() ([]manager.Tag, error) {
	var (
		tags, cacheTags []manager.Tag
	)

	cache := cache_service.Tag{
		State: t.State,

		PageNum:  t.PageNum,
		PageSize: t.PageSize,
	}
	key := cache.GetTagsKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheTags)
			return cacheTags, nil
		}
	}

	tags, err := manager.GetTags(t.PageNum, t.PageSize, t.getMaps())
	if err != nil {
		return nil, err
	}

	gredis.Set(key, tags, 3600)
	return tags, nil
}

func (t *Tag) Export() (string, error) {
	tags, err := t.GetAll()
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

	for _, v := range tags {
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
	filename := "tags-" + time + export.EXT

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

func (t *Tag) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0

	if t.Name != "" {
		maps["name"] = t.Name
	}
	if t.State >= 0 {
		maps["state"] = t.State
	}

	return maps
}
