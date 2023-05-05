package handler

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"crud-golang/db"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type JWTClaims struct {
	ID int `json:"id"`
	jwt.StandardClaims
}

var users []User

func hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Register(c echo.Context) error {
	var user User
	err := c.Bind(&user)
	if err != nil {
		log.Fatal(err)
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		log.Fatal(err)
	}

	user.Password = string(hashedPassword)

	var count int64
	if err := db.GetDB().Model(&User{}).Where("email = ?", user.Email).Count(&count).Error; err != nil {
		log.Fatal(err)
	}

	if count > 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Email already registered"})
	}

	if err := db.GetDB().Create(&user).Error; err != nil {
		log.Fatal(err)
	}

	return c.JSON(http.StatusOK, user)
}

func Login(c echo.Context) error {
	var user User
	err := c.Bind(&user)
	if err != nil {
		log.Fatal(err)
	}

	result := db.GetDB().Find(&users)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Failed to get users",
		})
	}

	for _, u := range users {
		if u.Email == user.Email && checkPasswordHash(user.Password, u.Password) {
			claims := &JWTClaims{
				ID: u.ID,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
				},
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

			secretKey := []byte(os.Getenv("JWT_SECRET"))
			tokenString, err := token.SignedString(secretKey)
			if err != nil {
				log.Fatal(err)
			}

			response := map[string]string{"token": tokenString}
			return c.JSON(http.StatusOK, response)
		}
	}

	return echo.ErrUnauthorized
}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing authorization token")
		}

		token, err := jwt.ParseWithClaims(tokenString[7:], &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization token")
		}

		claims, ok := token.Claims.(*JWTClaims)
		if !ok || !token.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization token")
		}

		user, err := getUserByID(claims.ID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get user data")
		}

		c.Set("user", user)
		return next(c)
	}
}

func GetUsers(c echo.Context) error {
	var users []User
	err := db.GetDB().Find(&users).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get users"})
	}

	// Hapus field password dari setiap user
	for i := range users {
		users[i].Password = ""
	}

	return c.JSON(http.StatusOK, users)
}

func GetUserProfile(c echo.Context) error {
	user := c.Get("user").(*User)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}

func getUserByID(id int) (*User, error) {
	user := new(User)
	if err := db.GetDB().First(user, id).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func Restricted(c echo.Context) error {
	response := map[string]string{
		"message": "Hello",
	}
	return c.JSON(http.StatusOK, response)
}
