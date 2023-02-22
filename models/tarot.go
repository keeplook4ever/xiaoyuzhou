package models

import "gorm.io/gorm"

type Tarot struct {
	Model
	ImgUrl        string `gorm:"column:img_url;not null;type:varchar(191)" json:"img_url"`                          // 图片链接
	Language      string `gorm:"column:language;not null;type:varchar(11)" json:"language" enums:"JP,EN,CH-S,CH-T"` // 语言
	Pos           string `gorm:"column:pos;not null;type:varchar(10)" json:"pos" enums:"正位,逆位"`                     // 塔罗正逆位
	CardName      string `gorm:"column:card_name;not null;type:varchar(100)" json:"card_name"`                      //卡牌名字
	KeyWord       string `gorm:"column:keyword;not null;type:varchar(190)" json:"keyword"`                          // 卡牌解读关键词
	Constellation string `gorm:"column:constellation;not null;type:varchar(190)" json:"constellation"`              //对应星座
	People        string `gorm:"column:people;not null;type:varchar(190)" json:"people"`                            //对应人物
	Element       string `gorm:"column:element;not null;type:varchar(190)" json:"element"`                          //对应元素
	Enhance       string `gorm:"column:enhance;not null;type:varchar(190)" json:"enhance"`                          // 加强牌
	AnalyzeOne    string `gorm:"column:analyze_one;not null;type:text" json:"analyze_one"`                          // 解析1
	AnalyzeTwo    string `gorm:"column:analyze_two;not null;type:text" json:"analyze_two"`                          // 解析2
	PosMeaning    string `gorm:"column:pos_meaning;not null;type:varchar(190)" json:"pos_meaning"`                  // 正逆位含义
	Love          string `gorm:"column:love;not null;type:varchar(190)" json:"love"`                                // 爱情婚姻
	Work          string `gorm:"column:work;not null;type:varchar(190)" json:"work"`                                // 事业学业
	Money         string `gorm:"column:money;not null;type:varchar(190)" json:"money"`                              // 人际财富
	Health        string `gorm:"column:health;not null;type:varchar(190)" json:"health"`                            // 健康生活
	Other         string `gorm:"column:other;not null;type:varchar(190)" json:"other"`                              // 其他
	AnswerOne     string `gorm:"column:answer_one;not null;type:text" json:"answer_one"`                            //回答1
	AnswerTwo     string `gorm:"column:answer_two;not null;type:text" json:"answer_two"`                            // 回答2
	AnswerThree   string `gorm:"column:answer_three;not null;type:text" json:"answer_three"`                        // 回答3
	AnswerFour    string `gorm:"column:answer_four;not null;type:text" json:"answer_four"`                          // 回答4
	AnswerFive    string `gorm:"column:answer_five;not null;type:text" json:"answer_five"`                          // 回答5
}

//type TarotAnswer struct {
//	Model
//	TarotId int    `gorm:"column:tarot_id;not null;type:tinyint(10)" json:"tarot_id"` // 塔罗牌id
//	Answer  string `gorm:"column:answer;not null;type:varchar(190)" json:"answer"`    //回答
//}

func AddTarot(data map[string]string) error {
	tarot := Tarot{
		ImgUrl:        data["img_url"],
		Language:      data["language"],
		Pos:           data["pos"],
		CardName:      data["card_name"],
		KeyWord:       data["keyword"],
		Constellation: data["constellation"], //对应星座
		People:        data["people"],        // 对应人物
		Element:       data["element"],       // 对应元素
		Enhance:       data["enhance"],       // 加强牌
		AnalyzeOne:    data["analyze_one"],   // 解析1
		AnalyzeTwo:    data["analyze_two"],   // 解析2
		PosMeaning:    data["pos_meaning"],   // 正逆位含义
		Love:          data["love"],          // 爱情婚姻
		Work:          data["work"],          // 事业学业
		Money:         data["money"],         // 人际财富
		Health:        data["health"],        // 健康生活
		Other:         data["other"],         // 其他
		AnswerOne:     data["answer_one"],    // 回答1
		AnswerTwo:     data["answer_two"],    // 回答2
		AnswerThree:   data["answer_three"],  // 回答3
		AnswerFour:    data["answer_four"],   // 回答4
		AnswerFive:    data["answer_five"],   // 回答5
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

func GetTarots(pageNum int, pageSize int, cond string, vals []interface{}) ([]Tarot, int64, error) {
	var tarots []Tarot
	var count int64
	Db.Model(&Tarot{}).Where(cond, vals...).Count(&count)
	err := Db.Where(cond, vals...).Order("created_at desc").Offset(pageNum).Limit(pageSize).Find(&tarots).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	return tarots, count, nil
}
