package article_service

import (
	"encoding/json"
	"xiaoyuzhou/models/manager"
	"xiaoyuzhou/service/manager/cache_service"

	"xiaoyuzhou/pkg/gredis"
	"xiaoyuzhou/pkg/logging"
)

type Article struct {
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
	ModifiedBy      string

	PageNum  int
	PageSize int
}

func (a *Article) Add() error {
	article := map[string]interface{}{
		"cate_id":          a.CategoryID,
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
	}

	if err := manager.AddArticle(article); err != nil {
		return err
	}

	return nil
}

func (a *Article) Edit() error {
	return manager.EditArticle(a.ID, map[string]interface{}{
		"cate_id":         a.CategoryID,
		"seo_title":       a.SeoTitle,
		"page_title":      a.PageTitle,
		"meta_desc":       a.MetaDesc,
		"content":         a.Content,
		"cover_image_url": a.CoverImageUrl,
		"state":           a.State,
		"modified_by":     a.ModifiedBy,
		"author_id":       a.AuthorId,
	})
}

func (a *Article) Get() (*manager.Article, error) {
	var cacheArticle *manager.Article

	cache := cache_service.Article{ID: a.ID}
	key := cache.GetArticleKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheArticle)
			return cacheArticle, nil
		}
	}

	article, err := manager.GetArticle(a.ID)
	if err != nil {
		return nil, err
	}

	gredis.Set(key, article, 3600)
	return article, nil
}

func (a *Article) GetAll() ([]*manager.Article, error) {
	var (
		articles, cacheArticles []*manager.Article
	)

	cache := cache_service.Article{
		TagID:    a.CategoryID,
		State:    a.State,
		AuthorId: a.AuthorId,
		PageNum:  a.PageNum,
		PageSize: a.PageSize,
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

func (a *Article) Delete() error {
	return manager.DeleteArticle(a.ID)
}

func (a *Article) ExistByID() (bool, error) {
	return manager.ExistArticleByID(a.ID)
}

func (a *Article) Count() (int, error) {
	return manager.GetArticleTotal(a.getMaps())
}

func (a *Article) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0
	if a.State != -1 {
		maps["state"] = a.State
	}
	if a.CategoryID != -1 {
		maps["tag_id"] = a.CategoryID
	}

	return maps
}
