package cache_service

import (
	"strconv"
	"strings"

	"xiaoyuzhou/pkg/e"
)

type Article struct {
	ID         int
	CategoryID int
	State      int
	AuthorId   int
	CreatedBy  string
	PageNum    int
	PageSize   int
}

func (a *Article) GetArticleKey() string {
	return e.CACHE_ARTICLE + "_" + strconv.Itoa(a.ID)
}

func (a *Article) GetArticlesKey() string {
	keys := []string{
		e.CACHE_ARTICLE,
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
