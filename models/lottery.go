package models

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"xiaoyuzhou/pkg/logging"
	"xiaoyuzhou/pkg/util"
)

type Lottery struct {
	Model
	MinScore    int     `gorm:"column:min_score;not null;type:tinyint(3)" json:"min_score"` // 最小分数
	MaxScore    int     `gorm:"column:max_score;not null;type:tinyint(3)" json:"max_score"` // 最大分数
	KeyWord     string  `gorm:"column:keyword;not null;type:varchar(50)" json:"keyword"`    // 运势文字
	Probability float32 `gorm:"column:probability;type:float" json:"probability"`           // 概率
	Type        string  `gorm:"column:type;type:varchar(10)" json:"type"`                   //枚举
}

type LotteryContent struct {
	Model
	Type    string `gorm:"column:type;not null;type:varchar(1)" json:"type"` //A-D 枚举
	Content string `gorm:"column:content;not null;type:text" json:"content"`
}

type LotteryDto struct {
	Score   int    `json:"score"`
	Keyword string `json:"keyword"`
	Content string `json:"content"`
}

type TypeAndProb struct {
	Type string  `json:"type"`
	Prob float32 `json:"prob"`
}

func (l *Lottery) makeLotteryWithContent() LotteryDto {
	score := util.GetScore(l.MinScore, l.MaxScore)
	content, _ := getOneRandLotteryContent(l.Type)
	return LotteryDto{
		Score:   score,
		Keyword: l.KeyWord,
		Content: content,
	}
}

func GetLotteries() ([]Lottery, error) {
	var lotteries []Lottery
	err := Db.Model(&Lottery{}).Find(&lotteries).Error
	if err != nil {
		return nil, err
	}
	return lotteries, nil
}

func GetLotteryContents(cond string, vals []interface{}) ([]LotteryContent, error) {
	var lotteryContents []LotteryContent
	err := Db.Model(&LotteryContent{}).Where(cond, vals...).Find(&lotteryContents).Error
	if err != nil {
		return nil, err
	}
	return lotteryContents, nil
}

func GetOneRandLottery() (*LotteryDto, error) {
	lotteries, err := GetLotteries()
	if err != nil {
		return nil, err
	}
	if len(lotteries) == 0 {
		return nil, errors.New("没有设置lottery")
	}
	lotteryChose := getOneLotteryWithProb(lotteries)
	lotteryWithContent := lotteryChose.makeLotteryWithContent()
	return &lotteryWithContent, nil
}

func getOneRandLotteryContent(tyPe string) (string, error) {
	var contents []string
	err := Db.Model(&LotteryContent{}).Pluck("content", &contents).Where("type = ?", tyPe).Error
	if err != nil {
		return "", err
	}
	// 随机取一个值
	rand.Seed(time.Now().Unix())
	return contents[rand.Intn(len(contents))], nil
}

func EditLottery(typE string, data map[string]interface{}) error {
	if err := Db.Model(&Lottery{}).Where("type = ?", typE).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

func AddLotteryContent(content, typE string) error {
	lc := LotteryContent{
		Content: content,
		Type:    typE,
	}
	err := Db.Model(&LotteryContent{}).Create(&lc).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateLotteryContent(id int, data interface{}) error {
	err := Db.Model(&LotteryContent{}).Where("id = ?", id).Updates(data).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteLotteryContent(id int) error {
	// 删除时需要确认是否该类型 >= 2
	var count int64
	Db.Model(&LotteryContent{}).Where("id = ?", id).Count(&count)
	if count >= 2 {
		err := Db.Where("id = ?", id).Delete(&LotteryContent{}).Error
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("该类型不够啦！")
}

func getOneLotteryWithProb(lotteries []Lottery) Lottery {
	// 随机算法
	var typeWithProb map[string]float32
	origIndexArray := make([]string, 0)

	var indexStringSlice [][]string // 类型字符按照概率生成的string切片
	// 生成100个1234下标，其中个数多少取决于其的概率值

	for i, l := range lotteries {
		typeWithProb[l.Type] = l.Probability
		indexString := strconv.Itoa(i)
		indexStringSlice[i] = strings.Split(strings.Repeat(indexString, int(l.Probability*100)), "")
		origIndexArray = append(origIndexArray, indexStringSlice[i]...)
	}

	logging.Debugf("count: %d || value: %v", len(origIndexArray), origIndexArray)
	// 将顺序排列的slice打乱，在rand.Intn(len)取下标，取随机
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(origIndexArray), func(i, j int) {
		origIndexArray[i], origIndexArray[j] = origIndexArray[j], origIndexArray[i]
	})
	logging.Debugf("count: %d || Shuffled: %v", len(origIndexArray), origIndexArray)

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(origIndexArray))
	indexChose := origIndexArray[n]

	intIndex, _ := strconv.Atoi(indexChose)
	return lotteries[intIndex]
}
