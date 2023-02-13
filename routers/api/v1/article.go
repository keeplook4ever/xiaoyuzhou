package v1

import (
	"net/http"
	"time"
	"xiaoyuzhou/models"
	"xiaoyuzhou/pkg/util"
	"xiaoyuzhou/service/article_service"
	"xiaoyuzhou/service/author_service"
	"xiaoyuzhou/service/category_service"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"

	"xiaoyuzhou/pkg/app"
	"xiaoyuzhou/pkg/e"
)

type AddArticleForm struct {
	CategoryID      int    `json:"category_id" binding:"required"`
	SeoTitle        string `json:"seo_title" binding:"required"`
	SeoUrl          string `json:"seo_url" binding:"required"`
	PageTitle       string `json:"page_title" binding:"required"`
	MetaDesc        string `json:"meta_desc" binding:"required"`
	RelatedArticles []int  `json:"related_articles"`
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
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
		return
	}

	// 判断是否类型存在
	categoryService := category_service.CategoryInput{ID: article.CategoryID}
	exists, err := categoryService.ExistByID()
	if err != nil {
		appG.Response(http.StatusOK, e.ErrorExistCategoryFail, nil)
		return
	}

	if !exists {
		appG.Response(http.StatusOK, e.ErrorNotExistCategory, nil)
		return
	}

	// 判断是否作者存在
	authorService := author_service.AuthorInput{ID: article.AuthorId}
	exists, err = authorService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorExistAuthorFail, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ErrorNotExistAuthor, nil)
		return
	}

	starNum := util.RandFromRange(100, 800)
	time.Sleep(time.Microsecond)
	readNum := util.RandFromRange(100, 800)

	// 点赞数要小于阅读数
	if starNum > readNum {
		util.SwapTwoInt(&starNum, &readNum)
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
		ReadNum:         readNum,
		StarNum:         starNum,
	}
	if err = articleService.Add(); err != nil {
		appG.Response(http.StatusOK, err.Error(), nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

type EditArticleForm struct {
	ID              int    `json:"id" binding:"required"`
	CategoryID      int    `json:"category_id"`
	SeoTitle        string `json:"seo_title"`
	SeoUrl          string `json:"seo_url"`
	PageTitle       string `json:"page_title"`
	MetaDesc        string `json:"meta_desc"`
	RelatedArticles []int  `json:"related_articles"`
	Content         string `json:"content"`
	AuthorId        int    `json:"author_id"`
	CoverImageUrl   string `json:"cover_image_url"`
	State           int    `json:"state" enums:"1,0" default:"1"` // 0表示禁用，1表示启用
	Language        string `json:"language" enums:"jp,zh,en" default:"jp"`
	UpdatedBy 		string `json:"updated_by"`
}

// EditArticle
// @Summary 修改文章
// @Produce  json
// @Param _ body EditArticleForm true "修改参数"
// @Param id path int true "文章ID"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Security ApiKeyAuth
// @Router /manager/articles/{id} [put]
// @Tags Manager
func EditArticle(c *gin.Context) {
	var (
		appG    = app.Gin{C: c}
		article = EditArticleForm{ID: com.StrTo(c.Param("id")).MustInt(),
			UpdatedBy: c.GetString("username")}
	)

	if err := c.ShouldBindJSON(&article); err != nil {
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
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
		appG.Response(http.StatusInternalServerError, e.ErrorCheckExistArticleFail, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ErrorNotExistArticle, nil)
		return
	}

	if article.CategoryID > 0 {
		// 判断类型是否存在
		tagService := category_service.CategoryInput{ID: article.CategoryID}
		exists, err = tagService.ExistByID()
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ErrorExistCategoryFail, nil)
			return
		}

		if !exists {
			appG.Response(http.StatusOK, e.ErrorNotExistCategory, nil)
			return
		}
	}

	if article.AuthorId > 0 {
		// 判断是否作者存在
		authorService := author_service.AuthorInput{ID: article.AuthorId}
		exists, err = authorService.ExistByID()
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ErrorExistAuthorFail, nil)
			return
		}
		if !exists {
			appG.Response(http.StatusOK, e.ErrorNotExistAuthor, nil)
			return
		}
	}

	err = articleService.Edit()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorEditArticleFail, nil)
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
		appG.Response(http.StatusOK, e.InvalidParams, nil)
		return
	}

	articleService := article_service.ArticleInput{ID: id}
	exists, err := articleService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorCheckExistArticleFail, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ErrorNotExistArticle, nil)
		return
	}

	err = articleService.Delete()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorDeleteArticleFail, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

type GetArticlesResponse struct {
	Lists []models.ArticleDto `json:"lists"`
	Count int64               `json:"total"` //符合条件的总数，不是单页数量
}

