package resources

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/omics"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

var testOmicsListerOpts = &nuke.ListerOpts{
	Logger: logrus.WithField("test", true),
}

type mockOmicsClient struct {
	mock.Mock
}

func (m *mockOmicsClient) ListReferenceStores(ctx context.Context, params *omics.ListReferenceStoresInput,
	_ ...func(*omics.Options)) (*omics.ListReferenceStoresOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*omics.ListReferenceStoresOutput), args.Error(1)
}

func (m *mockOmicsClient) ListReferences(ctx context.Context, params *omics.ListReferencesInput,
	_ ...func(*omics.Options)) (*omics.ListReferencesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*omics.ListReferencesOutput), args.Error(1)
}

func (m *mockOmicsClient) ListRunGroups(ctx context.Context, params *omics.ListRunGroupsInput,
	_ ...func(*omics.Options)) (*omics.ListRunGroupsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*omics.ListRunGroupsOutput), args.Error(1)
}

func (m *mockOmicsClient) ListSequenceStores(ctx context.Context, params *omics.ListSequenceStoresInput,
	_ ...func(*omics.Options)) (*omics.ListSequenceStoresOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*omics.ListSequenceStoresOutput), args.Error(1)
}

func (m *mockOmicsClient) ListReadSets(ctx context.Context, params *omics.ListReadSetsInput,
	_ ...func(*omics.Options)) (*omics.ListReadSetsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*omics.ListReadSetsOutput), args.Error(1)
}

func (m *mockOmicsClient) ListWorkflows(ctx context.Context, params *omics.ListWorkflowsInput,
	_ ...func(*omics.Options)) (*omics.ListWorkflowsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*omics.ListWorkflowsOutput), args.Error(1)
}

func (m *mockOmicsClient) DeleteReferenceStore(ctx context.Context, params *omics.DeleteReferenceStoreInput,
	_ ...func(*omics.Options)) (*omics.DeleteReferenceStoreOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*omics.DeleteReferenceStoreOutput), args.Error(1)
}

func (m *mockOmicsClient) DeleteReference(ctx context.Context, params *omics.DeleteReferenceInput,
	_ ...func(*omics.Options)) (*omics.DeleteReferenceOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*omics.DeleteReferenceOutput), args.Error(1)
}

func (m *mockOmicsClient) DeleteRunGroup(ctx context.Context, params *omics.DeleteRunGroupInput,
	_ ...func(*omics.Options)) (*omics.DeleteRunGroupOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*omics.DeleteRunGroupOutput), args.Error(1)
}

func (m *mockOmicsClient) DeleteSequenceStore(ctx context.Context, params *omics.DeleteSequenceStoreInput,
	_ ...func(*omics.Options)) (*omics.DeleteSequenceStoreOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*omics.DeleteSequenceStoreOutput), args.Error(1)
}

func (m *mockOmicsClient) BatchDeleteReadSet(ctx context.Context, params *omics.BatchDeleteReadSetInput,
	_ ...func(*omics.Options)) (*omics.BatchDeleteReadSetOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*omics.BatchDeleteReadSetOutput), args.Error(1)
}

func (m *mockOmicsClient) DeleteWorkflow(ctx context.Context, params *omics.DeleteWorkflowInput,
	_ ...func(*omics.Options)) (*omics.DeleteWorkflowOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*omics.DeleteWorkflowOutput), args.Error(1)
}
