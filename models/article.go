package models

import (
	"gorm.io/gorm"
	"strings"
	"xiaoyuzhou/pkg/util"
)

type Article struct {
	Model               // gorm.Model 包含了ID，CreatedAt， UpdatedAt， DeletedAt
	CategoryID int      `gorm:"column:category_id;type:int" json:"category_id"` // 默认外键
	Category   Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"` // 一个文章属于一个类型

	SeoTitle        string `gorm:"column:seo_title;not null;unique;type:varchar(100)" json:"seo_title"`
	SeoUrl          string `gorm:"column:seo_url;not null;unique;type:varchar(191)" json:"seo_url"`
	PageTitle       string `gorm:"column:page_title;not null;unique;type:varchar(100)" json:"page_title"`
	MetaDesc        string `gorm:"column:meta_desc;not null;type:varchar(100)" json:"meta_desc"`
	RelatedArticles string `gorm:"column:related_articles;type:varchar(20)" json:"related_articles"`
	Content         string `gorm:"column:content;not null;type:text" json:"content"`
	AuthorId        int    `gorm:"column:author_id;not null;type:int" json:"author_id"`
	Author          Author `gorm:"foreignKey:AuthorId" json:"author,omitempty"` // 一个文章属于一个作者
	CoverImageUrl   string `gorm:"column:cover_image_url;not null;type:varchar(191)" json:"cover_image_url"`
	State           int    `gorm:"column:state;not null;type:tinyint(1)" json:"state"`
	Language        string `gorm:"column:language;not null;type:varchar(2)" json:"language"`
	CreatedBy       string `gorm:"column:created_by;not null;type:varchar(50)" json:"created_by"`
	UpdatedBy       string `gorm:"column:updated_by;not null;type:varchar(50)" json:"updated_by"`
}

type ArticleDto struct {
	ID              uint   `json:"id,omitempty"`
	CategoryID      uint   `json:"category_id,omitempty"`
	CategoryName    string `json:"category_name,omitempty"`
	SeoTitle        string `json:"seo_title,omitempty"`
	SeoUrl          string `json:"seo_url,omitempty"`
	PageTitle       string `json:"page_title,omitempty"`
	MetaDesc        string `json:"meta_desc,omitempty"`
	RelatedArticles []int  `json:"related_articles,omitempty"`
	Content         string `json:"content,omitempty"`
	AuthorID        uint   `json:"author_id,omitempty"`
	AuthorName      string `json:"author_name,omitempty"`
	CoverImageUrl   string `json:"cover_image_url,omitempty"`
	State           int    `json:"state,omitempty"`
	Language        string `json:"language,omitempty"`
	CreatedAt       int    `json:"created_at,omitempty"`
	CreatedBy       string `json:"created_by,omitempty"`
	UpdatedAt       int    `json:"updated_at,omitempty"`
	UpdatedBy       string `json:"updated_by,omitempty"`
	StarNum         int    `json:"star_num,omitempty"` // 点赞数
	ReadNum         int    `json:"read_num,omitempty"` // 阅读数
}

// ToArticleDto 从数据库结构抽取前端需要的字段返回
func (itself *Article) ToArticleDto(hasContent bool) ArticleDto {
	content := ""
	if hasContent {
		content = itself.Content
	}
	return ArticleDto{
		ID:              itself.ID,
		CategoryID:      itself.Category.ID,
		CategoryName:    itself.Category.Name,
		SeoUrl:          itself.SeoUrl,
		SeoTitle:        itself.SeoTitle,
		PageTitle:       itself.PageTitle,
		MetaDesc:        itself.MetaDesc,
		RelatedArticles: util.String2Int(strings.Split(itself.RelatedArticles, ",")),
		Content:         content,
		AuthorID:        itself.Author.ID,
		AuthorName:      itself.Author.Name,
		CoverImageUrl:   itself.CoverImageUrl,
		State:           itself.State,
		Language:        itself.Language,
		CreatedAt:       itself.CreatedAt,
		UpdatedAt:       itself.UpdatedAt,
		CreatedBy:       itself.CreatedBy,
		UpdatedBy:       itself.UpdatedBy,
		StarNum:         util.RandFromRange(300, 500),
		ReadNum:         util.RandFromRange(900, 1400),
	}
}

// ExistArticleByID checks if an article exists based on ID
func ExistArticleByID(id int) (bool, error) {
	var article Article
	err := Db.Model(&Article{}).Select("id").Where("id = ? ", id).Find(&article).Error
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
func GetArticles(pageNum int, pageSize int, cond string, vals []interface{}) ([]ArticleDto, int64, error) {
	var articles []Article
	var count int64
	Db.Model(&Article{}).Where(cond, vals...).Count(&count)
	err := Db.Preload("Category").Preload("Author").Where(cond, vals...).Offset(pageNum).Limit(pageSize).Find(&articles).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}

	resp := make([]ArticleDto, len(articles))

	for i, aa := range articles {
		resp[i] = aa.ToArticleDto(true)
	}
	return resp, count, nil
}

// GetArticleByID Get a single article based on ID
func GetArticleByID(id int) (*ArticleDto, error) {
	var article Article
	err := Db.Where("id = ? ", id).Find(&article).Error
	if err != nil {
		return nil, err
	}
	resp := article.ToArticleDto(true)
	return &resp, nil
}

func GetArticleBySeoUrl(url string) (*ArticleDto, error) {
	var article Article
	err := Db.Where("seo_url = ? ", url).First(&article).Error
	if err != nil {
		return nil, err
	}
	resp := article.ToArticleDto(true)
	return &resp, nil
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

func GetLatestArticle(cnt int) ([]ArticleDto, error) {
	var articles []Article
	var count int64
	Db.Model(&Article{}).Count(&count)
	if count < int64(cnt) {
		cnt = int(count)
	}
	err := Db.Preload("Category").Preload("Author").Order("created_at desc").Limit(cnt).Find(&articles).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	resp := make([]ArticleDto, 0)
	for _, art := range articles {
		resp = append(resp, art.ToArticleDto(false))
	}
	return resp, nil
}
