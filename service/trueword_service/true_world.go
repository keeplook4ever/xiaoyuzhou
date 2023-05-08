package trueword_service

import (
	"errors"
	"fmt"
	"xiaoyuzhou/models"
	"xiaoyuzhou/pkg/logging"
	"xiaoyuzhou/pkg/util"
)

type TrueWordInput struct {
	Id        int      // id
	Lang      string   // 语言
	WordList  []string // 真言列表
	PageNum   int      // 页数
	PageSize  int      // 页大小
	UpdatedBy string   // 更新人
	CreatedBy string   // 创建人
}

func (t *TrueWordInput) Add() error {
	return models.AddTrueWord(t.Lang, t.WordList)
}

func (t *TrueWordInput) Edit() error {
	// 判断是否存在

	if !models.TrueWordExistById(t.Id) {
		return errors.New("记录不存在")
	}

	return models.EditTrueWord(t.Id, t.getMaps())
}

func (t *TrueWordInput) Delete() error {
	return models.DeleteTrueWord(t.Id)
}

func (t *TrueWordInput) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	if t.Id > 0 {
		maps["id"] = t.Id
	}
	if t.Lang != "" {
		maps["language"] = t.Lang
	}
	if t.WordList != nil {
		maps["word"] = t.WordList[0]
	}
	return maps
}

func (t *TrueWordInput) Get() ([]models.TrueWord, int64, error) {
	cond, vals, err := util.SqlWhereBuild(t.getMaps(), "and")
	if err != nil {
		logging.Error(fmt.Sprintf("Error Get(): %s", err.Error()))
		return nil, 0, err
	}
	return models.GetTrueWord(t.PageNum, t.PageSize, cond, vals)
}
