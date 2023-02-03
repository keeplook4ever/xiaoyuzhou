package routers

import (
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"xiaoyuzhou/middleware/cors"
	"xiaoyuzhou/middleware/jwt"
	"xiaoyuzhou/routers/api/v1"

	"github.com/gin-gonic/gin"

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

	r.GET("/health", v1.HealthCheck)
	// 用户使用接口
	apiPlayerV1 := r.Group("/api/v1/player/")
	{

		// 获取日签相关
		apiPlayerV1.GET("/lottery", v1.GetLotteryForUser)

		// 塔罗牌相关
		apiPlayerV1.Group("/tarot")
		{

		}

		// 文章相关
		{
			// 获取文章
		}

		// 星座相关
		apiPlayerV1.Group("/constellation")
		{

		}

		// PayPal相关
		paypalV1 := apiPlayerV1.Group("/paypal")
		{
			//创建PayPal支付订单
			paypalV1.POST("/checkout/order", v1.CreatePayPalOrder)

			//捕获PayPal订单
			paypalV1.POST("/capture/order/:order_id", v1.CapturePayPalOrder)
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

		//获取运势Lottery配置表
		apiManagerV1.GET("/lottery", v1.GetLotteryForManager)
		////添加运势Lottery配置表
		//apiManagerV1.POST("/lottery", v1.AddLotteryType)
		//修改Lottery
		apiManagerV1.PUT("/lottery", v1.EditLottery)

		//添加具体运势内容LotteryContent
		apiManagerV1.POST("/lottery-content", v1.AddLotteryContent)
		//修改运势内容表LotteryContent
		apiManagerV1.PUT("/lottery-content/:id", v1.EditLotteryContent)
		//获取运势内容表LotteryContent
		apiManagerV1.GET("/lottery-content", v1.GetLotteryContentForManager)
		//删除LotteryContent
		apiManagerV1.DELETE("/lottery-content/:id", v1.DeleteLotteryContent)

		//添加今日好运相关内容
		apiManagerV1.POST("/lucky", v1.AddLucky)

		//上传今日好运excel文件
		apiManagerV1.POST("/lucky/upload", v1.UploadLucky)
		//修改今日好运内容
		apiManagerV1.PUT("/lucky", v1.EditLucky)
		//删除今日好运内容
		apiManagerV1.DELETE("/lucky", v1.DeleteLucky)
		//获取今日好运
		apiManagerV1.GET("/lucky", v1.GetLucky)
	}

	return r
}
