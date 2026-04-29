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

func Test_Mock_VPCLatticeListener_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockVPCLatticeClient)

	mockClient.
		On("ListServices", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListServicesOutput{
				Items: []lattice.ServiceSummary{
					{Id: ptr.String("svc-1"), Name: ptr.String("service-one")},
				},
			}, nil,
		)

	mockClient.
		On("ListListeners", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListListenersOutput{
				Items: []lattice.ListenerSummary{
					{
						Id:       ptr.String("l-1"),
						Arn:      ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:listener/l-1"),
						Name:     ptr.String("listener-1"),
						Port:     ptr.Int32(80),
						Protocol: lattice.ListenerProtocolHttp,
					},
				},
			}, nil,
		)

	mockClient.
		On("ListTagsForResource", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListTagsForResourceOutput{
				Tags: map[string]string{"env": "test"},
			}, nil,
		)

	lister := &VPCLatticeListenerLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testVPCLatticeListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	listener := resources[0].(*VPCLatticeListener)
	assertions.Equal("listener-1", *listener.Name)
	assertions.Equal("svc-1", *listener.ServiceID)
	assertions.Equal("service-one", *listener.ServiceName)
	assertions.Equal(int32(80), *listener.Port)

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeListener_List_Empty_NoServices(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockVPCLatticeClient)

	mockClient.
		On("ListServices", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListServicesOutput{
				Items: []lattice.ServiceSummary{},
			}, nil,
		)

	lister := &VPCLatticeListenerLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testVPCLatticeListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeListener_List_Empty_NoListeners(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockVPCLatticeClient)

	mockClient.
		On("ListServices", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListServicesOutput{
				Items: []lattice.ServiceSummary{
					{Id: ptr.String("svc-1"), Name: ptr.String("service-one")},
				},
			}, nil,
		)

	mockClient.
		On("ListListeners", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListListenersOutput{
				Items: []lattice.ListenerSummary{},
			}, nil,
		)

	lister := &VPCLatticeListenerLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testVPCLatticeListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeListener_List_MultiPage_Services(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockVPCLatticeClient)

	// 101 services across 2 pages, each with 1 listener = 101 listeners total
	firstPageServices := make([]lattice.ServiceSummary, 100)
	for i := range firstPageServices {
		firstPageServices[i] = lattice.ServiceSummary{
			Id:   ptr.String(fmt.Sprintf("svc-%d", i)),
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
				Items:     firstPageServices,
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
					{Id: ptr.String("svc-100"), Name: ptr.String("service-100")},
				},
			}, nil,
		).
		Once()

	mockClient.
		On("ListListeners", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListListenersOutput{
				Items: []lattice.ListenerSummary{
					{
						Id:       ptr.String("l-1"),
						Arn:      ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:listener/l-1"),
						Name:     ptr.String("listener"),
						Port:     ptr.Int32(443),
						Protocol: lattice.ListenerProtocolHttps,
					},
				},
			}, nil,
		)

	mockClient.
		On("ListTagsForResource", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListTagsForResourceOutput{
				Tags: map[string]string{},
			}, nil,
		)

	lister := &VPCLatticeListenerLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testVPCLatticeListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 101)

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeListener_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockVPCLatticeClient)

	listener := &VPCLatticeListener{
		svc:       mockClient,
		ServiceID: ptr.String("svc-1"),
		ID:        ptr.String("l-1"),
	}

	mockClient.
		On(
			"DeleteListener",
			mock.Anything,
			&vpclattice.DeleteListenerInput{
				ServiceIdentifier:  listener.ServiceID,
				ListenerIdentifier: listener.ID,
			},
		).
		Return(&vpclattice.DeleteListenerOutput{}, nil)

	err := listener.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeListener_Properties(t *testing.T) {
	assertions := assert.New(t)

	listener := VPCLatticeListener{
		ID:          ptr.String("listener-12345"),
		ARN:         ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:listener/listener-12345"),
		Name:        ptr.String("my-listener"),
		ServiceID:   ptr.String("svc-111"),
		ServiceName: ptr.String("my-service"),
		Port:        ptr.Int32(443),
		Protocol:    ptr.String("HTTPS"),
		Tags:        map[string]string{"Team": "networking"},
	}

	properties := listener.Properties()

	assertions.Equal("listener-12345", properties.Get("ID"))
	assertions.Equal("my-listener", properties.Get("Name"))
	assertions.Equal("svc-111", properties.Get("ServiceID"))
	assertions.Equal("443", properties.Get("Port"))
	assertions.Equal("HTTPS", properties.Get("Protocol"))
	assertions.Equal("networking", properties.Get("tag:Team"))
}

func Test_Mock_VPCLatticeListener_String(t *testing.T) {
	assertions := assert.New(t)

	listener := VPCLatticeListener{
		ServiceName: ptr.String("my-service"),
		Name:        ptr.String("my-listener"),
	}

	assertions.Equal("my-service -> my-listener", listener.String())
}
