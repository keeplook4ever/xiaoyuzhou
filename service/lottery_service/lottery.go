package lottery_service

import (
	"xiaoyuzhou/models"
	"xiaoyuzhou/pkg/util"
)

type LotteryInput struct {
	MaxScore    int     `json:"max_score"`
	MinScore    int     `json:"min_score"`
	KeyWord     string  `json:"keyword"`
	Probability float32 `json:"probability"`
	Type        string  `json:"type"` //A-D 枚举四大类
}
type LotteryContentInput struct {
	Content string `json:"content"`
	Type    string `json:"type"`
	ID      int    `json:"id"`
}

func GetLotteryForPlayer() (*models.LotteryDto, error) {
	//TODO 按照相应概率抽取文字和分数
	return models.GetOneRandLottery()
}

func GetLuckyForPlayer() (*models.LuckyTodayDto, error) {
	return models.GetOneRandomLuckyToday()
}

func (l *LotteryInput) Edit() error {
	return models.EditLottery(l.Type, l.getMaps())
}

func (lc *LotteryContentInput) Add() error {
	return models.AddLotteryContent(lc.Content, lc.Type)
}

func (lc *LotteryContentInput) Update() error {
	return models.UpdateLotteryContent(lc.ID, lc.getMaps())
}

func (lc *LotteryContentInput) Delete() error {
	return models.DeleteLotteryContent(lc.ID)
}

func GetLotteryForManager() ([]models.Lottery, error) {
	return models.GetLotteries()
}

func (l *LotteryContentInput) GetLotteryContentForManager() ([]models.LotteryContent, error) {
	cond, vals, err := util.SqlWhereBuild(l.getMaps(), "and")
	if err != nil {
		return nil, err
	}
	return models.GetLotteryContents(cond, vals)
}

func (lc *LotteryContentInput) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	if lc.Type != "" {
		maps["type"] = lc.Type
	}
	if lc.ID > 0 {
		maps["id"] = lc.ID
	}

	if lc.Content != "" {
		maps["content"] = lc.Content
	}
	return maps
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
