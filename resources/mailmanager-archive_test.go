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

type TestMailManagerArchiveSuite struct {
	suite.Suite
	svc *mailmanager.Client
}

func (suite *TestMailManagerArchiveSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = mailmanager.NewFromConfig(cfg)
}

func (suite *TestMailManagerArchiveSuite) TestList() {
	a := assert.New(suite.T())
	lister := MailManagerArchiveLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testMailManagerListerOpts)
	a.NoError(err)
	_ = resources
}

func TestMailManagerArchiveIntegration(t *testing.T) {
	suite.Run(t, new(TestMailManagerArchiveSuite))
}
