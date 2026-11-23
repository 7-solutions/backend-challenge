package users_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chaihanij/backend-challenge/app/layers/deliveries/http/api/users"
	"github.com/chaihanij/backend-challenge/app/layers/deliveries/http/middlewares"
	_mocksAuth "github.com/chaihanij/backend-challenge/app/mocks/auth"
	_mocksUsersUseCase "github.com/chaihanij/backend-challenge/app/mocks/users"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandler_DeleteUser(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	const userID = "507f1f77bcf86cd799439011"
	executeWithRequest := func(useCase *_mocksUsersUseCase.MockUseCase, id string) *httptest.ResponseRecorder {
		res := httptest.NewRecorder()
		engine := gin.New()
		engine.Use(middlewares.ErrorHandler())
		users.New(engine, useCase, _mocksAuth.NewMockAuth())
		engine.ServeHTTP(res, httptest.NewRequest(http.MethodDelete, "/api/users/"+id, nil))
		return res
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		useCase := _mocksUsersUseCase.NewMockUseCase(ctrl)
		useCase.EXPECT().DeleteUser(gomock.Any(), userID).Return(nil)

		res := executeWithRequest(useCase, userID)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.Contains(t, res.Body.String(), `"code":200`)
		assert.Contains(t, res.Body.String(), `"message":"OK"`)
	})

	t.Run("invalid user ID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		useCase := _mocksUsersUseCase.NewMockUseCase(ctrl)

		res := executeWithRequest(useCase, "not-a-hex-id")

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	t.Run("use case error", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		useCase := _mocksUsersUseCase.NewMockUseCase(ctrl)
		useCase.EXPECT().DeleteUser(gomock.Any(), userID).Return(errors.New("delete failed"))

		res := executeWithRequest(useCase, userID)

		assert.Equal(t, http.StatusInternalServerError, res.Code)
	})
}
