package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/lightsail"
	lightsailtypes "github.com/aws/aws-sdk-go-v2/service/lightsail/types"
)

func Test_Mock_LightsailDiskSnapshot_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockLightsailClient)

	mockClient.On("GetDiskSnapshots", mock.Anything, mock.Anything).
		Return(&lightsail.GetDiskSnapshotsOutput{
			DiskSnapshots: []lightsailtypes.DiskSnapshot{
				{Name: ptr.String("test-value"), FromDiskName: ptr.String("test-value"), SizeInGb: ptr.Int32(10)},
			},
		}, nil)

	lister := &LightsailDiskSnapshotLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testLightsailListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*LightsailDiskSnapshot)
	a.Equal("test-value", *r.DiskSnapshotName)
	mockClient.AssertExpectations(t)
}

func Test_Mock_LightsailDiskSnapshot_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockLightsailClient)

	mockClient.On("GetDiskSnapshots", mock.Anything, mock.Anything).
		Return(&lightsail.GetDiskSnapshotsOutput{
			DiskSnapshots: []lightsailtypes.DiskSnapshot{},
		}, nil)

	lister := &LightsailDiskSnapshotLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testLightsailListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_LightsailDiskSnapshot_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockLightsailClient)

	r := &LightsailDiskSnapshot{
		svc:              mockClient,
		DiskSnapshotName: ptr.String("test-disksnapshotname"),
	}

	mockClient.On("DeleteDiskSnapshot", mock.Anything,
		&lightsail.DeleteDiskSnapshotInput{
			DiskSnapshotName: r.DiskSnapshotName,
		}).Return(&lightsail.DeleteDiskSnapshotOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_LightsailDiskSnapshot_Properties(t *testing.T) {
	a := assert.New(t)
	r := &LightsailDiskSnapshot{
		DiskSnapshotName: ptr.String("test-disksnapshotname"),
		FromDiskName:     ptr.String("test-fromdiskname"),
	}
	props := r.Properties()
	a.Equal("test-disksnapshotname", props.Get("DiskSnapshotName"))
	a.Equal("test-fromdiskname", props.Get("FromDiskName"))
}

func Test_Mock_LightsailDiskSnapshot_String(t *testing.T) {
	a := assert.New(t)
	r := &LightsailDiskSnapshot{
		DiskSnapshotName: ptr.String("test-disksnapshotname"),
	}
	a.Equal("test-disksnapshotname", r.String())
}
