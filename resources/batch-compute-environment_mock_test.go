package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/batch"
	batchtypes "github.com/aws/aws-sdk-go-v2/service/batch/types"
)

type mockBatchComputeEnvironmentClient struct {
	mock.Mock
}

func (m *mockBatchComputeEnvironmentClient) DescribeComputeEnvironments(
	ctx context.Context, params *batch.DescribeComputeEnvironmentsInput,
	_ ...func(*batch.Options)) (*batch.DescribeComputeEnvironmentsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*batch.DescribeComputeEnvironmentsOutput), args.Error(1)
}

func (m *mockBatchComputeEnvironmentClient) DeleteComputeEnvironment(
	ctx context.Context, params *batch.DeleteComputeEnvironmentInput,
	_ ...func(*batch.Options)) (*batch.DeleteComputeEnvironmentOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*batch.DeleteComputeEnvironmentOutput), args.Error(1)
}

func (m *mockBatchComputeEnvironmentClient) UpdateComputeEnvironment(
	ctx context.Context, params *batch.UpdateComputeEnvironmentInput,
	_ ...func(*batch.Options)) (*batch.UpdateComputeEnvironmentOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*batch.UpdateComputeEnvironmentOutput), args.Error(1)
}

// --- BatchComputeEnvironment Tests ---

func Test_Mock_BatchComputeEnvironment_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockBatchComputeEnvironmentClient)

	mockClient.
		On("DescribeComputeEnvironments", mock.Anything, mock.Anything).
		Return(&batch.DescribeComputeEnvironmentsOutput{
			ComputeEnvironments: []batchtypes.ComputeEnvironmentDetail{
				{
					ComputeEnvironmentName: ptr.String("test-ce-1"),
					Status:                 batchtypes.CEStatusValid,
					State:                  batchtypes.CEStateEnabled,
					Type:                   batchtypes.CETypeManaged,
					Tags:                   map[string]string{"env": "test"},
				},
				{
					ComputeEnvironmentName: ptr.String("test-ce-deleted"),
					Status:                 batchtypes.CEStatusDeleted,
					State:                  batchtypes.CEStateDisabled,
					Type:                   batchtypes.CETypeManaged,
				},
			},
		}, nil)

	lister := &BatchComputeEnvironmentLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testBatchListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	ce := resources[0].(*BatchComputeEnvironment)
	a.Equal("test-ce-1", *ce.Name)
	a.Equal("VALID", ce.Status)
	a.Equal("ENABLED", ce.State)

	mockClient.AssertExpectations(t)
}

func Test_Mock_BatchComputeEnvironment_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockBatchComputeEnvironmentClient)

	ce := &BatchComputeEnvironment{
		svc:    mockClient,
		Name:   ptr.String("test-ce-1"),
		Status: "VALID",
		State:  "DISABLED",
	}

	mockClient.
		On("DeleteComputeEnvironment", mock.Anything, &batch.DeleteComputeEnvironmentInput{
			ComputeEnvironment: ce.Name,
		}).
		Return(&batch.DeleteComputeEnvironmentOutput{}, nil)

	err := ce.Remove(context.TODO())
	a.NoError(err)
	mockClient.AssertExpectations(t)
}

func Test_Mock_BatchComputeEnvironment_Remove_Enabled(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockBatchComputeEnvironmentClient)

	ce := &BatchComputeEnvironment{
		svc:    mockClient,
		Name:   ptr.String("test-ce-enabled"),
		Status: "VALID",
		State:  "ENABLED",
	}

	mockClient.
		On("UpdateComputeEnvironment", mock.Anything, &batch.UpdateComputeEnvironmentInput{
			ComputeEnvironment: ce.Name,
			State:              batchtypes.CEStateDisabled,
		}).
		Return(&batch.UpdateComputeEnvironmentOutput{}, nil)

	mockClient.
		On("DeleteComputeEnvironment", mock.Anything, &batch.DeleteComputeEnvironmentInput{
			ComputeEnvironment: ce.Name,
		}).
		Return(&batch.DeleteComputeEnvironmentOutput{}, nil)

	err := ce.Remove(context.TODO())
	a.NoError(err)
	mockClient.AssertExpectations(t)
}

func Test_Mock_BatchComputeEnvironment_Filter_Deleting(t *testing.T) {
	a := assert.New(t)

	ce := &BatchComputeEnvironment{
		Name:   ptr.String("test-ce-deleting"),
		Status: string(batchtypes.CEStatusDeleting),
	}

	err := ce.Filter()
	a.Error(err)
	a.Contains(err.Error(), "already deleting")
}

func Test_Mock_BatchComputeEnvironment_Filter_Valid(t *testing.T) {
	a := assert.New(t)

	ce := &BatchComputeEnvironment{
		Name:   ptr.String("test-ce-valid"),
		Status: string(batchtypes.CEStatusValid),
	}

	err := ce.Filter()
	a.NoError(err)
}

