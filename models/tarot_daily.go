package models

import (
	"errors"
	"gorm.io/gorm"
	"math/rand"
	"time"
	"xiaoyuzhou/pkg/util"
)

type DailyTarot struct {
	Model
	ImgUrl    string `gorm:"column:img_url;not null;type:varchar(191)" json:"img_url"`                   // 图片链接
	Language  string `gorm:"column:language;not null;type:varchar(11)" json:"language" enums:"jp,zh,en"` // 语言
	CardName  string `gorm:"column:card_name;not null;type:varchar(100)" json:"card_name"`               // 卡牌名字
	Analyze   string `gorm:"column:analyze;not null;type:text" json:"analyze"`                           // 解读
	LoveList  string `gorm:"column:love_list;not null;type:text" json:"love_list"`                       // 爱情列表
	WorkList  string `gorm:"column:work_list;not null;type:text" json:"work_list"`                       // 工作列表
	CreatedBy string `gorm:"column:created_by;not null;type:varchar(50)" json:"created_by"`              // 创建者
	UpdatedBy string `gorm:"column:updated_by;not null;type:varchar(50)" json:"updated_by"`              // 更新者
}

type DailyTarotDto struct {
	Id        uint     `json:"id"`         // id
	ImgUrl    string   `json:"img_url"`    // 图片
	CardName  string   `json:"card_name"`  // 名称
	Language  string   `json:"language"`   // 语言
	Analyze   string   `json:"analyze"`    // 解析
	LoveList  []string `json:"love_list"`  // 爱情
	WorkList  []string `json:"work_list"`  // 工作
	CreatedBy string   `json:"created_by"` // 创建者
	UpdatedBy string   `json:"updated_by"` // 更新者
	CreatedAt int      `json:"created_at"` // 创建时间
	UpdatedAt int      `json:"updated_at"` // 更新时间
}

func (d *DailyTarot) ToDailyTarotDto() DailyTarotDto {
	return DailyTarotDto{
		Id:        d.ID,
		ImgUrl:    d.ImgUrl,
		CardName:  d.CardName,
		Language:  d.Language,
		Analyze:   d.Analyze,
		LoveList:  util.String2StringSlice(d.LoveList),
		WorkList:  util.String2StringSlice(d.WorkList),
		CreatedBy: d.CreatedBy,
		UpdatedBy: d.UpdatedBy,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

func GetOneRandDailyTarot(lang string) (*DailyTarotDto, error) {
	var tarots []DailyTarotDto
	tarots, num, err := GetDailyTarots(0, 10000, "language = ?", []interface{}{lang})
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

func GetDailyTarots(pageNum int, pageSize int, cond string, vals []interface{}) ([]DailyTarotDto, int64, error) {
	var tarots []DailyTarot
	var count int64
	Db.Model(&DailyTarot{}).Where(cond, vals...).Count(&count)
	err := Db.Where(cond, vals...).Order("created_at desc").Offset(pageNum).Limit(pageSize).Find(&tarots).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	var resp []DailyTarotDto
	for _, v := range tarots {
		resp = append(resp, v.ToDailyTarotDto())
	}
	return resp, count, nil
}

func AddDailyTarot(data map[string]string) error {
	dailyTarot := DailyTarot{
		ImgUrl:    data["img_url"],
		Language:  data["language"],
		CardName:  data["card_name"],
		Analyze:   data["analyze"],
		LoveList:  data["love_list"],
		WorkList:  data["work_list"],
		CreatedBy: data["created_by"], // 创建者
		UpdatedBy: data["updated_by"], // 更新者(默认同创建者)
	}
	if err := Db.Create(&dailyTarot).Error; err != nil {
		return err
	}
	return nil
}

func ExistDailyTarotByID(id int) (bool, error) {
	var t DailyTarot
	err := Db.Model(&DailyTarot{}).Select("id").Where("id = ? ", id).First(&t).Error
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

// EditDailyTarot modify a single dailyTarot
func EditDailyTarot(id int, data interface{}) error {
	if err := Db.Model(&DailyTarot{}).Where("id = ? ", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}
