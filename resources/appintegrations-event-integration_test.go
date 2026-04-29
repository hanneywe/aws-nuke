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
	"github.com/aws/aws-sdk-go-v2/service/appintegrations"
	aitypes "github.com/aws/aws-sdk-go-v2/service/appintegrations/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type TestAppIntegrationsEventIntegrationSuite struct {
	suite.Suite
	svc  *appintegrations.Client
	cfg  config.Config
	name *string
}

func (suite *TestAppIntegrationsEventIntegrationSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.cfg = cfg
	suite.svc = appintegrations.NewFromConfig(cfg)

	suite.name = ptr.String(fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano()))

	_, err = suite.svc.CreateEventIntegration(ctx, &appintegrations.CreateEventIntegrationInput{
		Name:           suite.name,
		EventBridgeBus: ptr.String("default"),
		EventFilter: &aitypes.EventFilter{
			Source: ptr.String("aws.partner/example.com"),
		},
	})
	if err != nil {
		suite.T().Fatalf("failed to create event integration: %v", err)
	}
}

func (suite *TestAppIntegrationsEventIntegrationSuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = suite.svc.DeleteEventIntegration(ctx, &appintegrations.DeleteEventIntegrationInput{
		Name: suite.name,
	})
}

func (suite *TestAppIntegrationsEventIntegrationSuite) TestList() {
	a := assert.New(suite.T())

	awsCfg := suite.cfg
	lister := AppIntegrationsEventIntegrationLister{}
	resources, err := lister.List(context.TODO(), &nuke.ListerOpts{
		Region: &nuke.Region{
			Name: "us-east-1",
		},
		Config: &awsCfg,
		Logger: logrus.WithField("test", "appintegrations-event-integration"),
	})
	a.Nil(err)
	a.Greater(len(resources), 0)
}

func (suite *TestAppIntegrationsEventIntegrationSuite) TestRemove() {
	a := assert.New(suite.T())

	ei := &AppIntegrationsEventIntegration{
		svc:  suite.svc,
		Name: suite.name,
	}

	err := ei.Remove(context.TODO())
	a.Nil(err)
}

func TestAppIntegrationsEventIntegrationIntegration(t *testing.T) {
	suite.Run(t, new(TestAppIntegrationsEventIntegrationSuite))
}
