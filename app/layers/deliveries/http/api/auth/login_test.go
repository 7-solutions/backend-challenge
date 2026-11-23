package auth_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	authapi "github.com/chaihanij/backend-challenge/app/layers/deliveries/http/api/auth"
	middlewares "github.com/chaihanij/backend-challenge/app/layers/deliveries/http/middlewares"
	"github.com/chaihanij/backend-challenge/app/layers/deliveries/http/response"
	mocksAuthUseCase "github.com/chaihanij/backend-challenge/app/mocks/auth"
	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
)

func TestHandler_Login(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	newRequest := func(t *testing.T, mockAuthUseCase *mocksAuthUseCase.MockUseCase, body string) *httptest.ResponseRecorder {
		t.Helper()

		engine := gin.New()
		engine.Use(middlewares.ErrorHandler())
		authapi.New(engine, mockAuthUseCase)

		request := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBufferString(body))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		engine.ServeHTTP(response, request)
		return response
	}

	t.Run("returns JWT token when credentials are valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAuthUseCase := mocksAuthUseCase.NewMockUseCase(ctrl)
		mockAuthUseCase.EXPECT().Login(gomock.Any(), "jane@example.com", "secret123").Return(stringPointer("jwt-token"), nil)

		res := newRequest(t, mockAuthUseCase, `{"email":"jane@example.com","password":"secret123"}`)

		if res.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d: %s", http.StatusOK, res.Code, res.Body.String())
		}

		var body response.SuccessResponse[string]
		if err := json.Unmarshal(res.Body.Bytes(), &body); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if body.Data != "jwt-token" {
			t.Fatalf("expected JWT token %q, got %q", "jwt-token", body.Data)
		}
	})

	t.Run("rejects invalid request without calling use case", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAuthUseCase := mocksAuthUseCase.NewMockUseCase(ctrl)
		mockAuthUseCase.EXPECT().Login(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		res := newRequest(t, mockAuthUseCase, `{"email":"bad-email","password":"123"}`)

		if res.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d: %s", http.StatusBadRequest, res.Code, res.Body.String())
		}
		if res.Body.Len() == 0 {
			t.Fatal("expected validation error response body")
		}
	})

	t.Run("returns error response when use case fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAuthUseCase := mocksAuthUseCase.NewMockUseCase(ctrl)
		mockAuthUseCase.EXPECT().Login(gomock.Any(), "jane@example.com", "secret123").Return(nil, errors.New("login failed"))

		res := newRequest(t, mockAuthUseCase, `{"email":"jane@example.com","password":"secret123"}`)

		if res.Code != http.StatusInternalServerError {
			t.Fatalf("expected status %d, got %d: %s", http.StatusInternalServerError, res.Code, res.Body.String())
		}
		if res.Body.Len() == 0 {
			t.Fatal("expected error response body")
		}
	})
}

func stringPointer(value string) *string {
	return &value
}
