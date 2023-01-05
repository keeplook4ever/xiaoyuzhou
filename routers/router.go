package routers

import (
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"xiaoyuzhou/middleware/cors"
	"xiaoyuzhou/middleware/jwt"
	"xiaoyuzhou/routers/api/v1/player"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	_ "xiaoyuzhou/docs"
	"xiaoyuzhou/pkg/export"
	"xiaoyuzhou/pkg/qrcode"
	"xiaoyuzhou/pkg/upload"
	"xiaoyuzhou/routers/api/v1/manager"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.Cors())

	r.StaticFS("/api/v1/export", http.Dir(export.GetExcelFullPath()))
	r.StaticFS("/api/v1/upload/images", http.Dir(upload.GetImageFullPath()))
	r.StaticFS("/api/v1/qrcode", http.Dir(qrcode.GetQrCodeFullPath()))

	r.POST("/api/v1/manager/user/auth", manager.GetAuth)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/api/v1/upload", manager.UploadImage)

	// 用户使用接口
	apiPlayerV1 := r.Group("/api/v1/player/")
	{

		// 日签相关
		lotteryV1 := apiPlayerV1.Group("/lottery")
		{
			lotteryV1.GET("", player.GetLottery)
		}

		// 塔罗牌相关
		apiPlayerV1.Group("/tarot")
		{

		}

		// 文章相关
		articlePlayer := apiPlayerV1.Group("/article")
		{
			// 获取文章
			articlePlayer.GET("", player.GetArticleForPlayer)
		}

		// 星座相关
		apiPlayerV1.Group("/constellation")
		{

		}

		// 支付相关
		apiPlayerV1.Group("/pay")
		{

		}
	}

	// 后台管理接口
	apiManagerV1 := r.Group("/api/v1/manager/")
	apiManagerV1.Use(jwt.JWT())
	{
		//获取类型列表
		apiManagerV1.GET("/category", manager.GetCategory)
		//新建类型
		apiManagerV1.POST("/category", manager.AddCategory)
		//更新指定类型
		apiManagerV1.PUT("/category/:id", manager.EditCategory)
		//删除指定类型
		apiManagerV1.DELETE("/category/:id", manager.DeleteCategory)

		//获取文章
		apiManagerV1.GET("/articles", manager.GetArticles)
		//新建文章
		apiManagerV1.POST("/articles", manager.AddArticle)
		//更新指定文章
		apiManagerV1.PUT("/articles/:id", manager.EditArticle)
		//删除指定文章
		apiManagerV1.DELETE("/articles/:id", manager.DeleteArticle)
		//生成文章海报
		apiManagerV1.POST("/articles/poster/generate", manager.GenerateArticlePoster)
		//上传文章图片
		apiManagerV1.POST("/articles/img", manager.UploadImage)

		//添加作者接口
		apiManagerV1.POST("/author", manager.AddAuthor)
		//修改作者信息
		apiManagerV1.PUT("/author/:id", manager.EditAuthor)
		//获取作者
		apiManagerV1.GET("/author", manager.GetAuthors)

		//添加用户
		apiManagerV1.POST("/user", manager.AddUser)
		//获取用户
		apiManagerV1.GET("/user", manager.GetUser)
		//获取当前登录用户信息
		apiManagerV1.GET("/user/info", manager.GetCurrentLoginUserInfo)

	}

	return r
}
