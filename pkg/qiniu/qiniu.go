package qiniu

import (
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"mime/multipart"
	"xiaoyuzhou/pkg/logging"
)

// 自定义返回值结构体
type MyPutRet struct {
	Key    string
	Hash   string
	Fsize  int
	Bucket string
	Name   string
}

func GetUpTools() (string, *storage.FormUploader, storage.PutRet) {
	bucket := "xiaoyuzhou-1"
	accessKey := "uYMK8yo6XKsRogBh5jb8A45WaW9pvq-xtiehnxzH"
	secretKey := "rmOnqYQgQ_2hc_skDfgfN-zu7h3gdyV5W7f06y1R"
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuadong
	// 是否使用https域名
	cfg.UseHTTPS = true
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	return upToken, formUploader, ret
}

func UploadLocalImg(path string) error {
	upToken, formUploader, ret := GetUpTools()

	key := ""
	// 可选配置
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "test img",
		},
	}

	err := formUploader.PutFile(context.Background(), &ret, upToken, key, path, &putExtra)
	if err != nil {
		logging.Error(err)
		return err
	}
	fmt.Println(ret.Key, ret.Hash)

	domain := "http://ro048sb6x.hd-bkt.clouddn.com"
	publicAccessURL := storage.MakePublicURL(domain, ret.Key)
	fmt.Println(publicAccessURL)

	return nil
}

func UploadImg(file multipart.File, size int64, md5Img string) (url string, err error) {
	key := md5Img
	upToken, formUploader, ret := GetUpTools()
	err = formUploader.Put(context.Background(), &ret, upToken, key, file, size, nil)
	if err != nil {
		return "", err
	}
	domain := "http://ro048sb6x.hd-bkt.clouddn.com"
	publicAccessURL := storage.MakePublicURL(domain, ret.Key)
	fmt.Println(publicAccessURL)

	return publicAccessURL, nil
}
