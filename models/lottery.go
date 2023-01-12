package models

import (
	"math/rand"
	"time"
	"xiaoyuzhou/pkg/util"
)

type Lottery struct {
	Model
	MinScore    int     `gorm:"column:min_score;not null;type:tinyint(3)" json:"min_score"` // 最小分数
	MaxScore    int     `gorm:"column:max_score;not null;type:tinyint(3)" json:"max_score"` // 最大分数
	KeyWord     string  `gorm:"column:keyword;not null;type:varchar(5)" json:"keyword"`     // 运势文字
	Probability float32 `gorm:"column:probability;type:float" json:"probability"`           // 概率
}

type LotteryContent struct {
	Model
	KeyWord string `gorm:"column:keyword;not null;type:varchar(5)" json:"key_word"`
	Content string `gorm:"column:content;not null;type:text" json:"content"`
}

type LotteryDto struct {
	Score   int    `json:"score"`
	Keyword string `json:"keyword"`
	Content string `json:"content"`
}

func (l *Lottery) ToLotteryDto() LotteryDto {
	score := util.GetScore(l.MinScore, l.MaxScore)
	content, _ := GetOneRandLotteryContent(l.KeyWord)
	return LotteryDto{
		Score:   score,
		Keyword: l.KeyWord,
		Content: content,
	}
}

func GetLotteries(cond string, vals []interface{}) ([]Lottery, error) {
	var lotteries []Lottery
	err := Db.Model(&Lottery{}).Where(cond, vals...).Find(&lotteries).Error
	if err != nil {
		return nil, err
	}
	return lotteries, nil
	//resp := make([]LotteryDto, 0)
	//for _, l := range lotteries {
	//	resp = append(resp, l.ToLotteryDto())
	//}
	//return resp, nil
}

func GetLotteryContents(cond string, vals []interface{}) ([]LotteryContent, error) {
	var lotteryContents []LotteryContent
	err := Db.Model(&LotteryContent{}).Where(cond, vals...).Find(&lotteryContents).Error
	if err != nil {
		return nil, err
	}
	return lotteryContents, nil
}

func GetOneRandLotteryContent(kw string) (string, error) {
	var contents []string
	err := Db.Model(&LotteryContent{}).Pluck("content", &contents).Where("keyword = ?", kw).Error
	if err != nil {
		return "", err
	}
	// 随机取一个值
	rand.Seed(time.Now().Unix())
	return contents[rand.Intn(len(contents))], nil
}

func AddLottery(minScore, maxScore int, keyWord string, probability float32) error {
	lott := Lottery{
		MinScore:    minScore,
		MaxScore:    maxScore,
		KeyWord:     keyWord,
		Probability: probability,
	}
	err := Db.Model(&Lottery{}).Create(&lott).Error
	if err != nil {
		return err
	}
	return nil
}

func EditLottery(id int, maps map[string]interface{}) error {
	if err := Db.Model(&Lottery{}).Where("id = ? ", id).Updates(maps).Error; err != nil {
		return err
	}
	return nil
}

func AddLotteryContent(keyWord, content string) error {
	lc := LotteryContent{
		KeyWord: keyWord,
		Content: content,
	}
	err := Db.Model(&LotteryContent{}).Create(&lc).Error
	if err != nil {
		return err
	}
	return nil
}
