package models

type TrueWord struct {
	Model
	Word     string `gorm:"column:word;type:varchar(191);not null" json:"word"`                                         // 真言
	Language string `gorm:"column:language;type:varchar(10);not null" json:"language" default:"jp" enums:"jp,zh,en,ts"` // 语言
}
