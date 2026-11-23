package users

import (
	swaggermodels "github.com/chaihanij/backend-challenge/app/layers/deliveries/http/swagger_docs/shared"
	"github.com/gin-gonic/gin"
)

var (
	_ = swaggermodels.ProblemDetailsResponse{}
)

// DeleteUserResponseJsonSwagger documents the success envelope returned when
// a user is deleted. The data field is omitted for a nil response payload.
type DeleteUserResponseJsonSwagger struct {
	Code      int         `json:"code" example:"200"`
	Message   string      `json:"message" example:"OK"`
	Meta      interface{} `json:"meta,omitempty"`
	Timestamp string      `json:"timestamp,omitempty"`
}

// DeleteUser
//
//	@Summary		Delete a user
//	@Description	Delete a user by ID.
//	@ID				users-delete
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header	string	true	"Bearer access token (format: Bearer <token>)"
//	@Param			id				path	string	true	"User ID (hexadecimal ObjectID)"
//	@Success		200				{object}	DeleteUserResponseJsonSwagger	"User deleted"
//	@Failure		400				{object}	swaggermodels.ProblemDetailsResponse	"Invalid user ID"
//	@Failure		401				{object}	swaggermodels.ProblemDetailsResponse	"Authentication is required"
//	@Failure		404				{object}	swaggermodels.ProblemDetailsResponse	"User not found"
//	@Failure		500				{object}	swaggermodels.ProblemDetailsResponse	"Failed to delete user"
//	@Router			/api/users/{id} [delete]
func DeleteUser(c *gin.Context) {}
