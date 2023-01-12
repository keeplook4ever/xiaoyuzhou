package models

import (
	"math/rand"
	"time"
	"xiaoyuzhou/pkg/util"
)

type Lottery struct {
	Model
	MinScore    int     `gorm:"column:min_score,not null" json:"min_score"` // 最小分数
	MaxScore    int     `gorm:"column:max_score,not null" json:"max_score"` // 最大分数
	KeyWord     string  `gorm:"column:keyword,not null" json:"keyword"`     // 运势文字
	Probability float32 `gorm:"column:probability" json:"probability"`      // 概率
}

type LotteryContent struct {
	Model
	KeyWord string `gorm:"column:keyword;not null"`
	Content string `gorm:"column:content;not null"`
}

type LotteryDto struct {
	Score   int    `json:"score"`
	Keyword string `json:"keyword"`
	Content string `json:"content"`
}

func (l *Lottery) ToLotteryDto() LotteryDto {
	score := util.GetScore(l.MinScore, l.MaxScore)
	content, _ := GetLotteryContent(l.KeyWord)
	return LotteryDto{
		Score:   score,
		Keyword: l.KeyWord,
		Content: content,
	}
}

func GetLotteryContent(keyword string) (string, error) {
	var contents []string
	err := Db.Model(&LotteryContent{}).Pluck("content", &contents).Where("keyword = %s", keyword).Error
	if err != nil {
		return "", err
	}
	// 随机取一个值
	rand.Seed(time.Now().Unix())
	return contents[rand.Intn(len(contents))], nil
}
