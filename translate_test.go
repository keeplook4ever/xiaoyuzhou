package main

import (
	"encoding/json"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tmt "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tmt/v20180321"
	"os"
	"testing"
)

func Test_Trans(t *testing.T) {
	credential := common.NewCredential(
		os.Getenv("TENCENTCLOUD_SECRET_ID"),
		os.Getenv("TENCENTCLOUD_SECRET_KEY"),
	)
	t.Log(credential)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "tmt.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := tmt.NewClient(credential, "ap-guangzhou", cpf)

	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := tmt.NewTextTranslateBatchRequest()
	request.Source = common.StringPtr("zh")
	request.Target = common.StringPtr("ja")
	request.ProjectId = common.Int64Ptr(0)
	request.SourceTextList = common.StringPtrs([]string{"你好", "白羊座很好"})

	// 返回的resp是一个TextTranslateBatchResponse的实例，与请求对象对应
	response, err := client.TextTranslateBatch(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		t.Error(err)
		return
	}
	if err != nil {
		t.Error(err)
		panic(err)
	}

	// 输出json格式的字符串回包
	//t.Log(response.ToJsonString())

	var result ResponseOfTranslated
	err = json.Unmarshal([]byte(response.ToJsonString()), &result)
	if err != nil {
		t.Error(err)
	}
	t.Log(result.Response.TargetTextList)

}

type ResponseOfTranslated struct {
	Response ResponseStruct `json:"Response"`
}

type ResponseStruct struct {
	Source         string   `json:"Source"`
	Target         string   `json:"Target"`
	TargetTextList []string `json:"TargetTextList"`
	RequestId      string   `json:"RequestId"`
}

func Test_Key(t *testing.T) {
	t.Log(os.Getenv("TENCENTCLOUD_SECRET_ID"))
	t.Log(os.Getenv("TENCENTCLOUD_SECRET_KEY"))
}
