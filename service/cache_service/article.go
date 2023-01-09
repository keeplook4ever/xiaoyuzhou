package cache_service

import (
	"strconv"
	"strings"

	"xiaoyuzhou/pkg/e"
)

type ArticleInput struct {
	ID         int
	CategoryID int
	State      int
	AuthorId   int
	CreatedBy  string
	PageNum    int
	PageSize   int
}

func (a *ArticleInput) GetArticleKey() string {
	return e.CacheArticle + "_" + strconv.Itoa(a.ID)
}

func (a *ArticleInput) GetArticlesKey() string {
	keys := []string{
		e.CacheArticle,
		"LIST",
	}

	if a.ID > 0 {
		keys = append(keys, strconv.Itoa(a.ID))
	}
	if a.CategoryID > 0 {
		keys = append(keys, strconv.Itoa(a.CategoryID))
	}
	if a.AuthorId > 0 {
		keys = append(keys, strconv.Itoa(a.AuthorId))
	}
	if a.State >= 0 {
		keys = append(keys, strconv.Itoa(a.State))
	}
	if a.PageNum > 0 {
		keys = append(keys, strconv.Itoa(a.PageNum))
	}
	if a.PageSize > 0 {
		keys = append(keys, strconv.Itoa(a.PageSize))
	}

	return strings.Join(keys, "_")
}
