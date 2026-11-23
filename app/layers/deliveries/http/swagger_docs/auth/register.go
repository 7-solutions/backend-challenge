package auth

import (
	swaggermodels "github.com/chaihanij/backend-challenge/app/layers/deliveries/http/swagger_docs/shared"
	"github.com/gin-gonic/gin"
)

var (
	_ = swaggermodels.CreateUserRequest{}
	_ = swaggermodels.UserResponseEnvelope{}
)

// Register
//
//	@Summary		Register
//	@Description	Create a user account. Registration is public; use the login endpoint to obtain a JWT.
//	@ID				auth-register
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body	swaggermodels.CreateUserRequest	true	"User details"
//	@Success		200	{object}	swaggermodels.UserResponseEnvelope	"User registered"
//	@Failure		400	{object}	swaggermodels.ProblemDetailsResponse	"Invalid registration request"
//	@Failure		409	{object}	swaggermodels.ProblemDetailsResponse	"Email already registered"
//	@Failure		500	{object}	swaggermodels.ProblemDetailsResponse	"Registration failed"
//	@Router			/api/auth/register [post]
func Register(c *gin.Context) {}
