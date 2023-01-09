package e

const (
	SUCCESS       = 200
	ERROR         = 500
	InvalidParams = 400

	ErrorExistCategory      = 10001
	ErrorExistCategoryFail  = 10002
	ErrorNotExistCategory   = 10003
	ErrorGetCategoriesFail  = 10004
	ErrorCountCategoryFail  = 10005
	ErrorAddCategoryFail    = 10006
	ErrorEditCategoryFail   = 10007
	ErrorDeleteCategoryFail = 10008
	ErrorExportCategoryFail = 10009
	ErrorImportCategoryFail = 10010

	ErrorNotExistArticle       = 10011
	ErrorCheckExistArticleFail = 10012
	ErrorAddArticleFail        = 10013
	ErrorDeleteArticleFail     = 10014
	ErrorEditArticleFail       = 10015
	ErrorCountArticleFail      = 10016
	ErrorGetArticlesFail       = 10017
	ErrorGetArticleFail        = 10018
	ErrorGenArticlePosterFail  = 10019

	ErrorAuthCheckTokenFail    = 20001
	ErrorAuthCheckTokenTimeout = 20002
	ErrorAuthToken             = 20003
	ErrorAuth                  = 20004

	ErrorUploadSaveImageFail    = 30001
	ErrorUploadCheckImageFail   = 30002
	ErrorUploadCheckImageFormat = 30003
	ErrorPutFileToQiniu         = 30004

	ErrorExistAuthorFail = 40001
	ErrorNotExistAuthor  = 40002
	ErrorExistAuthor     = 40003
	ErrorAddAuthorFail   = 40004
	ErrorEditAuthorFail  = 40005
	ErrorGetAuthorFail   = 40006
	ErrorCountAuthorFail = 40007

	ErrorCheckExistUser  = 50001
	ErrorUserHasExist    = 50002
	ErrorCreatUser       = 50003
	ErrorGetUserInfoFail = 50004
	ErrorGetUserFail     = 50005

	ErrorGetLotteryFail   = 60001
	ErrorGetLuckytodyFail = 60002
)
