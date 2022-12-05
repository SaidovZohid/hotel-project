package api

import (
	v1 "github.com/SaidovZohid/hotel-project/api/v1"
	"github.com/SaidovZohid/hotel-project/config"
	"github.com/SaidovZohid/hotel-project/storage"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type RouteOptions struct {
	Cfg      *config.Config
	Stgr     storage.StorageI
	InMemory storage.InMemoryStorageI
}

// New @title           Swagger for hotel api
// @version         2.0
// @description     This is a hotel service api.
// @host      		localhost:8080
// @BasePath  		/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @Security ApiKeyAuth
func New(opt *RouteOptions) *gin.Engine {
	router := gin.Default()

	handler := v1.New(&v1.HandlerV1{
		Cfg:      opt.Cfg,
		Strg:     opt.Stgr,
		InMemory: opt.InMemory,
	})
	api := router.Group("/v1")
	api.POST("/auth/register", handler.Register)
	api.POST("auth/verify", handler.Verify)
	api.POST("/auth/login", handler.Login)

	api.POST("/hotels", handler.AuthMiddleWare, handler.CreateHotel)
	api.GET("/hotels/:id", handler.GetHotel)
	api.PUT("/hotels/:id", handler.AuthMiddleWare, handler.UpdateHotel)
	api.DELETE("/hotels/:id", handler.AuthMiddleWare, handler.DeleteHotel)
	api.GET("/hotels", handler.GetAllHotels)

	api.POST("/rooms", handler.AuthMiddleWare, handler.CreateRoom)
	// api.GET("/hotels/:id", handler.GetRoom)
	// api.PUT("/hotels/:id", handler.UpdateRoom)
	// api.DELETE("/hotels/:id", handler.DeleteRoom)
	// api.GET("/hotels", handler.GetAllRooms)


	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}
