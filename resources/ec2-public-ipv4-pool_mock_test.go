package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func Test_Mock_EC2PublicIpv4Pool_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.On("DescribePublicIpv4Pools", mock.Anything, mock.Anything).
		Return(&ec2.DescribePublicIpv4PoolsOutput{
			PublicIpv4Pools: []ec2types.PublicIpv4Pool{
				{
					PoolId: ptr.String("ipv4pool-ec2-12345"),
					Tags: []ec2types.Tag{
						{Key: ptr.String("Name"), Value: ptr.String("test-pool")},
					},
				},
			},
		}, nil)

	lister := &EC2PublicIpv4PoolLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	pool := resources[0].(*EC2PublicIpv4Pool)
	assertions.Equal("ipv4pool-ec2-12345", *pool.PoolID)
	assertions.Equal("test-pool", pool.Tags["Name"])
	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2PublicIpv4Pool_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.On("DescribePublicIpv4Pools", mock.Anything, mock.Anything).
		Return(&ec2.DescribePublicIpv4PoolsOutput{
			PublicIpv4Pools: []ec2types.PublicIpv4Pool{},
		}, nil)

	lister := &EC2PublicIpv4PoolLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2PublicIpv4Pool_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEC2Client)

	pool := &EC2PublicIpv4Pool{
		svc:    mockClient,
		PoolID: ptr.String("ipv4pool-ec2-12345"),
	}

	mockClient.On("DeletePublicIpv4Pool", mock.Anything, &ec2.DeletePublicIpv4PoolInput{
		PoolId: pool.PoolID,
	}).Return(&ec2.DeletePublicIpv4PoolOutput{}, nil)

	err := pool.Remove(context.TODO())
	assertions.NoError(err)
	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2PublicIpv4Pool_Properties(t *testing.T) {
	assertions := assert.New(t)

	pool := EC2PublicIpv4Pool{
		PoolID: ptr.String("ipv4pool-ec2-12345"),
		Tags:   map[string]string{"Name": "test-pool"},
	}

	properties := pool.Properties()
	assertions.Equal("ipv4pool-ec2-12345", properties.Get("PoolId"))
	assertions.Equal("test-pool", properties.Get("tag:Name"))
}

func Test_Mock_EC2PublicIpv4Pool_String(t *testing.T) {
	assertions := assert.New(t)
	pool := EC2PublicIpv4Pool{PoolID: ptr.String("ipv4pool-ec2-12345")}
	assertions.Equal("ipv4pool-ec2-12345", pool.String())
}
