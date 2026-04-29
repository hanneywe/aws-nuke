//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/migrationhubrefactorspaces"
)

type TestMigrationHubRefactorSpacesEnvironmentSuite struct {
	suite.Suite
	svc *migrationhubrefactorspaces.Client
}

func (suite *TestMigrationHubRefactorSpacesEnvironmentSuite) SetupSuite() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}
	suite.svc = migrationhubrefactorspaces.NewFromConfig(cfg)
}

func (suite *TestMigrationHubRefactorSpacesEnvironmentSuite) TestList() {
	a := assert.New(suite.T())
	lister := MigrationHubRefactorSpacesEnvironmentLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testMigrationHubRefactorSpacesListerOpts)
	a.NoError(err)
	_ = resources
}

func TestMigrationHubRefactorSpacesEnvironmentIntegration(t *testing.T) {
	suite.Run(t, new(TestMigrationHubRefactorSpacesEnvironmentSuite))
}
