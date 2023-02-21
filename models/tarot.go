package models

type Tarot struct {
	Model
	ImgUrl        string `gorm:"" json:"img_url"`                  // 图片链接
	Language      string `json:"language" enums:"JP,EN,CH-S,CH-T"` // 语言
	Pos           string `json:"pos" enums:"正位,逆位"`                // 塔罗正逆位
	CardName      string `json:"card_name"`                        //卡牌名字
	KeyWord       string `json:"keyword"`                          // 卡牌解读关键词
	Constellation string `json:"constellation"`                    //对应星座
	People        string `json:"people"`                           //对应人物
	Element       string `json:"element"`                          //对应元素
	Enhance       string `json:"enhance"`                          // 加强牌
	AnalyzeOne    string `json:"analyze_one"`                      // 解析1
	AnalyzeTwo    string `json:"analyze_two"`                      // 解析2
	PosMeaning    string `json:"pos_meaning"`                      // 正逆位含义
	Love          string `json:"love"`                             // 爱情婚姻
	Work          string `json:"work"`                             // 事业学业
	Money         string `json:"money"`                            // 人际财富
	Health        string `json:"health"`                           // 健康生活
	Other         string `json:"other"`                            // 其他
	AnswerOne     string `json:"answer_one"`                       //回答1
	AnswerTwo     string `json:"answer_two"`                       // 回答2
	AnswerThree   string `json:"answer_three"`                     // 回答3
	AnswerFour    string `json:"answer_four"`                      // 回答4
	AnswerFive    string `json:"answer_five"`                      // 回答5
}
