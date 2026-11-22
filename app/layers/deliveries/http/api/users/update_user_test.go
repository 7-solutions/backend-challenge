package users_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chaihanij/backend-challenge/app/entities"
	"github.com/chaihanij/backend-challenge/app/layers/deliveries/http/api/users"
	"github.com/chaihanij/backend-challenge/app/layers/deliveries/http/middlewares"
	mocksAuth "github.com/chaihanij/backend-challenge/app/mocks/auth"
	mocksUsers "github.com/chaihanij/backend-challenge/app/mocks/users"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandler_UpdateUser(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	const userID = "507f1f77bcf86cd799439011"
	const body = `{"email":"tester@example.com","name":"Tester"}`
	execute := func(useCase *mocksUsers.MockUseCase, id, requestBody string) *httptest.ResponseRecorder {
		res := httptest.NewRecorder()
		engine := gin.New()
		engine.Use(middlewares.ErrorHandler())
		users.New(engine, useCase, mocksAuth.NewMockAuth())
		req := httptest.NewRequest(http.MethodPut, "/api/users/"+id, bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")
		engine.ServeHTTP(res, req)
		return res
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		useCase := mocksUsers.NewMockUseCase(ctrl)
		useCase.EXPECT().UpdateUser(gomock.Any(), &entities.User{ID: userID, Email: "tester@example.com", Name: "Tester"}).Return(&entities.User{ID: userID, Email: "tester@example.com", Name: "Tester"}, nil)

		res := execute(useCase, userID, body)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Contains(t, res.Body.String(), `"id":"`+userID+`"`)
		assert.Contains(t, res.Body.String(), `"name":"Tester"`)
	})

	t.Run("invalid user ID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		useCase := mocksUsers.NewMockUseCase(ctrl)

		res := execute(useCase, "invalid", body)
		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	t.Run("invalid request body", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		useCase := mocksUsers.NewMockUseCase(ctrl)

		res := execute(useCase, userID, `{"email":"invalid","name":"Tester"}`)
		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	t.Run("use case error", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		useCase := mocksUsers.NewMockUseCase(ctrl)
		useCase.EXPECT().UpdateUser(gomock.Any(), &entities.User{ID: userID, Email: "tester@example.com", Name: "Tester"}).Return(nil, errors.New("update failed"))

		res := execute(useCase, userID, body)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
	})
}
