package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/gamelift"
	gamelifttypes "github.com/aws/aws-sdk-go-v2/service/gamelift/types"
)

func Test_Mock_GameLiftVpcPeeringAuthorization_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGameLiftV2Client)
	mockClient.On("DescribeVpcPeeringAuthorizations", mock.Anything, mock.Anything).
		Return(&gamelift.DescribeVpcPeeringAuthorizationsOutput{
			VpcPeeringAuthorizations: []gamelifttypes.VpcPeeringAuthorization{
				{GameLiftAwsAccountId: ptr.String("123456789012"), PeerVpcId: ptr.String("vpc-12345")},
			},
		}, nil)
	lister := &GameLiftVpcPeeringAuthorizationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGameLiftV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	a.Equal("123456789012:vpc-12345", resources[0].(*GameLiftVpcPeeringAuthorization).String())
	mockClient.AssertExpectations(t)
}

func Test_Mock_GameLiftVpcPeeringAuthorization_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGameLiftV2Client)
	mockClient.On("DescribeVpcPeeringAuthorizations", mock.Anything, mock.Anything).
		Return(&gamelift.DescribeVpcPeeringAuthorizationsOutput{VpcPeeringAuthorizations: []gamelifttypes.VpcPeeringAuthorization{}}, nil)
	lister := &GameLiftVpcPeeringAuthorizationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGameLiftV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_GameLiftVpcPeeringAuthorization_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGameLiftV2Client)
	r := &GameLiftVpcPeeringAuthorization{
		svc:                  mockClient,
		GameLiftAwsAccountID: ptr.String("123456789012"),
		PeerVpcID:            ptr.String("vpc-12345"),
	}
	mockClient.On("DeleteVpcPeeringAuthorization", mock.Anything, &gamelift.DeleteVpcPeeringAuthorizationInput{
		GameLiftAwsAccountId: r.GameLiftAwsAccountID, PeerVpcId: r.PeerVpcID,
	}).Return(&gamelift.DeleteVpcPeeringAuthorizationOutput{}, nil)
	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_GameLiftVpcPeeringAuthorization_Properties(t *testing.T) {
	a := assert.New(t)
	r := GameLiftVpcPeeringAuthorization{
		GameLiftAwsAccountID: ptr.String("123456789012"),
		PeerVpcID:            ptr.String("vpc-12345"),
	}
	a.Equal("123456789012", r.Properties().Get("GameLiftAwsAccountId"))
	a.Equal("vpc-12345", r.Properties().Get("PeerVpcId"))
}

func Test_Mock_GameLiftVpcPeeringAuthorization_String(t *testing.T) {
	a := assert.New(t)
	r := &GameLiftVpcPeeringAuthorization{
		GameLiftAwsAccountID: ptr.String("123456789012"),
		PeerVpcID:            ptr.String("vpc-12345"),
	}
	a.Equal("123456789012:vpc-12345", r.String())
}
