package models

type Question struct {
	Model
	Uid      string `gorm:"column:uid;not null;type:varchar(190)" json:"uid"`           // 用户id
	Question string `gorm:"column:question;not null;type:varchar(190)" json:"question"` // 用户问题
	UUid     string `gorm:"column:uuid;not null;type:varchar(191)" json:"uuid"`         // 唯一记录
}

func InsertQuestionToDB(uid, ques, uuid string) error {
	if err := Db.Create(&Question{
		Uid:      uid,
		Question: ques,
		UUid:     uuid,
	}).Error; err != nil {
		return err
	}
	return nil
}

func GetQuestionByUUID(uuid, uid string) (string, error) {
	// 取15分钟前的时间戳
	var qus Question
	if err := Db.Model(&Question{}).Where("uid = ? and uuid =", uid, uuid).First(&qus).Error; err != nil {
		return "", err
	}
	return qus.Question, nil
}
