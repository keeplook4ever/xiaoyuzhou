package player

type Lottery struct {
	// 日签id
	Id string `json:"id,required"`
	// 日签分数
	Score int `json:"score,required"`
	// 日签关键字
	Keyword string `json:"keyword,required"`
	// 日签内容
	Content string `json:"content,required"`
}
