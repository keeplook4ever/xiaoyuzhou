package util

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"xiaoyuzhou/pkg/setting"
)

// Setup Initialize the util
func Setup() {
	jwtSecret = []byte(setting.AppSetting.JwtSecret)
}

// String2Int 字符串切片转数字切片
func String2Int(strArr []string) []int {
	res := make([]int, len(strArr))

	for index, val := range strArr {
		res[index], _ = strconv.Atoi(val)
	}
	return res
}

type NullType byte

const (
	_ NullType = iota
	// IsNull the same as `is null`
	IsNull
	// IsNotNull the same as `is not null`
	IsNotNull
)

// SqlWhereBuild sql build where
func SqlWhereBuild(where map[string]interface{}, connect string) (whereSQL string, vals []interface{}, err error) {
	for k, v := range where {
		ks := strings.Split(k, " ")
		if len(ks) > 2 {
			return "", nil, fmt.Errorf("Error in query condition: %s. ", k)
		}

		if whereSQL != "" {
			whereSQL += fmt.Sprintf(" %s ", strings.ToUpper(connect))
		}
		strings.Join(ks, ",")
		switch len(ks) {
		case 1:
			//fmt.Println(reflect.TypeOf(v))
			switch v := v.(type) {
			case NullType:
				if v == IsNotNull {
					whereSQL += fmt.Sprint(k, " IS NOT NULL")
				} else {
					whereSQL += fmt.Sprint(k, " IS NULL")
				}
			default:
				whereSQL += fmt.Sprint(k, "=?")
				vals = append(vals, v)
			}
			break
		case 2:
			k = ks[0]
			switch ks[1] {
			case "=":
				whereSQL += fmt.Sprint(k, "=?")
				vals = append(vals, v)
				break
			case ">":
				whereSQL += fmt.Sprint(k, ">?")
				vals = append(vals, v)
				break
			case ">=":
				whereSQL += fmt.Sprint(k, ">=?")
				vals = append(vals, v)
				break
			case "<":
				whereSQL += fmt.Sprint(k, "<?")
				vals = append(vals, v)
				break
			case "<=":
				whereSQL += fmt.Sprint(k, "<=?")
				vals = append(vals, v)
				break
			case "!=":
				whereSQL += fmt.Sprint(k, "!=?")
				vals = append(vals, v)
				break
			case "<>":
				whereSQL += fmt.Sprint(k, "!=?")
				vals = append(vals, v)
				break
			case "in":
				whereSQL += fmt.Sprint(k, " in (?)")
				vals = append(vals, v)
				break
			case "like":
				whereSQL += fmt.Sprint(k, " like ?")
				vals = append(vals, v)
			}
			break
		}
	}
	return
}

// RandFromRange 从最小最大获取一个中间随机数
func RandFromRange(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(max-min) + min
	return n
}

// SwapTwoInt 交换两个变量的值
func SwapTwoInt(a *int, b *int) {
	*a, *b = *b, *a
}

func StringToIntSlice(str string) (intSlice []int) {
	if str == "" {
		return
	}
	str = strings.Trim(str, "[]")
	strList := strings.Split(str, ",")
	if len(strList) == 0 {
		return
	}
	for _, item := range strList {
		if item == "" {
			continue
		}
		val, err := strconv.Atoi(item)
		if err != nil {
			continue
		}
		intSlice = append(intSlice, val)
	}
	return
}

func StringSlice2String(ori []string) *string {
	if ori == nil {
		return nil
	}
	tarBytes, err := json.Marshal(ori)
	if err != nil {
		return nil
	}
	str := string(tarBytes)
	return &str
}

func String2StringSlice(ori string) []string {
	var obj []string
	err := json.Unmarshal([]byte(ori), &obj)
	if err != nil {
		return nil
	}
	return obj
}

func IfInSlice(sSl []string, targ string) bool {
	for _, v := range sSl {
		if targ == v {
			return true
		}
	}
	return false
}
