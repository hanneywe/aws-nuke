//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/mediapackagev2"
)

type TestMediaPackageV2ChannelSuite struct {
	suite.Suite
	svc *mediapackagev2.Client
}

func (suite *TestMediaPackageV2ChannelSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = mediapackagev2.NewFromConfig(cfg)
}

func (suite *TestMediaPackageV2ChannelSuite) TestList() {
	a := assert.New(suite.T())
	lister := MediaPackageV2ChannelLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testMediaPackageV2ListerOpts)
	a.NoError(err)
	_ = resources
}

func TestMediaPackageV2ChannelIntegration(t *testing.T) {
	suite.Run(t, new(TestMediaPackageV2ChannelSuite))
}
