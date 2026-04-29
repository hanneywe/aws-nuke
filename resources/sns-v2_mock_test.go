package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/sns"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockSNSV2Client struct {
	mock.Mock
}

func (m *mockSNSV2Client) ListTopics(
	ctx context.Context, params *sns.ListTopicsInput,
	_ ...func(*sns.Options),
) (*sns.ListTopicsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sns.ListTopicsOutput), args.Error(1)
}

func (m *mockSNSV2Client) GetDataProtectionPolicy(
	ctx context.Context, params *sns.GetDataProtectionPolicyInput,
	_ ...func(*sns.Options),
) (*sns.GetDataProtectionPolicyOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sns.GetDataProtectionPolicyOutput), args.Error(1)
}

func (m *mockSNSV2Client) PutDataProtectionPolicy(
	ctx context.Context, params *sns.PutDataProtectionPolicyInput,
	_ ...func(*sns.Options),
) (*sns.PutDataProtectionPolicyOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sns.PutDataProtectionPolicyOutput), args.Error(1)
}

var testSNSV2ListerOpts = &nuke.ListerOpts{}
