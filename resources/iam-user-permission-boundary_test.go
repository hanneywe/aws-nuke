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

type TestIAMUserPermissionBoundarySuite struct {
	suite.Suite
	svc      *iamv2.Client
	userName *string
}

func (suite *TestIAMUserPermissionBoundarySuite) SetupSuite() {
	suite.userName = ptr.String(fmt.Sprintf("aws-nuke-test-pb-%d", time.Now().UnixNano()))

	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}

	suite.svc = iamv2.NewFromConfig(cfg)
}

func (suite *TestIAMUserPermissionBoundarySuite) TearDownSuite() {
	ctx := context.TODO()
	_, _ = suite.svc.DeleteUserPermissionsBoundary(ctx, &iamv2.DeleteUserPermissionsBoundaryInput{
		UserName: suite.userName,
	})
	_, _ = suite.svc.DeleteUser(ctx, &iamv2.DeleteUserInput{
		UserName: suite.userName,
	})
}

func (suite *TestIAMUserPermissionBoundarySuite) TestList() {
	a := assert.New(suite.T())

	lister := IAMUserPermissionBoundaryLister{}
	resources, err := lister.List(context.TODO(), testListerOpts)
	a.Nil(err)
	_ = resources
}

func TestIAMUserPermissionBoundaryIntegration(t *testing.T) {
	suite.Run(t, new(TestIAMUserPermissionBoundarySuite))
}
