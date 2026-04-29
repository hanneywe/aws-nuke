//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/mailmanager"
)

type TestMailManagerAddonSubscriptionSuite struct {
	suite.Suite
	svc *mailmanager.Client
}

func (suite *TestMailManagerAddonSubscriptionSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = mailmanager.NewFromConfig(cfg)
}

func (suite *TestMailManagerAddonSubscriptionSuite) TestList() {
	a := assert.New(suite.T())
	lister := MailManagerAddonSubscriptionLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testMailManagerListerOpts)
	a.NoError(err)
	_ = resources
}

func TestMailManagerAddonSubscriptionIntegration(t *testing.T) {
	suite.Run(t, new(TestMailManagerAddonSubscriptionSuite))
}
