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

func TestHandler_FindUser(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	execute := func(useCase *mocksUsers.MockUseCase, query string) *httptest.ResponseRecorder {
		res := httptest.NewRecorder()
		engine := gin.New()
		engine.Use(middlewares.ErrorHandler())
		users.New(engine, useCase, mocksAuth.NewMockAuth())
		req := httptest.NewRequest(http.MethodGet, "/api/users"+query, nil)
		engine.ServeHTTP(res, req)
		return res
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		useCase := mocksUsers.NewMockUseCase(ctrl)
		useCase.EXPECT().FindUser(gomock.Any(), gomock.Any()).Return([]entities.User{{ID: "id-1", Name: "Tester", Email: "tester@example.com"}}, nil)

		res := execute(useCase, "?name=Tester&email=tester%40example.com")
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Contains(t, res.Body.String(), `"email":"tester@example.com"`)
		assert.Contains(t, res.Body.String(), `"name":"Tester"`)
	})

	t.Run("invalid query", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		useCase := mocksUsers.NewMockUseCase(ctrl)

		res := execute(useCase, "?email=invalid")
		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	t.Run("use case error", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		useCase := mocksUsers.NewMockUseCase(ctrl)
		useCase.EXPECT().FindUser(gomock.Any(), gomock.Any()).Return(nil, errors.New("find failed"))

		res := execute(useCase, "")
		assert.Equal(t, http.StatusInternalServerError, res.Code)
	})
}
