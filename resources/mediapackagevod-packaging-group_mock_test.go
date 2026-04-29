package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/mediapackagevod"
	mediapackagevodtypes "github.com/aws/aws-sdk-go-v2/service/mediapackagevod/types"
)

func Test_Mock_MediaPackageVODPackagingGroup_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMediaPackageVODClient)
	mockClient.On("ListPackagingGroups", mock.Anything, mock.Anything).
		Return(&mediapackagevod.ListPackagingGroupsOutput{
			PackagingGroups: []mediapackagevodtypes.PackagingGroup{
				{Id: ptr.String("my-group"), Arn: ptr.String("arn:aws:mediapackage-vod:us-east-1:123456789012:packaging-groups/my-group")},
			},
		}, nil)
	lister := &MediaPackageVODPackagingGroupLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testMediaPackageVODListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	pg := resources[0].(*MediaPackageVODPackagingGroup)
	a.Equal("my-group", *pg.ID)
	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaPackageVODPackagingGroup_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMediaPackageVODClient)
	mockClient.On("ListPackagingGroups", mock.Anything, mock.Anything).
		Return(&mediapackagevod.ListPackagingGroupsOutput{PackagingGroups: []mediapackagevodtypes.PackagingGroup{}}, nil)
	lister := &MediaPackageVODPackagingGroupLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testMediaPackageVODListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaPackageVODPackagingGroup_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockMediaPackageVODClient)
	pg := &MediaPackageVODPackagingGroup{svc: mockClient, ID: ptr.String("my-group")}
	mockClient.On("DeletePackagingGroup", mock.Anything, &mediapackagevod.DeletePackagingGroupInput{Id: pg.ID}).
		Return(&mediapackagevod.DeletePackagingGroupOutput{}, nil)
	a.NoError(pg.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_MediaPackageVODPackagingGroup_Properties(t *testing.T) {
	a := assert.New(t)
	pg := MediaPackageVODPackagingGroup{
		ID:  ptr.String("my-group"),
		Arn: ptr.String("arn:aws:mediapackage-vod:us-east-1:123456789012:packaging-groups/my-group"),
	}
	a.Equal("my-group", pg.Properties().Get("Id"))
}

func Test_Mock_MediaPackageVODPackagingGroup_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("my-group", (&MediaPackageVODPackagingGroup{ID: ptr.String("my-group")}).String())
}
