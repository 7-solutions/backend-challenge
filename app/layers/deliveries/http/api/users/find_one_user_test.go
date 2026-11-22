package users_test

import (
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

func TestHandler_FindOneUser(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	const userID = "507f1f77bcf86cd799439011"
	execute := func(useCase *mocksUsers.MockUseCase, id string) *httptest.ResponseRecorder {
		res := httptest.NewRecorder()
		engine := gin.New()
		engine.Use(middlewares.ErrorHandler())
		users.New(engine, useCase, mocksAuth.NewMockAuth())
		engine.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/api/users/"+id, nil))
		return res
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		useCase := mocksUsers.NewMockUseCase(ctrl)
		useCase.EXPECT().FindOneUser(gomock.Any(), userID).Return(&entities.User{ID: userID, Name: "Tester", Email: "tester@example.com"}, nil)

		res := execute(useCase, userID)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Contains(t, res.Body.String(), `"id":"`+userID+`"`)
		assert.Contains(t, res.Body.String(), `"email":"tester@example.com"`)
	})

	t.Run("invalid user ID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		useCase := mocksUsers.NewMockUseCase(ctrl)

		res := execute(useCase, "invalid")
		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	t.Run("use case error", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		useCase := mocksUsers.NewMockUseCase(ctrl)
		useCase.EXPECT().FindOneUser(gomock.Any(), userID).Return(nil, errors.New("find failed"))

		res := execute(useCase, userID)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
	})
}
