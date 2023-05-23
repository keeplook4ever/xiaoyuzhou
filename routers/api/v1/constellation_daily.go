package v1

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"xiaoyuzhou/pkg/app"
	"xiaoyuzhou/pkg/e"
	"xiaoyuzhou/pkg/gredis"
	"xiaoyuzhou/pkg/logging"
	"xiaoyuzhou/pkg/tencent"
	"xiaoyuzhou/pkg/util"
	"xiaoyuzhou/pkg/xingzuoapi"
	"xiaoyuzhou/service/cache_service"
)

// 占卜页-星座

// GetDailyConstellation
// @Summary 获取星座运势
// @Param name query string true "星座名" Enums(baiyang,jinniu,shuangzi,juxie,shizi,chunv,tiancheng,tianxie,sheshou,mojie,shuiping,shuangyu)
// @Param language query string true "语言" Enums(zh,jp,en,tc)
// @Success 200 {object} RespOfConstellation
// @Failure 500 {object} app.Response
// @Tags Player
// @Router /player/constellation [get]
func GetDailyConstellation(c *gin.Context) {
	appG := app.Gin{C: c}
	name := c.Query("name")
	cstlList := []string{"baiyang", "jinniu", "shuangzi", "juxie", "shizi", "chunv", "tiancheng", "tianxie", "sheshou", "mojie", "shuiping", "shuangyu"}
	vaLid := util.IfInSlice(cstlList, name)
	if !vaLid {
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
		return
	}
	var dataS xingZuoApiRes

	cache := cache_service.XingzuoInput{
		Name: name,
	}
	key := cache.GetXingzuoKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Error(err)
		} else {
			err = json.Unmarshal([]byte(data), &dataS)
			if err != nil {
				logging.Error(err)
			}
		}
	}

	// dataS.ShowApiId == "" 表示上述获取缓存失败或者无缓存
	if dataS.ShowApiId == "" {
		showapi_appid := 1370755                           //要替换成自己的
		showapi_sign := "bd37a3ceb71a40a9bfd7ad19085ec725" //要替换成自己的
		res := xingzuoapi.ShowapiRequest("http://route.showapi.com/872-1", showapi_appid, showapi_sign)
		res.AddTextPara("star", "baiyang")
		res.AddTextPara("needTomorrow", "1")
		res.AddTextPara("needWeek", "1")
		res.AddTextPara("needMonth", "1")
		res.AddTextPara("needYear", "0")
		dateNow := time.Now().Format("0102")
		res.AddTextPara("date", dateNow)

		result, err := res.Post()
		if err != nil {
			appG.Response(http.StatusOK, "后台星座获取失败", nil)
			return
		}
		if err := json.Unmarshal([]byte(result), &dataS); err != nil {
			appG.Response(http.StatusOK, "星座数据解析失败", nil)
			return
		}
		if dataS.ShowApiError != "" {
			appG.Response(http.StatusOK, dataS.ShowApiError, nil)
			return
		}
		gredis.Set(key, dataS, 43200) // 12小时
	}

	waitForTransList := make([]string, 0)
	//today
	waitForTransList = append(waitForTransList, dataS.ShowApiBody.DayContext.LuckyColor)
	waitForTransList = append(waitForTransList, dataS.ShowApiBody.DayContext.Grxz)
	waitForTransList = append(waitForTransList, dataS.ShowApiBody.DayContext.LuckyDirection)
	waitForTransList = append(waitForTransList, dataS.ShowApiBody.DayContext.GeneralTxt)
	//tomorrow
	waitForTransList = append(waitForTransList, dataS.ShowApiBody.Tomorrow.LuckyColor)
	waitForTransList = append(waitForTransList, dataS.ShowApiBody.Tomorrow.Grxz)
	waitForTransList = append(waitForTransList, dataS.ShowApiBody.Tomorrow.LuckyDirection)
	waitForTransList = append(waitForTransList, dataS.ShowApiBody.Tomorrow.GeneralTxt)
	//week
	waitForTransList = append(waitForTransList, dataS.ShowApiBody.WeekC.LuckyColor)
	waitForTransList = append(waitForTransList, dataS.ShowApiBody.WeekC.Grxz)
	waitForTransList = append(waitForTransList, GetRandomDirection())
	waitForTransList = append(waitForTransList, dataS.ShowApiBody.WeekC.GeneralTxt)
	//month
	waitForTransList = append(waitForTransList, GetRandomColor())
	waitForTransList = append(waitForTransList, dataS.ShowApiBody.MonthC.Grxz)
	waitForTransList = append(waitForTransList, dataS.ShowApiBody.MonthC.LuckyDirection)
	waitForTransList = append(waitForTransList, dataS.ShowApiBody.MonthC.GeneralTxt)

	resAfterTrans := make([]string, len(waitForTransList))
	sourceL := c.Query("language")
	if sourceL == "zh" || sourceL == "tc" { // 繁体中文没有，用简体中文代替
		copy(resAfterTrans, waitForTransList)
	} else {
		err, resAfterTransG := tencent.TranslateTextList(waitForTransList, sourceL)
		resAfterTrans = resAfterTransG
		if err != nil {
			appG.Response(http.StatusOK, "翻译失败", nil)
			return
		}
	}

	// 获取需要的星座数据
	todayC := RespOfOneConstellation{
		LoveScore:       getScoreFromXingzuo(dataS.ShowApiBody.DayContext.LoveStar),
		WorkScore:       getScoreFromXingzuo(dataS.ShowApiBody.DayContext.WorkStar),
		MoneyScore:      getScoreFromXingzuo(dataS.ShowApiBody.DayContext.MoneyStar),
		HealthScore:     getScoreFromXingzuo(-1),
		SummaryScore:    getScoreFromXingzuo(dataS.ShowApiBody.DayContext.SummaryStar),
		LuckyColor:      resAfterTrans[0],
		LuckyColorOri:   waitForTransList[0],
		LuckyXingzuo:    resAfterTrans[1],
		LuckyXingzuoOri: GetXingzuoFromZh(waitForTransList[1]),
		LuckyDirection:  resAfterTrans[2],
		SummaryTxt:      resAfterTrans[3],
	}
	tomorrowC := RespOfOneConstellation{
		LoveScore:       getScoreFromXingzuo(dataS.ShowApiBody.Tomorrow.LoveStar),
		WorkScore:       getScoreFromXingzuo(dataS.ShowApiBody.Tomorrow.WorkStar),
		MoneyScore:      getScoreFromXingzuo(dataS.ShowApiBody.Tomorrow.MoneyStar),
		HealthScore:     getScoreFromXingzuo(-1),
		SummaryScore:    getScoreFromXingzuo(dataS.ShowApiBody.Tomorrow.SummaryStar),
		LuckyColor:      resAfterTrans[4],
		LuckyColorOri:   waitForTransList[4],
		LuckyXingzuo:    resAfterTrans[5],
		LuckyXingzuoOri: GetXingzuoFromZh(waitForTransList[5]),
		LuckyDirection:  resAfterTrans[6],
		SummaryTxt:      resAfterTrans[7],
	}
	weekC := RespOfOneConstellation{
		LoveScore:       getScoreFromXingzuo(dataS.ShowApiBody.WeekC.LoveStar),
		WorkScore:       getScoreFromXingzuo(dataS.ShowApiBody.WeekC.WorkStar),
		MoneyScore:      getScoreFromXingzuo(dataS.ShowApiBody.WeekC.MoneyStar),
		HealthScore:     getScoreFromXingzuo(-1),
		SummaryScore:    getScoreFromXingzuo(dataS.ShowApiBody.WeekC.SummaryStar),
		LuckyColor:      resAfterTrans[8],
		LuckyColorOri:   waitForTransList[8],
		LuckyXingzuo:    resAfterTrans[9],
		LuckyXingzuoOri: GetXingzuoFromZh(waitForTransList[9]),
		LuckyDirection:  resAfterTrans[10],
		SummaryTxt:      resAfterTrans[11],
	}
	monthC := RespOfOneConstellation{
		LoveScore:       getScoreFromXingzuo(dataS.ShowApiBody.MonthC.LoveStar),
		WorkScore:       getScoreFromXingzuo(dataS.ShowApiBody.MonthC.WorkStar),
		MoneyScore:      getScoreFromXingzuo(dataS.ShowApiBody.MonthC.MoneyStar),
		HealthScore:     getScoreFromXingzuo(-1),
		SummaryScore:    getScoreFromXingzuo(dataS.ShowApiBody.MonthC.SummaryStar),
		LuckyColor:      resAfterTrans[12],
		LuckyColorOri:   waitForTransList[12],
		LuckyXingzuo:    resAfterTrans[13],
		LuckyXingzuoOri: GetXingzuoFromZh(waitForTransList[13]),
		LuckyDirection:  resAfterTrans[14],
		SummaryTxt:      resAfterTrans[15],
	}
	resp := RespOfConstellation{
		Name:     name,
		Today:    todayC,
		Tomorrow: tomorrowC,
		Week:     weekC,
		Month:    monthC,
	}
	appG.Response(http.StatusOK, e.SUCCESS, resp)
}

