package tarot_service

import (
	"xiaoyuzhou/models"
	"xiaoyuzhou/pkg/util"
)

type TarotInput struct {
	ID            int    // 塔罗牌ID
	ImgUrl        string // 图片链接
	Language      string // 语言
	Pos           string // 塔罗正逆位
	CardName      string // 卡牌名字
	KeyWord       string // 卡牌解读关键词
	Constellation string // 对应星座
	People        string // 对应人物
	Element       string // 对应元素
	Enhance       string // 加强牌
	AnalyzeOne    string // 解析1
	AnalyzeTwo    string // 解析2
	PosMeaning    string // 正逆位含义
	Love          string // 爱情婚姻
	Work          string // 事业学业
	Money         string // 人际财富
	Health        string // 健康生活
	Other         string // 其他
	AnswerOne     string // 回答1
	AnswerTwo     string // 回答2
	AnswerThree   string // 回答3
	AnswerFour    string // 回答4
	AnswerFive    string // 回答5
	PageNum       int
	PageSize      int
	CreatedBy     string
	UpdatedBy     string
}

func (t *TarotInput) Add() error {
	dbData := map[string]string{
		"img_url":       t.ImgUrl,
		"language":      t.Language,
		"pos":           t.Pos,
		"card_name":     t.CardName,
		"keyword":       t.KeyWord,
		"constellation": t.Constellation,
		"people":        t.People,
		"element":       t.Element,
		"enhance":       t.Enhance,
		"analyze_one":   t.AnalyzeOne,
		"analyze_two":   t.AnalyzeTwo,
		"pos_meaning":   t.PosMeaning,
		"love":          t.Love,
		"work":          t.Work,
		"money":         t.Money,
		"health":        t.Health,
		"other":         t.Other,
		"answer_one":    t.AnswerOne,
		"answer_two":    t.AnswerTwo,
		"answer_three":  t.AnswerThree,
		"answer_four":   t.AnswerFour,
		"answer_five":   t.AnswerFive,
		"created_by":    t.CreatedBy,
		"updated_by":    t.UpdatedBy,
	}
	return models.AddTarot(dbData)
}

func (t *TarotInput) ExistByID() (bool, error) {
	return models.ExistTarotByID(t.ID)
}
func (t *TarotInput) Edit() error {

	data := make(map[string]interface{})
	if t.ImgUrl != "" {
		data["img_url"] = t.ImgUrl
	}
	if t.Pos != "" {
		data["pos"] = t.Pos
	}
	if t.Language != "" {
		data["language"] = t.Language
	}
	if t.CardName != "" {
		data["card_name"] = t.CardName
	}
	if t.KeyWord != "" {
		data["keyword"] = t.KeyWord
	}
	if t.Constellation != "" {
		data["constellation"] = t.Constellation
	}
	if t.People != "" {
		data["people"] = t.People
	}
	if t.Element != "" {
		data["element"] = t.Element
	}
	if t.Enhance != "" {
		data["enhance"] = t.Enhance
	}
	if t.AnalyzeOne != "" {
		data["analyze_one"] = t.AnalyzeOne
	}
	if t.AnalyzeTwo != "" {
		data["analyze_two"] = t.AnalyzeTwo
	}
	if t.PosMeaning != "" {
		data["pos_meaning"] = t.PosMeaning
	}
	if t.Love != "" {
		data["love"] = t.Love
	}
	if t.Work != "" {
		data["work"] = t.Work
	}
	if t.Money != "" {
		data["money"] = t.Money
	}
	if t.Health != "" {
		data["health"] = t.Health
	}
	if t.Other != "" {
		data["other"] = t.Other
	}
	if t.AnswerOne != "" {
		data["answer_one"] = t.AnswerOne
	}
	if t.AnswerTwo != "" {
		data["answer_two"] = t.AnswerTwo
	}
	if t.AnswerThree != "" {
		data["answer_three"] = t.AnswerThree
	}
	if t.AnswerFour != "" {
		data["answer_four"] = t.AnswerFour
	}
	if t.AnswerFive != "" {
		data["answer_five"] = t.AnswerFive
	}
	if t.UpdatedBy != "" {
		data["updated_by"] = t.UpdatedBy
	}
	return models.EditTarot(t.ID, data)
}

func (t *TarotInput) Get() ([]models.Tarot, int64, error) {
	cond, vals, err := util.SqlWhereBuild(t.getMaps(), "and")
	if err != nil {
		return nil, 0, err
	}
	tarots, count, err := models.GetTarots(t.PageNum, t.PageSize, cond, vals)
	if err != nil {
		return nil, 0, err
	}
	return tarots, count, nil
}

func (t *TarotInput) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	if t.ID > 0 {
		maps["id"] = t.ID
	}
	if t.CardName != "" {
		maps["card_name like"] = "%" + t.CardName + "%"
	}
	if t.Pos != "" {
		maps["pos"] = t.Pos
	}
	if t.Constellation != "" {
		maps["constellation like"] = "%" + t.Constellation + "%"
	}
	if t.KeyWord != "" {
		maps["keyword like"] = "%" + t.KeyWord + "%"
	}
	if t.Enhance != "" {
		maps["enhance like"] = "%" + t.Enhance + "%"
	}
	if t.Element != "" {
		maps["element"] = "%" + t.Element + "%"
	}
	return maps
}
