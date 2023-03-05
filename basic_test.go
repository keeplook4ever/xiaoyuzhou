package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
	"xiaoyuzhou/pkg/util"
)

func Test_Basic(t *testing.T) {
	ori := []string{"a", "b", "b", "c"}
	s, _ := json.Marshal(ori)
	str := string(s)
	fmt.Println(reflect.TypeOf(str))

	var obj []string
	err := json.Unmarshal([]byte(str), &obj)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(obj, reflect.TypeOf(obj))

	ts15 := int(time.Now().UnixMilli())
	fmt.Println(ts15)
}

func TestSlice(t *testing.T) {
	or := "[\"1\",\"2\",3]"

	str := strings.Trim(or, "[]\"'")
	t.Log(str)
	res := util.StringToIntSlice(or)
	t.Log(res)
}