type xingZuoApiRes struct {
	ShowApiCode   int           `json:"showapi_res_code"`
	ShowApiFeeNum int           `json:"showapi_fee_num"`
	ShowApiError  string        `json:"showapi_res_error"`
	ShowApiId     string        `json:"showapi_res_id"`
	ShowApiBody   BodyOfShowApi `json:"showapi_res_body"`
}

type BodyOfShowApi struct {
	Star       string       `json:"star"`
	RetCode    int          `json:"ret_code"`
	DayContext dayContext   `json:"day"`
	Tomorrow   dayContext   `json:"tomorrow"`
	MonthC     monthContent `json:"month"`
	WeekC      weekContext  `json:"week"`
	//YearC        YearContext `json:"year"`
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

type weekContext struct {
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

type monthContent struct {
	LoveTxt        string `json:"love_txt"`
	WorkTxt        string `json:"work_txt"`
	WorkStar       int    `json:"work_star"`
	MoneyStar      int    `json:"money_star"`
	MonthAdvantage string `json:"month_advantage"`
	Yfxz           string `json:"yfxz"`
	LuckyDirection string `json:"lucky_direction"`
	LoveStar       int    `json:"love_star"`
	MonthWeekness  string `json:"month_weekness"`
	SummaryStar    int    `json:"summary_star"`
	Time           string `json:"time"`
	MoneyTxt       string `json:"money_txt"`
	GeneralTxt     string `json:"general_txt"`
	Grxz           string `json:"grxz"`
	Xrxz           string `json:"xrxz"`
	LuckyNum       string `json:"lucky_num"`
}

type RespOfOneConstellation struct {
	LoveScore       int    `json:"love_score"`        // 爱情分
	WorkScore       int    `json:"work_score"`        // 事业分
	MoneyScore      int    `json:"money_score"`       // 金钱分
	HealthScore     int    `json:"health_score"`      // 健康分
	SummaryScore    int    `json:"summary_score"`     // 总体分
	LuckyColor      string `json:"lucky_color"`       // 幸运色
	LuckyColorOri   string `json:"lucky_color_ori"`   // 幸运色中文
	LuckyDirection  string `json:"lucky_direction"`   // 幸运方向
	LuckyXingzuo    string `json:"lucky_xingzuo"`     // 幸运星座
	LuckyXingzuoOri string `json:"lucky_xingzuo_ori"` // 星座枚举
	SummaryTxt      string `json:"summary_txt"`       // 运势概览
}

type RespOfConstellation struct {
	Name     string                 `json:"name"` // 星座名称
	Today    RespOfOneConstellation `json:"today"`
	Tomorrow RespOfOneConstellation `json:"tomorrow"`
	Week     RespOfOneConstellation `json:"week"`
	Month    RespOfOneConstellation `json:"month"`
}

// getScoreFromXingzuo 根据星座返回数字生成10-100的数
func getScoreFromXingzuo(star int) int {
	switch star {
	case 1:
		return util.RandFromRange(10, 20)
	case 2:
		return util.RandFromRange(21, 50)
	case 3:
		return util.RandFromRange(51, 70)
	case 4:
		return util.RandFromRange(71, 99)
	}
	return util.RandFromRange(50, 78)
}

// 颜色矩阵

//姜黄色
//青金石蓝
//绿色
//棕色
//稻草黄
//紫色
//玫红色
//蘑菇灰
//天蓝色
//玫瑰红
//橘色
//樱桃红
//橙色
//粉色
//琥珀棕
//银灰色
//墨绿色
//茄皮紫
//酒红色
//梅子青
//酸橙绿
//咖啡色
//松石绿
//土黄色
//卡其色
//琥珀橙

// 方向矩阵

// 东北方
// 东南方
// 西南方
// 西北方
// 正南方
// 正北方
// 正东方
// 正西方

func GetRandomDirection() string {
	directionList := []string{"东北方", "东南方", "西南方", "西北方", "正南方", "正北方", "正东方", "正西方", "茄皮紫"}
	ix := util.RandFromRange(0, len(directionList))
	return directionList[ix]
}

func GetRandomColor() string {
	colorList := []string{"姜黄色", "青金石蓝", "绿色", "棕色", "稻草黄", "紫色", "玫红色", "蘑菇灰", "天蓝色", "玫瑰红", "橘色", "樱桃红", "橙色", "粉色", "琥珀棕", "银灰色", "墨绿色", "酒红色", "梅子青", "酸橙绿", "咖啡色", "松石绿", "土黄色", "卡其色", "琥珀橙"}
	ix := util.RandFromRange(0, len(colorList))
	return colorList[ix]
}

func GetXingzuoFromZh(zh string) string {
	switch zh {
	case "白羊座":
		return "baiyang"
	case "金牛座":
		return "jinniu"
	case "双子座":
		return "shuangzi"
	case "巨蟹座":
		return "juxie"
	case "狮子座":
		return "shizi"
	case "处女座":
		return "chunv"
	case "天秤座":
		return "tiancheng"
	case "天蝎座":
		return "tianxie"
	case "射手座":
		return "shezhou"
	case "摩羯座":
		return "mojie"
	case "水瓶座":
		return "shuiping"
	case "双鱼座":
		return "shuangyu"
	}
	return "jinniu"
}
