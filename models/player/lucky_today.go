package player

type LuckyToday struct {
	// 今日好运id
	Id int `json:"id,required"`
	// 今日好运咒语
	Spell string `json:"spell,required"`
	// 今日适宜
	Todo string `json:"todo,required"`
	// 今日幸运之歌
	Song string `json:"song,required"`
}
