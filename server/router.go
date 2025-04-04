package server

import (
	"github.com/FXAZfung/image-board/server/handles"
	"github.com/FXAZfung/image-board/server/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Init(router *gin.Engine) {
	// 处理跨域
	Cors(router)

	// API 文档
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 直接访问图片相关路由（无需认证）
	imagesGroup := router.Group("/images")
	{
		imagesGroup.GET("/image/:name", handles.GetImageByName)
		imagesGroup.GET("/image/random", handles.GetRandomImage)
		imagesGroup.GET("/thumbnail/:name", handles.GetThumbnailByName)
	}

	api := router.Group("/api")

	// 登录
	authApi := api.Group("/auth")
	{
		authApi.POST("/login", handles.Login)
		authApi.POST("/logout", handles.Logout).Use(middleware.AuthMiddleware)
	}

	// 图片
	imageApi := api.Group("/image")
	{
		imageApiAuth := imageApi.Group("").Use(middleware.AuthMiddleware)
		{
			imageApiAuth.POST("/upload", handles.UploadImage)
			imageApiAuth.POST("/delete", handles.DeleteImage)
			imageApiAuth.POST("/tag/add", handles.AddTagToImage)
			imageApiAuth.POST("/tag/remove", handles.RemoveTagFromImage)
		}
		imageApi.POST("/list", handles.ListImages)
		imageApi.GET("/count", handles.GetImageCount)
		imageApi.POST("/tag/list", handles.GetImagesByTag)
	}

	// 设置
	settingApi := api.Group("/setting")
	{
		settingApi.GET("", handles.PublicSettings)
		settingApiAuth := settingApi.Group("").Use(middleware.AuthMiddleware)
		{
			//settingApiAuth.GET("/name", handles.GetSetting)
			settingApiAuth.POST("/save", handles.SaveSettings)
			settingApiAuth.GET("/list", handles.ListSettings)
			//settingApiAuth.DELETE("/delete", handles.DeleteSetting)
		}
	}
	// 标签
	tagApi := api.Group("/tag")
	{
		tagApi.POST("/list", handles.ListTags)
		tagApi.GET("/popular", handles.MostPopularTags)
		tagApi.GET("/search", handles.SearchTags)
		tagApi.GET("/image/:image_id", handles.GetTagsByImage)
		tagApi.GET("/name", handles.GetTagByName)
		tagApi.GET("/:id", handles.GetTagByID)
	}
}

// Cors 跨域配置
func Cors(router *gin.Engine) {
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	})
}
