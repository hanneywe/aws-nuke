package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/networkmanager"
	nmtypes "github.com/aws/aws-sdk-go-v2/service/networkmanager/types"
)

func Test_Mock_NetworkManagerTGWRegistration_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockNetworkManagerClient)

	mockClient.
		On("DescribeGlobalNetworks", mock.Anything, mock.Anything).
		Return(&networkmanager.DescribeGlobalNetworksOutput{
			GlobalNetworks: []nmtypes.GlobalNetwork{
				{GlobalNetworkId: ptr.String("gn-123")},
			},
		}, nil)

	mockClient.
		On("GetTransitGatewayRegistrations", mock.Anything, mock.Anything).
		Return(&networkmanager.GetTransitGatewayRegistrationsOutput{
			TransitGatewayRegistrations: []nmtypes.TransitGatewayRegistration{
				{
					GlobalNetworkId:   ptr.String("gn-123"),
					TransitGatewayArn: ptr.String("arn:aws:ec2:us-west-2:123456789012:transit-gateway/tgw-abc"),
					State:             &nmtypes.TransitGatewayRegistrationStateReason{Code: nmtypes.TransitGatewayRegistrationStateAvailable},
				},
			},
		}, nil)

	lister := &NetworkManagerTransitGatewayRegistrationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testNetworkManagerListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	reg := resources[0].(*NetworkManagerTransitGatewayRegistration)
	a.Equal("gn-123", *reg.GlobalNetworkID)
	a.Contains(*reg.TransitGatewayArn, "tgw-abc")
	a.Equal("AVAILABLE", reg.State)
	mockClient.AssertExpectations(t)
}

func Test_Mock_NetworkManagerTGWRegistration_List_EmptyGlobalNetworks(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockNetworkManagerClient)

	mockClient.
		On("DescribeGlobalNetworks", mock.Anything, mock.Anything).
		Return(&networkmanager.DescribeGlobalNetworksOutput{
			GlobalNetworks: []nmtypes.GlobalNetwork{},
		}, nil)

	lister := &NetworkManagerTransitGatewayRegistrationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testNetworkManagerListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_NetworkManagerTGWRegistration_List_NoRegistrations(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockNetworkManagerClient)

	mockClient.
		On("DescribeGlobalNetworks", mock.Anything, mock.Anything).
		Return(&networkmanager.DescribeGlobalNetworksOutput{
			GlobalNetworks: []nmtypes.GlobalNetwork{
				{GlobalNetworkId: ptr.String("gn-123")},
			},
		}, nil)

	mockClient.
		On("GetTransitGatewayRegistrations", mock.Anything, mock.Anything).
		Return(&networkmanager.GetTransitGatewayRegistrationsOutput{
			TransitGatewayRegistrations: []nmtypes.TransitGatewayRegistration{},
		}, nil)

	lister := &NetworkManagerTransitGatewayRegistrationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testNetworkManagerListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_NetworkManagerTGWRegistration_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockNetworkManagerClient)

	reg := &NetworkManagerTransitGatewayRegistration{
		svc:               mockClient,
		GlobalNetworkID:   ptr.String("gn-123"),
		TransitGatewayArn: ptr.String("arn:aws:ec2:us-west-2:123456789012:transit-gateway/tgw-abc"),
		State:             "AVAILABLE",
	}

	mockClient.
		On("DeregisterTransitGateway", mock.Anything, &networkmanager.DeregisterTransitGatewayInput{
			GlobalNetworkId:   reg.GlobalNetworkID,
			TransitGatewayArn: reg.TransitGatewayArn,
		}).
		Return(&networkmanager.DeregisterTransitGatewayOutput{}, nil)

	err := reg.Remove(context.TODO())
	a.NoError(err)
	mockClient.AssertExpectations(t)
}

func Test_Mock_NetworkManagerTGWRegistration_Filter_Deleting(t *testing.T) {
	a := assert.New(t)
	reg := &NetworkManagerTransitGatewayRegistration{
		TransitGatewayArn: ptr.String("arn:tgw"),
		State:             "DELETING",
	}
	a.Error(reg.Filter())
}

func Test_Mock_NetworkManagerTGWRegistration_Filter_Deleted(t *testing.T) {
	a := assert.New(t)
	reg := &NetworkManagerTransitGatewayRegistration{
		TransitGatewayArn: ptr.String("arn:tgw"),
		State:             "DELETED",
	}
	a.Error(reg.Filter())
}

func Test_Mock_NetworkManagerTGWRegistration_Filter_Available(t *testing.T) {
	a := assert.New(t)
	reg := &NetworkManagerTransitGatewayRegistration{
		TransitGatewayArn: ptr.String("arn:tgw"),
		State:             "AVAILABLE",
	}
	a.NoError(reg.Filter())
}

func Test_Mock_NetworkManagerTGWRegistration_Properties(t *testing.T) {
	a := assert.New(t)
	reg := &NetworkManagerTransitGatewayRegistration{
		GlobalNetworkID:   ptr.String("gn-123"),
		TransitGatewayArn: ptr.String("arn:tgw"),
		State:             "AVAILABLE",
	}
	props := reg.Properties()
	a.Equal("gn-123", props.Get("GlobalNetworkID"))
	a.Equal("arn:tgw", props.Get("TransitGatewayArn"))
	a.Equal("AVAILABLE", props.Get("State"))
}

func Test_Mock_NetworkManagerTGWRegistration_String(t *testing.T) {
	a := assert.New(t)
	reg := &NetworkManagerTransitGatewayRegistration{
		TransitGatewayArn: ptr.String("arn:aws:ec2:us-west-2:123456789012:transit-gateway/tgw-abc"),
	}
	a.Equal("arn:aws:ec2:us-west-2:123456789012:transit-gateway/tgw-abc", reg.String())
}
