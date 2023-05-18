package models

import (
	"errors"
	"gorm.io/gorm"
	"math/rand"
	"time"
	"xiaoyuzhou/pkg/util"
)

type Tarot struct {
	Model
	ImgUrl        string `gorm:"column:img_url;not null;type:varchar(191)" json:"img_url"`                   // 图片链接
	Language      string `gorm:"column:language;not null;type:varchar(11)" json:"language" enums:"jp,zh,en"` // 语言
	Pos           string `gorm:"column:pos;not null;type:varchar(10)" json:"pos" enums:"up,down"`            // 塔罗正逆位
	CardName      string `gorm:"column:card_name;not null;type:varchar(100)" json:"card_name"`               // 卡牌名字
	KeyWord       string `gorm:"column:keyword;not null;type:varchar(190)" json:"keyword"`                   // 卡牌解读关键词
	Constellation string `gorm:"column:constellation;not null;type:varchar(190)" json:"constellation"`       // 对应星座
	People        string `gorm:"column:people;not null;type:varchar(190)" json:"people"`                     // 对应人物
	Element       string `gorm:"column:element;not null;type:varchar(190)" json:"element"`                   // 对应元素
	Enhance       string `gorm:"column:enhance;not null;type:varchar(190)" json:"enhance"`                   // 加强牌
	AnalyzeOne    string `gorm:"column:analyze_one;not null;type:text" json:"analyze_one"`                   // 解析1
	AnalyzeTwo    string `gorm:"column:analyze_two;not null;type:text" json:"analyze_two"`                   // 解析2
	PosMeaning    string `gorm:"column:pos_meaning;not null;type:varchar(190)" json:"pos_meaning"`           // 正逆位含义
	Love          string `gorm:"column:love;not null;type:varchar(190)" json:"love"`                         // 爱情婚姻
	Work          string `gorm:"column:work;not null;type:varchar(190)" json:"work"`                         // 事业学业
	Money         string `gorm:"column:money;not null;type:varchar(190)" json:"money"`                       // 人际财富
	Health        string `gorm:"column:health;not null;type:varchar(190)" json:"health"`                     // 健康生活
	Other         string `gorm:"column:other;not null;type:varchar(190)" json:"other"`                       // 其他
	AnswerList    string `gorm:"column:answer_list;not null;type:text" json:"answer_list"`                   // 回答列表
	LuckyNumber   string `gorm:"column:lucky_number;not null;type:varchar(50)" json:"lucky_number"`          // 幸运数字
	Saying        string `gorm:"column:saying;not null;type:varchar(191)" json:"saying"`                     // 名言
	CreatedBy     string `gorm:"column:created_by;not null;type:varchar(50)" json:"created_by"`              // 创建者
	UpdatedBy     string `gorm:"column:updated_by;not null;type:varchar(50)" json:"updated_by"`              // 更新者
	Status        string `gorm:"column:status;not null;type:varchar(20);default:on" json:"status"`           // 状态:on默认启用,off关闭，需手动开启才可被抽取
}

type TarotDto struct {
	TarotId       uint     `json:"tarot_id"`      // 塔罗牌id
	ImgUrl        string   `json:"img_url"`       // 图片链接
	Language      string   `json:"language"`      // 语言
	Pos           string   `json:"pos"`           // 塔罗正逆位
	CardName      string   `json:"card_name"`     // 卡牌名字
	KeyWord       string   `json:"keyword"`       // 卡牌解读关键词
	Constellation string   `json:"constellation"` // 对应星座
	People        string   `json:"people"`        // 对应人物
	Element       string   `json:"element"`       // 对应元素
	Enhance       string   `json:"enhance"`       // 加强牌
	AnalyzeOne    string   `json:"analyze_one"`   // 解析1
	AnalyzeTwo    string   `json:"analyze_two"`   // 解析2
	PosMeaning    string   `json:"pos_meaning"`   // 正逆位含义
	Love          string   `json:"love"`          // 爱情婚姻
	Work          string   `json:"work"`          // 事业学业
	Money         string   `json:"money"`         // 人际财富
	Health        string   `json:"health"`        // 健康生活
	Other         string   `json:"other"`         // 其他
	AnswerList    []string `json:"answer_list"`   // 回答列表
	LuckyNumber   string   `json:"lucky_number"`  // 幸运数字
	Saying        string   `json:"saying"`        // 名言
	CreatedBy     string   `json:"created_by"`    // 创建者
	UpdatedBy     string   `json:"updated_by"`    // 更新者
	CreatedAt     int      `json:"created_at"`    // 创建时间
	UpdatedAt     int      `json:"updated_at"`    // 更新时间
	Status        string   `json:"status"`        // 状态
}

