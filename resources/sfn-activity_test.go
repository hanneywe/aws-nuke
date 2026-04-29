//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
)

type TestSFNActivitySuite struct {
	suite.Suite
	svc *sfn.Client
}

func (suite *TestSFNActivitySuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = sfn.NewFromConfig(cfg)
}

func (suite *TestSFNActivitySuite) TestList() {
	a := assert.New(suite.T())
	lister := SFNActivityLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testSFNv2ListerOpts)
	a.NoError(err)
	_ = resources
}

func TestSFNActivityIntegration(t *testing.T) {
	suite.Run(t, new(TestSFNActivitySuite))
}
