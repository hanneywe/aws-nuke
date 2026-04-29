package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/efs"
	efstypes "github.com/aws/aws-sdk-go-v2/service/efs/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

var testEFSV2ListerOpts = &nuke.ListerOpts{
	Logger: logrus.WithField("test", true),
}

func Test_Mock_EFSAccessPoint_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEFSV2Client)

	mockClient.On("DescribeAccessPoints", mock.Anything, mock.Anything).
		Return(&efs.DescribeAccessPointsOutput{
			AccessPoints: []efstypes.AccessPointDescription{
				{
					AccessPointId:  ptr.String("fsap-12345"),
					AccessPointArn: ptr.String("arn:aws:elasticfilesystem:us-east-1:123456789012:access-point/fsap-12345"),
					Name:           ptr.String("test-ap"),
					LifeCycleState: efstypes.LifeCycleStateAvailable,
				},
			},
		}, nil)

	lister := &EFSAccessPointLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testEFSV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	accessPoint := resources[0].(*EFSAccessPoint)
	assertions.Equal("fsap-12345", *accessPoint.AccessPointID)
	assertions.Equal("test-ap", *accessPoint.Name)
	mockClient.AssertExpectations(t)
}

func Test_Mock_EFSAccessPoint_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEFSV2Client)

	mockClient.On("DescribeAccessPoints", mock.Anything, mock.Anything).
		Return(&efs.DescribeAccessPointsOutput{
			AccessPoints: []efstypes.AccessPointDescription{},
		}, nil)

	lister := &EFSAccessPointLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testEFSV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_EFSAccessPoint_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEFSV2Client)

	accessPoint := &EFSAccessPoint{
		svc:           mockClient,
		AccessPointID: ptr.String("fsap-12345"),
	}

	mockClient.On("DeleteAccessPoint", mock.Anything, &efs.DeleteAccessPointInput{
		AccessPointId: accessPoint.AccessPointID,
	}).Return(&efs.DeleteAccessPointOutput{}, nil)

	err := accessPoint.Remove(context.TODO())
	assertions.NoError(err)
	mockClient.AssertExpectations(t)
}

func Test_Mock_EFSAccessPoint_Properties(t *testing.T) {
	assertions := assert.New(t)

	accessPoint := EFSAccessPoint{
		AccessPointID:  ptr.String("fsap-12345"),
		Name:           ptr.String("test-ap"),
		LifeCycleState: ptr.String("available"),
	}

	properties := accessPoint.Properties()
	assertions.Equal("fsap-12345", properties.Get("AccessPointId"))
	assertions.Equal("test-ap", properties.Get("Name"))
}

func Test_Mock_EFSAccessPoint_String(t *testing.T) {
	assertions := assert.New(t)
	accessPoint := EFSAccessPoint{AccessPointID: ptr.String("fsap-12345")}
	assertions.Equal("fsap-12345", accessPoint.String())
}

func Test_Mock_EFSAccessPoint_Filter(t *testing.T) {
	assertions := assert.New(t)

	deletingAP := EFSAccessPoint{LifeCycleState: ptr.String("deleting")}
	assertions.Error(deletingAP.Filter())

	deletedAP := EFSAccessPoint{LifeCycleState: ptr.String("deleted")}
	assertions.Error(deletedAP.Filter())

	availableAP := EFSAccessPoint{LifeCycleState: ptr.String("available")}
	assertions.NoError(availableAP.Filter())
}
