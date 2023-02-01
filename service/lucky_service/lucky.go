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

func (lk *LuckyInputContent) Get() (string, interface{}, int64, error) {
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
	if TodoRows != nil {
		var TodoData []string
		for _, row := range TodoRows {
			for _, cell := range row {
				if cell != "" {
					TodoData = append(TodoData, cell)
				}
			}
		}
		if err = models.AddLucky(TodoData, "todo"); err != nil {
			fiErrString += err.Error()
		}
	}
	if SpellRows != nil {
		var SpellData []string
		for _, row := range SpellRows {
			for _, cell := range row {
				if cell != "" {
					SpellData = append(SpellData, cell)
				}
			}
		}
		if err = models.AddLucky(SpellData, "spell"); err != nil {
			fiErrString += err.Error()
		}
	}
	if SongRows != nil {
		var SongData []string
		for _, row := range SongRows {
			for _, cell := range row {
				if cell != "" {
					SongData = append(SongData, cell)
				}
			}
		}
		if err = models.AddLucky(SongData, "song"); err != nil {
			fiErrString += err.Error()
		}
	}

	if fiErrString != "" {
		return errors.New(fiErrString)
	}

	return nil
}
