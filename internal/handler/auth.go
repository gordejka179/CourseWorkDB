package handler

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gordejka179/CourseWorkDB/internal/models"
	_ "github.com/gordejka179/CourseWorkDB/internal/usecase"
)

var jwtSecret = "secret-key"

type LoginForm struct {
	Email string `json:"username"`
	Password string `json:"password"`
	Role string `json:"role"`
}

type RegistrationForm struct {
	Email string `json:"email"`
	PassportSeries string `json:"passportseries"`
    PassportNumber string `json:"passportnumber"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Patronymic  string `json:"patronimyc"`
	Password  string `json:"password"`
}


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
        email, ok := claims["email"].(string)
        if !ok {
            c.Redirect(http.StatusSeeOther, "/auth/login")
            c.Abort()
            return
        }
        role, ok := claims["role"].(string)
        if !ok {
            c.Redirect(http.StatusSeeOther, "/auth/login")
            c.Abort()
            return
        }
        c.Set("email", email)
        c.Set("role", role)
        c.Next()
    }
}

// Создание токена
func generateToken(email, role string) (string, error) {
    token := jwt.New(jwt.SigningMethodHS256)
    claims := token.Claims.(jwt.MapClaims)
    claims["email"] = email
    claims["role"] = role
    claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
    tokenString, err := token.SignedString([]byte(jwtSecret))
    if err != nil {
        return "", err
    }
    return tokenString, nil
}

// Логинимся
func (h *Handler) signIn(c *gin.Context) {
    if isAuth(c) {
        c.Redirect(http.StatusSeeOther, "/home")
        return
    }

    switch c.Request.Method {
    case http.MethodGet:
        c.HTML(http.StatusOK, "login.tmpl", gin.H{})

    case http.MethodPost:
        var form LoginForm
        if err := c.BindJSON(&form); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
            return
        }

        var valid bool
        var role string

		hash := md5.Sum([]byte(form.Password))
		passwordHash := hex.EncodeToString(hash[:])

        if form.Role == "reader" {
            valid, _ = h.service.CheckReaderCredentials(form.Email, passwordHash)
            role = "reader"
        } else {
            valid, _ = h.service.CheckLibrarianCredentials(form.Email, passwordHash)
            role = "librarian"
        }

        if !valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверные почта или пароль"})
            return
        }

        token, err := generateToken(form.Email, role)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токена"})
            return
        }

        http.SetCookie(c.Writer, &http.Cookie{
            Name:     "token",
            Path:     "/",
            Value:    token,
            Expires:  time.Now().Add(time.Hour * 24),
            HttpOnly: true,
        })

        c.JSON(http.StatusOK, gin.H{"message": "Успешный вход", "role": role})

    default:
        c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Метод не разрешён"})
    }
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

			exist, err := h.service.CheckIfReaderExists(form.Email)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err})
			}

			if exist {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Пользователь с такой почтой уже есть"})
				return
			}

			hash := md5.Sum([]byte(form.Password))

			passwordHash := hex.EncodeToString(hash[:])


			err = h.service.CreateReader(
				&models.Reader{
					ReaderId:     0,
					Email:        form.Email,
					LibraryCard:  "",
					PassportSeries: form.PassportSeries,
    				PassportNumber: form.PassportNumber,
					FirstName:    form.FirstName,
					LastName:     form.LastName,
					Patronymic:    form.FirstName,
					PasswordHash: passwordHash,
				})



			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err})
			}


			c.Redirect(302, "/auth/login")
		}
	default:
		{
			c.JSON(http.StatusBadRequest, "No such router for this method")
		}
	}
}



