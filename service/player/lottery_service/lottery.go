package lottery_service

type Lottery struct {
	Id         int
	Score      int
	Keyword    string
	Content    string
	LuckySpell string
	LuckyTodo  string
	LuckySong  string
}

func GetLottery() (Lottery, error) {
	l := Lottery{
		Id: 1,
	}
	return l, nil
}
