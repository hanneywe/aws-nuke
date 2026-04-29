package resources

import (
	"context"
	"fmt"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func Test_Mock_EC2IPAM_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeIpams", mock.Anything, mock.Anything).
		Return(
			&ec2.DescribeIpamsOutput{
				Ipams: []ec2types.Ipam{
					{
						IpamId:  ptr.String("ipam-11111111111111111"),
						OwnerId: ptr.String("123456789012"),
						State:   ec2types.IpamStateCreateComplete,
						Tags: []ec2types.Tag{
							{Key: ptr.String("Name"), Value: ptr.String("test-ipam")},
						},
					},
				},
			}, nil,
		)

	lister := &EC2IPAMLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	ipam := resources[0].(*EC2IPAM)
	assertions.Equal("ipam-11111111111111111", *ipam.IpamID)
	assertions.Equal("123456789012", *ipam.OwnerID)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2IPAM_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeIpams", mock.Anything, mock.Anything).
		Return(
			&ec2.DescribeIpamsOutput{
				Ipams: []ec2types.Ipam{},
			}, nil,
		)

	lister := &EC2IPAMLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2IPAM_List_MultiPage(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	firstPageItems := make([]ec2types.Ipam, 100)
	for i := range firstPageItems {
		firstPageItems[i] = ec2types.Ipam{
			IpamId:  ptr.String(fmt.Sprintf("ipam-%d", i)),
			OwnerId: ptr.String("123456789012"),
			State:   ec2types.IpamStateCreateComplete,
		}
	}

	mockClient.
		On(
			"DescribeIpams",
			mock.Anything,
			mock.MatchedBy(func(input *ec2.DescribeIpamsInput) bool {
				return input.NextToken == nil
			}),
		).
		Return(
			&ec2.DescribeIpamsOutput{
				Ipams:     firstPageItems,
				NextToken: ptr.String("page2"),
			}, nil,
		).
		Once()

	mockClient.
		On(
			"DescribeIpams",
			mock.Anything,
			mock.MatchedBy(func(input *ec2.DescribeIpamsInput) bool {
				return input.NextToken != nil
			}),
		).
		Return(
			&ec2.DescribeIpamsOutput{
				Ipams: []ec2types.Ipam{
					{
						IpamId:  ptr.String("ipam-100"),
						OwnerId: ptr.String("123456789012"),
						State:   ec2types.IpamStateCreateComplete,
					},
				},
			}, nil,
		).
		Once()

	lister := &EC2IPAMLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 101)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2IPAM_Remove_WithCascade(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	ipam := &EC2IPAM{
		svc:    mockClient,
		IpamID: ptr.String("ipam-11111111111111111"),
	}

	mockClient.
		On(
			"DeleteIpam",
			mock.Anything,
			&ec2.DeleteIpamInput{
				IpamId:  ipam.IpamID,
				Cascade: aws.Bool(true),
			},
		).
		Return(&ec2.DeleteIpamOutput{}, nil)

	err := ipam.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2IPAM_Properties(t *testing.T) {
	assertions := assert.New(t)

	ipam := EC2IPAM{
		IpamID:  ptr.String("ipam-11111111111111111"),
		OwnerID: ptr.String("123456789012"),
		State:   "create-complete",
		Tags: []ec2types.Tag{
			{Key: ptr.String("Environment"), Value: ptr.String("production")},
		},
	}

	properties := ipam.Properties()

	assertions.Equal("ipam-11111111111111111", properties.Get("IpamId"))
	assertions.Equal("123456789012", properties.Get("OwnerId"))
	assertions.Equal("create-complete", properties.Get("State"))
	assertions.Equal("production", properties.Get("tag:Environment"))
}

func Test_Mock_EC2IPAM_Properties_EmptyTags(t *testing.T) {
	assertions := assert.New(t)

	ipam := EC2IPAM{
		IpamID:  ptr.String("ipam-99999999999999999"),
		OwnerID: ptr.String("111111111111"),
		State:   "create-complete",
		Tags:    []ec2types.Tag{},
	}

	properties := ipam.Properties()

	assertions.Equal("ipam-99999999999999999", properties.Get("IpamId"))
	assertions.Equal("create-complete", properties.Get("State"))
}

func Test_Mock_EC2IPAM_String(t *testing.T) {
	assertions := assert.New(t)

	ipam := EC2IPAM{
		IpamID: ptr.String("ipam-11111111111111111"),
	}

	assertions.Equal("ipam-11111111111111111", ipam.String())
}

func Test_Mock_EC2IPAM_Filter_ExcludesDeletedState(t *testing.T) {
	assertions := assert.New(t)

	ipam := EC2IPAM{
		IpamID: ptr.String("ipam-deleted"),
		State:  string(ec2types.IpamStateDeleteComplete),
	}

	err := ipam.Filter()
	assertions.Error(err)
	assertions.Contains(err.Error(), "already deleted")
}

func Test_Mock_EC2IPAM_Filter_PassesActiveState(t *testing.T) {
	assertions := assert.New(t)

	ipam := EC2IPAM{
		IpamID: ptr.String("ipam-active"),
		State:  string(ec2types.IpamStateCreateComplete),
	}

	err := ipam.Filter()
	assertions.NoError(err)
}
