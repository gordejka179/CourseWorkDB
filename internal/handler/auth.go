package handler

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gordejka179/CourseWorkDB/internal/repository"
	"github.com/gordejka179/CourseWorkDB/pkg"
)

const (
	clientID       = ""
	clientSecret   = ""
	redirectURL    = "http://localhost:8080/auth/callback"
	sessionName    = "mysession"
	accessTokenKey = "access_token"
)

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegistrationForm struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	LastName  string `json:"lastname"`
	FirstName string `json:"firstname"`
	Email     string `json:"email"`
	Address   string `json:"address"`
}

var (
	jwtSecret = "secret-key"
)


// проверка, залогинился ли уже пользователь
func isAuth(c *gin.Context) bool {
	cookie, err := c.Cookie("token")
	if err != nil {
		return false
	}

	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return false
	}
	_, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	return true
}

// Функция для middleware
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("token")
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/auth/login")
			c.Abort()
			return
		}

		token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			c.Redirect(http.StatusSeeOther, "/auth/login")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.Redirect(http.StatusSeeOther, "/auth/login")
			c.Abort()
			return
		}

		username, ok := claims["username"].(string)
		if !ok {
			c.Redirect(http.StatusSeeOther, "/auth/login")
			c.Abort()
			return
		}

		role, err := repository.GetUserRole(username)
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/auth/login")
			c.Abort()
			return
		}

		c.Set("username", username)
		c.Set("role", role)

		c.Next()
	}
}

// Логируемся
func (h *Handler) signIn(c *gin.Context) {
	// проверить авторизован ли уже пользователь
	if isAuth(c) {
		c.Redirect(http.StatusSeeOther, "/home")
		c.Abort()
	}

	requestMethod := c.Request.Method
	switch requestMethod {
	case "GET":
		{
			c.HTML(http.StatusOK, "login.tmpl", gin.H{})
		}
	case "POST":
		{
			var form LoginForm
			if err := c.BindJSON(&form); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
				return
			}

			valid, _ := checkLoginAndPassword(form.Username, form.Password)
			if valid {
				token, err := pkg.GenerateToken(form.Username)

				if err != nil {
					c.JSON(http.StatusUnauthorized, "Ошибка генерации токена jwt")
				}

				cookie := http.Cookie{
					Name:     "token",
					Path:     "/",
					Value:    token,
					Expires:  time.Now().Add(time.Hour * 24),
					HttpOnly: true,
				}
				http.SetCookie(c.Writer, &cookie)

				c.JSON(http.StatusOK, gin.H{"message": "Успешный вход"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверные логин или пароль"})
			}
		}
	default:
		{
			c.JSON(http.StatusBadRequest, "No such router for this method")
		}
	}
}

func checkLoginAndPassword(username, password string) (bool, string) {
	str := password
	hash := md5.Sum([]byte(str))

	hashString := hex.EncodeToString(hash[:])

	exist, err := repository.AuthenticateUser(username, hashString)

	if err != nil {
		fmt.Println("error")
		return false, ""
	}

	if !exist {
		return false, ""
	}

	role, err := repository.GetUserRole(username)

	if err != nil {
		fmt.Println("error")
		return false, ""
	}

	return true, role
}


// Регистрация
func (h *Handler) signUp(c *gin.Context) {
	requestMethod := c.Request.Method
	switch requestMethod {
	case "GET":
		{
			c.HTML(http.StatusOK, "registration.tmpl", gin.H{})
		}
	case "POST":
		{
			var form RegistrationForm
			if err := c.BindJSON(&form); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
				return
			}

			exist, err := repository.CheckIfUserExists(form.Username)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err})
			}

			if exist {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Пользователь с таким никнеймом уже есть"})
				return
			}

			hash := md5.Sum([]byte(form.Password))

			hashPassword := hex.EncodeToString(hash[:])

			repository.CreateUser(repository.User{
				Username:     form.Username,
				PasswordHash: hashPassword,
				LastName:     form.LastName,
				FirstName:    form.FirstName,
				Email:        form.Email,
			})

			c.Redirect(302, "/auth/login")
		}
	default:
		{
			c.JSON(http.StatusBadRequest, "No such router for this method")
		}
	}
}

