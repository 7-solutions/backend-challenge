package errors_test

import (
	"errors"
	"testing"

	httperrors "github.com/chaihanij/backend-challenge/app/layers/deliveries/http/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type HelpersTestSuite struct {
	suite.Suite
}

func (suite *HelpersTestSuite) SetupTest() {
	httperrors.BaseUrl = "https://api.example.com"
}

// Test for firstErr helper function
func (suite *HelpersTestSuite) TestFirstErr_WithEmptySlice() {
	// Note: firstErr is unexported, so we'll test it indirectly through factory functions
	// or we can create a test that uses the factory which uses firstErr
	result := httperrors.NewBadRequest("test detail")
	assert.NotNil(suite.T(), result)
}

func (suite *HelpersTestSuite) TestFirstErr_WithMultipleErrors() {
	err1 := errors.New("error 1")
	err2 := errors.New("error 2")

	// Test that the first error is used
	result := httperrors.NewBadRequest("test detail", err1, err2)
	assert.NotNil(suite.T(), result)

	// The wrapped error should match the first one
	wrappedErr := result.Unwrap()
	assert.Equal(suite.T(), err1, wrappedErr)
}

// Test for problem type generation
func (suite *HelpersTestSuite) TestProblemType_GeneratesCorrectUrl() {
	result := httperrors.NewBadRequest("test")

	// The JSON should contain the base URL and problem path
	assert.NotNil(suite.T(), result)
	jsonStr := string(result.JSON())
	assert.Contains(suite.T(), jsonStr, "https://api.example.com")
	assert.Contains(suite.T(), jsonStr, "/problems/")
	assert.Contains(suite.T(), jsonStr, "bad-request")
}

func (suite *HelpersTestSuite) TestProblemType_DifferentProblems() {
	testCases := []struct {
		name            string
		factory         func(string) httperrors.ProblemDetails
		expectedProblem string
		expectedStatus  int
		expectedTitle   string
	}{
		{
			name:            "BadRequest",
			factory:         func(d string) httperrors.ProblemDetails { return httperrors.NewBadRequest(d) },
			expectedProblem: "bad-request",
			expectedStatus:  400,
			expectedTitle:   "Bad Request",
		},
		{
			name:            "Unauthorized",
			factory:         func(d string) httperrors.ProblemDetails { return httperrors.NewUnauthorized(d) },
			expectedProblem: "unauthorized",
			expectedStatus:  401,
			expectedTitle:   "Unauthorized",
		},
		{
			name:            "Forbidden",
			factory:         func(d string) httperrors.ProblemDetails { return httperrors.NewForbidden(d) },
			expectedProblem: "forbidden",
			expectedStatus:  403,
			expectedTitle:   "Forbidden",
		},
		{
			name:            "NotFound",
			factory:         func(d string) httperrors.ProblemDetails { return httperrors.NewNotFound(d) },
			expectedProblem: "not-found",
			expectedStatus:  404,
			expectedTitle:   "Not Found",
		},
		{
			name:            "Internal",
			factory:         func(d string) httperrors.ProblemDetails { return httperrors.NewInternal(d) },
			expectedProblem: "internal-server-error",
			expectedStatus:  500,
			expectedTitle:   "Internal Server Error",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			result := tc.factory("detail")
			require.NotNil(suite.T(), result)

			// Check status and title
			assert.Equal(suite.T(), tc.expectedStatus, result.GetStatus())
			assert.Equal(suite.T(), tc.expectedTitle, result.GetTitle())
			assert.Equal(suite.T(), "detail", result.GetDetail())

			// Check JSON contains expected problem path
			jsonStr := string(result.JSON())
			assert.Contains(suite.T(), jsonStr, tc.expectedProblem)
			assert.Contains(suite.T(), jsonStr, httperrors.BaseUrl)
		})
	}
}

func TestHelpersTestSuite(t *testing.T) {
	suite.Run(t, new(HelpersTestSuite))
}
