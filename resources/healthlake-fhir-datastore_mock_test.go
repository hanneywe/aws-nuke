package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/healthlake"
	healthlaketypes "github.com/aws/aws-sdk-go-v2/service/healthlake/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

var testHealthLakeListerOpts = &nuke.ListerOpts{
	Logger: logrus.WithField("test", true),
}

func Test_Mock_HealthLakeFHIRDatastore_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockHealthLakeClient)

	mockClient.
		On("ListFHIRDatastores", mock.Anything, mock.Anything).
		Return(
			&healthlake.ListFHIRDatastoresOutput{
				DatastorePropertiesList: []healthlaketypes.DatastoreProperties{
					{
						DatastoreId:     ptr.String("ds-12345"),
						DatastoreName:   ptr.String("test-datastore"),
						DatastoreStatus: healthlaketypes.DatastoreStatusActive,
					},
				},
			}, nil,
		)

	lister := &HealthLakeFHIRDatastoreLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testHealthLakeListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	datastore := resources[0].(*HealthLakeFHIRDatastore)
	assertions.Equal("ds-12345", *datastore.DatastoreID)
	assertions.Equal("test-datastore", *datastore.DatastoreName)
	assertions.Equal("ACTIVE", *datastore.Status)

	mockClient.AssertExpectations(t)
}

func Test_Mock_HealthLakeFHIRDatastore_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockHealthLakeClient)

	mockClient.
		On("ListFHIRDatastores", mock.Anything, mock.Anything).
		Return(
			&healthlake.ListFHIRDatastoresOutput{
				DatastorePropertiesList: []healthlaketypes.DatastoreProperties{},
			}, nil,
		)

	lister := &HealthLakeFHIRDatastoreLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testHealthLakeListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_HealthLakeFHIRDatastore_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockHealthLakeClient)

	datastore := &HealthLakeFHIRDatastore{
		svc:           mockClient,
		DatastoreID:   ptr.String("ds-12345"),
		DatastoreName: ptr.String("test-datastore"),
		Status:        ptr.String("ACTIVE"),
	}

	mockClient.
		On(
			"DeleteFHIRDatastore",
			mock.Anything,
			&healthlake.DeleteFHIRDatastoreInput{
				DatastoreId: datastore.DatastoreID,
			},
		).
		Return(&healthlake.DeleteFHIRDatastoreOutput{}, nil)

	err := datastore.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_HealthLakeFHIRDatastore_Properties(t *testing.T) {
	assertions := assert.New(t)

	datastore := HealthLakeFHIRDatastore{
		DatastoreID:   ptr.String("ds-12345"),
		DatastoreName: ptr.String("test-datastore"),
		Status:        ptr.String("ACTIVE"),
	}

	properties := datastore.Properties()

	assertions.Equal("ds-12345", properties.Get("DatastoreID"))
	assertions.Equal("test-datastore", properties.Get("DatastoreName"))
	assertions.Equal("ACTIVE", properties.Get("Status"))
}

func Test_Mock_HealthLakeFHIRDatastore_String(t *testing.T) {
	assertions := assert.New(t)

	datastore := HealthLakeFHIRDatastore{
		DatastoreID:   ptr.String("ds-12345"),
		DatastoreName: ptr.String("test-datastore"),
	}

	assertions.Equal("ds-12345", datastore.String())
}

func Test_Mock_HealthLakeFHIRDatastore_Filter(t *testing.T) {
	assertions := assert.New(t)

	// DELETING status should be filtered
	deletingDatastore := HealthLakeFHIRDatastore{
		DatastoreID:   ptr.String("ds-deleting"),
		DatastoreName: ptr.String("deleting-datastore"),
		Status:        ptr.String("DELETING"),
	}
	err := deletingDatastore.Filter()
	assertions.Error(err)
	assertions.Contains(err.Error(), "DELETING")

	// DELETED status should be filtered
	deletedDatastore := HealthLakeFHIRDatastore{
		DatastoreID:   ptr.String("ds-deleted"),
		DatastoreName: ptr.String("deleted-datastore"),
		Status:        ptr.String("DELETED"),
	}
	err = deletedDatastore.Filter()
	assertions.Error(err)
	assertions.Contains(err.Error(), "DELETED")

	// ACTIVE status should not be filtered
	activeDatastore := HealthLakeFHIRDatastore{
		DatastoreID:   ptr.String("ds-active"),
		DatastoreName: ptr.String("active-datastore"),
		Status:        ptr.String("ACTIVE"),
	}
	err = activeDatastore.Filter()
	assertions.NoError(err)
}
