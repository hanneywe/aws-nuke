package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockCloudFormationClient struct {
	mock.Mock
}

func (m *mockCloudFormationClient) ListGeneratedTemplates(
	ctx context.Context,
	params *cloudformation.ListGeneratedTemplatesInput,
	_ ...func(*cloudformation.Options),
) (*cloudformation.ListGeneratedTemplatesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cloudformation.ListGeneratedTemplatesOutput), args.Error(1)
}

func (m *mockCloudFormationClient) DeleteGeneratedTemplate(
	ctx context.Context,
	params *cloudformation.DeleteGeneratedTemplateInput,
	_ ...func(*cloudformation.Options),
) (*cloudformation.DeleteGeneratedTemplateOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cloudformation.DeleteGeneratedTemplateOutput), args.Error(1)
}

var testCloudFormationListerOpts = &nuke.ListerOpts{}
