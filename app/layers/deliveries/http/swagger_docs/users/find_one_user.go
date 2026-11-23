package users

import (
	swaggermodels "github.com/chaihanij/backend-challenge/app/layers/deliveries/http/swagger_docs/shared"
	"github.com/gin-gonic/gin"
)

var (
	_ = swaggermodels.UserResponseEnvelope{}
	_ = swaggermodels.ProblemDetailsResponse{}
)

// FindOneUser
//
//	@Summary		Find one user
//	@Description	Retrieves a single user by their ID.
//	@ID				users-find-one
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							true	"Bearer access token (format: Bearer <token>)"
//	@Param			id				path		string							true	"User ID"
//	@Success		200				{object}	swaggermodels.UserResponseEnvelope	"Successfully retrieved the user"
//	@Failure		400				{object}	swaggermodels.ProblemDetailsResponse	"The user ID is invalid"
//	@Failure		401				{object}	swaggermodels.ProblemDetailsResponse	"Authentication is required"
//	@Failure		404				{object}	swaggermodels.ProblemDetailsResponse	"User not found"
//	@Failure		500				{object}	swaggermodels.ProblemDetailsResponse	"Failed to retrieve the user"
//	@Router			/api/users/{id} [get]
func FindOneUser(c *gin.Context) {
}
