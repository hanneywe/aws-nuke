package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/backupgateway"
	bgtypes "github.com/aws/aws-sdk-go-v2/service/backupgateway/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

var testBackupGatewayListerOpts = &nuke.ListerOpts{
	Logger: logrus.WithField("test", true),
}

func Test_Mock_BackupGatewayHypervisor_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockBackupGatewayClient)

	mockClient.
		On("ListHypervisors", mock.Anything, mock.Anything).
		Return(
			&backupgateway.ListHypervisorsOutput{
				Hypervisors: []bgtypes.Hypervisor{
					{
						HypervisorArn: ptr.String("arn:aws:backup-gateway:us-east-1:123456789012:hypervisor/hv-1234567890"),
						Name:          ptr.String("test-hypervisor"),
					},
				},
			}, nil,
		)

	lister := &BackupGatewayHypervisorLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testBackupGatewayListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	hv := resources[0].(*BackupGatewayHypervisor)
	assertions.Equal("arn:aws:backup-gateway:us-east-1:123456789012:hypervisor/hv-1234567890", *hv.HypervisorArn)
	assertions.Equal("test-hypervisor", *hv.Name)

	mockClient.AssertExpectations(t)
}

func Test_Mock_BackupGatewayHypervisor_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockBackupGatewayClient)

	mockClient.
		On("ListHypervisors", mock.Anything, mock.Anything).
		Return(
			&backupgateway.ListHypervisorsOutput{
				Hypervisors: []bgtypes.Hypervisor{},
			}, nil,
		)

	lister := &BackupGatewayHypervisorLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testBackupGatewayListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_BackupGatewayHypervisor_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockBackupGatewayClient)

	hv := &BackupGatewayHypervisor{
		svc:           mockClient,
		HypervisorArn: ptr.String("arn:aws:backup-gateway:us-east-1:123456789012:hypervisor/hv-1234567890"),
	}

	mockClient.
		On("DeleteHypervisor", mock.Anything, &backupgateway.DeleteHypervisorInput{
			HypervisorArn: hv.HypervisorArn,
		}).
		Return(&backupgateway.DeleteHypervisorOutput{}, nil)

	err := hv.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_BackupGatewayHypervisor_Properties(t *testing.T) {
	assertions := assert.New(t)

	hv := BackupGatewayHypervisor{
		HypervisorArn: ptr.String("arn:aws:backup-gateway:us-east-1:123456789012:hypervisor/hv-1234567890"),
		Name:          ptr.String("test-hypervisor"),
	}

	properties := hv.Properties()

	assertions.Equal("arn:aws:backup-gateway:us-east-1:123456789012:hypervisor/hv-1234567890", properties.Get("HypervisorArn"))
	assertions.Equal("test-hypervisor", properties.Get("Name"))
}

func Test_Mock_BackupGatewayHypervisor_String(t *testing.T) {
	assertions := assert.New(t)

	hv := BackupGatewayHypervisor{
		HypervisorArn: ptr.String("arn:aws:backup-gateway:us-east-1:123456789012:hypervisor/hv-1234567890"),
	}

	assertions.Equal("arn:aws:backup-gateway:us-east-1:123456789012:hypervisor/hv-1234567890", hv.String())
}
