package models

import (
	"gorm.io/gorm"
	"strings"
	"xiaoyuzhou/pkg/util"
)

type Article struct {
	Model               // gorm.Model 包含了ID，CreatedAt， UpdatedAt， DeletedAt
	CategoryID int      `gorm:"column:category_id;type:int" json:"category_id"`  // 默认外键
	Category   Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"` // 一个文章属于一个类型

	SeoTitle        string `gorm:"column:seo_title;not null;unique;type:varchar(50)" json:"seo_title"`
	SeoUrl          string `gorm:"column:seo_url;not null;unique;type:varchar(50)" json:"seo_url"`
	PageTitle       string `gorm:"column:page_title;not null;unique;type:varchar(50)" json:"page_title"`
	MetaDesc        string `gorm:"column:meta_desc;not null;type:varchar(50)" json:"meta_desc"`
	RelatedArticles string `gorm:"column:related_articles;type:varchar(20)" json:"related_articles"`
	Content         string `gorm:"column:content;not null;type:text" json:"content"`
	AuthorId        int    `gorm:"column:author_id;not null;type:int" json:"author_id"`
	Author          Author `gorm:"foreignKey:AuthorId" json:"author,omitempty"` // 一个文章属于一个作者
	CoverImageUrl   string `gorm:"column:cover_image_url;not null;type:varchar(50)" json:"cover_image_url"`
	State           int    `gorm:"column:state;not null;type:tinyint(1)" json:"state"`
	Language        string `gorm:"column:language;not null;type:varchar(2)" json:"language"`
	CreatedBy       string `gorm:"column:created_by;not null;type:varchar(50)" json:"created_by"`
	UpdatedBy       string `gorm:"column:updated_by;not null;type:varchar(50)" json:"updated_by"`
}

type ArticleDto struct {
	ID              uint   `json:"id"`
	CategoryID      uint   `json:"category_id"`
	CategoryName    string `json:"category_name"`
	SeoTitle        string `json:"seo_title"`
	SeoUrl          string `json:"seo_url"`
	PageTitle       string `json:"page_title"`
	MetaDesc        string `json:"meta_desc"`
	RelatedArticles []int  `json:"related_articles"`
	Content         string `json:"content"`
	AuthorID        uint   `json:"author_id"`
	AuthorName      string `json:"author_name"`
	CoverImageUrl   string `json:"cover_image_url"`
	State           int    `json:"state"`
	Language        string `json:"language"`
	CreatedAt       int    `json:"created_at"`
	CreatedBy       string `json:"created_by"`
	UpdatedAt       int    `json:"updated_at"`
	UpdatedBy       string `json:"updated_by"`
}

// ToArticleDto 从数据库结构抽取前端需要的字段返回
func (itself *Article) ToArticleDto() ArticleDto {
	return ArticleDto{
		ID:              itself.ID,
		CategoryID:      itself.Category.ID,
		CategoryName:    itself.Category.Name,
		SeoUrl:          itself.SeoUrl,
		SeoTitle:        itself.SeoTitle,
		PageTitle:       itself.PageTitle,
		MetaDesc:        itself.MetaDesc,
		RelatedArticles: util.String2Int(strings.Split(itself.RelatedArticles, ",")),
		Content:         itself.Content,
		AuthorID:        itself.Author.ID,
		AuthorName:      itself.Author.Name,
		CoverImageUrl:   itself.CoverImageUrl,
		State:           itself.State,
		Language:        itself.Language,
		CreatedAt:       itself.CreatedAt,
		UpdatedAt:       itself.UpdatedAt,
		CreatedBy:       itself.CreatedBy,
		UpdatedBy:       itself.UpdatedBy,
	}
}

// ExistArticleByID checks if an article exists based on ID
func ExistArticleByID(id int) (bool, error) {
	var article Article
	err := Db.Model(&Article{}).Select("id").Where("id = ? ", id).First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if article.ID > 0 {
		return true, nil
	}

	return false, nil
}

// GetArticleTotal gets the total number of articles based on the constraints
func GetArticleTotal(cond string, vals []interface{}) (int64, error) {
	var count int64
	if err := Db.Model(&Article{}).Where(cond, vals...).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// GetArticles gets a list of articles based on paging constraints
func GetArticles(pageNum int, pageSize int, cond string, vals []interface{}) ([]ArticleDto, error) {
	var articles []Article

	err := Db.Preload("Category").Preload("Author").Where(cond, vals...).Offset(pageNum).Limit(pageSize).Find(&articles).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	resp := make([]ArticleDto, len(articles))

	for i, aa := range articles {
		resp[i] = aa.ToArticleDto()
	}
	return resp, nil
}

// GetArticle Get a single article based on ID
func GetArticle(id int) (*Article, error) {
	var article Article
	err := Db.Where("id = ? ", id).First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &article, nil
}

// EditArticle modify a single article
func EditArticle(id int, data interface{}) error {
	if err := Db.Model(&Article{}).Where("id = ? ", id).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

// AddArticle add a single article
func AddArticle(data map[string]interface{}) error {
	article := Article{
		CategoryID:      data["category_id"].(int),
		SeoTitle:        data["seo_title"].(string),
		SeoUrl:          data["seo_url"].(string),
		PageTitle:       data["page_title"].(string),
		MetaDesc:        data["meta_desc"].(string),
		RelatedArticles: data["related_articles"].(string),
		Content:         data["content"].(string),
		AuthorId:        data["author_id"].(int),
		CoverImageUrl:   data["cover_image_url"].(string),
		State:           data["state"].(int),
		Language:        data["language"].(string),
		CreatedBy:       data["created_by"].(string),
		UpdatedBy:       data["updated_by"].(string),
	}
	if err := Db.Create(&article).Error; err != nil {
		return err
	}

	return nil
}

// DeleteArticle delete a single article
func DeleteArticle(id int) error {
	if err := Db.Where("id = ?", id).Delete(Article{}).Error; err != nil {
		return err
	}

	return nil
}

// CleanAllArticle clear all article
func CleanAllArticle() error {
	if err := Db.Unscoped().Delete(&Article{}).Error; err != nil {
		return err
	}

	return nil
}
