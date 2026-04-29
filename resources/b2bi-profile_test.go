//go:build integration

package resources

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/b2bi"
	b2bitypes "github.com/aws/aws-sdk-go-v2/service/b2bi/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type TestB2BIProfileSuite struct {
	suite.Suite
	svc       *b2bi.Client
	cfg       config.Config
	profileId *string
}

func (suite *TestB2BIProfileSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.cfg = cfg
	suite.svc = b2bi.NewFromConfig(cfg)

	name := fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano())

	resp, err := suite.svc.CreateProfile(ctx, &b2bi.CreateProfileInput{
		Name:         ptr.String(name),
		BusinessName: ptr.String("Test Business"),
		Phone:        ptr.String("+1234567890"),
		Logging:      b2bitypes.LoggingDisabled,
	})
	if err != nil {
		suite.T().Fatalf("failed to create B2BI profile: %v", err)
	}
	suite.profileId = resp.ProfileId
}

func (suite *TestB2BIProfileSuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = suite.svc.DeleteProfile(ctx, &b2bi.DeleteProfileInput{
		ProfileId: suite.profileId,
	})
}

func (suite *TestB2BIProfileSuite) TestList() {
	a := assert.New(suite.T())

	awsCfg := suite.cfg
	lister := B2BIProfileLister{}
	resources, err := lister.List(context.TODO(), &nuke.ListerOpts{
		Region: &nuke.Region{
			Name: "us-east-1",
		},
		Config: &awsCfg,
		Logger: logrus.WithField("test", "b2bi-profile"),
	})
	a.Nil(err)
	a.Greater(len(resources), 0)
}

func (suite *TestB2BIProfileSuite) TestRemove() {
	a := assert.New(suite.T())

	profile := &B2BIProfile{
		svc:       suite.svc,
		ProfileId: suite.profileId,
	}

	err := profile.Remove(context.TODO())
	a.Nil(err)
}

func TestB2BIProfileIntegration(t *testing.T) {
	suite.Run(t, new(TestB2BIProfileSuite))
}
