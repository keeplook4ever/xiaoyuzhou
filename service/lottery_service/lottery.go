package lottery_service

import (
	"fmt"
	"xiaoyuzhou/models"
	"xiaoyuzhou/pkg/logging"
	"xiaoyuzhou/pkg/util"
)

type LotteryInput struct {
	MaxScore    int
	MinScore    int
	KeyWord     string
	Probability float32
	Type        string //A-D 枚举四大类
}
type LotteryContentInput struct {
	Content  string
	Type     string
	ID       int
	PageNum  int
	PageSize int
	Language string
}

func GetLotteryForPlayer(uid, language string) (*models.LotteryDto, error) {
	cc, err := models.GetOneRandLottery(language)
	if err != nil {
		logging.Error(fmt.Sprintf("Eroor get one rand lottery: %s", err.Error()))
		return nil, err
	}
	return cc, nil
}

func GetLuckyForPlayer(language string) (*models.LuckyTodayDto, error) {
	return models.GetOneRandomLuckyToday(language)
}

func (l *LotteryInput) Edit() error {
	return models.EditLottery(l.Type, l.getMapsEdit())
}

func (lc *LotteryContentInput) Add() error {
	return models.AddLotteryContent(lc.Content, lc.Type, lc.Language)
}

func (lc *LotteryContentInput) Update() error {
	return models.UpdateLotteryContent(lc.ID, lc.getMapsEdit())
}

func (lc *LotteryContentInput) Delete() error {
	return models.DeleteLotteryContent(lc.ID)
}

func GetLotteryForManager() ([]models.Lottery, int64, error) {
	return models.GetLotteries("")
}

func (l *LotteryContentInput) GetLotteryContentForManager() ([]models.LotteryContent, int64, error) {
	cond, vals, err := util.SqlWhereBuild(l.getMapsGet(), "and")
	if err != nil {
		return nil, 0, err
	}
	return models.GetLotteryContents(l.PageNum, l.PageSize, cond, vals)
}

func (lc *LotteryContentInput) getMapsEdit() map[string]interface{} {
	maps := make(map[string]interface{})
	if lc.Type != "" {
		maps["type"] = lc.Type
	}
	if lc.ID > 0 {
		maps["id"] = lc.ID
	}
	if lc.Language != "" {
		maps["language"] = lc.Language
	}
	if lc.Content != "" {
		maps["content"] = lc.Content
	}
	return maps
}

func (lc *LotteryContentInput) getMapsGet() map[string]interface{} {
	maps := make(map[string]interface{})
	if lc.Type != "" {
		maps["type like"] = "%" + lc.Type + "%"
	}
	if lc.ID > 0 {
		maps["id"] = lc.ID
	}
	if lc.Content != "" {
		maps["content like"] = "%" + lc.Content + "%"
	}
	if lc.Language != "" {
		maps["language"] = lc.Language
	}
	return maps
}

func (l *LotteryInput) getMapsEdit() map[string]interface{} {
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
