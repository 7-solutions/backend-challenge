package auth_test

import (
	"bytes"
	"encoding/json"
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

	executeWithRequest := func(mockAuthUseCase *mocksAuthUseCase.MockUseCase, jsonBody []byte) *httptest.ResponseRecorder {
		res := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		engine := gin.New()

		authapi.New(engine, mockAuthUseCase)

		req, err := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Fatalf("failed to build request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		engine.ServeHTTP(res, req)
		engine.Use(middlewares.RequestId())
		engine.Use(middlewares.CorrelationId())
		engine.Use(middlewares.TraceId())
		engine.Use(middlewares.ErrorHandler())
		return res
	}

	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		email := "email@test.com"
		password := "P@ssw0rd"
		body := []byte(`{"email":"email@test.com","password":"P@ssw0rd"}`)
		token := "mocked-token"

		mockAuthUseCase := mocksAuthUseCase.NewMockUseCase(ctrl)
		mockAuthUseCase.EXPECT().Login(gomock.Any(), email, password).Return(&token, nil).Times(1)

		res := executeWithRequest(mockAuthUseCase, body)

		if res.Code != http.StatusOK {
			t.Fatalf("expected status code 200, got %d", res.Code)
		}

		var responseBody response.SuccessResponse[string]
		if err := json.Unmarshal(res.Body.Bytes(), &responseBody); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		if responseBody.Data != token {
			t.Fatalf("expected data to be %s, got %s", token, responseBody.Data)
		}
	})

	t.Run("does not call login use case when request payload is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthUseCase := mocksAuthUseCase.NewMockUseCase(ctrl)
		mockAuthUseCase.EXPECT().Login(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		res := executeWithRequest(mockAuthUseCase, []byte(`{"email":"invalid-email","password":"123"}`))

		if res.Body.Len() != 0 {
			t.Fatalf("expected empty body on error path, got %s", res.Body.String())
		}
		t.Log("response body:", string(res.Body.Bytes()))
		var responseBody map[string]interface{}
		_ = json.Unmarshal(res.Body.Bytes(), &responseBody)
		t.Log("response body:", responseBody)
	})

	t.Run("returns no success response when login use case returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		email := "email@test.com"
		password := "P@ssw0rd"
		body := []byte(`{"email":"email@test.com","password":"P@ssw0rd"}`)

		mockAuthUseCase := mocksAuthUseCase.NewMockUseCase(ctrl)
		mockAuthUseCase.EXPECT().Login(gomock.Any(), email, password).Return(nil, gin.Error{})

		res := executeWithRequest(mockAuthUseCase, body)

		if res.Body.Len() != 0 {
			t.Fatalf("expected empty body on error path, got %s", res.Body.String())
		}
	})

}
