//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/inspector2"
)

type TestInspector2CodeSecurityScanConfigurationSuite struct {
	suite.Suite
	svc *inspector2.Client
}

func (suite *TestInspector2CodeSecurityScanConfigurationSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = inspector2.NewFromConfig(cfg)
}

func (suite *TestInspector2CodeSecurityScanConfigurationSuite) TestList() {
	a := assert.New(suite.T())
	lister := Inspector2CodeSecurityScanConfigurationLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testInspector2V2ListerOpts)
	a.NoError(err)
	_ = resources
}

func TestInspector2CodeSecurityScanConfigurationIntegration(t *testing.T) {
	suite.Run(t, new(TestInspector2CodeSecurityScanConfigurationSuite))
}
