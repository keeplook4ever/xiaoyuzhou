package article_service

import (
	"encoding/json"
	"strings"
	"xiaoyuzhou/models"
	"xiaoyuzhou/pkg/util"
)

type ArticleInput struct {
	ID              int
	CategoryID      int
	SeoTitle        string
	SeoUrl          string
	PageTitle       string
	MetaDesc        string
	RelatedArticles []int
	Content         string
	AuthorId        int
	CoverImageUrl   string
	State           int
	Language        string
	UpdatedBy       string
	UpdatedAt       int
	CreatedBy       string
	CreatedAt       int
	PageNum         int
	PageSize        int
}

func (a *ArticleInput) Add() error {
	marshalData, _ := json.Marshal(a.RelatedArticles)
	relatedA := strings.Trim(string(marshalData), "[]")

	article := map[string]interface{}{
		"category_id":      a.CategoryID,
		"seo_title":        a.SeoTitle,
		"seo_url":          a.SeoUrl,
		"page_title":       a.PageTitle,
		"meta_desc":        a.MetaDesc,
		"related_articles": relatedA,
		"content":          a.Content,
		"author_id":        a.AuthorId,
		"cover_image_url":  a.CoverImageUrl,
		"state":            a.State,
		"language":         a.Language,
		"created_by":       a.CreatedBy,
		"updated_by":       a.UpdatedBy,
	}

	return models.AddArticle(article)
}

func (a *ArticleInput) Edit() error {
	data := make(map[string]interface{})

	marshalData, _ := json.Marshal(a.RelatedArticles)
	relatedA := strings.Trim(string(marshalData), "[]")

	data["updated_by"] = a.UpdatedBy
	data["state"] = a.State
	if a.CategoryID > 0 {
		data["category_id"] = a.CategoryID
	}
	if a.AuthorId > 0 {
		data["author_id"] = a.AuthorId
	}
	if a.SeoTitle != "" {
		data["seo_title"] = a.SeoTitle
	}
	if a.SeoUrl != "" {
		data["seo_url"] = a.SeoUrl
	}
	if a.PageTitle != "" {
		data["page_title"] = a.PageTitle
	}
	if a.MetaDesc != "" {
		data["meta_desc"] = a.MetaDesc
	}
	if a.Content != "" {
		data["content"] = a.Content
	}
	if a.CoverImageUrl != "" {
		data["cover_image_url"] = a.CoverImageUrl
	}
	if len(a.RelatedArticles) > 0 {
		data["related_articles"] = relatedA
	}
	return models.EditArticle(a.ID, data)
}

func (a *ArticleInput) Get() ([]models.ArticleDto, int64, error) {
	var (
		articles []models.ArticleDto
		//cacheArticles []manager.ArticleDto
	)

	//cache := cache_service.ArticleInput{
	//	ID:         a.ID,
	//	CreatedBy:  a.CreatedBy,
	//	CategoryID: a.CategoryID,
	//	State:      a.State,
	//	AuthorId:   a.AuthorId,
	//	PageNum:    a.PageNum,
	//	PageSize:   a.PageSize,
	//}
	//key := cache.GetArticlesKey()
	//if gredis.Exists(key) {
	//	data, err := gredis.Get(key)
	//	if err != nil {
	//		logging.Info(err)
	//	} else {
	//		json.Unmarshal(data, &cacheArticles)
	//		return cacheArticles, nil
	//	}
	//}
	cond, vals, err := util.SqlWhereBuild(a.getMaps(), "and")
	if err != nil {
		return nil, 0, err
	}
	articles, count, err := models.GetArticles(a.PageNum, a.PageSize, cond, vals)
	if err != nil {
		return nil, 0, err
	}

	//gredis.Set(key, articles, 3600)
	return articles, count, nil
}

func (a *ArticleInput) Delete() error {
	return models.DeleteArticle(a.ID)
}

func (a *ArticleInput) ExistByID() (bool, error) {
	return models.ExistArticleByID(a.ID)
}

func (a *ArticleInput) Count() (int64, error) {
	cond, vals, err := util.SqlWhereBuild(a.getMaps(), "and")
	if err != nil {
		return 0, err
	}
	return models.GetArticleTotal(cond, vals)
}

func (a *ArticleInput) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	marshalData, _ := json.Marshal(a.RelatedArticles)
	relatedA := strings.Trim(string(marshalData), "[]")
	if a.State != -1 {
		maps["state"] = a.State
	}
	if a.CategoryID != -1 {
		maps["category_id"] = a.CategoryID
	}
	if a.ID > 0 {
		maps["id"] = a.ID
	}
	if a.CreatedBy != "" {
		maps["created_by like"] = "%" + a.CreatedBy + "%"
	}
	if a.AuthorId > 0 {
		maps["author_id"] = a.AuthorId
	}
	if a.SeoTitle != "" {
		maps["seo_title like"] = "%" + a.SeoTitle + "%"
	}
	if a.SeoUrl != "" {
		maps["seo_url"] = a.SeoUrl
	}
	if a.PageTitle != "" {
		maps["page_title like"] = "%" + a.PageTitle + "%"
	}
	if a.MetaDesc != "" {
		maps["meta_desc like"] = "%" + a.MetaDesc + "%"
	}
	if a.CoverImageUrl != "" {
		maps["cover_image_url"] = a.CoverImageUrl
	}
	if len(a.RelatedArticles) > 0 {
		maps["related_articles"] = relatedA
	}
	return maps
}


func GetArticleForPlayer(cnt int) ([]models.ArticleDto, error) {
	return models.GetLatestArticle(cnt)
}

func GetSpecificArticleForPlayer(id int) (*models.ArticleDto, error) {
	return models.GetArticleByID(id)
}