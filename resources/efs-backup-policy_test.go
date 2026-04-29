//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/efs"
)

type TestEFSBackupPolicySuite struct {
	suite.Suite
	svc *efs.Client
}

func (suite *TestEFSBackupPolicySuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = efs.NewFromConfig(cfg)
}

func (suite *TestEFSBackupPolicySuite) TestList() {
	a := assert.New(suite.T())
	lister := EFSBackupPolicyLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testEFSV2ListerOpts)
	a.NoError(err)
	_ = resources
}

func TestEFSBackupPolicyIntegration(t *testing.T) {
	suite.Run(t, new(TestEFSBackupPolicySuite))
}
