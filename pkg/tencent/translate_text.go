package tencent

import (
	"encoding/json"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tmt "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tmt/v20180321"
	"os"
)

func TranslateText(txt string) (error, *string) {
	credential := common.NewCredential(
		os.Getenv("TENCENTCLOUD_SECRET_ID"),
		os.Getenv("TENCENTCLOUD_SECRET_KEY"),
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "tmt.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	reGion := os.Getenv("TENCENTCLOUD_REGION")
	client, _ := tmt.NewClient(credential, reGion, cpf)

	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := tmt.NewTextTranslateRequest()
	request.Source = common.StringPtr("zh")
	request.Target = common.StringPtr("ja")
	request.ProjectId = common.Int64Ptr(0)
	request.SourceText = common.StringPtr(txt)

	// 返回的resp是一个TextTranslateResponse的实例，与请求对象对应
	response, err := client.TextTranslate(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return err, nil
	}
	if err != nil {
		return err, nil
	}

	// 输出json格式的字符串回包
	//t.Log(response.ToJsonString())

	var result ResponseOfTranslated
	err = json.Unmarshal([]byte(response.ToJsonString()), &result)
	if err != nil {
		return err, nil
	}
	return nil, &result.Response.TargetText
}

type ResponseOfTranslated struct {
	Response ResponseStruct `json:"Response"`
}

type ResponseStruct struct {
	Source     string `json:"Source"`
	Target     string `json:"Target"`
	TargetText string `json:"TargetText"`
	RequestId  string `json:"RequestId"`
}
