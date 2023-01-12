package lottery_service

import (
	"xiaoyuzhou/models"
	"xiaoyuzhou/pkg/util"
)

type LotteryInput struct {
	ID          int     `json:"id"`
	MinScore    int     `json:"min_score"`   // 最小分数
	MaxScore    int     `json:"max_score"`   // 最大分数
	KeyWord     string  `json:"keyword"`     // 运势文字
	Probability float32 `json:"probability"` // 概率
}
type LotteryContentInput struct {
	KeyWord string `json:"key_word"`
	Content string `json:"content"`
}

func GetLotteryForPlayer() (models.LotteryDto, error) {
	return models.LotteryDto{}, nil
}

func GetLuckyForPlayer() (models.LuckyTodayDto, error) {
	return models.LuckyTodayDto{}, nil
}

func (l *LotteryInput) Add() error {
	return models.AddLottery(l.MinScore, l.MaxScore, l.KeyWord, l.Probability)
}

func (l *LotteryInput) Edit() error {
	return models.EditLottery(l.ID, l.getMaps())
}

func (lc *LotteryContentInput) Add() error {
	return models.AddLotteryContent(lc.KeyWord, lc.Content)
}

func (l *LotteryInput) GetLotteryForManager() ([]models.Lottery, error) {
	cond, vals, err := util.SqlWhereBuild(l.getMaps(), "and")
	if err != nil {
		return nil, err
	}
	return models.GetLotteries(cond, vals)
}

func (l *LotteryInput) GetLotteryContentForManager() ([]models.LotteryContent, error) {
	cond, vals, err := util.SqlWhereBuild(l.getMaps(), "and")
	if err != nil {
		return nil, err
	}
	return models.GetLotteryContents(cond, vals)
}

func (l *LotteryInput) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	if l.MaxScore > 0 && l.MaxScore < 100 {
		maps["max_score"] = l.MaxScore
	}
	if l.MinScore > 0 && l.MinScore < 100 {
		maps["min_score"] = l.MinScore
	}
	if l.KeyWord != "" {
		maps["keyword"] = l.KeyWord
	}
	if l.Probability != 0.0 {
		maps["probability"] = l.Probability
	}
	return maps
}
