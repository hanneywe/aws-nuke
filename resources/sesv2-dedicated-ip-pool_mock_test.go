package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/sesv2"
)

func Test_Mock_SESv2DedicatedIpPool_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSESv2Client)
	mockClient.On("ListDedicatedIpPools", mock.Anything, mock.Anything).
		Return(&sesv2.ListDedicatedIpPoolsOutput{DedicatedIpPools: []string{"my-pool"}}, nil)
	lister := &SESv2DedicatedIPPoolLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSESv2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	pool := resources[0].(*SESv2DedicatedIPPool)
	a.Equal("my-pool", *pool.PoolName)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SESv2DedicatedIpPool_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSESv2Client)
	mockClient.On("ListDedicatedIpPools", mock.Anything, mock.Anything).
		Return(&sesv2.ListDedicatedIpPoolsOutput{DedicatedIpPools: []string{}}, nil)
	lister := &SESv2DedicatedIPPoolLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSESv2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SESv2DedicatedIpPool_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSESv2Client)
	pool := &SESv2DedicatedIPPool{svc: mockClient, PoolName: ptr.String("my-pool")}
	mockClient.On("DeleteDedicatedIpPool", mock.Anything, &sesv2.DeleteDedicatedIpPoolInput{PoolName: pool.PoolName}).
		Return(&sesv2.DeleteDedicatedIpPoolOutput{}, nil)
	a.NoError(pool.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_SESv2DedicatedIpPool_Properties(t *testing.T) {
	a := assert.New(t)
	pool := SESv2DedicatedIPPool{PoolName: ptr.String("my-pool")}
	a.Equal("my-pool", pool.Properties().Get("PoolName"))
}

func Test_Mock_SESv2DedicatedIpPool_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("my-pool", (&SESv2DedicatedIPPool{PoolName: ptr.String("my-pool")}).String())
}
