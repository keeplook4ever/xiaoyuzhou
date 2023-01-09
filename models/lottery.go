package models

type Lottery struct {
	Model
	Score   int    `gorm:"score,not null" json:"score"` // 日签分数
	Keyword string `gorm:"word,not null" json:"word"`   // 日签关键字
	Content string `gorm:"content" json:"content"`      // 日签内容
}

type LotteryDto struct {
	Score   int    `json:"score"`
	Keyword string `json:"keyword"`
	Content string `json:"content"`
}

func (l *Lottery) ToLotteryDto() LotteryDto {
	return LotteryDto{
		Score:   l.Score,
		Keyword: l.Keyword,
		Content: l.Content,
	}
}
