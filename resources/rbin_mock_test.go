package resources

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/rbin"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

var testRbinListerOpts = &nuke.ListerOpts{
	Logger: logrus.WithField("test", true),
}

type mockRbinClient struct {
	mock.Mock
}

func (m *mockRbinClient) ListRules(ctx context.Context, params *rbin.ListRulesInput,
	_ ...func(*rbin.Options)) (*rbin.ListRulesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*rbin.ListRulesOutput), args.Error(1)
}

func (m *mockRbinClient) GetRule(ctx context.Context, params *rbin.GetRuleInput,
	_ ...func(*rbin.Options)) (*rbin.GetRuleOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*rbin.GetRuleOutput), args.Error(1)
}

func (m *mockRbinClient) DeleteRule(ctx context.Context, params *rbin.DeleteRuleInput,
	_ ...func(*rbin.Options)) (*rbin.DeleteRuleOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*rbin.DeleteRuleOutput), args.Error(1)
}

func (m *mockRbinClient) UnlockRule(ctx context.Context, params *rbin.UnlockRuleInput,
	_ ...func(*rbin.Options)) (*rbin.UnlockRuleOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*rbin.UnlockRuleOutput), args.Error(1)
}
