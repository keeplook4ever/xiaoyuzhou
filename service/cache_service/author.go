package cache_service

import (
	"strconv"
	"strings"
	"xiaoyuzhou/pkg/e"
)

func (t *CategoryInput) GetAuthorsKey() string {
	keys := []string{
		e.CacheCategory,
		"LIST",
	}

	if t.Name != "" {
		keys = append(keys, t.Name)
	}

	if t.PageNum > 0 {
		keys = append(keys, strconv.Itoa(t.PageNum))
	}
	if t.PageSize > 0 {
		keys = append(keys, strconv.Itoa(t.PageSize))
	}

	return strings.Join(keys, "_")
}