func Test_Mock_BatchComputeEnvironment_Filter_Invalid(t *testing.T) {
	a := assert.New(t)

	ce := &BatchComputeEnvironment{
		Name:   ptr.String("test-ce-invalid"),
		Status: string(batchtypes.CEStatusInvalid),
	}

	err := ce.Filter()
	a.NoError(err)
}

func Test_Mock_BatchComputeEnvironment_Properties(t *testing.T) {
	a := assert.New(t)

	ce := &BatchComputeEnvironment{
		Name:   ptr.String("test-ce-1"),
		Status: "VALID",
		State:  "ENABLED",
		Type:   "MANAGED",
		Tags:   map[string]string{"env": "test"},
	}

	props := ce.Properties()
	a.Equal("test-ce-1", props.Get("Name"))
	a.Equal("VALID", props.Get("Status"))
	a.Equal("ENABLED", props.Get("State"))
	a.Equal("MANAGED", props.Get("Type"))
	a.Equal("test", props.Get("tag:env"))
}

func Test_Mock_BatchComputeEnvironment_String(t *testing.T) {
	a := assert.New(t)

	ce := &BatchComputeEnvironment{
		Name: ptr.String("test-ce-1"),
	}

	a.Equal("test-ce-1", ce.String())
}

func Test_Mock_BatchComputeEnvironment_HandleWait_Gone(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockBatchComputeEnvironmentClient)

	ce := &BatchComputeEnvironment{
		svc:  mockClient,
		Name: ptr.String("test-ce-1"),
	}

	mockClient.
		On("DescribeComputeEnvironments", mock.Anything, &batch.DescribeComputeEnvironmentsInput{
			ComputeEnvironments: []string{"test-ce-1"},
		}).
		Return(&batch.DescribeComputeEnvironmentsOutput{
			ComputeEnvironments: []batchtypes.ComputeEnvironmentDetail{},
		}, nil)

	err := ce.HandleWait(context.TODO())
	a.NoError(err)
	mockClient.AssertExpectations(t)
}

func Test_Mock_BatchComputeEnvironment_HandleWait_Deleted(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockBatchComputeEnvironmentClient)

	ce := &BatchComputeEnvironment{
		svc:  mockClient,
		Name: ptr.String("test-ce-1"),
	}

	mockClient.
		On("DescribeComputeEnvironments", mock.Anything, &batch.DescribeComputeEnvironmentsInput{
			ComputeEnvironments: []string{"test-ce-1"},
		}).
		Return(&batch.DescribeComputeEnvironmentsOutput{
			ComputeEnvironments: []batchtypes.ComputeEnvironmentDetail{
				{
					ComputeEnvironmentName: ptr.String("test-ce-1"),
					Status:                 batchtypes.CEStatusDeleted,
				},
			},
		}, nil)

	err := ce.HandleWait(context.TODO())
	a.NoError(err)
	mockClient.AssertExpectations(t)
}

func Test_Mock_BatchComputeEnvironment_HandleWait_Deleting(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockBatchComputeEnvironmentClient)

	ce := &BatchComputeEnvironment{
		svc:  mockClient,
		Name: ptr.String("test-ce-1"),
	}

	mockClient.
		On("DescribeComputeEnvironments", mock.Anything, &batch.DescribeComputeEnvironmentsInput{
			ComputeEnvironments: []string{"test-ce-1"},
		}).
		Return(&batch.DescribeComputeEnvironmentsOutput{
			ComputeEnvironments: []batchtypes.ComputeEnvironmentDetail{
				{
					ComputeEnvironmentName: ptr.String("test-ce-1"),
					Status:                 batchtypes.CEStatusDeleting,
				},
			},
		}, nil)

	err := ce.HandleWait(context.TODO())
	a.Error(err)
	a.Contains(err.Error(), "waiting for compute environment to delete")
	mockClient.AssertExpectations(t)
}

func Test_Mock_BatchComputeEnvironment_HandleWait_Invalid(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockBatchComputeEnvironmentClient)

	ce := &BatchComputeEnvironment{
		svc:  mockClient,
		Name: ptr.String("test-ce-invalid"),
	}

	mockClient.
		On("DescribeComputeEnvironments", mock.Anything, &batch.DescribeComputeEnvironmentsInput{
			ComputeEnvironments: []string{"test-ce-invalid"},
		}).
		Return(&batch.DescribeComputeEnvironmentsOutput{
			ComputeEnvironments: []batchtypes.ComputeEnvironmentDetail{
				{
					ComputeEnvironmentName: ptr.String("test-ce-invalid"),
					Status:                 batchtypes.CEStatusInvalid,
				},
			},
		}, nil)

	mockClient.
		On("DeleteComputeEnvironment", mock.Anything, &batch.DeleteComputeEnvironmentInput{
			ComputeEnvironment: ptr.String("test-ce-invalid"),
		}).
		Return(&batch.DeleteComputeEnvironmentOutput{}, nil)

	err := ce.HandleWait(context.TODO())
	a.Error(err)
	a.Contains(err.Error(), "retrying delete")
	mockClient.AssertExpectations(t)
}

