package util

import "xiaoyuzhou/pkg/setting"

// Setup Initialize the util
func Setup() {
	jwtSecret = []byte(setting.AppSetting.JwtSecret)
}
