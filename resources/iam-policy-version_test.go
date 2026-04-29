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

type TestIAMPolicyVersionSuite struct {
	suite.Suite
	svc        *iamv2.Client
	policyName *string
	policyArn  *string
}

func (suite *TestIAMPolicyVersionSuite) SetupSuite() {
	suite.policyName = ptr.String(fmt.Sprintf("aws-nuke-test-pv-%d", time.Now().UnixNano()))

	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}

	suite.svc = iamv2.NewFromConfig(cfg)
}

func (suite *TestIAMPolicyVersionSuite) TearDownSuite() {
	ctx := context.TODO()
	if suite.policyArn != nil {
		_, _ = suite.svc.DeletePolicy(ctx, &iamv2.DeletePolicyInput{
			PolicyArn: suite.policyArn,
		})
	}
}

func (suite *TestIAMPolicyVersionSuite) TestList() {
	a := assert.New(suite.T())

	lister := IAMPolicyVersionLister{}
	resources, err := lister.List(context.TODO(), testListerOpts)
	a.Nil(err)
	_ = resources
}

func TestIAMPolicyVersionIntegration(t *testing.T) {
	suite.Run(t, new(TestIAMPolicyVersionSuite))
}
