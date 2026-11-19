package errors_test

import (
	"testing"

	httperrors "github.com/chaihanij/backend-challenge/app/layers/deliveries/http/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BaseUrlTestSuite struct {
	suite.Suite
}

func (suite *BaseUrlTestSuite) SetupTest() {
	httperrors.BaseUrl = ""
}

func (suite *BaseUrlTestSuite) TestSetBaseUrl_WithValidUrl() {
	url := "https://api.example.com"
	httperrors.SetBaseUrl(url)

	assert.Equal(suite.T(), url, httperrors.BaseUrl)
}

func (suite *BaseUrlTestSuite) TestSetBaseUrl_WithCustomUrl() {
	url := "https://custom.domain.com/errors"
	httperrors.SetBaseUrl(url)

	assert.Equal(suite.T(), url, httperrors.BaseUrl)
}

func (suite *BaseUrlTestSuite) TestSetBaseUrl_WithEmptyString() {
	initialUrl := "https://initial.com"
	httperrors.BaseUrl = initialUrl

	httperrors.SetBaseUrl("")

	assert.Equal(suite.T(), "", httperrors.BaseUrl)
}

func TestBaseUrlTestSuite(t *testing.T) {
	suite.Run(t, new(BaseUrlTestSuite))
}
