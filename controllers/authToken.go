package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
)

func generateAuthToken() (string, error) {
	// Buat request untuk login
	loginRequest := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    "aimrzki@gmail.com",
		Password: "admin123",
	}

	// Marshal request ke JSON
	requestBody, err := json.Marshal(loginRequest)
	if err != nil {
		return "", err
	}

	// Buat HTTP request
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Buat konteks Echo
	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Panggil handler LoginUserController untuk mendapatkan token
	if err := LoginUserController(c); err != nil {
		return "", err
	}

	// Periksa status kode HTTP
	if rec.Code != http.StatusOK {
		return "", fmt.Errorf("Login failed with status code %d", rec.Code)
	}

	// Ambil token dari respons
	var response map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		return "", err
	}

	token, ok := response["token"].(string)
	if !ok {
		return "", errors.New("Token not found in response")
	}

	return token, nil
}
