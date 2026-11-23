package auth_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chaihanij/backend-challenge/app/entities"
	authapi "github.com/chaihanij/backend-challenge/app/layers/deliveries/http/api/auth"
	"github.com/chaihanij/backend-challenge/app/layers/deliveries/http/middlewares"
	"github.com/chaihanij/backend-challenge/app/layers/deliveries/http/response"
	mocksAuthUseCase "github.com/chaihanij/backend-challenge/app/mocks/auth"
	mocksUsersUseCase "github.com/chaihanij/backend-challenge/app/mocks/users"
	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
)

func TestHandler_Register(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	newRequest := func(t *testing.T, mockUsersUseCase *mocksUsersUseCase.MockUseCase, body string) *httptest.ResponseRecorder {
		t.Helper()
		engine := gin.New()
		engine.Use(middlewares.ErrorHandler())
		authapi.New(engine, mocksAuthUseCase.NewMockUseCase(gomock.NewController(t)), mockUsersUseCase)
		request := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBufferString(body))
		request.Header.Set("Content-Type", "application/json")
		result := httptest.NewRecorder()
		engine.ServeHTTP(result, request)
		return result
	}

	t.Run("creates user and omits password from response", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockUsersUseCase := mocksUsersUseCase.NewMockUseCase(ctrl)
		input := &entities.User{Email: "jane@example.com", Name: "Jane Doe", Password: "secret123"}
		created := &entities.User{ID: "65f1a4d76a9c0d1234567890", Email: input.Email, Name: input.Name}
		mockUsersUseCase.EXPECT().CreateUser(gomock.Any(), input).Return(created, nil)

		result := newRequest(t, mockUsersUseCase, `{"email":"jane@example.com","name":"Jane Doe","password":"secret123"}`)
		if result.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d: %s", http.StatusOK, result.Code, result.Body.String())
		}
		var body response.SuccessResponse[struct {
			ID       string `json:"id"`
			Email    string `json:"email"`
			Name     string `json:"name"`
			Password string `json:"password"`
		}]
		if err := json.Unmarshal(result.Body.Bytes(), &body); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if body.Data.ID != created.ID || body.Data.Email != created.Email || body.Data.Name != created.Name {
			t.Fatalf("unexpected registration response data: %+v", body.Data)
		}
		if body.Data.Password != "" {
			t.Fatal("registration response must not expose a password")
		}
	})

	t.Run("rejects invalid request without calling use case", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockUsersUseCase := mocksUsersUseCase.NewMockUseCase(ctrl)
		mockUsersUseCase.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)

		result := newRequest(t, mockUsersUseCase, `{"email":"bad-email","name":"J","password":"123"}`)
		if result.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d: %s", http.StatusBadRequest, result.Code, result.Body.String())
		}
	})

	t.Run("returns error response when use case fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockUsersUseCase := mocksUsersUseCase.NewMockUseCase(ctrl)
		input := &entities.User{Email: "jane@example.com", Name: "Jane Doe", Password: "secret123"}
		mockUsersUseCase.EXPECT().CreateUser(gomock.Any(), input).Return(nil, errors.New("create failed"))

		result := newRequest(t, mockUsersUseCase, `{"email":"jane@example.com","name":"Jane Doe","password":"secret123"}`)
		if result.Code != http.StatusInternalServerError {
			t.Fatalf("expected status %d, got %d: %s", http.StatusInternalServerError, result.Code, result.Body.String())
		}
	})
}
