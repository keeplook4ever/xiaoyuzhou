package models

// QuestionAndTarot : 记录问题、用户id、用户抽取塔罗牌、时间戳信息的表
type QuestionAndTarot struct {
	Model
	Uid      string `gorm:"column:uid;not null;type:varchar(190)" json:"uid"`             // 用户id
	Question string `gorm:"column:question;not null;type:varchar(190)" json:"question"`   // 用户问题
	RecordId string `gorm:"column:record_id;not null;type:varchar(191)" json:"record_id"` // 唯一记录
	TarotId  uint   `gorm:"column:tarot_id;not null;type:tinyint(10)" json:"tarot_id"`    // 输入问题获取到的塔罗牌ID
}

func InsertQuestionToDB(uid, ques, rid string, tId uint) error {
	if err := Db.Create(&QuestionAndTarot{
		Uid:      uid,
		Question: ques,
		RecordId: rid,
		TarotId:  tId,
	}).Error; err != nil {
		return err
	}
	return nil
}

func GetQuestionByUUID(rid, uid string) (string, error) {
	// 取15分钟前的时间戳
	var qus QuestionAndTarot
	if err := Db.Model(&QuestionAndTarot{}).Where("uid = ? and record_id =", uid, rid).First(&qus).Error; err != nil {
		return "", err
	}
	return qus.Question, nil
}
