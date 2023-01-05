package article_service

import (
	"encoding/json"
	"time"
	"xiaoyuzhou/models/manager"
	"xiaoyuzhou/pkg/gredis"
	"xiaoyuzhou/pkg/logging"
	"xiaoyuzhou/service/manager/cache_service"
)

type ArticleInput struct {
	ID              int
	CategoryID      int
	SeoTitle        string
	SeoUrl          string
	PageTitle       string
	MetaDesc        string
	RelatedArticles string
	Content         string
	AuthorId        int
	CoverImageUrl   string
	State           int
	Language        string
	UpdatedBy       string
	UpdatedAt       time.Time
	CreatedBy       string
	CreatedAt       time.Time
	PageNum         int
	PageSize        int
}

func (a *ArticleInput) Add() error {
	article := map[string]interface{}{
		"category_id":      a.CategoryID,
		"seo_title":        a.SeoTitle,
		"seo_url":          a.SeoUrl,
		"page_title":       a.PageTitle,
		"meta_desc":        a.MetaDesc,
		"related_articles": a.RelatedArticles,
		"content":          a.Content,
		"author_id":        a.AuthorId,
		"cover_image_url":  a.CoverImageUrl,
		"state":            a.State,
		"language":         a.Language,
		"created_by":       a.CreatedBy,
		"updated_by":       a.UpdatedBy,
	}

	return manager.AddArticle(article)
}

func (a *ArticleInput) Edit() error {
	return manager.EditArticle(a.ID, map[string]interface{}{
		"category_id":     a.CategoryID,
		"seo_title":       a.SeoTitle,
		"page_title":      a.PageTitle,
		"meta_desc":       a.MetaDesc,
		"content":         a.Content,
		"cover_image_url": a.CoverImageUrl,
		"state":           a.State,
		"updated_by":      a.UpdatedBy,
		"updated_at":      a.UpdatedAt,
		"author_id":       a.AuthorId,
	})
}

type ArticleReturn struct {
	Id              int    `json:"id"`
	CategoryID      int    `json:"category_id"`
	CategoryName    string `json:"category_name"`
	SeoTitle        string `json:"seo_title"`
	SeoUrl          string `json:"seo_url"`
	PageTitle       string `json:"page_title"`
	MetaDesc        string `json:"meta_desc"`
	RelatedArticles string `json:"related_articles"`
	Content         string `json:"content"`
	AuthorId        int    `json:"author_id"`
	CoverImageUrl   string `json:"cover_image_url"`
	State           int    `json:"state"`
	Language        string `json:"language"`
	ModifiedBy      string `json:"modified_by"`
}

func (a *ArticleInput) Get() ([]manager.ArticleDto, error) {
	var (
		articles      []manager.ArticleDto
		cacheArticles []manager.ArticleDto
	)

	cache := cache_service.ArticleInput{
		ID:         a.ID,
		CreatedBy:  a.CreatedBy,
		CategoryID: a.CategoryID,
		State:      a.State,
		AuthorId:   a.AuthorId,
		PageNum:    a.PageNum,
		PageSize:   a.PageSize,
	}
	key := cache.GetArticlesKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheArticles)
			return cacheArticles, nil
		}
	}

	articles, err := manager.GetArticles(a.PageNum, a.PageSize, a.getMaps())
	if err != nil {
		return nil, err
	}

	gredis.Set(key, articles, 3600)
	return articles, nil
}

func (a *ArticleInput) Delete() error {
	return manager.DeleteArticle(a.ID)
}

func (a *ArticleInput) ExistByID() (bool, error) {
	return manager.ExistArticleByID(a.ID)
}

func (a *ArticleInput) Count() (int64, error) {
	return manager.GetArticleTotal(a.getMaps())
}

func (a *ArticleInput) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
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
		maps["created_by"] = a.CreatedBy
	}
	if a.AuthorId > 0 {
		maps["author_id"] = a.AuthorId
	}

	return maps
}
