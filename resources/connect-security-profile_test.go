//go:build integration

package resources

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/connect"
	connecttypes "github.com/aws/aws-sdk-go-v2/service/connect/types"
)

type TestConnectSecurityProfileSuite struct {
	suite.Suite
	svc               *connect.Client
	instanceId        *string
	securityProfileId *string
}

func (suite *TestConnectSecurityProfileSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}

	suite.svc = connect.NewFromConfig(cfg)

	// Create a test instance first
	alias := fmt.Sprintf("aws-nuke-test-sp-%d", time.Now().UnixNano())
	instResp, err := suite.svc.CreateInstance(ctx, &connect.CreateInstanceInput{
		InstanceAlias:          &alias,
		IdentityManagementType: connecttypes.DirectoryTypeConnectManaged,
		InboundCallsEnabled:    boolPtr(true),
		OutboundCallsEnabled:   boolPtr(true),
	})
	if err != nil {
		suite.T().Fatalf("failed to create test instance: %v", err)
	}
	suite.instanceId = instResp.Id

	// Create a test security profile
	profileName := fmt.Sprintf("aws-nuke-test-profile-%d", time.Now().UnixNano())
	profResp, err := suite.svc.CreateSecurityProfile(ctx, &connect.CreateSecurityProfileInput{
		InstanceId:          suite.instanceId,
		SecurityProfileName: &profileName,
	})
	if err != nil {
		suite.T().Fatalf("failed to create test security profile: %v", err)
	}
	suite.securityProfileId = profResp.SecurityProfileId
}

func (suite *TestConnectSecurityProfileSuite) TearDownSuite() {
	ctx := context.TODO()
	if suite.securityProfileId != nil && suite.instanceId != nil {
		_, _ = suite.svc.DeleteSecurityProfile(ctx, &connect.DeleteSecurityProfileInput{
			InstanceId:        suite.instanceId,
			SecurityProfileId: suite.securityProfileId,
		})
	}
	if suite.instanceId != nil {
		_, _ = suite.svc.DeleteInstance(ctx, &connect.DeleteInstanceInput{
			InstanceId: suite.instanceId,
		})
	}
}

func (suite *TestConnectSecurityProfileSuite) TestList() {
	a := assert.New(suite.T())

	lister := ConnectSecurityProfileLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testConnectListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (suite *TestConnectSecurityProfileSuite) TestRemove() {
	a := assert.New(suite.T())

	resource := ConnectSecurityProfile{
		svc:               suite.svc,
		InstanceId:        suite.instanceId,
		SecurityProfileId: suite.securityProfileId,
	}

	err := resource.Remove(context.TODO())
	a.NoError(err)
	suite.securityProfileId = nil
}

func TestConnectSecurityProfileIntegration(t *testing.T) {
	suite.Run(t, new(TestConnectSecurityProfileSuite))
}
