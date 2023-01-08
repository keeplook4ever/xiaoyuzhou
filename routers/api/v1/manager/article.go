package manager

import (
	"net/http"
	"xiaoyuzhou/models/manager"
	"xiaoyuzhou/pkg/setting"
	"xiaoyuzhou/pkg/util"
	"xiaoyuzhou/service/manager/article_service"
	"xiaoyuzhou/service/manager/author_service"
	"xiaoyuzhou/service/manager/category_service"

	"github.com/astaxie/beego/validation"
	"github.com/boombuler/barcode/qr"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"

	"xiaoyuzhou/pkg/app"
	"xiaoyuzhou/pkg/e"
	"xiaoyuzhou/pkg/qrcode"
)

type AddArticleForm struct {
	CategoryID      int    `json:"category_id" binding:"required"`
	SeoTitle        string `json:"seo_title" binding:"required"`
	SeoUrl          string `json:"seo_url" binding:"required"`
	PageTitle       string `json:"page_title" binding:"required"`
	MetaDesc        string `json:"meta_desc" binding:"required"`
	RelatedArticles string `json:"related_articles" binding:"required"`
	Content         string `json:"content" binding:"required"`
	AuthorId        int    `json:"author_id"  binding:"required"`
	CoverImageUrl   string `json:"cover_image_url" binding:"required"`
	State           int    `json:"state" binding:"required" enums:"1,0" default:"1"` // 0表示禁用，1表示启用
	Language        string `json:"language" binding:"required" enums:"jp,zh,en" default:"jp"`
}

// AddArticle
// @Summary 添加文章
// @Produce  json
// @Param _ body AddArticleForm true "文章详情"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Security ApiKeyAuth
// @Router /manager/articles [post]
// @Tags Manager
func AddArticle(c *gin.Context) {
	var (
		appG    = app.Gin{C: c}
		article AddArticleForm
	)
	if err := c.ShouldBindJSON(&article); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	// 判断是否类型存在
	categoryService := category_service.CategoryInput{ID: article.CategoryID}
	exists, err := categoryService.ExistByID()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_EXIST_CATEGORY_FAIL, nil)
		return
	}

	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_CATEGORY, nil)
		return
	}

	// 判断是否作者存在
	authorService := author_service.AuthorInput{ID: article.AuthorId}
	exists, err = authorService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_AUTHOR_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_AUTHOR, nil)
		return
	}

	articleService := article_service.ArticleInput{
		CategoryID:      article.CategoryID,
		SeoTitle:        article.SeoTitle,
		SeoUrl:          article.SeoUrl,
		PageTitle:       article.PageTitle,
		MetaDesc:        article.MetaDesc,
		RelatedArticles: article.RelatedArticles,
		Content:         article.Content,
		AuthorId:        article.AuthorId,
		CoverImageUrl:   article.CoverImageUrl,
		State:           article.State,
		Language:        article.Language,
		CreatedBy:       c.GetString("username"), // 根据登录态获取
		UpdatedBy:       c.GetString("username"),
	}
	if err = articleService.Add(); err != nil {
		appG.Response(http.StatusOK, e.ERROR_ADD_ARTICLE_FAIL, err.Error())
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

type EditArticleForm struct {
	ID              int    `form:"id" binding:"required"`
	CategoryID      int    `form:"category_id"`
	SeoTitle        string `form:"seo_title"`
	SeoUrl          string `form:"seo_url"`
	PageTitle       string `form:"page_title"`
	MetaDesc        string `form:"meta_desc"`
	AuthorId        int    `form:"author_id"`
	Content         string `form:"content"`
	UpdatedBy       string `form:"updated_by" binding:"required"`
	CoverImageUrl   string `form:"cover_image_url"`
	State           int    `form:"state" enums:"0,1"`
	RelatedArticles string `form:"related_articles"`
}

// EditArticle
// @Summary 修改文章
// @Produce  json
// @Param id path int true "ID"
// @Param category_id formData int false "Category ID"
// @Param page_title formData string false "Page Title"
// @Param seo_title formData string false "SEO Title"
// @Param seo_url formData string false "SEO URL"
// @Param related_articles formData string false "Related Articles"
// @Param meta_desc formData string false "Desc"
// @Param author_id formData int false "Author ID"
// @Param content formData string false "Content"
// @Param cover_image_url formData string false "Cover img URL"
// @Param state formData int false "State" default(1)
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
//@Security ApiKeyAuth
// @Router /manager/articles/{id} [put]
// @Tags Manager
func EditArticle(c *gin.Context) {
	var (
		appG    = app.Gin{C: c}
		article = EditArticleForm{ID: com.StrTo(c.Param("id")).MustInt(),
			UpdatedBy: c.GetString("username")}
	)

	if err := c.ShouldBind(&article); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	articleService := article_service.ArticleInput{
		ID:              article.ID,
		CategoryID:      article.CategoryID,
		PageTitle:       article.PageTitle,
		SeoTitle:        article.SeoTitle,
		SeoUrl:          article.SeoUrl,
		MetaDesc:        article.MetaDesc,
		Content:         article.Content,
		CoverImageUrl:   article.CoverImageUrl,
		RelatedArticles: article.RelatedArticles,
		UpdatedBy:       c.GetString("username"), // 后端获取，通过登录态
		AuthorId:        article.AuthorId,
		State:           article.State,
	}
	exists, err := articleService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	if article.CategoryID > 0 {
		// 判断类型是否存在
		tagService := category_service.CategoryInput{ID: article.CategoryID}
		exists, err = tagService.ExistByID()
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_CATEGORY_FAIL, nil)
			return
		}

		if !exists {
			appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_CATEGORY, nil)
			return
		}
	}

	if article.AuthorId > 0 {
		// 判断是否作者存在
		authorService := author_service.AuthorInput{ID: article.AuthorId}
		exists, err = authorService.ExistByID()
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_AUTHOR_FAIL, nil)
			return
		}
		if !exists {
			appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_AUTHOR, nil)
			return
		}
	}

	err = articleService.Edit()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EDIT_ARTICLE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// DeleteArticle
