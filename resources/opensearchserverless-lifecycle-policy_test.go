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
	"github.com/aws/aws-sdk-go-v2/service/opensearchserverless"
	osstypes "github.com/aws/aws-sdk-go-v2/service/opensearchserverless/types"
)

type TestOpenSearchServerlessLifecyclePolicySuite struct {
	suite.Suite
	svc        *opensearchserverless.Client
	policyName *string
}

func (s *TestOpenSearchServerlessLifecyclePolicySuite) SetupSuite() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		s.T().Fatalf("failed to load config: %v", err)
	}

	s.svc = opensearchserverless.NewFromConfig(cfg)

	s.policyName = ptr.String(fmt.Sprintf("aws-nuke-test-%d", time.Now().UnixNano()))
	_, err = s.svc.CreateLifecyclePolicy(ctx, &opensearchserverless.CreateLifecyclePolicyInput{
		Name:   s.policyName,
		Type:   osstypes.LifecyclePolicyTypeRetention,
		Policy: ptr.String(`{"Rules":[{"ResourceType":"index","Resource":["index/test-collection/*"],"MinIndexRetention":"24h"}]}`),
	})
	if err != nil {
		s.T().Fatalf("failed to create lifecycle policy: %v", err)
	}
}

func (s *TestOpenSearchServerlessLifecyclePolicySuite) TearDownSuite() {
	ctx := context.TODO()
	if s.policyName != nil {
		_, _ = s.svc.DeleteLifecyclePolicy(ctx, &opensearchserverless.DeleteLifecyclePolicyInput{
			Name: s.policyName,
			Type: osstypes.LifecyclePolicyTypeRetention,
		})
	}
}

func (s *TestOpenSearchServerlessLifecyclePolicySuite) TestList() {
	a := assert.New(s.T())
	lister := &OpenSearchServerlessLifecyclePolicyLister{svc: s.svc}
	resources, err := lister.List(context.TODO(), testOpenSearchServerlessListerOpts)
	a.NoError(err)
	a.Greater(len(resources), 0)
}

func (s *TestOpenSearchServerlessLifecyclePolicySuite) TestRemove() {
	a := assert.New(s.T())
	policy := &OpenSearchServerlessLifecyclePolicy{
		svc:  s.svc,
		Name: s.policyName,
		Type: osstypes.LifecyclePolicyTypeRetention,
	}
	a.NoError(policy.Remove(context.TODO()))
}

func TestOpenSearchServerlessLifecyclePolicyIntegration(t *testing.T) {
	suite.Run(t, new(TestOpenSearchServerlessLifecyclePolicySuite))
}
