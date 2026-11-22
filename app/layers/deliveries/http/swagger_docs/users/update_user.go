package users

import (
	swaggermodels "github.com/chaihanij/backend-challenge/app/layers/deliveries/http/swagger_docs/shared"
	"github.com/gin-gonic/gin"
)

var (
	_ = swaggermodels.UpdateUserRequest{}
	_ = swaggermodels.UserResponseEnvelope{}
	_ = swaggermodels.ProblemDetailsResponse{}
)

// UpdateUser
//
//	@Summary		Update a user
//	@Description	Update the email and name for a user by ID.
//	@ID				users-update
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header	string							true	"Bearer access token (format: Bearer <token>)"
//	@Param			id				path	string							true	"User ID (hexadecimal ObjectID)"
//	@Param			request			body	swaggermodels.UpdateUserRequest	true	"Fields to update"
//	@Success		200				{object}	swaggermodels.UserResponseEnvelope	"User updated"
//	@Failure		400				{object}	swaggermodels.ProblemDetailsResponse	"Invalid user ID or request body"
//	@Failure		401				{object}	swaggermodels.ProblemDetailsResponse	"Authentication is required"
//	@Failure		404				{object}	swaggermodels.ProblemDetailsResponse	"User not found"
//	@Failure		500				{object}	swaggermodels.ProblemDetailsResponse	"Failed to update user"
//	@Router			/api/users/{id} [put]
func UpdateUser(c *gin.Context) {}
