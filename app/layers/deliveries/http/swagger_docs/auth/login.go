package auth

import (
	_ "github.com/chaihanij/backend-challenge/app/layers/deliveries/http/api/auth/models"
	swaggermodels "github.com/chaihanij/backend-challenge/app/layers/deliveries/http/swagger_docs/shared"
	"github.com/gin-gonic/gin"
)

var (
	_ = swaggermodels.ProblemDetailsResponse{}
)

// LoginRequestJsonSwagger documents the JSON body accepted by login.
type LoginRequestJsonSwagger struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"secret123"`
}

// LoginMetaJsonSwagger documents the response metadata envelope.
type LoginMetaJsonSwagger struct {
	RequestID     string `json:"requestId,omitempty"`
	CorrelationID string `json:"correlationId,omitempty"`
	TraceID       string `json:"traceId,omitempty"`
}

// LoginResponseJsonSwagger documents the success envelope returned by login.
type LoginResponseJsonSwagger struct {
	Code      int                   `json:"code" example:"200"`
	Message   string                `json:"message" example:"OK"`
	Data      string                `json:"data" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	Meta      *LoginMetaJsonSwagger `json:"meta,omitempty"`
	Timestamp string                `json:"timestamp,omitempty"`
}

// Login
//
//	@Summary		Login
//	@Description	Authenticate with an email and password to receive a JWT access token.
//	@ID				auth-login
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body	LoginRequestJsonSwagger	true	"Login credentials"
//	@Success		200		{object}	LoginResponseJsonSwagger	"Login successful"
//	@Failure		400		{object}	swaggermodels.ProblemDetailsResponse	"Invalid login request"
//	@Failure		401		{object}	swaggermodels.ProblemDetailsResponse	"Invalid password"
//	@Failure		404		{object}	swaggermodels.ProblemDetailsResponse	"User not found"
//	@Failure		500		{object}	swaggermodels.ProblemDetailsResponse	"Failed to authenticate"
//	@Router			/api/auth/login [post]
func Login(c *gin.Context) {}
