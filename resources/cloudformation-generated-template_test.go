//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
)

type TestCloudFormationGeneratedTemplateSuite struct {
	suite.Suite
	svc *cloudformation.Client
}

func (suite *TestCloudFormationGeneratedTemplateSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = cloudformation.NewFromConfig(cfg)
}

func (suite *TestCloudFormationGeneratedTemplateSuite) TestList() {
	a := assert.New(suite.T())
	lister := CloudFormationGeneratedTemplateLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testCloudFormationListerOpts)
	a.NoError(err)
	_ = resources
}

func TestCloudFormationGeneratedTemplateIntegration(t *testing.T) {
	suite.Run(t, new(TestCloudFormationGeneratedTemplateSuite))
}
