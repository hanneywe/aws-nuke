package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/connect"
	connecttypes "github.com/aws/aws-sdk-go-v2/service/connect/types"
)

func Test_Mock_ConnectUseCase_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConnectClient)

	mockClient.
		On("ListInstances", mock.Anything, mock.Anything).
		Return(&connect.ListInstancesOutput{
			InstanceSummaryList: []connecttypes.InstanceSummary{
				{Id: ptr.String("i-12345")},
			},
		}, nil)

	mockClient.
		On("ListIntegrationAssociations", mock.Anything, mock.Anything).
		Return(&connect.ListIntegrationAssociationsOutput{
			IntegrationAssociationSummaryList: []connecttypes.IntegrationAssociationSummary{
				{IntegrationAssociationId: ptr.String("ia-12345")},
			},
		}, nil)

	mockClient.
		On("ListUseCases", mock.Anything, mock.Anything).
		Return(&connect.ListUseCasesOutput{
			UseCaseSummaryList: []connecttypes.UseCase{
				{UseCaseId: ptr.String("uc-12345")},
			},
		}, nil)

	lister := &ConnectUseCaseLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testConnectListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	uc := resources[0].(*ConnectUseCase)
	a.Equal("uc-12345", *uc.UseCaseID)
	a.Equal("ia-12345", *uc.IntegrationAssociationID)
	a.Equal("i-12345", *uc.InstanceID)

	mockClient.AssertExpectations(t)
}

func Test_Mock_ConnectUseCase_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConnectClient)

	mockClient.
		On("ListInstances", mock.Anything, mock.Anything).
		Return(&connect.ListInstancesOutput{
			InstanceSummaryList: []connecttypes.InstanceSummary{},
		}, nil)

	lister := &ConnectUseCaseLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testConnectListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_ConnectUseCase_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConnectClient)

	uc := &ConnectUseCase{
		svc:                      mockClient,
		InstanceID:               ptr.String("i-12345"),
		IntegrationAssociationID: ptr.String("ia-12345"),
		UseCaseID:                ptr.String("uc-12345"),
	}

	mockClient.
		On("DeleteUseCase", mock.Anything, &connect.DeleteUseCaseInput{
			InstanceId:               uc.InstanceID,
			IntegrationAssociationId: uc.IntegrationAssociationID,
			UseCaseId:                uc.UseCaseID,
		}).
		Return(&connect.DeleteUseCaseOutput{}, nil)

	err := uc.Remove(context.TODO())
	a.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_ConnectUseCase_Properties(t *testing.T) {
	a := assert.New(t)

	uc := ConnectUseCase{
		InstanceID:               ptr.String("i-12345"),
		IntegrationAssociationID: ptr.String("ia-12345"),
		UseCaseID:                ptr.String("uc-12345"),
	}

	props := uc.Properties()
	a.Equal("i-12345", props.Get("InstanceId"))
	a.Equal("ia-12345", props.Get("IntegrationAssociationId"))
	a.Equal("uc-12345", props.Get("UseCaseId"))
}

func Test_Mock_ConnectUseCase_String(t *testing.T) {
	a := assert.New(t)

	uc := ConnectUseCase{
		UseCaseID: ptr.String("uc-12345"),
	}

	a.Equal("uc-12345", uc.String())
}
