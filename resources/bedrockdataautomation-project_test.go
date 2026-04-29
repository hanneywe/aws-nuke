//go:build integration

package resources

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockdataautomation"
)

type TestBedrockDataAutomationProjectSuite struct {
	suite.Suite
	svc        *bedrockdataautomation.Client
	projectArn *string
}

func (suite *TestBedrockDataAutomationProjectSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}

	suite.svc = bedrockdataautomation.NewFromConfig(cfg)

	resp, err := suite.svc.CreateDataAutomationProject(ctx, &bedrockdataautomation.CreateDataAutomationProjectInput{
		ProjectName: ptrString(fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano())),
	})
	if err != nil {
		suite.T().Fatalf("failed to create test project: %v", err)
	}
	suite.projectArn = resp.ProjectArn
}

func (suite *TestBedrockDataAutomationProjectSuite) TearDownSuite() {
	ctx := context.TODO()
	if suite.projectArn != nil {
		_, _ = suite.svc.DeleteDataAutomationProject(ctx, &bedrockdataautomation.DeleteDataAutomationProjectInput{
			ProjectArn: suite.projectArn,
		})
	}
}

func (suite *TestBedrockDataAutomationProjectSuite) TestList() {
	a := assert.New(suite.T())

	lister := BedrockDataAutomationProjectLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testBedrockDataAutomationListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (suite *TestBedrockDataAutomationProjectSuite) TestRemove() {
	a := assert.New(suite.T())

	resource := BedrockDataAutomationProject{
		svc:        suite.svc,
		ProjectArn: suite.projectArn,
	}

	err := resource.Remove(context.TODO())
	a.NoError(err)
	suite.projectArn = nil
}

func TestBedrockDataAutomationProjectIntegration(t *testing.T) {
	suite.Run(t, new(TestBedrockDataAutomationProjectSuite))
}

func ptrString(s string) *string {
	return &s
}
