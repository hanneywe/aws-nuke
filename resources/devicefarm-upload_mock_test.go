package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/devicefarm"
	devicefarmtypes "github.com/aws/aws-sdk-go-v2/service/devicefarm/types"
)

func Test_Mock_DeviceFarmUpload_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockDeviceFarmClient)

	mockClient.
		On("ListProjects", mock.Anything, mock.Anything).
		Return(&devicefarm.ListProjectsOutput{
			Projects: []devicefarmtypes.Project{
				{
					Arn:  ptr.String("arn:aws:devicefarm:us-west-2:123456789012:project:project-1"),
					Name: ptr.String("my-project"),
				},
			},
		}, nil)

	mockClient.
		On("ListUploads", mock.Anything, mock.Anything).
		Return(&devicefarm.ListUploadsOutput{
			Uploads: []devicefarmtypes.Upload{
				{
					Arn:  ptr.String("arn:aws:devicefarm:us-west-2:123456789012:upload:upload-1"),
					Name: ptr.String("my-upload.apk"),
				},
			},
		}, nil)

	lister := &DeviceFarmUploadLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testDeviceFarmListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	upload := resources[0].(*DeviceFarmUpload)
	assertions.Equal("arn:aws:devicefarm:us-west-2:123456789012:upload:upload-1", *upload.Arn)
	assertions.Equal("my-upload.apk", *upload.Name)

	mockClient.AssertExpectations(t)
}

func Test_Mock_DeviceFarmUpload_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockDeviceFarmClient)

	mockClient.
		On("ListProjects", mock.Anything, mock.Anything).
		Return(&devicefarm.ListProjectsOutput{
			Projects: []devicefarmtypes.Project{},
		}, nil)

	lister := &DeviceFarmUploadLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testDeviceFarmListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_DeviceFarmUpload_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockDeviceFarmClient)

	upload := &DeviceFarmUpload{
		svc:  mockClient,
		Arn:  ptr.String("arn:aws:devicefarm:us-west-2:123456789012:upload:upload-1"),
		Name: ptr.String("my-upload.apk"),
	}

	mockClient.
		On("DeleteUpload", mock.Anything, &devicefarm.DeleteUploadInput{
			Arn: upload.Arn,
		}).
		Return(&devicefarm.DeleteUploadOutput{}, nil)

	err := upload.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_DeviceFarmUpload_Properties(t *testing.T) {
	assertions := assert.New(t)

	upload := DeviceFarmUpload{
		Arn:  ptr.String("arn:aws:devicefarm:us-west-2:123456789012:upload:upload-1"),
		Name: ptr.String("my-upload.apk"),
	}

	properties := upload.Properties()

	assertions.Equal("arn:aws:devicefarm:us-west-2:123456789012:upload:upload-1", properties.Get("Arn"))
	assertions.Equal("my-upload.apk", properties.Get("Name"))
}

func Test_Mock_DeviceFarmUpload_String(t *testing.T) {
	assertions := assert.New(t)

	upload := DeviceFarmUpload{
		Name: ptr.String("my-upload.apk"),
	}

	assertions.Equal("my-upload.apk", upload.String())
}
