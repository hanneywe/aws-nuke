package resources

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Mock_EC2VpcEncryptionControl_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeVpcEncryptionControls", mock.Anything, mock.Anything).
		Return(
			&ec2.DescribeVpcEncryptionControlsOutput{
				VpcEncryptionControls: []ec2types.VpcEncryptionControl{
					{
						VpcId:                  ptr.String("vpc-11111111111111111"),
						VpcEncryptionControlId: ptr.String("vec-22222222222222222"),
					},
				},
			}, nil,
		)

	lister := &EC2VpcEncryptionControlLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	control := resources[0].(*EC2VpcEncryptionControl)
	assertions.Equal("vpc-11111111111111111", *control.VpcID)
	assertions.Equal("vec-22222222222222222", *control.VpcEncryptionControlID)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2VpcEncryptionControl_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeVpcEncryptionControls", mock.Anything, mock.Anything).
		Return(
			&ec2.DescribeVpcEncryptionControlsOutput{
				VpcEncryptionControls: []ec2types.VpcEncryptionControl{},
			}, nil,
		)

	lister := &EC2VpcEncryptionControlLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2VpcEncryptionControl_List_MultiPage(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	firstPageItems := make([]ec2types.VpcEncryptionControl, 100)
	for index := range firstPageItems {
		firstPageItems[index] = ec2types.VpcEncryptionControl{
			VpcId:                  ptr.String(fmt.Sprintf("vpc-%d", index)),
			VpcEncryptionControlId: ptr.String(fmt.Sprintf("vec-%d", index)),
		}
	}

	mockClient.
		On(
			"DescribeVpcEncryptionControls",
			mock.Anything,
			mock.MatchedBy(func(input *ec2.DescribeVpcEncryptionControlsInput) bool {
				return input.NextToken == nil
			}),
		).
		Return(
			&ec2.DescribeVpcEncryptionControlsOutput{
				VpcEncryptionControls: firstPageItems,
				NextToken:             ptr.String("page2"),
			}, nil,
		).
		Once()

	mockClient.
		On(
			"DescribeVpcEncryptionControls",
			mock.Anything,
			mock.MatchedBy(func(input *ec2.DescribeVpcEncryptionControlsInput) bool {
				return input.NextToken != nil
			}),
		).
		Return(
			&ec2.DescribeVpcEncryptionControlsOutput{
				VpcEncryptionControls: []ec2types.VpcEncryptionControl{
					{
						VpcId:                  ptr.String("vpc-100"),
						VpcEncryptionControlId: ptr.String("vec-100"),
					},
				},
			}, nil,
		).
		Once()

	lister := &EC2VpcEncryptionControlLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 101)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2VpcEncryptionControl_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	control := &EC2VpcEncryptionControl{
		svc:                    mockClient,
		VpcID:                  ptr.String("vpc-11111111111111111"),
		VpcEncryptionControlID: ptr.String("vec-22222222222222222"),
	}

	mockClient.
		On(
			"DeleteVpcEncryptionControl",
			mock.Anything,
			&ec2.DeleteVpcEncryptionControlInput{
				VpcEncryptionControlId: control.VpcEncryptionControlID,
			},
		).
		Return(&ec2.DeleteVpcEncryptionControlOutput{}, nil)

	err := control.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2VpcEncryptionControl_Properties(t *testing.T) {
	assertions := assert.New(t)

	control := EC2VpcEncryptionControl{
		VpcID:                  ptr.String("vpc-11111111111111111"),
		VpcEncryptionControlID: ptr.String("vec-22222222222222222"),
	}

	properties := control.Properties()

	assertions.Equal("vpc-11111111111111111", properties.Get("VpcId"))
	assertions.Equal("vec-22222222222222222", properties.Get("VpcEncryptionControlId"))
}

func Test_Mock_EC2VpcEncryptionControl_String(t *testing.T) {
	assertions := assert.New(t)

	control := EC2VpcEncryptionControl{
		VpcEncryptionControlID: ptr.String("vec-22222222222222222"),
	}

	assertions.Equal("vec-22222222222222222", control.String())
}
