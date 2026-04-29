//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/emr"
)

type TestEMRBlockPublicAccessConfigurationSuite struct {
	suite.Suite
	svc *emr.Client
}

func (suite *TestEMRBlockPublicAccessConfigurationSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = emr.NewFromConfig(cfg)
}

func (suite *TestEMRBlockPublicAccessConfigurationSuite) TestList() {
	a := assert.New(suite.T())
	lister := EMRBlockPublicAccessConfigurationLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testEMRV2ListerOpts)
	a.NoError(err)
	_ = resources
}

func TestEMRBlockPublicAccessConfigurationIntegration(t *testing.T) {
	suite.Run(t, new(TestEMRBlockPublicAccessConfigurationSuite))
}
