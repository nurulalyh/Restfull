package service

import (
	"errors"
	"test/features/users"
	"test/features/users/mocks"
	helper "test/helper/mocks"

	"github.com/stretchr/testify/mock"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	generator := helper.NewGeneratorInterface(t)
	jwt := helper.NewJWTInterface(t)
	data := mocks.NewUserDataInterface(t)
	service := New(data, generator, jwt)
	newUser := users.User{
		Nama:     "dida",
		HP:       "123",
		Password: "didadejan123",
	}

	t.Run("Success insert", func(t *testing.T) {
		generator.On("GenerateUUID").Return("randomUUID", nil).Once()
		newUser.ID = "randomUUID"
		data.On("Insert", newUser).Return(&newUser, nil).Once()

		result, err := service.Register(newUser)
		assert.Nil(t, err)
		assert.Equal(t, newUser.ID, result.ID)
		assert.Equal(t, newUser.Nama, result.Nama)
		generator.AssertExpectations(t)
		data.AssertExpectations(t)
	})

	t.Run("Generate failed", func(t *testing.T) {
		generator.On("GenerateUUID").Return("", errors.New("some error on generator")).Once()

		result, err := service.Register(newUser)
		assert.Error(t, err)
		assert.EqualError(t, err, "id generator failed")
		assert.Nil(t, result)
		generator.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	generator := helper.NewGeneratorInterface(t)
	j := helper.NewJWTInterface(t)
	data := mocks.NewUserDataInterface(t)
	service := New(data, generator, j)
	userData := users.User{
		ID:       "randomUserID",
		Nama:     "dida",
		HP:       "123",
		Password: "didadejan123",
	}

	t.Run("success login", func(t *testing.T) {
		jwtResult := map[string]any{"access_token": "randomAccessToken"}
		data.On("Login", userData.HP, userData.Password).Return(&userData, nil)
		j.On("GenerateJWT", mock.Anything).Return(jwtResult)
		result, err := service.Login(userData.HP, userData.Password)

		data.AssertExpectations(t)
		j.AssertExpectations(t)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "dida", result.Nama)
		assert.Equal(t, jwtResult, result.Access)
	})
}
