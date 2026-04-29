//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
)

type TestAPIGatewayV2PortalSuite struct {
	suite.Suite
	svc *apigatewayv2.Client
}

func (s *TestAPIGatewayV2PortalSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}
	s.svc = apigatewayv2.NewFromConfig(cfg)
}

func (s *TestAPIGatewayV2PortalSuite) TestListPortals() {
	a := assert.New(s.T())
	lister := APIGatewayV2PortalLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testAPIGatewayV2PortalListerOpts)
	a.NoError(err)
	_ = resources
}

func (s *TestAPIGatewayV2PortalSuite) TestListPortalProducts() {
	a := assert.New(s.T())
	lister := APIGatewayV2PortalProductLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testAPIGatewayV2PortalListerOpts)
	a.NoError(err)
	_ = resources
}

func (s *TestAPIGatewayV2PortalSuite) TestListProductPages() {
	a := assert.New(s.T())
	lister := APIGatewayV2ProductPageLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testAPIGatewayV2PortalListerOpts)
	a.NoError(err)
	_ = resources
}

func (s *TestAPIGatewayV2PortalSuite) TestListProductRestEndpointPages() {
	a := assert.New(s.T())
	lister := APIGatewayV2ProductRestEndpointPageLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testAPIGatewayV2PortalListerOpts)
	a.NoError(err)
	_ = resources
}

func TestAPIGatewayV2PortalIntegration(t *testing.T) {
	suite.Run(t, new(TestAPIGatewayV2PortalSuite))
}
