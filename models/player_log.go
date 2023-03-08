package models

// PlayerLotteryLog
// 记录用户抽取的运势签
type PlayerLotteryLog struct {
	Model
	Uid       string `gorm:"column:uid;type:varchar(191)" json:"uid"`            //用户ID
	Timestamp int    `gorm:"column:timestamp;type:tinyint(20)" json:"timestamp"` //事件发生时间
	Score     int    `gorm:"column:score;type:tinyint(3)" json:"score"`          //运势分数
	Keyword   string `gorm:"column:keyword;type:varchar(20)" json:"keyword"`     //运势关键字
	Content   string `gorm:"column:content;type:varchar(191)" json:"content"`    //运势内容
}

func CreatPlayerLotteryLog(uid string, ts int, score int, kw string, cont string) error {
	newPlayerLotteryLog := PlayerLotteryLog{
		Uid:       uid,
		Timestamp: ts,
		Score:     score,
		Keyword:   kw,
		Content:   cont,
	}
	if err := Db.Create(&newPlayerLotteryLog).Error; err != nil {
		return err
	}
	return nil

}
