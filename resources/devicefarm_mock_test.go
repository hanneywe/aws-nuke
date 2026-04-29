package resources

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/devicefarm"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

var testDeviceFarmListerOpts = &nuke.ListerOpts{
	Logger: logrus.WithField("test", true),
}

type mockDeviceFarmClient struct {
	mock.Mock
}

func (m *mockDeviceFarmClient) ListInstanceProfiles(ctx context.Context, params *devicefarm.ListInstanceProfilesInput,
	_ ...func(*devicefarm.Options)) (*devicefarm.ListInstanceProfilesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*devicefarm.ListInstanceProfilesOutput), args.Error(1)
}

func (m *mockDeviceFarmClient) ListTestGridProjects(ctx context.Context, params *devicefarm.ListTestGridProjectsInput,
	_ ...func(*devicefarm.Options)) (*devicefarm.ListTestGridProjectsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*devicefarm.ListTestGridProjectsOutput), args.Error(1)
}

func (m *mockDeviceFarmClient) ListProjects(ctx context.Context, params *devicefarm.ListProjectsInput,
	_ ...func(*devicefarm.Options)) (*devicefarm.ListProjectsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*devicefarm.ListProjectsOutput), args.Error(1)
}

func (m *mockDeviceFarmClient) ListUploads(ctx context.Context, params *devicefarm.ListUploadsInput,
	_ ...func(*devicefarm.Options)) (*devicefarm.ListUploadsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*devicefarm.ListUploadsOutput), args.Error(1)
}

func (m *mockDeviceFarmClient) DeleteInstanceProfile(ctx context.Context, params *devicefarm.DeleteInstanceProfileInput,
	_ ...func(*devicefarm.Options)) (*devicefarm.DeleteInstanceProfileOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*devicefarm.DeleteInstanceProfileOutput), args.Error(1)
}

func (m *mockDeviceFarmClient) DeleteTestGridProject(ctx context.Context, params *devicefarm.DeleteTestGridProjectInput,
	_ ...func(*devicefarm.Options)) (*devicefarm.DeleteTestGridProjectOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*devicefarm.DeleteTestGridProjectOutput), args.Error(1)
}

func (m *mockDeviceFarmClient) DeleteUpload(ctx context.Context, params *devicefarm.DeleteUploadInput,
	_ ...func(*devicefarm.Options)) (*devicefarm.DeleteUploadOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*devicefarm.DeleteUploadOutput), args.Error(1)
}
