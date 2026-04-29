package resources

import (
	"context"
	"fmt"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func Test_Mock_EC2IPAMPool_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeIpamPools", mock.Anything, mock.Anything).
		Return(
			&ec2.DescribeIpamPoolsOutput{
				IpamPools: []ec2types.IpamPool{
					{
						IpamPoolId:    ptr.String("ipam-pool-11111111111111111"),
						IpamScopeArn:  ptr.String("arn:aws:ec2::123456789012:ipam-scope/ipam-scope-aaa"),
						AddressFamily: ec2types.AddressFamilyIpv4,
						State:         ec2types.IpamPoolStateCreateComplete,
						Tags: []ec2types.Tag{
							{Key: ptr.String("Name"), Value: ptr.String("test-pool")},
						},
					},
				},
			}, nil,
		)

	lister := &EC2IPAMPoolLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	pool := resources[0].(*EC2IPAMPool)
	assertions.Equal("ipam-pool-11111111111111111", *pool.IpamPoolID)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2IPAMPool_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeIpamPools", mock.Anything, mock.Anything).
		Return(
			&ec2.DescribeIpamPoolsOutput{
				IpamPools: []ec2types.IpamPool{},
			}, nil,
		)

	lister := &EC2IPAMPoolLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2IPAMPool_List_MultiPage(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	firstPageItems := make([]ec2types.IpamPool, 100)
	for i := range firstPageItems {
		firstPageItems[i] = ec2types.IpamPool{
			IpamPoolId:    ptr.String(fmt.Sprintf("ipam-pool-%d", i)),
			IpamScopeArn:  ptr.String(fmt.Sprintf("arn:aws:ec2::123456789012:ipam-scope/scope-%d", i)),
			AddressFamily: ec2types.AddressFamilyIpv4,
			State:         ec2types.IpamPoolStateCreateComplete,
		}
	}

	mockClient.
		On(
			"DescribeIpamPools",
			mock.Anything,
			mock.MatchedBy(func(input *ec2.DescribeIpamPoolsInput) bool {
				return input.NextToken == nil
			}),
		).
		Return(
			&ec2.DescribeIpamPoolsOutput{
				IpamPools: firstPageItems,
				NextToken: ptr.String("page2"),
			}, nil,
		).
		Once()

	mockClient.
		On(
			"DescribeIpamPools",
			mock.Anything,
			mock.MatchedBy(func(input *ec2.DescribeIpamPoolsInput) bool {
				return input.NextToken != nil
			}),
		).
		Return(
			&ec2.DescribeIpamPoolsOutput{
				IpamPools: []ec2types.IpamPool{
					{
						IpamPoolId:    ptr.String("ipam-pool-100"),
						IpamScopeArn:  ptr.String("arn:aws:ec2::123456789012:ipam-scope/scope-100"),
						AddressFamily: ec2types.AddressFamilyIpv6,
						State:         ec2types.IpamPoolStateCreateComplete,
					},
				},
			}, nil,
		).
		Once()

	lister := &EC2IPAMPoolLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 101)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2IPAMPool_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	pool := &EC2IPAMPool{
		svc:        mockClient,
		IpamPoolID: ptr.String("ipam-pool-11111111111111111"),
	}

	mockClient.
		On(
			"DeleteIpamPool",
			mock.Anything,
			&ec2.DeleteIpamPoolInput{
				IpamPoolId: pool.IpamPoolID,
			},
		).
		Return(&ec2.DeleteIpamPoolOutput{}, nil)

	err := pool.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2IPAMPool_Properties(t *testing.T) {
	assertions := assert.New(t)

	pool := EC2IPAMPool{
		IpamPoolID:    ptr.String("ipam-pool-11111111111111111"),
		IpamScopeID:   ptr.String("arn:aws:ec2::123456789012:ipam-scope/ipam-scope-aaa"),
		AddressFamily: "ipv4",
		State:         "create-complete",
		Tags: []ec2types.Tag{
			{Key: ptr.String("Environment"), Value: ptr.String("production")},
		},
	}

	properties := pool.Properties()

	assertions.Equal("ipam-pool-11111111111111111", properties.Get("IpamPoolId"))
	assertions.Equal("ipv4", properties.Get("AddressFamily"))
	assertions.Equal("create-complete", properties.Get("State"))
	assertions.Equal("production", properties.Get("tag:Environment"))
}

func Test_Mock_EC2IPAMPool_Properties_EmptyTags(t *testing.T) {
	assertions := assert.New(t)

	pool := EC2IPAMPool{
		IpamPoolID:    ptr.String("ipam-pool-99999999999999999"),
		IpamScopeID:   ptr.String("arn:aws:ec2::123456789012:ipam-scope/scope-zzz"),
		AddressFamily: "ipv6",
		State:         "create-complete",
		Tags:          []ec2types.Tag{},
	}

	properties := pool.Properties()

	assertions.Equal("ipam-pool-99999999999999999", properties.Get("IpamPoolId"))
	assertions.Equal("ipv6", properties.Get("AddressFamily"))
}

func Test_Mock_EC2IPAMPool_String(t *testing.T) {
	assertions := assert.New(t)

	pool := EC2IPAMPool{
		IpamPoolID: ptr.String("ipam-pool-11111111111111111"),
	}

	assertions.Equal("ipam-pool-11111111111111111", pool.String())
}

func Test_Mock_EC2IPAMPool_Filter_ExcludesDeletedState(t *testing.T) {
	assertions := assert.New(t)

	pool := EC2IPAMPool{
		IpamPoolID: ptr.String("ipam-pool-deleted"),
		State:      string(ec2types.IpamPoolStateDeleteComplete),
	}

	err := pool.Filter()
	assertions.Error(err)
	assertions.Contains(err.Error(), "already deleted")
}

func Test_Mock_EC2IPAMPool_Filter_PassesActiveState(t *testing.T) {
	assertions := assert.New(t)

	pool := EC2IPAMPool{
		IpamPoolID: ptr.String("ipam-pool-active"),
		State:      string(ec2types.IpamPoolStateCreateComplete),
	}

	err := pool.Filter()
	assertions.NoError(err)
}
