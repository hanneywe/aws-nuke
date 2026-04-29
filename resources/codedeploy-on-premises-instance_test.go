//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/codedeploy"
)

type TestCodeDeployOnPremisesInstanceSuite struct {
	suite.Suite
	svc *codedeploy.Client
}

func (suite *TestCodeDeployOnPremisesInstanceSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = codedeploy.NewFromConfig(cfg)
}

func (suite *TestCodeDeployOnPremisesInstanceSuite) TestList() {
	a := assert.New(suite.T())
	lister := CodeDeployOnPremisesInstanceLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testCodeDeployV2ListerOpts)
	a.NoError(err)
	_ = resources
}

func TestCodeDeployOnPremisesInstanceIntegration(t *testing.T) {
	suite.Run(t, new(TestCodeDeployOnPremisesInstanceSuite))
}
