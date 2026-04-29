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

	"github.com/aws/aws-sdk-go-v2/config"
	iamv2 "github.com/aws/aws-sdk-go-v2/service/iam"
)

type TestIAMRolePermissionBoundarySuite struct {
	suite.Suite
	svc      *iamv2.Client
	roleName *string
}

func (suite *TestIAMRolePermissionBoundarySuite) SetupSuite() {
	suite.roleName = ptr.String(fmt.Sprintf("aws-nuke-test-pb-%d", time.Now().UnixNano()))

	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}

	suite.svc = iamv2.NewFromConfig(cfg)
}

func (suite *TestIAMRolePermissionBoundarySuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = suite.svc.DeleteRolePermissionsBoundary(ctx, &iamv2.DeleteRolePermissionsBoundaryInput{
		RoleName: suite.roleName,
	})
	_, _ = suite.svc.DeleteRole(ctx, &iamv2.DeleteRoleInput{
		RoleName: suite.roleName,
	})
}

func (suite *TestIAMRolePermissionBoundarySuite) TestList() {
	a := assert.New(suite.T())

	lister := IAMRolePermissionBoundaryLister{}
	resources, err := lister.List(context.TODO(), testListerOpts)
	a.Nil(err)
	_ = resources
}

func TestIAMRolePermissionBoundaryIntegration(t *testing.T) {
	suite.Run(t, new(TestIAMRolePermissionBoundarySuite))
}
