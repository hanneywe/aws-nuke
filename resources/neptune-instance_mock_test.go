package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go/service/neptune" //nolint:staticcheck

	libsettings "github.com/ekristen/libnuke/pkg/settings"
)

type mockNeptuneInstanceClient struct {
	mock.Mock
}

func (m *mockNeptuneInstanceClient) DescribeDBInstances(
	input *neptune.DescribeDBInstancesInput,
) (*neptune.DescribeDBInstancesOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*neptune.DescribeDBInstancesOutput), args.Error(1)
}

func (m *mockNeptuneInstanceClient) ListTagsForResource(
	input *neptune.ListTagsForResourceInput,
) (*neptune.ListTagsForResourceOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*neptune.ListTagsForResourceOutput), args.Error(1)
}

func (m *mockNeptuneInstanceClient) ModifyDBCluster(input *neptune.ModifyDBClusterInput) (*neptune.ModifyDBClusterOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*neptune.ModifyDBClusterOutput), args.Error(1)
}

func (m *mockNeptuneInstanceClient) ModifyDBInstance(input *neptune.ModifyDBInstanceInput) (*neptune.ModifyDBInstanceOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*neptune.ModifyDBInstanceOutput), args.Error(1)
}

func (m *mockNeptuneInstanceClient) DeleteDBInstance(input *neptune.DeleteDBInstanceInput) (*neptune.DeleteDBInstanceOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*neptune.DeleteDBInstanceOutput), args.Error(1)
}

func Test_Mock_NeptuneInstance_Remove_SkipDisableDeletionProtection_Clustered(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockNeptuneInstanceClient)

	settings := &libsettings.Setting{}
	settings.Set("DisableDeletionProtection", true)
	settings.Set("DisableClusterDeletionProtection", false)

	instance := &NeptuneInstance{
		svc:       mockClient,
		settings:  settings,
		ID:        ptr.String("my-instance"),
		ClusterID: ptr.String("my-cluster"),
	}

	// ModifyDBInstance should NOT be called because instance is part of a cluster
	mockClient.
		On("DeleteDBInstance", mock.Anything).
		Return(&neptune.DeleteDBInstanceOutput{}, nil)

	err := instance.Remove(context.TODO())
	a.NoError(err)

	// Verify ModifyDBInstance was never called
	mockClient.AssertNotCalled(t, "ModifyDBInstance", mock.Anything)
	mockClient.AssertExpectations(t)
}

func Test_Mock_NeptuneInstance_Remove_DisableDeletionProtection_Standalone(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockNeptuneInstanceClient)

	settings := &libsettings.Setting{}
	settings.Set("DisableDeletionProtection", true)
	settings.Set("DisableClusterDeletionProtection", false)

	instance := &NeptuneInstance{
		svc:      mockClient,
		settings: settings,
		ID:       ptr.String("my-standalone-instance"),
		// ClusterID is nil — standalone instance
	}

	// ModifyDBInstance SHOULD be called for standalone instances
	mockClient.
		On("ModifyDBInstance", &neptune.ModifyDBInstanceInput{
			DBInstanceIdentifier: ptr.String("my-standalone-instance"),
			DeletionProtection:   ptr.Bool(false),
		}).
		Return(&neptune.ModifyDBInstanceOutput{}, nil)

	mockClient.
		On("DeleteDBInstance", mock.Anything).
		Return(&neptune.DeleteDBInstanceOutput{}, nil)

	err := instance.Remove(context.TODO())
	a.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_NeptuneInstance_Remove_Basic(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockNeptuneInstanceClient)

	settings := &libsettings.Setting{}

	instance := &NeptuneInstance{
		svc:      mockClient,
		settings: settings,
		ID:       ptr.String("my-instance"),
	}

	mockClient.
		On("DeleteDBInstance", &neptune.DeleteDBInstanceInput{
			DBInstanceIdentifier: ptr.String("my-instance"),
			SkipFinalSnapshot:    ptr.Bool(true),
		}).
		Return(&neptune.DeleteDBInstanceOutput{}, nil)

	err := instance.Remove(context.TODO())
	a.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_NeptuneInstance_Properties(t *testing.T) {
	a := assert.New(t)

	instance := NeptuneInstance{
		ID:        ptr.String("my-instance"),
		ClusterID: ptr.String("my-cluster"),
		Name:      ptr.String("my-db"),
		Status:    ptr.String("available"),
	}

	props := instance.Properties()
	a.Equal("my-instance", props.Get("ID"))
	a.Equal("my-cluster", props.Get("ClusterID"))
	a.Equal("my-db", props.Get("Name"))
	a.Equal("available", props.Get("Status"))
}

func Test_Mock_NeptuneInstance_String(t *testing.T) {
	a := assert.New(t)

	instance := NeptuneInstance{
		ID: ptr.String("my-instance"),
	}

	a.Equal("my-instance", instance.String())
}

func Test_Mock_NeptuneInstance_Filter_Deleting(t *testing.T) {
	a := assert.New(t)

	instance := NeptuneInstance{
		Status: ptr.String("deleting"),
	}

	err := instance.Filter()
	a.Error(err)
	a.Contains(err.Error(), "already deleting")
}

func Test_Mock_NeptuneInstance_Filter_Available(t *testing.T) {
	a := assert.New(t)

	instance := NeptuneInstance{
		Status: ptr.String("available"),
	}

	err := instance.Filter()
	a.NoError(err)
}
