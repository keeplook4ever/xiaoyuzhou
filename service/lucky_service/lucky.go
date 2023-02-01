package lucky_service

import (
	"errors"
	"github.com/xuri/excelize/v2"
	"io"
	"xiaoyuzhou/models"
)

type LuckyInputContent struct {
	PageNum  int
	PageSize int
	Type     string `enums:"spell,todo,song"`
	Lists    []string
}

func (lk *LuckyInputContent) Add() error {
	return models.AddLucky(lk.Lists, lk.Type)
}

func (lk *LuckyInputContent) Get() (string, interface{}, int, error) {
	return models.GetLuckys(lk.Type, lk.PageNum, lk.PageSize)
}

func Delete(xtype string, idSlice []int) error {
	return models.DeleteLucky(xtype, idSlice)
}

func EditLucky(xtype string, id int, data string) error {
	return models.EditLucky(xtype, id, data)
}

func Import(r io.Reader) error {
	fiErrString := "" //最终错误信息，包含三个类别的
	xlsx, err := excelize.OpenReader(r)
	if err != nil {
		return err
	}
	TodoRows, _ := xlsx.GetRows("todo")
	SpellRows, _ := xlsx.GetRows("spell")
	SongRows, _ := xlsx.GetRows("song")
	for _, row := range TodoRows {
		var data []string
		for _, cell := range row {
			data = append(data, cell)
		}
		if err = models.AddLucky(data, "todo"); err != nil {
			fiErrString += err.Error()
		}
	}

	for _, row := range SpellRows {
		var data []string
		for _, cell := range row {
			data = append(data, cell)
		}
		if err = models.AddLucky(data, "spell"); err != nil {
			fiErrString += err.Error()
		}
	}

	for _, row := range SongRows {
		var data []string
		for _, cell := range row {
			data = append(data, cell)
		}
		if err = models.AddLucky(data, "song"); err != nil {
			fiErrString += err.Error()
		}
	}
	if fiErrString != "" {
		return errors.New(fiErrString)
	}

	return nil
}
