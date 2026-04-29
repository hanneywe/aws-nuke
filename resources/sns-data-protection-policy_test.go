//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type TestSNSDataProtectionPolicySuite struct {
	suite.Suite
	svc *sns.Client
}

func (suite *TestSNSDataProtectionPolicySuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = sns.NewFromConfig(cfg)
}

func (suite *TestSNSDataProtectionPolicySuite) TestList() {
	a := assert.New(suite.T())
	lister := SNSDataProtectionPolicyLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testSNSV2ListerOpts)
	a.NoError(err)
	_ = resources
}

func TestSNSDataProtectionPolicyIntegration(t *testing.T) {
	suite.Run(t, new(TestSNSDataProtectionPolicySuite))
}
