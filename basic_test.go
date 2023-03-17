package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
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

func TestTime(t *testing.T) {
	//year := time.Now().Format("2006")
	//month := time.Now().Format("01")
	//day := time.Now().Format("02")
	//fmt.Println(year, month, day)
	//str := year[2:] + month + day
	//t.Log(str)
	ts := time.Now().Unix()
	t.Log(ts)
	ls := strconv.FormatInt(ts, 10)[5:]
	t.Log(ls)

	amount := 21.85
	ta := int(amount)
	t.Log(ta)
}
