package tarot_service

import (
	"fmt"
	"math/rand"
	"time"
	"xiaoyuzhou/models"
	"xiaoyuzhou/pkg/logging"
	"xiaoyuzhou/pkg/util"
)

type DailyTarotInput struct {
	Id        int      // id
	ImgUrl    string   // 图片链接
	Language  string   // 语言
	CardName  string   // 卡牌名字
	Analyze   string   // 解读
	LoveList  []string // 爱情列表
	WorkList  []string // 工作列表
	CreatedBy string   // 创建者
	UpdatedBy string   // 更新者
	PageNum   int      // 第几页
	PageSize  int      // 每页数量
}

func (d *DailyTarotInput) Add() error {
	loves := util.StringSlice2String(d.LoveList)
	lovesValue := ""
	if loves != nil {
		lovesValue = *loves
	}
	works := util.StringSlice2String(d.WorkList)
	worksValue := ""
	if works != nil {
		worksValue = *works
	}
	dbData := map[string]string{
		"img_url":    d.ImgUrl,
		"language":   d.Language,
		"card_name":  d.CardName,
		"analyze":    d.Analyze,
		"love_list":  lovesValue,
		"work_list":  worksValue,
		"created_by": d.CreatedBy,
		"updated_by": d.UpdatedBy,
	}
	return models.AddDailyTarot(dbData)
}

func (d *DailyTarotInput) ExistByID() (bool, error) {
	return models.ExistDailyTarotByID(d.Id)
}

func (d *DailyTarotInput) Edit() error {
	data := make(map[string]string)
	if d.ImgUrl != "" {
		data["img_url"] = d.ImgUrl
	}
	if d.Language != "" {
		data["language"] = d.Language
	}
	if d.CardName != "" {
		data["card_name"] = d.CardName
	}
	if d.Analyze != "" {
		data["analyze"] = d.Analyze
	}
	if d.WorkList != nil {
		works := util.StringSlice2String(d.WorkList)
		if works != nil {
			data["work_list"] = *works
		}
	}
	if d.LoveList != nil {
		loves := util.StringSlice2String(d.LoveList)
		if loves != nil {
			data["love_list"] = *loves
		}
	}
	if d.UpdatedBy != "" {
		data["updated_by"] = d.UpdatedBy
	}
	return models.EditDailyTarot(d.Id, data)
}

func (d *DailyTarotInput) Get() ([]models.DailyTarotDto, int64, error) {
	cond, vals, err := util.SqlWhereBuild(d.getMaps(), "and")
	if err != nil {
		return nil, 0, err
	}
	tarots, count, err := models.GetDailyTarots(d.PageNum, d.PageSize, cond, vals)
	if err != nil {
		return nil, 0, err
	}
	return tarots, count, nil
}

func (d *DailyTarotInput) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	if d.Id > 0 {
		maps["id"] = d.Id
	}
	if d.CardName != "" {
		maps["card_name like"] = "%" + d.CardName + "%"
	}
	if d.Language != "" {
		maps["language"] = d.Language
	}
	return maps
}

func GetDailyTarotFree(uid, language string) (*models.DailyTarotDto, error) {
	randDailyTarot, err := models.GetOneRandDailyTarot(language)
	if err != nil {
		logging.Error(fmt.Sprintf("Get GetDailyTarotFree Error, [uid:%s] %s", uid, err.Error()))
		return nil, err
	}

	// 随机取一个work和love
	rand.Seed(time.Now().Unix())
	love := make([]string, 0)
	love = append(love, randDailyTarot.LoveList[rand.Intn(len(randDailyTarot.LoveList))])
	work := make([]string, 0)
	work = append(work, randDailyTarot.WorkList[rand.Intn(len(randDailyTarot.WorkList))])
	randDailyTarot.WorkList = work
	randDailyTarot.LoveList = love
	return randDailyTarot, nil
}
