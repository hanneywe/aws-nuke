//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/panorama"
)

type TestPanoramaPackageSuite struct {
	suite.Suite
	svc *panorama.Client
}

func (suite *TestPanoramaPackageSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = panorama.NewFromConfig(cfg)
}

func (suite *TestPanoramaPackageSuite) TestList() {
	a := assert.New(suite.T())
	lister := PanoramaPackageLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testPanoramaListerOpts)
	a.NoError(err)
	_ = resources
}

func TestPanoramaPackageIntegration(t *testing.T) {
	suite.Run(t, new(TestPanoramaPackageSuite))
}
