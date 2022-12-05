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

	api.POST("/auth/forgot-password", handler.ForgotPassword)
	api.POST("/auth/verify-forgot-password", handler.VerifyForgotPassword)
	api.POST("/auth/update-password", handler.AuthMiddleWare, handler.UpdatePassword)

	api.POST("/hotels", handler.AuthMiddleWare, handler.CreateHotel)
	api.GET("/hotels/:id", handler.GetHotel)
	api.PUT("/hotels/:id", handler.AuthMiddleWare, handler.UpdateHotel)
	api.DELETE("/hotels/:id", handler.AuthMiddleWare, handler.DeleteHotel)
	api.GET("/hotels", handler.GetAllHotels)

	api.POST("/rooms", handler.AuthMiddleWare, handler.CreateRoom)
	api.GET("/rooms/:id", handler.GetRoom)
	api.PUT("/rooms/:id", handler.UpdateRoom)
	api.DELETE("/rooms/:id", handler.AuthMiddleWare, handler.DeleteRoom)
	api.GET("/rooms", handler.GetAllRooms)
	api.GET("/rooms/available/:id", handler.GetAllHotelRoomsAvailable)

	api.POST("/users", handler.AuthMiddleWare, handler.CreateUser)
	api.GET("/users/:id", handler.GetUser)
	api.GET("/users/profile", handler.AuthMiddleWare, handler.GetProfileUser)
	api.PUT("/users/:id", handler.AuthMiddleWare, handler.UpdateUser)
	api.DELETE("/users/:id", handler.AuthMiddleWare ,handler.DeleteUser)
	api.GET("/users", handler.GetAllUser)

	api.POST("/bookings", handler.AuthMiddleWare, handler.CreateBooking)
	api.GET("/bookings/:id", handler.GetBooking)
	api.PUT("/bookings/:id", handler.AuthMiddleWare, handler.UpdateBooking)
	api.DELETE("/bookings/:id", handler.AuthMiddleWare ,handler.DeleteBooking)
	api.GET("/bookings", handler.AuthMiddleWare, handler.GetAllBooking)
	api.GET("/bookings/hotel/:id", handler.AuthMiddleWare, handler.GetAllHotelsBooking)


	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}
