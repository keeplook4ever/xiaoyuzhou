package routers

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"xiaoyuzhou/middleware/cors"
	"xiaoyuzhou/middleware/jwt"
	"xiaoyuzhou/routers/api/v1"

	swaggerFiles "github.com/swaggo/files"
	_ "xiaoyuzhou/docs"
	"xiaoyuzhou/pkg/export"
	"xiaoyuzhou/pkg/qrcode"
	"xiaoyuzhou/pkg/upload"
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

	r.POST("/api/v1/manager/user/auth", v1.GetAuth)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//r.POST("/api/v1/upload", v1.UploadImage)

	// 用户使用接口
	apiPlayerV1 := r.Group("/api/v1/player/")
	{

		// 日签相关
		lotteryV1 := apiPlayerV1.Group("/lottery")
		{
			lotteryV1.GET("", v1.GetLottery)
		}

		// 塔罗牌相关
		apiPlayerV1.Group("/tarot")
		{

		}

		// 文章相关
		{
			// 获取文章
			//articlePlayer.GET("", player.GetArticleForPlayer)
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
		// 获取S3上传权限token
		apiManagerV1.POST("/s3/token", v1.GetS3Token)

		//获取类型列表
		apiManagerV1.GET("/category", v1.GetCategory)
		//新建类型
		apiManagerV1.POST("/category", v1.AddCategory)
		//更新指定类型
		apiManagerV1.PUT("/category/:id", v1.EditCategory)
		//删除指定类型
		apiManagerV1.DELETE("/category/:id", v1.DeleteCategory)

		//获取文章
		apiManagerV1.GET("/articles", v1.GetArticles)
		//新建文章
		apiManagerV1.POST("/articles", v1.AddArticle)
		//更新指定文章
		apiManagerV1.PUT("/articles/:id", v1.EditArticle)
		//删除指定文章
		apiManagerV1.DELETE("/articles/:id", v1.DeleteArticle)
		//生成文章海报
		apiManagerV1.POST("/articles/poster/generate", v1.GenerateArticlePoster)
		////上传文章图片
		//apiManagerV1.POST("/articles/img", v1.UploadImage)

		//添加作者接口
		apiManagerV1.POST("/author", v1.AddAuthor)
		//修改作者信息
		apiManagerV1.PUT("/author/:id", v1.EditAuthor)
		//获取作者
		apiManagerV1.GET("/author", v1.GetAuthors)

		//添加用户
		apiManagerV1.POST("/user", v1.AddUser)
		//获取用户
		apiManagerV1.GET("/user", v1.GetUser)
		//获取当前登录用户信息
		apiManagerV1.GET("/user/info", v1.GetCurrentLoginUserInfo)

	}

	return r
}
