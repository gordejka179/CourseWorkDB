package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gordejka179/CourseWorkDB/internal/usecase"
)

type Handler struct {
	service *usecase.Service
}

func NewUserHandler(service *usecase.Service) *Handler {
    return &Handler{service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.LoadHTMLGlob("internal/web/templates/*")
	router.Static("/static", "./internal/web/static")
	
	auth := router.Group("/auth")
	{
		auth.GET("/registration", h.signUp)
		auth.POST("/registration", h.signUp)
		auth.GET("/login", h.signIn)
		auth.POST("/login", h.signIn)
		auth.GET("/logout", h.signOut)		
	}


	//нужна авторизация
	protected := router.Group("/")
	protected.Use(authMiddleware())
	{
		protected.GET("/home", h.home)
		protected.GET("/search", h.search)
		protected.POST("/searchBook", h.searchBook)
		protected.POST("/reserve", h.reserve)
		protected.GET("/getCurrentBookings", h.getCurrentBookings)
		protected.GET("/makeLoan", h.makeLoan)

		protected.GET("/makeLoan", h.getLoanedBooks)
	}

	return router
}