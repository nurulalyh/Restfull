package controllers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserController(t *testing.T) {
	e := echo.New()

	reqBody := `{
		"name": "Lee Mujin",
		"email": "morilla_lmj@gmail.com",
		"password": "28122000"
	}`
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, CreateUserController(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestLoginUserController(t *testing.T) {
	e := echo.New()
	reqBody := `{
		"email": "morilla_lmj@gmail.com",
		"password": "28122000"
	}`
	req := httptest.NewRequest(http.MethodPost, "/users/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, LoginUserController(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestGetUsersController(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	token, err := generateAuthToken()
	if err != nil {
		t.Fatalf("Failed to generate auth token: %v", err)
	}

	// Kasus 1: Token valid
	req.Header.Set("Authorization", "Bearer "+token)
	if assert.NoError(t, GetUsersController(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

	}

	// Kasus 2: Token tidak ada (Unauthorized)
	req.Header.Del("Authorization")
	if assert.NoError(t, GetUsersController(c)) {
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	}

	// Kasus 3: Token tidak valid (Unauthorized)
	req.Header.Set("Authorization", "InvalidTokenFormat")
	if assert.NoError(t, GetUsersController(c)) {
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	}

	// Kasus 4: Token tidak valid (Unauthorized)
	req.Header.Set("Authorization", "Bearer InvalidToken")
	if assert.NoError(t, GetUsersController(c)) {
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	}
}

func TestGetUserController(t *testing.T) {
	e := echo.New()

	// Kasus 1: Pengguna yang sah mendapatkan data pengguna mereka sendiri
	t.Run("ValidUserOwnData", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		token, err := generateAuthToken()
		if err != nil {
			t.Fatalf("Failed to generate auth token: %v", err)
		}

		req.Header.Set("Authorization", "Bearer "+token)

		if assert.NoError(t, GetUserController(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	// Kasus 2: Pengguna yang sah mencoba mendapatkan data pengguna lain (Forbidden)
	t.Run("ValidUserOtherUserData", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/2", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("2")

		token, err := generateAuthToken()
		if err != nil {
			t.Fatalf("Failed to generate auth token: %v", err)
		}

		req.Header.Set("Authorization", "Bearer "+token)

		if assert.NoError(t, GetUserController(c)) {
			assert.Equal(t, http.StatusForbidden, rec.Code)
		}
	})

	// Kasus 3: Token tidak valid (Unauthorized)
	t.Run("InvalidToken", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		req.Header.Set("Authorization", "InvalidToken")

		if assert.NoError(t, GetUserController(c)) {
			assert.Equal(t, http.StatusUnauthorized, rec.Code)
		}
	})

	// Kasus 4: Pengguna tidak memiliki akses token (Unauthorized)
	t.Run("NoToken", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		if assert.NoError(t, GetUserController(c)) {
			assert.Equal(t, http.StatusUnauthorized, rec.Code)
		}
	})

	// Kasus 5: ID yang salah (Bad Request)
	t.Run("InvalidID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/invalid_id", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("invalid_id")

		token, err := generateAuthToken()
		if err != nil {
			t.Fatalf("Failed to generate auth token: %v", err)
		}

		req.Header.Set("Authorization", "Bearer "+token)

		if assert.NoError(t, GetUserController(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})
}

func TestUpdateUserController(t *testing.T) {
	e := echo.New()

	// Kasus 1: Pengguna yang sah berhasil memperbarui data pengguna mereka sendiri
	t.Run("ValidUserUpdateOwnData", func(t *testing.T) {
		reqBody := `{
			"name": "Lee Mujin",
			"email": "morilla_lmj@gmail.com",
			"password": "vocal20"
		}`

		req := httptest.NewRequest(http.MethodPut, "/users/1", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		token, err := generateAuthToken()
		if err != nil {
			t.Fatalf("Failed to generate auth token: %v", err)
		}

		req.Header.Set("Authorization", "Bearer "+token)

		if assert.NoError(t, UpdateUserController(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	// Kasus 2: Pengguna yang sah mencoba memperbarui data pengguna lain (Forbidden)
	t.Run("ValidUserUpdateOtherUserData", func(t *testing.T) {
		reqBody := `{
			"name": "Lee Mujin",
			"email": "morilla_lmj@gmail.com",
			"password": "mujinservice"
		}`

		req := httptest.NewRequest(http.MethodPut, "/users/2", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("2")

		token, err := generateAuthToken()
		if err != nil {
			t.Fatalf("Failed to generate auth token: %v", err)
		}

		req.Header.Set("Authorization", "Bearer "+token)

		// Assertions
		if assert.NoError(t, UpdateUserController(c)) {
			assert.Equal(t, http.StatusForbidden, rec.Code)
		}
	})

	// Kasus 3: Token tidak valid (Unauthorized)
	t.Run("InvalidToken", func(t *testing.T) {
		reqBody := `{
			"name": "Lee Mujin",
			"email": "morilla_lmj@gmail.com",
			"password": "mujiniee"
		}`

		req := httptest.NewRequest(http.MethodPut, "/users/1", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		req.Header.Set("Authorization", "InvalidToken")

		if assert.NoError(t, UpdateUserController(c)) {
			assert.Equal(t, http.StatusUnauthorized, rec.Code)
		}
	})

	// Kasus 4: Pengguna tidak memiliki akses token (Unauthorized)
	t.Run("NoToken", func(t *testing.T) {
		reqBody := `{
			"name": "Lee Mujin",
			"email": "morilla_lmj@gmail.com",
			"password": "2mujin"
		}`

		req := httptest.NewRequest(http.MethodPut, "/users/1", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		if assert.NoError(t, UpdateUserController(c)) {
			assert.Equal(t, http.StatusUnauthorized, rec.Code)
		}
	})

	// Kasus 5: ID yang salah (Bad Request)
	t.Run("InvalidID", func(t *testing.T) {
		reqBody := `{
			"name": "Lee Mujin",
			"email": "morilla_lmj@gmail.com",
			"password": "2000"
		}`

		req := httptest.NewRequest(http.MethodPut, "/users/invalid_id", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("invalid_id")

		token, err := generateAuthToken()
		if err != nil {
			t.Fatalf("Failed to generate auth token: %v", err)
		}

		req.Header.Set("Authorization", "Bearer "+token)

		if assert.NoError(t, UpdateUserController(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})
}

func TestDeleteUserController(t *testing.T) {
	e := echo.New()

	// Kasus 1: Pengguna yang sah berhasil menghapus akun
	t.Run("ValidUserDeleteUser", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		token, err := generateAuthToken()
		if err != nil {
			t.Fatalf("Failed to generate auth token: %v", err)
		}

		req.Header.Set("Authorization", "Bearer "+token)

		if assert.NoError(t, DeleteUserController(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	// Kasus 2: Token tidak valid (Unauthorized)
	t.Run("InvalidToken", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		req.Header.Set("Authorization", "InvalidToken")

		if assert.NoError(t, DeleteUserController(c)) {
			assert.Equal(t, http.StatusUnauthorized, rec.Code)
		}
	})

	// Kasus 3: Pengguna tidak memiliki akses token (Unauthorized)
	t.Run("NoToken", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		if assert.NoError(t, DeleteUserController(c)) {
			assert.Equal(t, http.StatusUnauthorized, rec.Code)
		}
	})

	// Kasus 4: ID pengguna yang salah (Not Found)
	t.Run("InvalidUserID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/users/999", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("999")

		token, err := generateAuthToken()
		if err != nil {
			t.Fatalf("Failed to generate auth token: %v", err)
		}

		req.Header.Set("Authorization", "Bearer "+token)

		if assert.NoError(t, DeleteUserController(c)) {
			assert.Equal(t, http.StatusNotFound, rec.Code)
		}
	})
}