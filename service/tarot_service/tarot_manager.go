package tarot_service

import (
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"io"
	"xiaoyuzhou/models"
	"xiaoyuzhou/pkg/logging"
	"xiaoyuzhou/pkg/util"
)

type TarotInput struct {
	ID            int      // 塔罗牌ID
	ImgUrl        string   // 图片链接
	Language      string   // 语言
	Pos           string   // 塔罗正逆位
	CardName      string   // 卡牌名字
	KeyWord       string   // 卡牌解读关键词
	Constellation string   // 对应星座
	People        string   // 对应人物
	Element       string   // 对应元素
	Enhance       string   // 加强牌
	AnalyzeOne    string   // 解析1
	AnalyzeTwo    string   // 解析2
	PosMeaning    string   // 正逆位含义
	Love          string   // 爱情婚姻
	Work          string   // 事业学业
	Money         string   // 人际财富
	Health        string   // 健康生活
	Other         string   // 其他
	LuckyNumber   string   // 幸运数字
	Saying        string   // 名言
	AnswerList    []string // 答案列表
	PageNum       int      // 分页偏移数
	PageSize      int      // 每页数量
	CreatedBy     string   // 创建人
	UpdatedBy     string   // 修改人
	Status        string   // 状态
}

func (t *TarotInput) Add() error {
	answers := util.StringSlice2String(t.AnswerList)
	answersValue := ""
	if answers != nil {
		answersValue = *answers
	}
	dbData := map[string]interface{}{
		"img_url":       t.ImgUrl,
		"language":      t.Language,
		"pos":           t.Pos,
		"card_name":     t.CardName,
		"keyword":       t.KeyWord,
		"constellation": t.Constellation,
		"people":        t.People,
		"element":       t.Element,
		"enhance":       t.Enhance,
		"analyze_one":   t.AnalyzeOne,
		"analyze_two":   t.AnalyzeTwo,
		"pos_meaning":   t.PosMeaning,
		"love":          t.Love,
		"work":          t.Work,
		"money":         t.Money,
		"health":        t.Health,
		"other":         t.Other,
		"answer_list":   answersValue,
		"lucky_number":  t.LuckyNumber,
		"saying":        t.Saying,
		"created_by":    t.CreatedBy,
		"updated_by":    t.UpdatedBy,
		"status":        t.Status,
	}
	return models.AddTarot(dbData)
}

func Import(r io.Reader) error {
	xlsx, err := excelize.OpenReader(r)
	if err != nil {
		return err
	}
	tarotData, _ := xlsx.GetRows("tarot")
	tarotList := make([]map[string]interface{}, 0)
	for i, row := range tarotData {
		fmt.Println(row)
		if i == 0 {
			continue
		}
		// 如果数据不够，则抛出异常
		if len(row) < 23 {
			return errors.New("格式不对")
		}

		answers := make([]string, 0)
		answers = append(answers, row[18], row[19], row[20], row[21], row[22])
		answersList := util.StringSlice2String(answers)
		tarotOne := map[string]interface{}{
			"img_url":       "",
			"language":      row[0],
			"pos":           row[1],
			"card_name":     row[2],
			"keyword":       row[3],
			"constellation": row[4],
			"people":        row[5],
			"element":       row[6],
			"enhance":       row[7],
			"saying":        row[8],
			"lucky_number":  row[9],
			"analyze_one":   row[10],
			"analyze_two":   row[11],
			"pos_meaning":   row[12],
			"love":          row[13],
			"work":          row[14],
			"money":         row[15],
			"health":        row[16],
			"other":         row[17],
			"answer_list":   *answersList,
			"created_by":    "admin",
			"updated_by":    "admin",
			"status":        "off", // 批量导入的由于没有图片，所以状态：关闭
		}
		tarotList = append(tarotList, tarotOne)
	}
	if err := models.Db.Model(&models.Tarot{}).Create(tarotList).Error; err != nil {
		logging.Error(fmt.Sprintf("批量导入塔罗失败: %s", err.Error()))
		return err
	}

	return err
}

func (t *TarotInput) ExistByID() (bool, error) {
	return models.ExistTarotByID(t.ID)
}
func (t *TarotInput) Edit() error {

	data := make(map[string]interface{})
	if t.ImgUrl != "" {
		data["img_url"] = t.ImgUrl
	}
	if t.Pos != "" {
		data["pos"] = t.Pos
	}
	if t.Language != "" {
		data["language"] = t.Language
	}
	if t.CardName != "" {
		data["card_name"] = t.CardName
	}
	if t.KeyWord != "" {
		data["keyword"] = t.KeyWord
	}
	if t.Constellation != "" {
		data["constellation"] = t.Constellation
	}
	if t.People != "" {
		data["people"] = t.People
	}
	if t.Element != "" {
		data["element"] = t.Element
	}
	if t.Enhance != "" {
		data["enhance"] = t.Enhance
	}
	if t.AnalyzeOne != "" {
		data["analyze_one"] = t.AnalyzeOne
	}
	if t.AnalyzeTwo != "" {
		data["analyze_two"] = t.AnalyzeTwo
	}
	if t.PosMeaning != "" {
		data["pos_meaning"] = t.PosMeaning
	}
	if t.Love != "" {
		data["love"] = t.Love
	}
	if t.Work != "" {
		data["work"] = t.Work
	}
	if t.Money != "" {
		data["money"] = t.Money
	}
	if t.Health != "" {
		data["health"] = t.Health
	}
	if t.Other != "" {
		data["other"] = t.Other
	}
	if t.AnswerList != nil {
		answers := util.StringSlice2String(t.AnswerList)
		if answers != nil {
			data["answer_list"] = *answers
		}
	}
	if t.Saying != "" {
		data["saying"] = t.Saying
	}
	if t.LuckyNumber != "" {
		data["lucky_number"] = t.LuckyNumber
	}
	if t.UpdatedBy != "" {
		data["updated_by"] = t.UpdatedBy
	}
	if t.Status != "" {
		data["status"] = t.Status
	}
	return models.EditTarot(t.ID, data)
}

func (t *TarotInput) Get() ([]models.TarotDto, int64, error) {
	cond, vals, err := util.SqlWhereBuild(t.getMaps(), "and")
	if err != nil {
		logging.Error(fmt.Sprintf("Error Get(): %s", err.Error()))
		return nil, 0, err
	}
	tarots, count, err := models.GetTarots(t.PageNum, t.PageSize, cond, vals)
	if err != nil {
		logging.Error(fmt.Sprintf("Error Get(): %s", err.Error()))
		return nil, 0, err
	}
	return tarots, count, nil
}

func (t *TarotInput) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	if t.ID > 0 {
		maps["id"] = t.ID
	}
	if t.CardName != "" {
		maps["card_name like"] = "%" + t.CardName + "%"
	}
	if t.Pos != "" {
		maps["pos"] = t.Pos
	}
	if t.Constellation != "" {
		maps["constellation like"] = "%" + t.Constellation + "%"
	}
	if t.KeyWord != "" {
		maps["keyword like"] = "%" + t.KeyWord + "%"
	}
	if t.Enhance != "" {
		maps["enhance like"] = "%" + t.Enhance + "%"
	}
	if t.Element != "" {
		maps["element"] = "%" + t.Element + "%"
	}
	if t.Language != "" {
		maps["language"] = t.Language
	}
	if t.Status != "" {
		maps["status"] = t.Status
	}
	return maps
}
