package e

var MsgFlags = map[int]string{
	SUCCESS:                         "ok",
	ERROR:                           "fail",
	INVALID_PARAMS:                  "请求参数错误",
	ERROR_EXIST_CATEGORY:            "已存在该类型名称",
	ERROR_EXIST_CATEGORY_FAIL:       "获取已存在类型失败",
	ERROR_NOT_EXIST_CATEGORY:        "该类型不存在",
	ERROR_GET_CATEGORYS_FAIL:        "获取所有类型失败",
	ERROR_COUNT_CATEGORY_FAIL:       "统计类型失败",
	ERROR_ADD_CATEGORY_FAIL:         "新增类型失败",
	ERROR_EDIT_CATEGORY_FAIL:        "修改类型失败",
	ERROR_DELETE_CATEGORY_FAIL:      "删除类型失败",
	ERROR_EXPORT_CATEGORY_FAIL:      "导出类型失败",
	ERROR_IMPORT_CATEGORY_FAIL:      "导入类型失败",
	ERROR_NOT_EXIST_ARTICLE:         "该文章不存在",
	ERROR_ADD_ARTICLE_FAIL:          "新增文章失败",
	ERROR_DELETE_ARTICLE_FAIL:       "删除文章失败",
	ERROR_CHECK_EXIST_ARTICLE_FAIL:  "检查文章是否存在失败",
	ERROR_EDIT_ARTICLE_FAIL:         "修改文章失败",
	ERROR_COUNT_ARTICLE_FAIL:        "统计文章失败",
	ERROR_GET_ARTICLES_FAIL:         "获取多个文章失败",
	ERROR_GET_ARTICLE_FAIL:          "获取单个文章失败",
	ERROR_GEN_ARTICLE_POSTER_FAIL:   "生成文章海报失败",
	ERROR_AUTH_CHECK_TOKEN_FAIL:     "Token鉴权失败",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT:  "Token已超时",
	ERROR_AUTH_TOKEN:                "Token生成失败",
	ERROR_AUTH:                      "Token错误",
	ERROR_UPLOAD_SAVE_IMAGE_FAIL:    "保存图片失败",
	ERROR_UPLOAD_CHECK_IMAGE_FAIL:   "检查图片失败",
	ERROR_UPLOAD_CHECK_IMAGE_FORMAT: "校验图片错误，图片格式或大小有问题",
	ERROR_PUT_FILE_TO_QINIU:         "上传图片到七牛云失败",
	ERROR_EXIST_AUTHOR_FAIL:         "获取作者失败",
	ERROR_NOT_EXIST_AUTHOR:          "该作者不存在",
	ERROR_EXIST_AUTHOR:              "该作者已存在",
	ERROR_ADD_AUTHOR_FAIL:           "添加作者失败",
	ERROR_EDIT_AUTHOR_FAIL:          "编辑作者失败",
	ERROR_GET_AUTHOR_FAIL:           "获取作者失败",
	ERROR_COUNT_AUTHOR_FAIL:         "统计作者失败",
	ERROR_CHECK_EXIST_USER:          "检查用户失败",
	ERROR_USER_HAS_EXIST:            "用户已经存在",
	ERROR_CREAT_USER:                "创建用户失败",
}

// GetMsg get error information based on Code
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
