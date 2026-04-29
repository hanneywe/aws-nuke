//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/mediaconnect"
)

type TestMediaConnectGatewaySuite struct {
	suite.Suite
	svc *mediaconnect.Client
}

func (s *TestMediaConnectGatewaySuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = mediaconnect.NewFromConfig(cfg)
}

func (s *TestMediaConnectGatewaySuite) TestList() {
	assertions := assert.New(s.T())

	lister := MediaConnectGatewayLister{}
	resources, err := lister.List(context.TODO(), testMediaConnectListerOpts)

	assertions.Nil(err)
	assertions.NotNil(resources)
}

func TestMediaConnectGatewayIntegration(t *testing.T) {
	suite.Run(t, new(TestMediaConnectGatewaySuite))
}
