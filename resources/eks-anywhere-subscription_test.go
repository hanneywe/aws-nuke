//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eks"
)

type TestEKSAnywhereSubscriptionSuite struct {
	suite.Suite
	svc *eks.Client
}

func (suite *TestEKSAnywhereSubscriptionSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = eks.NewFromConfig(cfg)
}

func (suite *TestEKSAnywhereSubscriptionSuite) TestList() {
	a := assert.New(suite.T())
	lister := EKSAnywhereSubscriptionLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testEKSv2ListerOpts)
	a.NoError(err)
	_ = resources
}

func TestEKSAnywhereSubscriptionIntegration(t *testing.T) {
	suite.Run(t, new(TestEKSAnywhereSubscriptionSuite))
}
