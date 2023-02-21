package v1

import "github.com/gin-gonic/gin"

type AddTarotForm struct {
	ImgUrl        string `json:"img_url" binding:"required"`                          // 图片链接
	Language      string `json:"language" enums:"JP,EN,CH-S,CH-T" binding:"required"` // 语言
	Pos           string `json:"pos" enums:"正位,逆位" binding:"required"`            // 塔罗正逆位
	CardName      string `json:"card_name" binding:"required"`                        //卡牌名字
	KeyWord       string `json:"keyword" binding:"required"`                          // 卡牌解读关键词
	Constellation string `json:"constellation" binding:"required"`                    //对应星座
	People        string `json:"people" binding:"required"`                           //对应人物
	Element       string `json:"element" binding:"required"`                          //对应元素
	Enhance       string `json:"enhance" binding:"required"`                          // 加强牌
	AnalyzeOne    string `json:"analyze_one" binding:"required"`                      // 解析1
	AnalyzeTwo    string `json:"analyze_two" binding:"required"`                      // 解析2
	PosMeaning    string `json:"pos_meaning" binding:"required"`                      // 正逆位含义
	Love          string `json:"love" binding:"required"`                             // 爱情婚姻
	Work          string `json:"work" binding:"required"`                             // 事业学业
	Money         string `json:"money" binding:"required"`                            // 人际财富
	Health        string `json:"health" binding:"required"`                           // 健康生活
	Other         string `json:"other" binding:"required"`                            // 其他
	AnswerOne     string `json:"answer_one" binding:"required"`                       //回答1
	AnswerTwo     string `json:"answer_two" binding:"required"`                       // 回答2
	AnswerThree   string `json:"answer_three" binding:"required"`                     // 回答3
	AnswerFour    string `json:"answer_four" binding:"required"`                      // 回答4
	AnswerFive    string `json:"answer_five" binding:"required"`                      // 回答5
}

// AddTarot
// @Summary 添加塔罗牌内容
// @Param _ body AddTarotForm true "参数"
// @Produce json
// @Accept json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /manager/tarot [post]
// @Security ApiKeyAuth
// @Tags Manager
func AddTarot(c *gin.Context) {

}
