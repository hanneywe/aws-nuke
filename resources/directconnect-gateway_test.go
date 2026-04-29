//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/directconnect"
)

type TestDirectConnectGatewaySuite struct {
	suite.Suite
	svc *directconnect.Client
}

func (suite *TestDirectConnectGatewaySuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}

	suite.svc = directconnect.NewFromConfig(cfg)
}

func (suite *TestDirectConnectGatewaySuite) TestList() {
	a := assert.New(suite.T())

	lister := DirectConnectGatewayLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testDirectConnectListerOpts)
	a.NoError(err)
	_ = resources
}

func TestDirectConnectGatewayIntegration(t *testing.T) {
	suite.Run(t, new(TestDirectConnectGatewaySuite))
}
