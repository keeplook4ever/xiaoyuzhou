package main

import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"os"
	"testing"
)

func Test_Trans(t *testing.T) {
	credential := common.NewCredential(
		os.Getenv("TENCENTCLOUD_SECRET_ID"),
		os.Getenv("TENCENTCLOUD_SECRET_KEY"),
	)
	t.Log(credential)
	// 非必要步骤
	// 实例化一个客户端配置对象，可以指定超时时间等配置
	cpf := profile.NewClientProfile()
	// SDK默认使用POST方法。
	// 如果你一定要使用GET方法，可以在这里设置。GET方法无法处理一些较大的请求。
	// 如非必要请不要修改默认设置。
	cpf.HttpProfile.ReqMethod = "POST"
	// SDK有默认的超时时间，如非必要请不要修改默认设置。
	// 如有需要请在代码中查阅以获取最新的默认值。
	cpf.HttpProfile.ReqTimeout = 30
	// SDK会自动指定域名。通常是不需要特地指定域名的，但是如果你访问的是金融区的服务，
	// 则必须手动指定域名，例如云服务器的上海金融区域名： cvm.ap-shanghai-fsi.tencentcloudapi.com
	cpf.HttpProfile.Endpoint = "tmt.tencentcloudapi.com"
	// SDK默认用TC3-HMAC-SHA256进行签名，它更安全但是会轻微降低性能。
	// 如非必要请不要修改默认设置。
	cpf.SignMethod = "TC3-HMAC-SHA256"
	// SDK 默认用 zh-CN 调用返回中文。此外还可以设置 en-US 返回全英文。
	// 但大部分产品或接口并不支持全英文的返回。
	// 如非必要请不要修改默认设置。
	cpf.Language = "en-US"

	t.Log(cpf)
}
