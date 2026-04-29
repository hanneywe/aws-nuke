//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type TestEC2RouteServerSuite struct {
	suite.Suite
	svc *ec2.Client
}

func (suite *TestEC2RouteServerSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = ec2.NewFromConfig(cfg)
}

func (suite *TestEC2RouteServerSuite) TestList() {
	a := assert.New(suite.T())
	lister := EC2RouteServerLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	a.NoError(err)
	_ = resources
}

func TestEC2RouteServerIntegration(t *testing.T) {
	suite.Run(t, new(TestEC2RouteServerSuite))
}
