package server

import (
	conf "github.com/FXAZfung/image-board/internal/config"
	"github.com/FXAZfung/image-board/server/common"
	"github.com/FXAZfung/image-board/server/handles"
	"github.com/FXAZfung/image-board/server/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Init(router *gin.Engine) {
	// 处理跨域
	Cors(router)

	common.SecretKey = []byte(conf.Conf.JwtSecret)

	api := router.Group("/api")
	// 不需要登录的接口
	public := api.Group("/public")
	public.POST("/login", handles.Login)
	public.GET("/images/:name", handles.GetImageByName)
	public.GET("/images", handles.ListImages)
	public.GET("/short/:short_link", handles.GetImageByShortLink)
	public.GET("/categories", handles.GetCategories)
	public.GET("/categories/:name", handles.GetCategoryByName)
	public.GET("/random", handles.GetRandomImage)
	public.GET("/settings", handles.PublicSettings)
	public.GET("/info", handles.GetInfo)

	// 需要登录的接口
	auth := api.Group("/auth")
	auth.Use(middleware.AuthMiddleware)
	auth.POST("/upload", handles.UploadImage)
	auth.POST("/categories", handles.CreateCategory)
	auth.GET("/logout", handles.Logout)

	// 私有接口
	private := api.Group("/private")
	private.Use(middleware.AuthMiddleware)
	private.GET("/setting", handles.GetSetting)
	private.POST("/setting", handles.SaveSettings)
	private.GET("/settings", handles.ListSettings)
	private.DELETE("/setting", handles.DeleteSetting)
	private.POST("/setting/token", handles.ResetToken)
	private.POST("/users", handles.ListUser)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

}

func Cors(r *gin.Engine) {
	config := cors.DefaultConfig()
	config.AllowOrigins = conf.Conf.Cors.AllowOrigins
	config.AllowHeaders = conf.Conf.Cors.AllowHeaders
	config.AllowMethods = conf.Conf.Cors.AllowMethods
	r.Use(cors.New(config))
}
