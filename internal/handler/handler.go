package handler

import (
	"github.com/gin-gonic/gin"
)

type Handler struct {
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.LoadHTMLGlob("web/templates/*")
	router.Static("/images", "web/static/images")
	router.Static("/styles", "web/static/styles")
	router.Static("/scripts", "web/static/scripts")

	auth := router.Group("/auth")
	{
		auth.GET("/login", h.signIn)
		auth.POST("/login", h.signIn)
		auth.GET("/registration", h.signUp)
		auth.POST("/registration", h.signUp)
	}

	view := router.Group("/view")
	view.Use(authMiddleware())

	return router
}