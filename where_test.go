package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"xiaoyuzhou/models"
	"xiaoyuzhou/pkg/util"
)

func TestWhereBuild(t *testing.T) {

	var author models.Author
	cond, vals, err := util.SqlWhereBuild(map[string]interface{}{
		"name like": "%q%",
		"age in":    []int{12, 19, 18},
		"id =":      1,
	}, "or")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(cond, reflect.TypeOf(cond))
	t.Log(vals, reflect.TypeOf(vals))

	models.Db.Where(cond, vals...).Find(&author)
	t.Log(author.ID, author.Name, author.Desc)

	s := fmt.Sprintf(" %s ", strings.ToUpper("and"))
	t.Log(s)
}
