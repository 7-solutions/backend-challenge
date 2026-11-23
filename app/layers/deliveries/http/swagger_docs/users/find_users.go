package users

import (
	swaggermodels "github.com/chaihanij/backend-challenge/app/layers/deliveries/http/swagger_docs/shared"
	"github.com/gin-gonic/gin"
)

var (
	_ = swaggermodels.UsersResponseEnvelope{}
	_ = swaggermodels.ProblemDetailsResponse{}
)

// FindUser
//
//	@Summary		Find users
//	@Description	Find users by optional email and name filters.
//	@ID				users-find
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header	string	true	"Bearer access token (format: Bearer <token>)"
//	@Param			email			query	string	false	"Filter by email address"
//	@Param			name			query	string	false	"Filter by user name"
//	@Success		200				{object}	swaggermodels.UsersResponseEnvelope	"Users found"
//	@Failure		400				{object}	swaggermodels.ProblemDetailsResponse	"Invalid filter parameters"
//	@Failure		401				{object}	swaggermodels.ProblemDetailsResponse	"Authentication is required"
//	@Failure		500				{object}	swaggermodels.ProblemDetailsResponse	"Failed to find users"
//	@Router			/api/users [get]
func FindUser(c *gin.Context) {}
