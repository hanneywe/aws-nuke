//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
)

type TestAPIGatewayGatewayResponseSuite struct {
	suite.Suite
	svc *apigateway.Client
}

func (suite *TestAPIGatewayGatewayResponseSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = apigateway.NewFromConfig(cfg)
}

func (suite *TestAPIGatewayGatewayResponseSuite) TestList() {
	a := assert.New(suite.T())
	lister := APIGatewayGatewayResponseLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testAPIGatewayV2ListerOpts)
	a.NoError(err)
	_ = resources
}

func TestAPIGatewayGatewayResponseIntegration(t *testing.T) {
	suite.Run(t, new(TestAPIGatewayGatewayResponseSuite))
}
