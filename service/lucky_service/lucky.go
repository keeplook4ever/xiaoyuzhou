package lucky_service

import (
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"io"
	"xiaoyuzhou/models"
	"xiaoyuzhou/pkg/logging"
)

type LuckyInputContent struct {
	PageNum  int
	PageSize int
	Language string
	Type     string `enums:"spell,todo,song"`
	Lists    []string
}

func (lk *LuckyInputContent) Add() error {
	return models.AddLucky(lk.Lists, lk.Type, lk.Language)
}

func (lk *LuckyInputContent) Get() (string, interface{}, int64, error) {
	return models.GetLuckys(lk.Type, lk.PageNum, lk.PageSize, lk.Language)
}

func Delete(xtype string, idSlice []int) error {
	return models.DeleteLucky(xtype, idSlice)
}

func EditLucky(xtype string, id int, data string, lang string) error {
	return models.EditLucky(xtype, id, data, lang)
}

func Import(r io.Reader) error {
	fiErrString := "" //最终错误信息，包含三个类别的
	xlsx, err := excelize.OpenReader(r)
	if err != nil {
		return err
	}
	TodoRows_jp, _ := xlsx.GetRows("todo-jp")
	SpellRows_jp, _ := xlsx.GetRows("spell-jp")
	SongRows_jp, _ := xlsx.GetRows("song-jp")

	TodoRows_zh, _ := xlsx.GetRows("todo-zh")
	SpellRows_zh, _ := xlsx.GetRows("spell-zh")
	SongsRows_zh, _ := xlsx.GetRows("song-zh")

	TodoRows_en, _ := xlsx.GetRows("todo-en")
	SpellRows_en, _ := xlsx.GetRows("spell-en")
	SongsRows_en, _ := xlsx.GetRows("song-en")

	TodoRows_ts, _ := xlsx.GetRows("todo-ts")
	SpellRows_ts, _ := xlsx.GetRows("spell-ts")
	SongsRows_ts, _ := xlsx.GetRows("song-ts")

	err_jp := AddLuckyByLangWithExcel(TodoRows_jp, "todo", "jp")
	err_jp1 := AddLuckyByLangWithExcel(SpellRows_jp, "spell", "jp")
	err_jp2 := AddLuckyByLangWithExcel(SongRows_jp, "song", "jp")

	err_zh := AddLuckyByLangWithExcel(TodoRows_zh, "todo", "zh")
	err_zh1 := AddLuckyByLangWithExcel(SpellRows_zh, "spell", "zh")
	err_zh2 := AddLuckyByLangWithExcel(SongsRows_zh, "song", "zh")

	err_en := AddLuckyByLangWithExcel(TodoRows_en, "todo", "en")
	err_en1 := AddLuckyByLangWithExcel(SpellRows_en, "spell", "en")
	err_en2 := AddLuckyByLangWithExcel(SongsRows_en, "song", "en")

	err_ts := AddLuckyByLangWithExcel(TodoRows_ts, "todo", "ts")
	err_ts1 := AddLuckyByLangWithExcel(SpellRows_ts, "spell", "ts")
	err_ts2 := AddLuckyByLangWithExcel(SongsRows_ts, "song", "ts")

	err_jps := ""
	err_zhs := ""
	err_ens := ""
	err_tss := ""
	if err_jp != nil || err_jp1 != nil || err_jp2 != nil {
		err_jps = "日本导入出错"
	}
	if err_zh != nil || err_zh1 != nil || err_zh2 != nil {
		err_zhs = "中文导入出错"
	}
	if err_en != nil || err_en1 != nil || err_en2 != nil {
		err_ens = "英文出错"
	}
	if err_ts != nil || err_ts1 != nil || err_ts2 != nil {
		err_tss = "繁体中文出错"
	}
	if err_jps != "" || err_zhs != "" || err_ens != "" || err_tss != "" {
		fiErrString = err_jps + err_zhs + err_ens + err_tss
	}

	if fiErrString != "" {
		logging.Error(fmt.Sprintf("Error import excel: %s", fiErrString))
		return errors.New(fiErrString)
	}

	return nil
}

func AddLuckyByLangWithExcel(dataRows [][]string, _type string, lan string) error {
	if dataRows != nil {
		var TodoData []string
		for _, row := range dataRows {
			for _, cell := range row {
				if cell != "" {
					TodoData = append(TodoData, cell)
				}
			}
		}
		if err := models.AddLucky(TodoData, _type, lan); err != nil {
			return err
		}
	}
	return nil
}
