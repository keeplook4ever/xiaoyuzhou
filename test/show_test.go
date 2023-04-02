package test

/**
SDK依赖第三方包，mahonia（字符编码处理）,uuid（生成uuid），需要预先获取依赖包
go env -w GO111MODULE=on
go get -u "github.com/axgle/mahonia"
go get -u "github.com/satori/go.uuid"
**/

import (
	"encoding/json"
	"reflect"
	"testing"
	"xiaoyuzhou/pkg/xingzuoapi"
)

type xingZuoApiRes struct {
	ShowApiCode   int           `json:"showapi_res_code"`
	ShowApiFeeNum int           `json:"showapi_fee_num"`
	ShowApiError  string        `json:"showapi_res_error"`
	ShowApiId     string        `json:"showapi_res_id"`
	ShowApiBody   BodyOfShowApi `json:"showapi_res_body"`
}

type BodyOfShowApi struct {
	Star         string      `json:"star"`
	RetCode      int         `json:"ret_code"`
	DayContext   dayContext  `json:"day"`
	Tomorrow     dayContext  `json:"tomorrow"`
	MonthContext interface{} `json:"month"`
	WeekC        WeekContext `json:"week"`
	YearC        YearContext `json:"year"`
}

type dayContext struct {
	LoveTxt        string `json:"love_txt"`
	WorkTxt        string `json:"work_txt"`
	WorkStar       int    `json:"work_star"`
	MoneyStar      int    `json:"money_star"`
	LuckyColor     string `json:"lucky_color"`
	LuckyTime      string `json:"lucky_time"`
	LoveStar       int    `json:"love_star"`
	LuckyDirection string `json:"lucky_direction"`
	SummaryStar    int    `json:"summary_star"`
	Time           string `json:"time"`
	MoneyTxt       string `json:"money_txt"`
	GeneralTxt     string `json:"general_txt"`
	Grxz           string `json:"grxz"`
	LuckyNum       string `json:"lucky_num"`
	DayNotice      string `json:"day_notice"`
}

type WeekContext struct {
	LoveTxt     string `json:"love_txt"`
	HealthTxt   string `json:"health_txt"`
	WorkTxt     string `json:"work_txt"`
	LuckyDay    string `json:"lucky_day"`
	WorkStar    int    `json:"work_star"`
	WeekNotice  string `json:"week_notice"`
	MoneyStar   int    `json:"money_star"`
	LuckyColor  string `json:"lucky_color"`
	LoveStar    int    `json:"love_star"`
	SummaryStar int    `json:"summary_star"`
	Time        string `json:"time"`
	MoneyTxt    string `json:"money_txt"`
	GeneralTxt  string `json:"general_txt"`
	Grxz        string `json:"grxz"`
	LuckyNum    string `json:"lucky_num"`
	Xrxz        string `json:"xrxz"`
	DayNotice   string `json:"day_notice"`
}

type YearContext struct {
	LoveTxt      string `json:"love_txt"`
	HealthTxt    string `json:"health_txt"`
	Time         string `json:"time"`
	WorkTxt      string `json:"work_txt"`
	MoneyTxt     string `json:"money_txt"`
	WorkIndex    string `json:"work_index"`
	MoneyIndex   string `json:"money_index"`
	GeneralTxt   string `json:"general_txt"`
	OneWord      string `json:"oneword"`
	GeneralIndex string `json:"general_index"`
	LoveIndex    string `json:"love_index"`
}

func Test_show(t *testing.T) {
	showapi_appid := 1370755                           //要替换成自己的
	showapi_sign := "bd37a3ceb71a40a9bfd7ad19085ec725" //要替换成自己的
	res := xingzuoapi.ShowapiRequest("http://route.showapi.com/872-1", showapi_appid, showapi_sign)
	res.AddTextPara("star", "baiyang")
	res.AddTextPara("needTomorrow", "1")
	res.AddTextPara("needWeek", "1")
	res.AddTextPara("needMonth", "1")
	res.AddTextPara("needYear", "1")
	res.AddTextPara("date", "0111")
	//fmt.Println(res.Post())
	result, err := res.Post()
	if err != nil {
		t.Error("Error")
	}

	//t.Log(result)
	t.Log(reflect.TypeOf(result))
	var dataS xingZuoApiRes
	if err := json.Unmarshal([]byte(result), &dataS); err != nil {
		t.Error(err)
	} else {
		t.Log(dataS)
	}
}