// @Summary 删除文章
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Security ApiKeyAuth
// @Router /manager/articles/{id} [delete]
// @Tags Manager
func DeleteArticle(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt()
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	articleService := article_service.ArticleInput{ID: id}
	exists, err := articleService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	err = articleService.Delete()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_ARTICLE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

type GetArticlesResponse struct {
	Lists []manager.ArticleDto `json:"lists"`
	Count int                  `json:"count"`
}

// GetArticles
// @Summary 获取文章
// @Produce  json
// @Param category_id query int false "Category ID"
// @Param author_id query int false "Author ID"
// @Param state query int false "State"
// @Param created_by query int false "CreatedBy"
// @Param id query int false "ID"
// @Success 200 {object} GetArticlesResponse
// @Failure 500 {object} app.Response
// @Security ApiKeyAuth
// @Router /manager/articles [get]
// @Tags Manager
func GetArticles(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}

	state := -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state")
	}

	tagId := -1
	if arg := c.Query("category_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()
		valid.Min(tagId, 1, "category_id")
	}

	authorId := -1
	if arg := c.Query("author_id"); arg != "" {
		authorId = com.StrTo(arg).MustInt()
		valid.Min(authorId, 1, "author_id")
	}

	createdBy := c.Query("created_by")

	id := com.StrTo(c.Query("id")).MustInt()

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	articleService := article_service.ArticleInput{
		ID:         id,
		CreatedBy:  createdBy,
		CategoryID: tagId,
		AuthorId:   authorId,
		State:      state,
		PageNum:    util.GetPage(c),
		PageSize:   setting.AppSetting.PageSize,
	}

	total, err := articleService.Count()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_COUNT_ARTICLE_FAIL, nil)
		return
	}

	articles, err := articleService.Get()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLES_FAIL, nil)
		return
	}

	data := make(map[string]interface{})
	data["lists"] = articles
	data["total"] = total

	appG.Response(http.StatusOK, e.SUCCESS, data)
}

const (
	QRCODE_URL = "https://github.com/EDDYCJY/blog#gin%E7%B3%BB%E5%88%97%E7%9B%AE%E5%BD%95"
)

func GenerateArticlePoster(c *gin.Context) {
	appG := app.Gin{C: c}
	article := &article_service.ArticleInput{}
	qr := qrcode.NewQrCode(QRCODE_URL, 300, 300, qr.M, qr.Auto)
	posterName := article_service.GetPosterFlag() + "-" + qrcode.GetQrCodeFileName(qr.URL) + qr.GetQrCodeExt()
	articlePoster := article_service.NewArticlePoster(posterName, article, qr)
	articlePosterBgService := article_service.NewArticlePosterBg(
		"bg.jpg",
		articlePoster,
		&article_service.Rect{
			X0: 0,
			Y0: 0,
			X1: 550,
			Y1: 700,
		},
		&article_service.Pt{
			X: 125,
			Y: 298,
		},
	)

	_, filePath, err := articlePosterBgService.Generate()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GEN_ARTICLE_POSTER_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"poster_url":      qrcode.GetQrCodeFullUrl(posterName),
		"poster_save_url": filePath + posterName,
	})
}
