package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/storagegateway"
	storagegatewaytypes "github.com/aws/aws-sdk-go-v2/service/storagegateway/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

var testStorageGatewayV2ListerOpts = &nuke.ListerOpts{
	Logger: logrus.WithField("test", true),
}

func Test_Mock_StorageGatewayTapePool_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockStorageGatewayV2Client)

	mockClient.On("ListTapePools", mock.Anything, mock.Anything).
		Return(&storagegateway.ListTapePoolsOutput{
			PoolInfos: []storagegatewaytypes.PoolInfo{
				{
					PoolARN:  ptr.String("arn:aws:storagegateway:us-east-1:123456789012:tape-pool/pool-12345"),
					PoolName: ptr.String("my-tape-pool"),
				},
			},
		}, nil)

	lister := &StorageGatewayTapePoolLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testStorageGatewayV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	tapePool := resources[0].(*StorageGatewayTapePool)
	assertions.Equal("arn:aws:storagegateway:us-east-1:123456789012:tape-pool/pool-12345", *tapePool.PoolARN)
	assertions.Equal("my-tape-pool", *tapePool.PoolName)
	mockClient.AssertExpectations(t)
}

func Test_Mock_StorageGatewayTapePool_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockStorageGatewayV2Client)

	mockClient.On("ListTapePools", mock.Anything, mock.Anything).
		Return(&storagegateway.ListTapePoolsOutput{
			PoolInfos: []storagegatewaytypes.PoolInfo{},
		}, nil)

	lister := &StorageGatewayTapePoolLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testStorageGatewayV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_StorageGatewayTapePool_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockStorageGatewayV2Client)

	tapePool := &StorageGatewayTapePool{
		svc:     mockClient,
		PoolARN: ptr.String("arn:aws:storagegateway:us-east-1:123456789012:tape-pool/pool-12345"),
	}

	mockClient.On("DeleteTapePool", mock.Anything, &storagegateway.DeleteTapePoolInput{
		PoolARN: tapePool.PoolARN,
	}).Return(&storagegateway.DeleteTapePoolOutput{}, nil)

	err := tapePool.Remove(context.TODO())
	assertions.NoError(err)
	mockClient.AssertExpectations(t)
}

func Test_Mock_StorageGatewayTapePool_Properties(t *testing.T) {
	assertions := assert.New(t)

	tapePool := StorageGatewayTapePool{
		PoolARN:  ptr.String("arn:aws:storagegateway:us-east-1:123456789012:tape-pool/pool-12345"),
		PoolName: ptr.String("my-tape-pool"),
	}

	properties := tapePool.Properties()
	assertions.Equal("arn:aws:storagegateway:us-east-1:123456789012:tape-pool/pool-12345", properties.Get("PoolARN"))
	assertions.Equal("my-tape-pool", properties.Get("PoolName"))
}

func Test_Mock_StorageGatewayTapePool_String(t *testing.T) {
	assertions := assert.New(t)
	tapePool := StorageGatewayTapePool{PoolARN: ptr.String("arn:aws:storagegateway:us-east-1:123456789012:tape-pool/pool-12345")}
	assertions.Equal("arn:aws:storagegateway:us-east-1:123456789012:tape-pool/pool-12345", tapePool.String())
}
