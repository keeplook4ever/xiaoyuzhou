package tarot_service

import (
	"time"
	"xiaoyuzhou/models"
	"xiaoyuzhou/pkg/logging"
)

// GetOneTarotPic 获取一张塔罗牌
func GetRandomOneTarot(uid string) (*models.TarotDto, error){
	randTarot, err := models.GetOneRandTarot()
	if err != nil {
		return nil, err
	}
	// 成功之后需要记录日志
	ts := int(time.Now().Unix())
	// Log TODO
	logging.Debugf("User %s Get one random tarot at %v", uid, ts)
	return randTarot, nil
}
