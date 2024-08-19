package router

import (
	"cors/internal/http/handler"
	"cors/internal/http/middleware"
	userusecase "cors/internal/usecase/user"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(user *userusecase.User) *gin.Engine {
	userHandler := handler.NewUser(user)
	router := gin.New()

	router.Use(middleware.TimingMiddleware)
	router.Use(middleware.Middleware)

	url := ginSwagger.URL("swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	userGroup := router.Group("/user")
	{
		userGroup.POST("/register", userHandler.Register)
		userGroup.POST("/login", userHandler.Login)
		userGroup.POST("/verify", userHandler.VerifyUser)
	}
	return router
}
