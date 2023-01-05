package cache_service

import (
	"strconv"
	"strings"

	"xiaoyuzhou/pkg/e"
)

type CategoryInput struct {
	ID    int
	Name  string
	State int

	PageNum  int
	PageSize int
}

func (t *CategoryInput) GetCategoryKey() string {
	keys := []string{
		e.CACHE_CATEGORY,
		"LIST",
	}

	if t.Name != "" {
		keys = append(keys, t.Name)
	}
	if t.State >= 0 {
		keys = append(keys, strconv.Itoa(t.State))
	}
	if t.PageNum > 0 {
		keys = append(keys, strconv.Itoa(t.PageNum))
	}
	if t.PageSize > 0 {
		keys = append(keys, strconv.Itoa(t.PageSize))
	}

	return strings.Join(keys, "_")
}
