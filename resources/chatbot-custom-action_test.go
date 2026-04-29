//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/chatbot"
)

type TestChatbotCustomActionSuite struct {
	suite.Suite
	svc *chatbot.Client
	arn *string
}

func (suite *TestChatbotCustomActionSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}

	suite.svc = chatbot.NewFromConfig(cfg)

	name := "aws-nuke-test-action"
	resp, err := suite.svc.CreateCustomAction(ctx, &chatbot.CreateCustomActionInput{
		ActionName: &name,
	})
	if err != nil {
		suite.T().Fatalf("failed to create test custom action: %v", err)
	}
	suite.arn = resp.CustomActionArn
}

func (suite *TestChatbotCustomActionSuite) TearDownSuite() {
	ctx := context.TODO()
	if suite.arn != nil {
		_, _ = suite.svc.DeleteCustomAction(ctx, &chatbot.DeleteCustomActionInput{
			CustomActionArn: suite.arn,
		})
	}
}

func (suite *TestChatbotCustomActionSuite) TestList() {
	a := assert.New(suite.T())

	lister := ChatbotCustomActionLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testChatbotListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (suite *TestChatbotCustomActionSuite) TestRemove() {
	a := assert.New(suite.T())

	resource := ChatbotCustomAction{
		svc:             suite.svc,
		CustomActionArn: suite.arn,
	}

	err := resource.Remove(context.TODO())
	a.NoError(err)
	suite.arn = nil
}

func TestChatbotCustomActionIntegration(t *testing.T) {
	suite.Run(t, new(TestChatbotCustomActionSuite))
}
