package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/datasync"
	datasynctypes "github.com/aws/aws-sdk-go-v2/service/datasync/types"
)

func Test_Mock_DataSyncLocation_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockDataSyncClient)

	mockClient.
		On("ListLocations", mock.Anything, mock.Anything).
		Return(&datasync.ListLocationsOutput{
			Locations: []datasynctypes.LocationListEntry{
				{
					LocationArn: ptr.String("arn:aws:datasync:us-east-1:123456789012:location/loc-12345"),
					LocationUri: ptr.String("s3://my-bucket"),
				},
			},
		}, nil)

	lister := &DataSyncLocationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testDataSyncListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	loc := resources[0].(*DataSyncLocation)
	a.Equal("arn:aws:datasync:us-east-1:123456789012:location/loc-12345", *loc.LocationArn)
	a.Equal("s3://my-bucket", *loc.LocationURI)

	mockClient.AssertExpectations(t)
}

func Test_Mock_DataSyncLocation_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockDataSyncClient)

	mockClient.
		On("ListLocations", mock.Anything, mock.Anything).
		Return(&datasync.ListLocationsOutput{
			Locations: []datasynctypes.LocationListEntry{},
		}, nil)

	lister := &DataSyncLocationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testDataSyncListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_DataSyncLocation_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockDataSyncClient)

	loc := &DataSyncLocation{
		svc:         mockClient,
		LocationArn: ptr.String("arn:aws:datasync:us-east-1:123456789012:location/loc-12345"),
	}

	mockClient.
		On("DeleteLocation", mock.Anything, &datasync.DeleteLocationInput{
			LocationArn: loc.LocationArn,
		}).
		Return(&datasync.DeleteLocationOutput{}, nil)

	err := loc.Remove(context.TODO())
	a.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_DataSyncLocation_Properties(t *testing.T) {
	a := assert.New(t)

	loc := DataSyncLocation{
		LocationArn: ptr.String("arn:aws:datasync:us-east-1:123456789012:location/loc-12345"),
		LocationURI: ptr.String("s3://my-bucket"),
	}

	props := loc.Properties()
	a.Equal("arn:aws:datasync:us-east-1:123456789012:location/loc-12345", props.Get("LocationArn"))
	a.Equal("s3://my-bucket", props.Get("LocationUri"))
}

func Test_Mock_DataSyncLocation_String(t *testing.T) {
	a := assert.New(t)

	loc := DataSyncLocation{
		LocationArn: ptr.String("arn:aws:datasync:us-east-1:123456789012:location/loc-12345"),
	}

	a.Equal("arn:aws:datasync:us-east-1:123456789012:location/loc-12345", loc.String())
}
