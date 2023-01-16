package models

import (
	"errors"
	"math/rand"
	"time"
	"xiaoyuzhou/pkg/util"
)

type Lottery struct {
	Model
	MinScore    int     `gorm:"column:min_score;not null;type:tinyint(3)" json:"min_score"` // 最小分数
	MaxScore    int     `gorm:"column:max_score;not null;type:tinyint(3)" json:"max_score"` // 最大分数
	KeyWord     string  `gorm:"column:keyword;not null;type:varchar(50)" json:"keyword"`    // 运势文字
	Probability float32 `gorm:"column:probability;type:float" json:"probability"`           // 概率
	Type        string  `gorm:"column:type;type:varchar(10)" json:"type"`                   //枚举
}

type LotteryContent struct {
	Model
	Type    string `gorm:"column:type;not null;type:varchar(1)" json:"type"` //A-D 枚举
	Content string `gorm:"column:content;not null;type:text" json:"content"`
}

type LotteryDto struct {
	Score   int    `json:"score"`
	Keyword string `json:"keyword"`
	Content string `json:"content"`
}

func (l *Lottery) ToLotteryDto() LotteryDto {
	score := util.GetScore(l.MinScore, l.MaxScore)
	content, _ := GetOneRandLotteryContent(l.Type)
	return LotteryDto{
		Score:   score,
		Keyword: l.KeyWord,
		Content: content,
	}
}

func GetLotteries() ([]Lottery, error) {
	var lotteries []Lottery
	err := Db.Model(&Lottery{}).Find(&lotteries).Error
	if err != nil {
		return nil, err
	}
	return lotteries, nil
}

func GetLotteryContents(cond string, vals []interface{}) ([]LotteryContent, error) {
	var lotteryContents []LotteryContent
	err := Db.Model(&LotteryContent{}).Where(cond, vals...).Find(&lotteryContents).Error
	if err != nil {
		return nil, err
	}
	return lotteryContents, nil
}

func GetOneRandLotteryContent(tyPe string) (string, error) {
	var contents []string
	err := Db.Model(&LotteryContent{}).Pluck("content", &contents).Where("type = ?", tyPe).Error
	if err != nil {
		return "", err
	}
	// 随机取一个值
	rand.Seed(time.Now().Unix())
	return contents[rand.Intn(len(contents))], nil
}

func EditLottery(typE string, data map[string]interface{}) error {
	if err := Db.Model(&Lottery{}).Where("type = ?", typE).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

func AddLotteryContent(content, typE string) error {
	lc := LotteryContent{
		Content: content,
		Type:    typE,
	}
	err := Db.Model(&LotteryContent{}).Create(&lc).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateLotteryContent(id int, data interface{}) error {
	err := Db.Model(&LotteryContent{}).Where("id = ?", id).Updates(data).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteLotteryContent(id int) error {
	// 删除时需要确认是否该类型 >= 2
	var count int64
	Db.Model(&LotteryContent{}).Where("id = ?", id).Count(&count)
	if count >= 2 {
		err := Db.Where("id = ?", id).Delete(&LotteryContent{}).Error
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("该类型不够啦！")
}
