package lottery_service

import (
	"xiaoyuzhou/models"
	"xiaoyuzhou/pkg/util"
)

type LotteryInput struct {
	KeyWordList     []string  `json:"keyword_list"`
	ScoreList       []int     `json:"score_list"`
	ProbabilityList []float32 `json:"probability_list"`
}
type LotteryContentInput struct {
	KeyWord string `json:"key_word"`
	Content string `json:"content"`
	ID      int    `json:"id"`
}

func GetLotteryForPlayer() (models.LotteryDto, error) {
	return models.LotteryDto{}, nil
}

func GetLuckyForPlayer() (models.LuckyTodayDto, error) {
	return models.LuckyTodayDto{}, nil
}

func (l *LotteryInput) Add() error {
	return models.AddLottery(l.KeyWordList, l.ScoreList, l.ProbabilityList)
}

func (l *LotteryInput) Edit() error {
	return models.EditLottery(l.KeyWordList, l.ScoreList, l.ProbabilityList)
}

func (lc *LotteryContentInput) Add() error {
	return models.AddLotteryContent(lc.KeyWord, lc.Content)
}

func (lc *LotteryContentInput) Update() error {
	return models.UpdateLotteryContent(lc.ID, lc.getMaps())
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
	if lc.KeyWord != "" {
		maps["keyword"] = lc.KeyWord
	}
	if lc.ID > 0 {
		maps["id"] = lc.ID
	}

	if lc.Content != "" {
		maps["content"] = lc.Content
	}
	return maps
}

//func (l *LotteryInput) getMaps() map[string]interface{} {
//	maps := make(map[string]interface{})
//	if l.MaxScore > 0 && l.MaxScore < 100 {
//		maps["max_score"] = l.MaxScore
//	}
//	if l.MinScore > 0 && l.MinScore < 100 {
//		maps["min_score"] = l.MinScore
//	}
//	if l.KeyWord != "" {
//		maps["keyword"] = l.KeyWord
//	}
//	if l.Probability != 0.0 {
//		maps["probability"] = l.Probability
//	}
//	return maps
//}
