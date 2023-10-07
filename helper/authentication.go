package helper

import (
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

var JwtSecret = []byte("secret-key")

func CreateJWTToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyPassword(storedPassword, providedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(providedPassword))
	return err == nil
}

func ExtractUserIDFromToken(c echo.Context) (float64, error) {
	authorization := c.Request().Header.Get("Authorization")
	if authorization == "" {
		authorization = c.QueryParam("token")
	}

	if authorization == "" {
		return 0, echo.NewHTTPError(http.StatusUnauthorized, "need authorization token")
	}

	token, err := jwt.Parse(authorization, func(token *jwt.Token) (interface{}, error) {
		return JwtSecret, nil
	})
	if err != nil || !token.Valid {
		return 0, echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, echo.NewHTTPError(http.StatusInternalServerError, "Error in token claims")
	}
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, echo.NewHTTPError(http.StatusInternalServerError, "Token does not match user ID")
	}

	return userID, nil
}

func IsAuthorized(c echo.Context, userID float64) bool {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return false
	}

	return int(userID) == id
}
