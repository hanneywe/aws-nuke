package resources

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/route53"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

var testRoute53ListerOpts = &nuke.ListerOpts{
	Logger: logrus.WithField("test", true),
}

type mockRoute53Client struct {
	mock.Mock
}

func (m *mockRoute53Client) ListCidrCollections(ctx context.Context, params *route53.ListCidrCollectionsInput,
	_ ...func(*route53.Options)) (*route53.ListCidrCollectionsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*route53.ListCidrCollectionsOutput), args.Error(1)
}

func (m *mockRoute53Client) DeleteCidrCollection(ctx context.Context, params *route53.DeleteCidrCollectionInput,
	_ ...func(*route53.Options)) (*route53.DeleteCidrCollectionOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*route53.DeleteCidrCollectionOutput), args.Error(1)
}

func (m *mockRoute53Client) ListReusableDelegationSets(ctx context.Context, params *route53.ListReusableDelegationSetsInput,
	_ ...func(*route53.Options)) (*route53.ListReusableDelegationSetsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*route53.ListReusableDelegationSetsOutput), args.Error(1)
}

func (m *mockRoute53Client) DeleteReusableDelegationSet(ctx context.Context, params *route53.DeleteReusableDelegationSetInput,
	_ ...func(*route53.Options)) (*route53.DeleteReusableDelegationSetOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*route53.DeleteReusableDelegationSetOutput), args.Error(1)
}
