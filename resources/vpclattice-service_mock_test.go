package resources

import (
	"context"
	"fmt"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/vpclattice"
	lattice "github.com/aws/aws-sdk-go-v2/service/vpclattice/types"
)

func Test_Mock_VPCLatticeService_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockVPCLatticeClient)

	mockClient.
		On("ListServices", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListServicesOutput{
				Items: []lattice.ServiceSummary{
					{
						Id:   ptr.String("svc-1"),
						Arn:  ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:service/svc-1"),
						Name: ptr.String("service-one"),
					},
				},
			}, nil,
		)

	mockClient.
		On("ListTagsForResource", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListTagsForResourceOutput{
				Tags: map[string]string{"team": "infra"},
			}, nil,
		)

	lister := &VPCLatticeServiceLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testVPCLatticeListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	service := resources[0].(*VPCLatticeService)
	assertions.Equal("service-one", *service.Name)
	assertions.Equal("infra", service.Tags["team"])

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeService_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockVPCLatticeClient)

	mockClient.
		On("ListServices", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListServicesOutput{
				Items: []lattice.ServiceSummary{},
			}, nil,
		)

	lister := &VPCLatticeServiceLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testVPCLatticeListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeService_List_MultiPage(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockVPCLatticeClient)

	firstPageItems := make([]lattice.ServiceSummary, 100)
	for i := range firstPageItems {
		firstPageItems[i] = lattice.ServiceSummary{
			Id:   ptr.String(fmt.Sprintf("svc-%d", i)),
			Arn:  ptr.String(fmt.Sprintf("arn:aws:vpc-lattice:us-east-1:123456789012:service/svc-%d", i)),
			Name: ptr.String(fmt.Sprintf("service-%d", i)),
		}
	}

	mockClient.
		On(
			"ListServices",
			mock.Anything,
			mock.MatchedBy(func(input *vpclattice.ListServicesInput) bool {
				return input.NextToken == nil
			}),
		).
		Return(
			&vpclattice.ListServicesOutput{
				Items:     firstPageItems,
				NextToken: ptr.String("page2"),
			}, nil,
		).
		Once()

	mockClient.
		On(
			"ListServices",
			mock.Anything,
			mock.MatchedBy(func(input *vpclattice.ListServicesInput) bool {
				return input.NextToken != nil
			}),
		).
		Return(
			&vpclattice.ListServicesOutput{
				Items: []lattice.ServiceSummary{
					{
						Id:   ptr.String("svc-100"),
						Arn:  ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:service/svc-100"),
						Name: ptr.String("service-100"),
					},
				},
			}, nil,
		).
		Once()

	mockClient.
		On("ListTagsForResource", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListTagsForResourceOutput{
				Tags: map[string]string{},
			}, nil,
		)

	lister := &VPCLatticeServiceLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testVPCLatticeListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 101)

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeService_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockVPCLatticeClient)

	service := &VPCLatticeService{
		svc: mockClient,
		ARN: ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:service/svc-1"),
	}

	mockClient.
		On(
			"DeleteService",
			mock.Anything,
			&vpclattice.DeleteServiceInput{
				ServiceIdentifier: service.ARN,
			},
		).
		Return(&vpclattice.DeleteServiceOutput{}, nil)

	err := service.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeService_Properties(t *testing.T) {
	assertions := assert.New(t)

	service := VPCLatticeService{
		ID:   ptr.String("svc-12345"),
		ARN:  ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:service/svc-12345"),
		Name: ptr.String("my-service"),
		Tags: map[string]string{"Environment": "production"},
	}

	properties := service.Properties()

	assertions.Equal("svc-12345", properties.Get("ID"))
	assertions.Equal("my-service", properties.Get("Name"))
	assertions.Equal("production", properties.Get("tag:Environment"))
}

func Test_Mock_VPCLatticeService_String(t *testing.T) {
	assertions := assert.New(t)

	service := VPCLatticeService{
		Name: ptr.String("my-service"),
	}

	assertions.Equal("my-service", service.String())
}
