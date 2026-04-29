//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/synthetics"
)

type TestSyntheticsGroupSuite struct {
	suite.Suite
	svc *synthetics.Client
}

func (suite *TestSyntheticsGroupSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = synthetics.NewFromConfig(cfg)
}

func (suite *TestSyntheticsGroupSuite) TestList() {
	a := assert.New(suite.T())
	lister := SyntheticsGroupLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testSyntheticsListerOpts)
	a.NoError(err)
	_ = resources
}

func TestSyntheticsGroupIntegration(t *testing.T) {
	suite.Run(t, new(TestSyntheticsGroupSuite))
}
