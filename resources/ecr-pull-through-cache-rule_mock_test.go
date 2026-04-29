package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	ecrtypes "github.com/aws/aws-sdk-go-v2/service/ecr/types"
)

func Test_Mock_ECRPullThroughCacheRule_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockECRv2Client)
	mockClient.On("DescribePullThroughCacheRules", mock.Anything, mock.Anything).
		Return(&ecr.DescribePullThroughCacheRulesOutput{
			PullThroughCacheRules: []ecrtypes.PullThroughCacheRule{
				{
					EcrRepositoryPrefix: ptr.String("ecr-public"),
					UpstreamRegistryUrl: ptr.String("public.ecr.aws"),
				},
			},
		}, nil)
	lister := &ECRPullThroughCacheRuleLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testECRv2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	r := resources[0].(*ECRPullThroughCacheRule)
	a.Equal("ecr-public", *r.EcrRepositoryPrefix)
	a.Equal("public.ecr.aws", *r.UpstreamRegistryURL)
	mockClient.AssertExpectations(t)
}

func Test_Mock_ECRPullThroughCacheRule_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockECRv2Client)
	mockClient.On("DescribePullThroughCacheRules", mock.Anything, mock.Anything).
		Return(&ecr.DescribePullThroughCacheRulesOutput{
			PullThroughCacheRules: []ecrtypes.PullThroughCacheRule{},
		}, nil)
	lister := &ECRPullThroughCacheRuleLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testECRv2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_ECRPullThroughCacheRule_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockECRv2Client)
	r := &ECRPullThroughCacheRule{
		svc:                 mockClient,
		EcrRepositoryPrefix: ptr.String("ecr-public"),
	}
	mockClient.On("DeletePullThroughCacheRule", mock.Anything, &ecr.DeletePullThroughCacheRuleInput{
		EcrRepositoryPrefix: r.EcrRepositoryPrefix,
	}).Return(&ecr.DeletePullThroughCacheRuleOutput{}, nil)
	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_ECRPullThroughCacheRule_Properties(t *testing.T) {
	a := assert.New(t)
	r := ECRPullThroughCacheRule{
		EcrRepositoryPrefix: ptr.String("ecr-public"),
		UpstreamRegistryURL: ptr.String("public.ecr.aws"),
	}
	props := r.Properties()
	a.Equal("ecr-public", props.Get("EcrRepositoryPrefix"))
	a.Equal("public.ecr.aws", props.Get("UpstreamRegistryUrl"))
}

func Test_Mock_ECRPullThroughCacheRule_String(t *testing.T) {
	a := assert.New(t)
	r := &ECRPullThroughCacheRule{EcrRepositoryPrefix: ptr.String("ecr-public")}
	a.Equal("ecr-public", r.String())
}
