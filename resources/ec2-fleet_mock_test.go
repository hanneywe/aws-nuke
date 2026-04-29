package resources

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Mock_EC2Fleet_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeFleets", mock.Anything, mock.Anything).
		Return(
			&ec2.DescribeFleetsOutput{
				Fleets: []ec2types.FleetData{
					{
						FleetId:    ptr.String("fleet-11111111111111111"),
						FleetState: ec2types.FleetStateCodeActive,
					},
				},
			}, nil,
		)

	lister := &EC2FleetLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	fleet := resources[0].(*EC2Fleet)
	assertions.Equal("fleet-11111111111111111", *fleet.FleetID)
	assertions.Equal(ec2types.FleetStateCodeActive, fleet.FleetState)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2Fleet_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeFleets", mock.Anything, mock.Anything).
		Return(
			&ec2.DescribeFleetsOutput{
				Fleets: []ec2types.FleetData{},
			}, nil,
		)

	lister := &EC2FleetLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2Fleet_List_MultiPage(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	firstPageItems := make([]ec2types.FleetData, 100)
	for index := range firstPageItems {
		firstPageItems[index] = ec2types.FleetData{
			FleetId:    ptr.String(fmt.Sprintf("fleet-%d", index)),
			FleetState: ec2types.FleetStateCodeActive,
		}
	}

	mockClient.
		On(
			"DescribeFleets",
			mock.Anything,
			mock.MatchedBy(func(input *ec2.DescribeFleetsInput) bool {
				return input.NextToken == nil
			}),
		).
		Return(
			&ec2.DescribeFleetsOutput{
				Fleets:    firstPageItems,
				NextToken: ptr.String("page2"),
			}, nil,
		).
		Once()

	mockClient.
		On(
			"DescribeFleets",
			mock.Anything,
			mock.MatchedBy(func(input *ec2.DescribeFleetsInput) bool {
				return input.NextToken != nil
			}),
		).
		Return(
			&ec2.DescribeFleetsOutput{
				Fleets: []ec2types.FleetData{
					{
						FleetId:    ptr.String("fleet-100"),
						FleetState: ec2types.FleetStateCodeActive,
					},
				},
			}, nil,
		).
		Once()

	lister := &EC2FleetLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 101)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2Fleet_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	fleet := &EC2Fleet{
		svc:     mockClient,
		FleetID: ptr.String("fleet-11111111111111111"),
	}

	mockClient.
		On(
			"DeleteFleets",
			mock.Anything,
			&ec2.DeleteFleetsInput{
				FleetIds:           []string{"fleet-11111111111111111"},
				TerminateInstances: aws.Bool(false),
			},
		).
		Return(&ec2.DeleteFleetsOutput{}, nil)

	err := fleet.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2Fleet_Properties(t *testing.T) {
	assertions := assert.New(t)

	fleet := EC2Fleet{
		FleetID:    ptr.String("fleet-11111111111111111"),
		FleetState: ec2types.FleetStateCodeActive,
	}

	properties := fleet.Properties()

	assertions.Equal("fleet-11111111111111111", properties.Get("FleetId"))
}

func Test_Mock_EC2Fleet_String(t *testing.T) {
	assertions := assert.New(t)

	fleet := EC2Fleet{
		FleetID: ptr.String("fleet-11111111111111111"),
	}

	assertions.Equal("fleet-11111111111111111", fleet.String())
}

func Test_Mock_EC2Fleet_Filter_ExcludesDeletedState(t *testing.T) {
	assertions := assert.New(t)

	fleet := EC2Fleet{
		FleetID:    ptr.String("fleet-deleted"),
		FleetState: ec2types.FleetStateCodeDeleted,
	}

	err := fleet.Filter()
	assertions.Error(err)
	assertions.Contains(err.Error(), "deleted")
}

func Test_Mock_EC2Fleet_Filter_ExcludesDeletedRunningState(t *testing.T) {
	assertions := assert.New(t)

	fleet := EC2Fleet{
		FleetID:    ptr.String("fleet-deleted-running"),
		FleetState: ec2types.FleetStateCodeDeletedRunning,
	}

	err := fleet.Filter()
	assertions.Error(err)
	assertions.Contains(err.Error(), "deleted_running")
}

func Test_Mock_EC2Fleet_Filter_ExcludesDeletedTerminatingState(t *testing.T) {
	assertions := assert.New(t)

	fleet := EC2Fleet{
		FleetID:    ptr.String("fleet-deleted-terminating"),
		FleetState: ec2types.FleetStateCodeDeletedTerminatingInstances,
	}

	err := fleet.Filter()
	assertions.Error(err)
	assertions.Contains(err.Error(), "deleted_terminating")
}

func Test_Mock_EC2Fleet_Filter_PassesActiveState(t *testing.T) {
	assertions := assert.New(t)

	fleet := EC2Fleet{
		FleetID:    ptr.String("fleet-active"),
		FleetState: ec2types.FleetStateCodeActive,
	}

	err := fleet.Filter()
	assertions.NoError(err)
}
