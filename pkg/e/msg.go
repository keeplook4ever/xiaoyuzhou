package e

var MsgFlags = map[int]string{
	SUCCESS:                     "ok",
	ERROR:                       "fail",
	InvalidParams:               "请求参数错误",
	ErrorExistCategory:          "已存在该类型名称",
	ErrorExistCategoryFail:      "获取已存在类型失败",
	ErrorNotExistCategory:       "该类型不存在",
	ErrorGetCategoriesFail:      "获取所有类型失败",
	ErrorCountCategoryFail:      "统计类型失败",
	ErrorAddCategoryFail:        "新增类型失败",
	ErrorEditCategoryFail:       "修改类型失败",
	ErrorDeleteCategoryFail:     "删除类型失败",
	ErrorExportCategoryFail:     "导出类型失败",
	ErrorImportCategoryFail:     "导入类型失败",
	ErrorNotExistArticle:        "该文章不存在",
	ErrorAddArticleFail:         "新增文章失败",
	ErrorDeleteArticleFail:      "删除文章失败",
	ErrorCheckExistArticleFail:  "检查文章是否存在失败",
	ErrorEditArticleFail:        "修改文章失败",
	ErrorCountArticleFail:       "统计文章失败",
	ErrorGetArticlesFail:        "获取多个文章失败",
	ErrorGetArticleFail:         "获取单个文章失败",
	ErrorGenArticlePosterFail:   "生成文章海报失败",
	ErrorAuthCheckTokenFail:     "Token鉴权失败",
	ErrorAuthCheckTokenTimeout:  "Token已超时",
	ErrorAuthToken:              "Token生成失败",
	ErrorAuth:                   "Token错误",
	ErrorUploadSaveImageFail:    "保存图片失败",
	ErrorUploadCheckImageFail:   "检查图片失败",
	ErrorUploadCheckImageFormat: "校验图片错误，图片格式或大小有问题",
	ErrorPutFileToQiniu:         "上传图片到七牛云失败",
	ErrorExistAuthorFail:        "获取作者失败",
	ErrorNotExistAuthor:         "该作者不存在",
	ErrorExistAuthor:            "该作者已存在",
	ErrorAddAuthorFail:          "添加作者失败",
	ErrorEditAuthorFail:         "编辑作者失败",
	ErrorGetAuthorFail:          "获取作者失败",
	ErrorCountAuthorFail:        "统计作者失败",
	ErrorCheckExistUser:         "检查用户失败",
	ErrorUserHasExist:           "用户已经存在",
	ErrorCreatUser:              "创建用户失败",
	ErrorGetUserInfoFail:        "获取当前登录用户信息失败",
	ErrorGetUserFail:            "获取用户列表失败",
	ErrorGetLotteryFail:         "获取运势失败",
	ErrorGetLuckytodyFail:       "获取每日好运失败",
}

// GetMsg get error information based on Code
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
