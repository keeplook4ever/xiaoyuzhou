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

// GetOneTarotByOrderAndUser 根据订单号和用户ID获取抽取的塔罗牌的解答
func GetOneTarotByOrderAndUser(orderId, uid string) (*models.TarotDto, string, error) {
	//ok, err := models.CheckOrderIfValid(orderId, uid)
	//if err != nil {
	//	return nil, err
	//}
	//if !ok {
	//	return nil, errors.New("订单合法性校验失败")
	//}
	// 根据订单号找到对应塔罗牌列表
	tarot, question, err := models.GetOneTarotFromOrder(orderId)

	if err != nil {
		return nil, "", err
	}
	return tarot, question, nil
}
