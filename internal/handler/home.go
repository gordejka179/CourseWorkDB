package handler

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func (h *Handler) home(c *gin.Context) {
    cookie, _ := c.Cookie("token")

    token, _ := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
        return []byte(jwtSecret), nil
    })

    claims, _ := token.Claims.(jwt.MapClaims)

    role, _ := claims["role"].(string)

	switch role {
    case "reader":
        c.HTML(http.StatusOK, "homeReader.tmpl", gin.H{
            "role": role,
        })
    case "librarian":
        c.HTML(http.StatusOK, "homeLibrarian.tmpl", gin.H{
            "role": role,
        })
	}
}