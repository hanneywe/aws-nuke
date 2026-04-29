//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
)

type TestSESv2DedicatedIpPoolSuite struct {
	suite.Suite
	svc *sesv2.Client
}

func (suite *TestSESv2DedicatedIpPoolSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = sesv2.NewFromConfig(cfg)
}

func (suite *TestSESv2DedicatedIpPoolSuite) TestList() {
	a := assert.New(suite.T())
	lister := SESv2DedicatedIpPoolLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testSESv2ListerOpts)
	a.NoError(err)
	_ = resources
}

func TestSESv2DedicatedIpPoolIntegration(t *testing.T) {
	suite.Run(t, new(TestSESv2DedicatedIpPoolSuite))
}