// GetArticles
// @Summary 获取文章
// @Produce  json
// @Param category_id query int false "Category ID"
// @Param author_id query int false "Author ID"
// @Param state query int false "State"
// @Param created_by query int false "CreatedBy"
// @Param id query int false "ID"
// @Param seo_title query string false "SEO Title"
// @Param seo_url query string false "SEO Url"
// @Param page_title query string false "Page Title"
// @Param meta_desc query string false "Meta Desc"
// @Param cover_image_url query string false "Cover Img URL"
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

	seoTitle := c.Query("seo_title")
	seoUrl := c.Query("seo_url")
	pageTitle := c.Query("page_title")
	metaDesc := c.Query("meta_desc")
	coverImageUrl := c.Query("cover_image_url")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
		return
	}

	articleService := article_service.ArticleInput{
		SeoTitle:      seoTitle,
		SeoUrl:        seoUrl,
		PageTitle:     pageTitle,
		MetaDesc:      metaDesc,
		CoverImageUrl: coverImageUrl,
		ID:            id,
		CreatedBy:     createdBy,
		CategoryID:    tagId,
		AuthorId:      authorId,
		State:         state,
		PageNum:       util.GetPage(c),
		PageSize:      util.GetPageSize(c),
	}

	// 获取带content的文章
	articles, count, err := articleService.Get(true)
	if err != nil {
		appG.Response(http.StatusOK, e.ErrorGetArticlesFail, nil)
		return
	}
	var res GetArticlesResponse
	res.Lists = articles
	res.Count = count

	appG.Response(http.StatusOK, e.SUCCESS, res)
}

// GetIndexArticleForPlayer
// @Summary 首页展示文章
// @Produce json
// @Success 200 {object} []models.ArticleDto
// @Failure 500 {object} app.Response
// @Tags Player
// @Router /player/articles/index [get]
func GetIndexArticleForPlayer(c *gin.Context) {
	appG := app.Gin{C: c}
	articleList, err := article_service.GetArticleForPlayer(4)
	if err != nil {
		appG.Response(http.StatusOK, "获取首页文章失败", nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, articleList)
}

// GetArticlesAll
// @Summary 获取文章(可传文章id)
// @Produce json
// @Param id_list query []int false "ID list"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Tags Player
// @Router /player/articles [get]
func GetArticlesAll(c *gin.Context) {
	appG := app.Gin{C: c}
	idList := c.Query("id_list")

	// -1 代表无此参数
	state := -1
	tagId := -1
	authorId := -1
	if idList == "" {
		articleService := article_service.ArticleInput{
			State:      state,
			CategoryID: tagId,
			AuthorId:   authorId,
			PageNum:    util.GetPage(c),
			PageSize:   util.GetPageSize(c),
		}
		// 获取不带content的文章
		articles, count, err := articleService.Get(false)
		if err != nil {
			appG.Response(http.StatusOK, e.ErrorGetArticlesFail, nil)
			return
		}
		var res GetArticlesResponse
		res.Lists = articles
		res.Count = count

		appG.Response(http.StatusOK, e.SUCCESS, res)
	} else {
		ids := util.StringToIntSlice(idList)
		article, err := article_service.GetSpecificArticleByIDs(ids, false)
		if err != nil {
			appG.Response(http.StatusOK, "获取文章失败", nil)
			return
		}

		var res GetArticlesResponse
		res.Lists = article
		res.Count = int64(len(article))
		appG.Response(http.StatusOK, e.SUCCESS, res)
	}

}

// GetSpecificArticleForPlayer
// @Summary 根据文章SEO URL获取特定文章给用户
// @Param seo_url query string true "SEO URL"
// @Accept json
// @Produce json
// @Success 200 {object} models.ArticleDto
// @Failure 500 {object} app.Response
// @Router /player/article [get]
// @Tags Player
func GetSpecificArticleForPlayer(c *gin.Context) {
	appG := app.Gin{
		C: c,
	}
	seoUrl := c.Query("seo_url")

	article, err := article_service.GetSpecificArticleBySeoUrl(seoUrl)
	if err != nil {
		appG.Response(http.StatusOK, err.Error(), nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, article)
}

// StarOneArticle
// @Summary 用户点赞文章API
// @Param id path int true "文章ID"
// @Param uid query string true "用户ID"
// @Accept json
// @Produce json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /player/article/star/{id} [put]
// @Tags Player
func StarOneArticle(c *gin.Context) {
	appG := app.Gin{C: c}
	id := com.StrTo(c.Param("id")).MustInt()
	uid := c.Query("uid")
	err := article_service.StarArticle(id, uid)
	if err != nil {
		appG.Response(http.StatusOK, "点赞失败", nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}


type StarStatusResp struct {
	Status int `json:"status" enums:"1,0"` // 1代表点赞过，0代表没点赞
}


// GetStarStatus
// @Summary 获取用户是否已经点赞此文章
// @Param id path int true "文章ID"
// @Param uid query string true "用户ID"
// @Accept json
// @Produce json
// @Success 200 {object} StarStatusResp
// @Failure 500 {object} app.Response
// @Router /player/article/star/{id} [get]
// @Tags Player
func GetStarStatus(c *gin.Context) {
	appG := app.Gin{C: c}
	id := com.StrTo(c.Param("id")).MustInt()
	uid := c.Query("uid")
	stared, err := article_service.GetArticleStarStatus(id, uid)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	if stared {
		appG.Response(http.StatusOK, e.SUCCESS, StarStatusResp{Status: 1})
	} else {
		appG.Response(http.StatusOK, e.SUCCESS, StarStatusResp{Status: 0})
	}
}