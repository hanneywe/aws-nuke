//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/swf"
)

type TestSWFDomainSuite struct {
	suite.Suite
	svc *swf.Client
}

func (suite *TestSWFDomainSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = swf.NewFromConfig(cfg)
}

func (suite *TestSWFDomainSuite) TestList() {
	a := assert.New(suite.T())
	lister := SWFDomainLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testSWFListerOpts)
	a.NoError(err)
	_ = resources
}

func TestSWFDomainIntegration(t *testing.T) {
	suite.Run(t, new(TestSWFDomainSuite))
}
