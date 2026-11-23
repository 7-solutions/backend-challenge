package users

import (
	swaggermodels "github.com/chaihanij/backend-challenge/app/layers/deliveries/http/swagger_docs/shared"
	"github.com/gin-gonic/gin"
)

var (
	_ = swaggermodels.CreateUserRequest{}
	_ = swaggermodels.UserResponseEnvelope{}
	_ = swaggermodels.ProblemDetailsResponse{}
)

// CreateUser
//
//	@Summary		Create a user
//	@Description	Create a user with an email, name, and password.
//	@ID				users-create
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header	string							true	"Bearer access token (format: Bearer <token>)"
//	@Param			request			body	swaggermodels.CreateUserRequest	true	"User details"
//	@Success		200				{object}	swaggermodels.UserResponseEnvelope	"User created"
//	@Failure		400				{object}	swaggermodels.ProblemDetailsResponse	"Invalid user details"
//	@Failure		401				{object}	swaggermodels.ProblemDetailsResponse	"Authentication is required"
//	@Failure		409				{object}	swaggermodels.ProblemDetailsResponse	"A user with this email already exists"
//	@Failure		500				{object}	swaggermodels.ProblemDetailsResponse	"Failed to create user"
//	@Router			/api/users [post]
func CreateUser(c *gin.Context) {}
