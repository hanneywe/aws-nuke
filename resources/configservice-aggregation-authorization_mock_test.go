package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/configservice"
	configtypes "github.com/aws/aws-sdk-go-v2/service/configservice/types"
)

func Test_Mock_ConfigServiceAggregationAuthorization_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConfigServiceClient)
	mockClient.On("DescribeAggregationAuthorizations", mock.Anything, mock.Anything).
		Return(&configservice.DescribeAggregationAuthorizationsOutput{
			AggregationAuthorizations: []configtypes.AggregationAuthorization{
				{AuthorizedAccountId: ptr.String("123456789012"), AuthorizedAwsRegion: ptr.String("us-east-1")},
			},
		}, nil)
	lister := &ConfigServiceAggregationAuthorizationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testConfigServiceListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	a.Equal("123456789012:us-east-1", resources[0].(*ConfigServiceAggregationAuthorization).String())
	mockClient.AssertExpectations(t)
}

func Test_Mock_ConfigServiceAggregationAuthorization_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConfigServiceClient)
	mockClient.On("DescribeAggregationAuthorizations", mock.Anything, mock.Anything).
		Return(&configservice.DescribeAggregationAuthorizationsOutput{AggregationAuthorizations: []configtypes.AggregationAuthorization{}}, nil)
	lister := &ConfigServiceAggregationAuthorizationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testConfigServiceListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_ConfigServiceAggregationAuthorization_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConfigServiceClient)
	r := &ConfigServiceAggregationAuthorization{
		svc:                 mockClient,
		AuthorizedAccountID: ptr.String("123456789012"),
		AuthorizedAwsRegion: ptr.String("us-east-1"),
	}
	mockClient.On("DeleteAggregationAuthorization", mock.Anything, &configservice.DeleteAggregationAuthorizationInput{
		AuthorizedAccountId: r.AuthorizedAccountID, AuthorizedAwsRegion: r.AuthorizedAwsRegion,
	}).Return(&configservice.DeleteAggregationAuthorizationOutput{}, nil)
	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_ConfigServiceAggregationAuthorization_Properties(t *testing.T) {
	a := assert.New(t)
	r := ConfigServiceAggregationAuthorization{
		AuthorizedAccountID: ptr.String("123456789012"),
		AuthorizedAwsRegion: ptr.String("us-east-1"),
	}
	a.Equal("123456789012", r.Properties().Get("AuthorizedAccountId"))
	a.Equal("us-east-1", r.Properties().Get("AuthorizedAwsRegion"))
}

func Test_Mock_ConfigServiceAggregationAuthorization_String(t *testing.T) {
	a := assert.New(t)
	r := &ConfigServiceAggregationAuthorization{
		AuthorizedAccountID: ptr.String("123456789012"),
		AuthorizedAwsRegion: ptr.String("us-east-1"),
	}
	a.Equal("123456789012:us-east-1", r.String())
}
