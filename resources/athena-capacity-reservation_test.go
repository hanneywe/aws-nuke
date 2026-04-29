//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/athena"
)

type TestAthenaCapacityReservationSuite struct {
	suite.Suite
	svc *athena.Client
}

func (suite *TestAthenaCapacityReservationSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = athena.NewFromConfig(cfg)
}

func (suite *TestAthenaCapacityReservationSuite) TestList() {
	a := assert.New(suite.T())
	lister := AthenaCapacityReservationLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testAthenaListerOpts)
	a.NoError(err)
	_ = resources
}

func TestAthenaCapacityReservationIntegration(t *testing.T) {
	suite.Run(t, new(TestAthenaCapacityReservationSuite))
}
