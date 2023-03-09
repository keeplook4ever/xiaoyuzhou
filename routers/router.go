package routers

import (
	ginSwagger "github.com/swaggo/gin-swagger"
	"xiaoyuzhou/middleware/cors"
	"xiaoyuzhou/middleware/jwt"
	"xiaoyuzhou/routers/api/v1"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	_ "xiaoyuzhou/docs"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.Cors())

	r.POST("/api/v1/manager/user/auth", v1.GetAuth)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/health", v1.HealthCheck)
	// 用户使用接口
	apiPlayerV1 := r.Group("/api/v1/player/")
	{

		// 获取日签相关
		apiPlayerV1.GET("/lottery", v1.GetLotteryForUser)

		{
			// 塔罗抽取接口：返回塔罗图片和名字
			apiPlayerV1.POST("/tarot/one", v1.GetTarotOne)

			// 获取单个塔罗答案接口：根据订单返回
			apiPlayerV1.GET("/tarot/one/answer", v1.GetTarotOneAnswer)

			// 添加webhook监听paypal事件
			apiPlayerV1.POST("/tarot/webhook/paypal", v1.ReceiveOrderEventsFromPayPal)
		}

		// 文章相关
		{
			// 获取全部文章，可根据id获取
			apiPlayerV1.GET("/articles", v1.GetArticlesAll)

			// 首页展示最新几篇文章
			apiPlayerV1.GET("/articles/index", v1.GetIndexArticleForPlayer)

			// 首页获取某个特定文章
			apiPlayerV1.GET("/article", v1.GetSpecificArticleForPlayer)

			// 点赞文章接口
			apiPlayerV1.PUT("/article/star/:id", v1.StarOneArticle)

			// 获取是否用户点赞该文章接口
			apiPlayerV1.GET("/article/star/:id", v1.GetStarStatus)
		}

		// 星座相关
		apiPlayerV1.Group("/constellation")
		{

		}

		// PayPal相关
		paypalV1 := apiPlayerV1.Group("/paypal")
		{
			//创建PayPal支付订单
			paypalV1.POST("/checkout/orders", v1.CreatePayPalOrder)

			//确认支付订单
			paypalV1.POST("/confirm/orders/:order_id", v1.ConfirmPayment)

			//捕获PayPal订单
			paypalV1.POST("/capture/orders/:order_id", v1.CapturePayPalOrder)

			//获取PayPal订单详情
			paypalV1.GET("/checkout/orders/:order_id", v1.GetPayPalOrderDetail)
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

		//创建塔罗牌
		apiManagerV1.POST("/tarot", v1.AddTarot)
		//修改塔罗牌
		apiManagerV1.PUT("/tarot/:id", v1.EditTarot)

		//获取塔罗牌
		apiManagerV1.GET("/tarot", v1.GetTarot)

		//创建塔罗牌价格
		apiManagerV1.POST("/tarot/price", v1.SetPrice)
		//获取塔罗牌价格
		apiManagerV1.GET("/tarot/price", v1.GetPrice)
		//修改塔罗牌价格
		apiManagerV1.PUT("/tarot/price", v1.UpdatePrice)

	}

	return r
}
