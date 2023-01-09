package lottery_service

import (
	"xiaoyuzhou/models"
)

func GetLottery() (models.LotteryDto, error) {
	return models.LotteryDto{}, nil
}

func GetLucky() (models.LuckyTodayDto, error) {
	return models.LuckyTodayDto{}, nil
}
