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
	router.Static("/images", "web/static/images")
	router.Static("/styles", "web/static/styles")
	router.Static("/scripts", "web/static/scripts")

	auth := router.Group("/auth")
	{
		auth.GET("/registration", h.signUp)
		auth.POST("/registration", h.signUp)
		auth.GET("/login", h.signIn)
		auth.POST("/login", h.signIn)
	}


	router.GET("/home", h.home)
	router.GET("/search", h.search)
	

	view := router.Group("/view")
	view.Use(authMiddleware())

	return router
}