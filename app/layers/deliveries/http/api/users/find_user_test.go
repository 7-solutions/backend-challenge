package users_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/chaihanij/backend-challenge/app/entities"
	"github.com/chaihanij/backend-challenge/app/layers/deliveries/http/api/users"
	mocsksAuth "github.com/chaihanij/backend-challenge/app/mocks/auth"
	mocksUsers "github.com/chaihanij/backend-challenge/app/mocks/users"
	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
)

func TestHandler_FindUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	executeWithRequest := func(mockUsersUseCase *mocksUsers.MockUseCase, query string) *httptest.ResponseRecorder {
		response := httptest.NewRecorder()
		_, ginEngine := gin.CreateTestContext(response)
		requestUrl := "/users"
		if query != "" {
			requestUrl += "?" + query
		}
		httpRequest, _ := http.NewRequest(http.MethodGet, requestUrl, nil)

		mockAuthMiddleware := mocsksAuth.NewMockAuth()

		users.New(ginEngine, mockUsersUseCase, mockAuthMiddleware)

		ginEngine.ServeHTTP(response, httpRequest)
		return response
	}

	t.Run("Success", func(t *testing.T) {
		query := url.Values{}
		query.Add("name", "Tester")
		query.Add("email", "tester@test.com")

		mockUsersUseCase := mocksUsers.NewMockUseCase(ctrl)

		mockUsersUseCase.EXPECT().FindUser(gomock.Any(), gomock.Any()).Return([]entities.User{{
			ID:    "id",
			Name:  "Tester",
			Email: "tester@test.com",
		}}, nil).Times(1)

		res := executeWithRequest(mockUsersUseCase, query.Encode())

		if res.Code != http.StatusOK {
			t.Fatalf("expected status code 200, got %d", res.Code)
		}
	})

	t.Run("Failure: Bad Request", func(t *testing.T) {

	})
}
