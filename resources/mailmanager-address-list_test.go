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

type TestMailManagerAddressListSuite struct {
	suite.Suite
	svc *mailmanager.Client
}

func (suite *TestMailManagerAddressListSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = mailmanager.NewFromConfig(cfg)
}

func (suite *TestMailManagerAddressListSuite) TestList() {
	a := assert.New(suite.T())
	lister := MailManagerAddressListLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testMailManagerListerOpts)
	a.NoError(err)
	_ = resources
}

func TestMailManagerAddressListIntegration(t *testing.T) {
	suite.Run(t, new(TestMailManagerAddressListSuite))
}
