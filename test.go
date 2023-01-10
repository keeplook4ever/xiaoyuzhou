package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"xiaoyuzhou/pkg/util"
)

func main() {
	var s = []int{1, 2, 3}
	v, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	a := strings.Trim(string(v), "[]")
	fmt.Println(a, reflect.TypeOf(a))

	mm := []string{a}
	fmt.Println(mm, reflect.TypeOf(mm))
	ss := util.String2Int(mm)
	fmt.Println(ss, reflect.TypeOf(ss))
}
