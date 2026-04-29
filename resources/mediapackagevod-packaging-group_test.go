//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/mediapackagevod"
)

type TestMediaPackageVODPackagingGroupSuite struct {
	suite.Suite
	svc *mediapackagevod.Client
}

func (suite *TestMediaPackageVODPackagingGroupSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = mediapackagevod.NewFromConfig(cfg)
}

func (suite *TestMediaPackageVODPackagingGroupSuite) TestList() {
	a := assert.New(suite.T())
	lister := MediaPackageVODPackagingGroupLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testMediaPackageVODListerOpts)
	a.NoError(err)
	_ = resources
}

func TestMediaPackageVODPackagingGroupIntegration(t *testing.T) {
	suite.Run(t, new(TestMediaPackageVODPackagingGroupSuite))
}
