package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/databasemigrationservice"
	dmstypes "github.com/aws/aws-sdk-go-v2/service/databasemigrationservice/types"
)

func Test_Mock_DMSDataProvider_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockDMSClient)
	mockClient.On("DescribeDataProviders", mock.Anything, mock.Anything).
		Return(&databasemigrationservice.DescribeDataProvidersOutput{
			DataProviders: []dmstypes.DataProvider{
				{DataProviderArn: ptr.String("arn:aws:dms:us-east-1:123456789012:data-provider:test"), DataProviderName: ptr.String("test-dp")},
			},
		}, nil)
	lister := &DMSDataProviderLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testDMSListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	a.Equal("test-dp", resources[0].(*DMSDataProvider).String())
	mockClient.AssertExpectations(t)
}

func Test_Mock_DMSDataProvider_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockDMSClient)
	mockClient.On("DescribeDataProviders", mock.Anything, mock.Anything).
		Return(&databasemigrationservice.DescribeDataProvidersOutput{DataProviders: []dmstypes.DataProvider{}}, nil)
	lister := &DMSDataProviderLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testDMSListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_DMSDataProvider_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockDMSClient)
	r := &DMSDataProvider{
		svc:              mockClient,
		DataProviderArn:  ptr.String("arn:aws:dms:us-east-1:123456789012:data-provider:test"),
		DataProviderName: ptr.String("test-dp"),
	}
	mockClient.On("DeleteDataProvider", mock.Anything, &databasemigrationservice.DeleteDataProviderInput{
		DataProviderIdentifier: r.DataProviderArn,
	}).Return(&databasemigrationservice.DeleteDataProviderOutput{}, nil)
	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_DMSDataProvider_Properties(t *testing.T) {
	a := assert.New(t)
	r := DMSDataProvider{
		DataProviderArn:  ptr.String("arn:aws:dms:us-east-1:123456789012:data-provider:test"),
		DataProviderName: ptr.String("test-dp"),
	}
	a.Equal("arn:aws:dms:us-east-1:123456789012:data-provider:test",
		r.Properties().Get("DataProviderArn"))
	a.Equal("test-dp", r.Properties().Get("DataProviderName"))
}

func Test_Mock_DMSDataProvider_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("test-dp", (&DMSDataProvider{DataProviderName: ptr.String("test-dp")}).String())
}
