package route

import (
	"github/SmoothieNoIce/dcard-backend-intern-assignment-2024/app/http/controllers"
	"github/SmoothieNoIce/dcard-backend-intern-assignment-2024/config"

	docs "github/SmoothieNoIce/dcard-backend-intern-assignment-2024/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup() *gin.Engine {
	var r *gin.Engine

	// Set Gin Mode
	if config.AppConfig.AppEnv == "local" {
		r = gin.Default()
		gin.SetMode(gin.DebugMode)
	} else if config.AppConfig.AppEnv == "production" {
		r = gin.New()
		r.Use(gin.Recovery())
		gin.SetMode(gin.ReleaseMode)
	} else {
		r = gin.Default()
		gin.SetMode(gin.DebugMode)
	}

	// Set Trusted
	r.SetTrustedProxies([]string{"discordservers.tw", "client:3000"})

	// Set a lower memory limit for multipart forms (default is 32 MiB)
	r.MaxMultipartMemory = 10 << 20 // 8 MiB

	apiv1 := r.Group("/api/v1")

	// === swagger ===
	docs.SwaggerInfo.Host = config.AppConfig.Host
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// === router ===
	apiv1.Use()
	{
		guest := apiv1.Group("/ad")
		guest.Use()
		{
			guest.GET("", controllers.GetAdvertistmentList)
			guest.POST("", controllers.AddAdvertistment)
		}
	}
	return r
}