func (t *Tarot) ToTarotDto() TarotDto {
	return TarotDto{
		TarotId:       t.ID,
		ImgUrl:        t.ImgUrl,
		Language:      t.Language,
		Pos:           t.Pos,
		CardName:      t.CardName,
		KeyWord:       t.KeyWord,
		Constellation: t.Constellation,
		People:        t.People,
		Element:       t.Element,
		Enhance:       t.Enhance,
		AnalyzeOne:    t.AnalyzeOne,
		AnalyzeTwo:    t.AnalyzeTwo,
		PosMeaning:    t.PosMeaning,
		Love:          t.Love,
		Work:          t.Work,
		Money:         t.Money,
		Health:        t.Health,
		Other:         t.Other,
		AnswerList:    util.String2StringSlice(t.AnswerList),
		LuckyNumber:   t.LuckyNumber,
		Saying:        t.Saying,
		CreatedBy:     t.CreatedBy,
		UpdatedBy:     t.UpdatedBy,
		CreatedAt:     t.CreatedAt,
		UpdatedAt:     t.UpdatedAt,
		Status:        t.Status,
	}
}

func AddTarot(data map[string]interface{}) error {
	tarot := Tarot{
		ImgUrl:        data["img_url"].(string),
		Language:      data["language"].(string),
		Pos:           data["pos"].(string),
		CardName:      data["card_name"].(string),
		KeyWord:       data["keyword"].(string),
		Constellation: data["constellation"].(string), //对应星座
		People:        data["people"].(string),        // 对应人物
		Element:       data["element"].(string),       // 对应元素
		Enhance:       data["enhance"].(string),       // 加强牌
		AnalyzeOne:    data["analyze_one"].(string),   // 解析1
		AnalyzeTwo:    data["analyze_two"].(string),   // 解析2
		PosMeaning:    data["pos_meaning"].(string),   // 正逆位含义
		Love:          data["love"].(string),          // 爱情婚姻
		Work:          data["work"].(string),          // 事业学业
		Money:         data["money"].(string),         // 人际财富
		Health:        data["health"].(string),        // 健康生活
		Other:         data["other"].(string),         // 其他
		LuckyNumber:   data["lucky_number"].(string),  // 幸运数字
		Saying:        data["saying"].(string),        // 名言
		AnswerList:    data["answer_list"].(string),   // 回答列表
		CreatedBy:     data["created_by"].(string),    // 创建者
		UpdatedBy:     data["updated_by"].(string),    // 更新者(默认同创建者)
		Status:        data["status"].(string),        // 状态
	}
	if err := Db.Create(&tarot).Error; err != nil {
		return err
	}

	return nil
}

// EditTarot modify a single tarot
func EditTarot(id int, data interface{}) error {
	if err := Db.Model(&Tarot{}).Where("id = ? ", id).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

func ExistTarotByID(id int) (bool, error) {
	var t Tarot
	err := Db.Model(&Tarot{}).Select("id").Where("id = ? ", id).First(&t).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	} else if err == gorm.ErrRecordNotFound {
		return false, nil
	}

	if t.ID > 0 {
		return true, nil
	}

	return false, nil
}

func GetTarots(pageNum int, pageSize int, cond string, vals []interface{}) ([]TarotDto, int64, error) {
	var tarots []Tarot
	var count int64
	Db.Model(&Tarot{}).Where(cond, vals...).Count(&count)
	err := Db.Where(cond, vals...).Order("created_at desc").Offset(pageNum).Limit(pageSize).Find(&tarots).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	var resp []TarotDto
	for _, v := range tarots {
		resp = append(resp, v.ToTarotDto())
	}
	return resp, count, nil
}

// GetOneRandTarot 获取一张随机塔罗
func GetOneRandTarot(lang string) (*TarotDto, error) {
	var tarots []TarotDto
	tarots, num, err := GetTarots(0, 10000, "language = ? and status = \"on\"", []interface{}{lang})
	if err != nil {
		return nil, err
	}
	if num == 0 {
		return nil, errors.New("No tarot")
	}
	rand.Seed(time.Now().Unix())
	resp := tarots[rand.Intn(int(num))]
	return &resp, nil
}

func GetThreeRandTarot() {

}
