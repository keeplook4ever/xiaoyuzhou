package tarot_service

import (
	"errors"
	"github.com/google/uuid"
	"time"
	"xiaoyuzhou/models"
	"xiaoyuzhou/pkg/logging"
)

// GetRandomOneTarot 获取一张塔罗牌
func GetRandomOneTarot(uid, question string) (*models.TarotDto, string, error) {

	// 生成唯一标识Id
	uuid_ := uuid.New().String() + uid

	// 将问题插入数据库
	err := InsertQuestionToDB(question, uid, uuid_)
	if err != nil {
		return nil, "", err
	}
	randTarot, err := models.GetOneRandTarot()
	if err != nil {
		return nil, "", err
	}
	// 成功之后需要记录日志
	ts := int(time.Now().Unix())
	// Log TODO
	logging.Debugf("User %s Get one random tarot at %v", uid, ts)
	return randTarot, uuid_, nil
}

// GetOneTarotByOrderAndUser 根据订单号和用户ID获取塔罗牌解答
func GetOneTarotByOrderAndUser(orderId, uid string) (*models.TarotDto, error) {
	ok, err := models.CheckOrderIfValid(orderId, uid)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("订单合法性校验失败")
	}
	// 根据订单号找到对应塔罗牌列表
	tarot, err := models.GetOneTarotFromOrder(orderId)

	if err != nil {
		return nil, err
	}
	return tarot, nil
}

func InsertQuestionToDB(qus, uid, uuid string) error {
	return models.InsertQuestionToDB(qus, uid, uuid)
}

func GetQuestionByUser(uid, uuid string) (string, error) {
	ques, err := models.GetQuestionByUUID(uid, uuid)
	if err != nil {
		return "", err
	}
	return ques, nil
}
