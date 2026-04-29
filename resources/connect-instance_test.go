//go:build integration

package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/connect"
	connecttypes "github.com/aws/aws-sdk-go-v2/service/connect/types"
)

type TestConnectInstanceSuite struct {
	suite.Suite
	svc *connect.Client
	id  *string
}

func (suite *TestConnectInstanceSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}

	suite.svc = connect.NewFromConfig(cfg)

	alias := "aws-nuke-test-instance"
	resp, err := suite.svc.CreateInstance(ctx, &connect.CreateInstanceInput{
		InstanceAlias:          &alias,
		IdentityManagementType: connecttypes.DirectoryTypeConnectManaged,
		InboundCallsEnabled:    boolPtr(true),
		OutboundCallsEnabled:   boolPtr(true),
	})
	if err != nil {
		suite.T().Fatalf("failed to create test instance: %v", err)
	}
	suite.id = resp.Id
}

func (suite *TestConnectInstanceSuite) TearDownSuite() {
	ctx := context.TODO()
	if suite.id != nil {
		_, _ = suite.svc.DeleteInstance(ctx, &connect.DeleteInstanceInput{
			InstanceId: suite.id,
		})
	}
}

func (suite *TestConnectInstanceSuite) TestList() {
	a := assert.New(suite.T())

	lister := ConnectInstanceLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testConnectListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (suite *TestConnectInstanceSuite) TestRemove() {
	a := assert.New(suite.T())

	resource := ConnectInstance{
		svc: suite.svc,
		Id:  suite.id,
	}

	err := resource.Remove(context.TODO())
	a.NoError(err)
	suite.id = nil
}

func TestConnectInstanceIntegration(t *testing.T) {
	suite.Run(t, new(TestConnectInstanceSuite))
}

func boolPtr(b bool) *bool {
	return &b
}
