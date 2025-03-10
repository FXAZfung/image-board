package server

import (
	"github.com/FXAZfung/image-board/server/handles"
	"github.com/FXAZfung/image-board/server/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Image Board API
// @version 1.0
// @description API for Image Board application
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

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

	// 公共API路由（无需认证）
	publicGroup := router.Group("/api/public")
	{
		// 登录相关
		publicGroup.POST("/login", handles.Login)

		// 公共设置
		publicGroup.GET("/settings", handles.PublicSettings)

		// 图片相关
		imageGroup := publicGroup.Group("/images")
		{
			imageGroup.GET("/:id", handles.GetImageByID)
			imageGroup.POST("", handles.ListImages)
			imageGroup.POST("/tag", handles.GetImagesByTag)
			imageGroup.GET("/count", handles.GetImageCount)
		}

		// Public tag routes
		tagGroup := publicGroup.Group("/tags")
		{
			tagGroup.POST("", handles.ListTags)
			tagGroup.GET("/:id", handles.GetTagByID)
			tagGroup.GET("/name", handles.GetTagByName)
			tagGroup.GET("/popular", handles.MostPopularTags)
			tagGroup.GET("/search", handles.SearchTags)
			tagGroup.GET("/image/:image_id", handles.GetTagsByImage)
		}
	}

	// 需要认证的API路由
	authGroup := router.Group("/api/auth")
	authGroup.Use(middleware.AuthMiddleware)
	{
		// 登出
		authGroup.GET("/logout", handles.Logout)

		// 图片上传和管理
		authGroup.POST("/upload", handles.UploadImage)

		// In the authenticated routes group
		authGroup.POST("/images/:id/tags", handles.AddTagToImage)

		// 图片操作
		authImageGroup := authGroup.Group("/images")
		{
			authImageGroup.PUT("/:id", handles.UpdateImage)
			authImageGroup.DELETE("/:id", handles.DeleteImage)

			// 标签操作
			authImageGroup.POST("/:id/tags", handles.AddTagsToImage)
			authImageGroup.DELETE("/:id/tags/:tag_id", handles.RemoveTagFromImage)
		}

		// Auth required tag routes
		authTagGroup := authGroup.Group("/tags")
		{
			authTagGroup.POST("", handles.CreateTag)
			authTagGroup.PUT("/:id", handles.UpdateTag)
			authTagGroup.DELETE("/:id", handles.DeleteTag)
		}
	}

	// 私有API路由（需要认证）
	privateGroup := router.Group("/api/private")
	privateGroup.Use(middleware.AuthMiddleware)
	{
		// 设置相关
		privateGroup.GET("/setting", handles.GetSetting)
		privateGroup.POST("/setting", handles.SaveSettings)
		privateGroup.GET("/settings", handles.ListSettings)
		privateGroup.DELETE("/setting", handles.DeleteSetting)
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