func Test_Mock_BatchComputeEnvironment_HandleWait_Valid(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockBatchComputeEnvironmentClient)

	ce := &BatchComputeEnvironment{
		svc:  mockClient,
		Name: ptr.String("test-ce-valid"),
	}

	mockClient.
		On("DescribeComputeEnvironments", mock.Anything, &batch.DescribeComputeEnvironmentsInput{
			ComputeEnvironments: []string{"test-ce-valid"},
		}).
		Return(&batch.DescribeComputeEnvironmentsOutput{
			ComputeEnvironments: []batchtypes.ComputeEnvironmentDetail{
				{
					ComputeEnvironmentName: ptr.String("test-ce-valid"),
					Status:                 batchtypes.CEStatusValid,
				},
			},
		}, nil)

	err := ce.HandleWait(context.TODO())
	a.Error(err)
	a.Contains(err.Error(), "waiting for compute environment to transition")
	mockClient.AssertExpectations(t)
}

// --- BatchComputeEnvironmentState Tests ---

func Test_Mock_BatchComputeEnvironmentState_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockBatchComputeEnvironmentClient)

	mockClient.
		On("DescribeComputeEnvironments", mock.Anything, mock.Anything).
		Return(&batch.DescribeComputeEnvironmentsOutput{
			ComputeEnvironments: []batchtypes.ComputeEnvironmentDetail{
				{
					ComputeEnvironmentName: ptr.String("test-ce-1"),
					Status:                 batchtypes.CEStatusValid,
					State:                  batchtypes.CEStateEnabled,
					Type:                   batchtypes.CETypeManaged,
				},
			},
		}, nil)

	lister := &BatchComputeEnvironmentStateLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testBatchListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	mockClient.AssertExpectations(t)
}

func Test_Mock_BatchComputeEnvironmentState_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockBatchComputeEnvironmentClient)

	state := &BatchComputeEnvironmentState{
		svc:    mockClient,
		Name:   ptr.String("test-ce-1"),
		State:  "ENABLED",
		Status: "VALID",
	}

	mockClient.
		On("UpdateComputeEnvironment", mock.Anything, &batch.UpdateComputeEnvironmentInput{
			ComputeEnvironment: state.Name,
			State:              batchtypes.CEStateDisabled,
		}).
		Return(&batch.UpdateComputeEnvironmentOutput{}, nil)

	err := state.Remove(context.TODO())
	a.NoError(err)
	mockClient.AssertExpectations(t)
}

func Test_Mock_BatchComputeEnvironmentState_Filter_AlreadyDisabled(t *testing.T) {
	a := assert.New(t)

	state := &BatchComputeEnvironmentState{
		Name:   ptr.String("test-ce-1"),
		State:  string(batchtypes.CEStateDisabled),
		Status: string(batchtypes.CEStatusValid),
	}

	err := state.Filter()
	a.Error(err)
	a.Contains(err.Error(), "already disabled")
}

func Test_Mock_BatchComputeEnvironmentState_Filter_Deleting(t *testing.T) {
	a := assert.New(t)

	state := &BatchComputeEnvironmentState{
		Name:   ptr.String("test-ce-1"),
		State:  string(batchtypes.CEStateEnabled),
		Status: string(batchtypes.CEStatusDeleting),
	}

	err := state.Filter()
	a.Error(err)
	a.Contains(err.Error(), "being deleted")
}

func Test_Mock_BatchComputeEnvironmentState_Filter_Enabled(t *testing.T) {
	a := assert.New(t)

	state := &BatchComputeEnvironmentState{
		Name:   ptr.String("test-ce-1"),
		State:  string(batchtypes.CEStateEnabled),
		Status: string(batchtypes.CEStatusValid),
	}

	err := state.Filter()
	a.NoError(err)
}

func Test_Mock_BatchComputeEnvironmentState_Properties(t *testing.T) {
	a := assert.New(t)

	state := &BatchComputeEnvironmentState{
		Name:   ptr.String("test-ce-1"),
		State:  "ENABLED",
		Status: "VALID",
		Type:   "MANAGED",
		Tags:   map[string]string{"env": "test"},
	}

	props := state.Properties()
	a.Equal("test-ce-1", props.Get("Name"))
	a.Equal("ENABLED", props.Get("State"))
	a.Equal("VALID", props.Get("Status"))
	a.Equal("test", props.Get("tag:env"))
}
