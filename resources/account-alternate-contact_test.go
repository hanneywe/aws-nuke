//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/account"
	accounttypes "github.com/aws/aws-sdk-go-v2/service/account/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type TestAccountAlternateContactSuite struct {
	suite.Suite
	svc *account.Client
	cfg config.Config
}

func (suite *TestAccountAlternateContactSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.cfg = cfg
	suite.svc = account.NewFromConfig(cfg)

	// Create a billing alternate contact for testing
	_, err = suite.svc.PutAlternateContact(ctx, &account.PutAlternateContactInput{
		AlternateContactType: accounttypes.AlternateContactTypeBilling,
		Name:                 ptr.String("Test Billing Contact"),
		EmailAddress:         ptr.String("billing-test@example.com"),
		PhoneNumber:          ptr.String("+1234567890"),
		Title:                ptr.String("Billing Manager"),
	})
	if err != nil {
		suite.T().Fatalf("failed to create alternate contact: %v", err)
	}
}

func (suite *TestAccountAlternateContactSuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = suite.svc.DeleteAlternateContact(ctx, &account.DeleteAlternateContactInput{
		AlternateContactType: accounttypes.AlternateContactTypeBilling,
	})
}

func (suite *TestAccountAlternateContactSuite) TestList() {
	a := assert.New(suite.T())

	awsCfg := suite.cfg
	lister := AccountAlternateContactLister{}
	resources, err := lister.List(context.TODO(), &nuke.ListerOpts{
		Region: &nuke.Region{
			Name: "us-east-1",
		},
		Config: &awsCfg,
		Logger: logrus.WithField("test", "account-alternate-contact"),
	})
	a.Nil(err)
	a.Greater(len(resources), 0)
}

func (suite *TestAccountAlternateContactSuite) TestRemove() {
	a := assert.New(suite.T())

	contact := &AccountAlternateContact{
		svc:                  suite.svc,
		AlternateContactType: ptr.String("BILLING"),
	}

	err := contact.Remove(context.TODO())
	a.Nil(err)
}

func TestAccountAlternateContactIntegration(t *testing.T) {
	suite.Run(t, new(TestAccountAlternateContactSuite))
}
