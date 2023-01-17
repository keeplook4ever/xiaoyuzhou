package models

// PlayerLotteryLog
// 记录用户抽取的运势签
type PlayerLotteryLog struct {
	ID        int    `gorm:"column:id;type:tinyint(20);primaryKey;autoIncrement;not null" json:"id"`
	Uid       string `gorm:"column:uid;type:varchar(191);unique" json:"uid"`     //用户ID
	Timestamp int    `gorm:"column:timestamp;type:tinyint(20)" json:"timestamp"` //事件发生时间
	Score     int    `gorm:"column:score;type:tinyint(3)" json:"score"`          //运势分数
	Keyword   string `gorm:"column:keyword;type:varchar(20)" json:"keyword"`     //运势关键字
	Content   string `gorm:"column:content;type:varchar(191)" json:"content"`    //运势内容
}

func CreatPlayerLotteryLog(uid string, ts int, score int, kw string, cont string) error {
	newPlayreLotteryLog := PlayerLotteryLog{
		Uid:       uid,
		Timestamp: ts,
		Score:     score,
		Keyword:   kw,
		Content:   cont,
	}
	if err := Db.Create(&newPlayreLotteryLog).Error; err != nil {
		return err
	}
	return nil

}
