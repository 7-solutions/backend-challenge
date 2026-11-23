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
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestHandler_FindUser(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	execute := func(t *testing.T, useCase *mocksUsers.MockUseCase, query string) *httptest.ResponseRecorder {
		t.Helper()
		engine := gin.New()
		engine.Use(middlewares.ErrorHandler())
		users.New(engine, useCase, mocksAuth.NewMockAuth())

		res := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/api/users"+query, nil)
		engine.ServeHTTP(res, request)
		return res
	}

	t.Run("returns users matching query filters", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		useCase := mocksUsers.NewMockUseCase(ctrl)
		useCase.EXPECT().FindUser(gomock.Any(), gomock.Cond(func(filter *entities.UserFilter) bool {
			return filter.Name != nil && *filter.Name == "Tester" &&
				filter.Email != nil && *filter.Email == "tester@example.com" && filter.ID == nil
		})).Return([]entities.User{{ID: "id-1", Name: "Tester", Email: "tester@example.com"}}, nil)

		res := execute(t, useCase, "?name=Tester&email=tester%40example.com")

		require.Equal(t, http.StatusOK, res.Code, res.Body.String())
		assert.Contains(t, res.Body.String(), `"id":"id-1"`)
		assert.Contains(t, res.Body.String(), `"email":"tester@example.com"`)
		assert.Contains(t, res.Body.String(), `"name":"Tester"`)
	})

	t.Run("allows searching without filters", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		useCase := mocksUsers.NewMockUseCase(ctrl)
		useCase.EXPECT().FindUser(gomock.Any(), gomock.Cond(func(filter *entities.UserFilter) bool {
			return filter.Name == nil && filter.Email == nil && filter.ID == nil
		})).Return([]entities.User{}, nil)

		res := execute(t, useCase, "")

		require.Equal(t, http.StatusOK, res.Code, res.Body.String())
		assert.Contains(t, res.Body.String(), `"data":null`)
	})

	t.Run("rejects invalid email without calling use case", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		useCase := mocksUsers.NewMockUseCase(ctrl)
		useCase.EXPECT().FindUser(gomock.Any(), gomock.Any()).Times(0)

		res := execute(t, useCase, "?email=invalid")

		require.Equal(t, http.StatusBadRequest, res.Code, res.Body.String())
		assert.NotEmpty(t, res.Body.String())
	})

	t.Run("rejects too short name without calling use case", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		useCase := mocksUsers.NewMockUseCase(ctrl)
		useCase.EXPECT().FindUser(gomock.Any(), gomock.Any()).Times(0)

		res := execute(t, useCase, "?name=ab")

		require.Equal(t, http.StatusBadRequest, res.Code, res.Body.String())
		assert.NotEmpty(t, res.Body.String())
	})

	t.Run("returns internal server error when use case fails", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		useCase := mocksUsers.NewMockUseCase(ctrl)
		useCase.EXPECT().FindUser(gomock.Any(), gomock.Any()).Return(nil, errors.New("find failed"))

		res := execute(t, useCase, "")

		require.Equal(t, http.StatusInternalServerError, res.Code, res.Body.String())
		assert.NotEmpty(t, res.Body.String())
	})
}
