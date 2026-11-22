package users_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chaihanij/backend-challenge/app/entities"
	"github.com/chaihanij/backend-challenge/app/layers/deliveries/http/api/users"
	"github.com/chaihanij/backend-challenge/app/layers/deliveries/http/middlewares"
	_mocksAuth "github.com/chaihanij/backend-challenge/app/mocks/auth"
	_mocksUsersUseCase "github.com/chaihanij/backend-challenge/app/mocks/users"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestHandler_CreateUser(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	executeWithRequest := func(useCase *_mocksUsersUseCase.MockUseCase, body []byte) *httptest.ResponseRecorder {
		res := httptest.NewRecorder()
		engine := gin.New()
		engine.Use(middlewares.ErrorHandler())
		users.New(engine, useCase, _mocksAuth.NewMockAuth())

		req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		engine.ServeHTTP(res, req)
		return res
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		useCase := _mocksUsersUseCase.NewMockUseCase(ctrl)
		input := &entities.User{Email: "email@test.com", Name: "Example User", Password: "P@ssw0rd"}
		created := &entities.User{ID: "user-123", Email: input.Email, Name: input.Name}
		useCase.EXPECT().CreateUser(gomock.Any(), input).Return(created, nil)

		body := []byte(`{"email":"email@test.com","name":"Example User","password":"P@ssw0rd"}`)
		res := executeWithRequest(useCase, body)

		assert.Equal(t, http.StatusOK, res.Code)
		var response struct {
			Data struct {
				ID    string `json:"id"`
				Email string `json:"email"`
				Name  string `json:"name"`
			} `json:"data"`
		}
		require.NoError(t, json.Unmarshal(res.Body.Bytes(), &response))
		assert.Equal(t, "user-123", response.Data.ID)
		assert.Equal(t, input.Email, response.Data.Email)
		assert.Equal(t, input.Name, response.Data.Name)
	})

	t.Run("invalid request", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		useCase := _mocksUsersUseCase.NewMockUseCase(ctrl)

		res := executeWithRequest(useCase, []byte(`{"email":"not-an-email"}`))

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	t.Run("use case error", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		useCase := _mocksUsersUseCase.NewMockUseCase(ctrl)
		input := &entities.User{Email: "email@test.com", Name: "Example User", Password: "P@ssw0rd"}
		useCase.EXPECT().CreateUser(gomock.Any(), input).Return(nil, errors.New("create failed"))

		body := []byte(`{"email":"email@test.com","name":"Example User","password":"P@ssw0rd"}`)
		res := executeWithRequest(useCase, body)

		assert.Equal(t, http.StatusInternalServerError, res.Code)
	})
}
