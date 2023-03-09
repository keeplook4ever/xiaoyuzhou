package tarot_service

import (
	"xiaoyuzhou/models"
)

// GetRandomOneTarot 获取随机一张塔罗牌
func GetRandomOneTarot() (*models.TarotDto, error) {

	randTarot, err := models.GetOneRandTarot()
	if err != nil {
		return nil, err
	}

	// 将问题与对应塔罗牌插入数据库
	//err = InsertQuestionToDB(question, uid, recordId, randTarot.TarotId)
	//if err != nil {
	//	return nil, err
	//}

	return randTarot, nil
}

// GetOneTarotByOrderAndUser 根据订单号获取抽取的塔罗牌的解答
func GetOneTarotByOrderAndUser(OriOrderId string) (*models.TarotDto, string, error) {
	// 根据订单号找到对应塔罗牌列表
	tarot, question, err := models.GetOneTarotFromOrder(OriOrderId)

	if err != nil {
		return nil, "", err
	}
	return tarot, question, nil
}
