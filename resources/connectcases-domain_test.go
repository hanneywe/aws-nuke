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
	"github.com/aws/aws-sdk-go-v2/service/connectcases"
)

type TestConnectCasesDomainSuite struct {
	suite.Suite
	svc      *connectcases.Client
	domainId *string
}

func (suite *TestConnectCasesDomainSuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		suite.T().Fatalf("failed to load config: %v", err)
	}

	suite.svc = connectcases.NewFromConfig(cfg)

	name := fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano())
	resp, err := suite.svc.CreateDomain(ctx, &connectcases.CreateDomainInput{
		Name: &name,
	})
	if err != nil {
		suite.T().Fatalf("failed to create test domain: %v", err)
	}
	suite.domainId = resp.DomainId
}

func (suite *TestConnectCasesDomainSuite) TearDownSuite() {
	ctx := context.TODO()
	if suite.domainId != nil {
		_, _ = suite.svc.DeleteDomain(ctx, &connectcases.DeleteDomainInput{
			DomainId: suite.domainId,
		})
	}
}

func (suite *TestConnectCasesDomainSuite) TestList() {
	a := assert.New(suite.T())

	lister := ConnectCasesDomainLister{svc: suite.svc}
	resources, err := lister.List(context.TODO(), testConnectCasesListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (suite *TestConnectCasesDomainSuite) TestRemove() {
	a := assert.New(suite.T())

	resource := ConnectCasesDomain{
		svc:      suite.svc,
		DomainId: suite.domainId,
	}

	err := resource.Remove(context.TODO())
	a.NoError(err)
	suite.domainId = nil
}

func TestConnectCasesDomainIntegration(t *testing.T) {
	suite.Run(t, new(TestConnectCasesDomainSuite))
}
