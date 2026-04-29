//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type TestSSMMaintenanceWindowTargetSuite struct {
	suite.Suite
	svc *ssm.Client
}

func (s *TestSSMMaintenanceWindowTargetSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}
	s.svc = ssm.NewFromConfig(cfg)
}

func (s *TestSSMMaintenanceWindowTargetSuite) TestList() {
	a := assert.New(s.T())
	lister := SSMMaintenanceWindowTargetLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testSSMV2ListerOpts)
	a.NoError(err)
	_ = resources
}

func TestSSMMaintenanceWindowTargetIntegration(t *testing.T) {
	suite.Run(t, new(TestSSMMaintenanceWindowTargetSuite))
}
