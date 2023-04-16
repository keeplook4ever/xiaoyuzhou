package tarot_service

import (
	"time"
	"xiaoyuzhou/models"
	"xiaoyuzhou/pkg/logging"
	"xiaoyuzhou/service/order_service"
)

// GetRandomOneTarot 获取随机一张塔罗牌, 并且创建塔罗抽取记录, 获得塔罗牌，订单号
func GetRandomOneTarot(uid, question, lang string) (*models.TarotDto, string, error) {

	randTarot, err := models.GetOneRandTarot(lang)
	if err != nil {
		logging.Debugf("Error: %s", err.Error())
		return nil, "", err
	}

	// 将问题与对应塔罗牌插入数据库
	idList := make([]uint, 0)
	idList = append(idList, randTarot.TarotId)
	ts := time.Now().Unix()
	// 创建记录
	err, orderId := order_service.CreateRecordWithNoOrder(uid, question, ts, idList)
	if err != nil {
		logging.Debugf("Error: %s", err.Error())
		return nil, "", err
	}

	return randTarot, orderId, nil
}

// GetOneTarotByOrderAndUser 根据订单号获取抽取的塔罗牌的解答
func GetOneTarotByOrderAndUser(OrderId string) (*models.TarotDto, string, int64, error) {
	// 根据订单号找到对应塔罗牌列表
	tarot, question, ts, err := models.GetOneTarotFromOrder(OrderId)

	if err != nil {
		return nil, "", 0, err
	}
	return tarot, question, ts, nil
}
