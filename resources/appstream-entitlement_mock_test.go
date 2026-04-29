package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/appstream"
	appstreamtypes "github.com/aws/aws-sdk-go-v2/service/appstream/types"
)

func Test_Mock_AppStreamEntitlement_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAppStreamClient)

	mockClient.On("DescribeStacks", mock.Anything, mock.Anything).
		Return(&appstream.DescribeStacksOutput{
			Stacks: []appstreamtypes.Stack{
				{Name: ptr.String("test-stackname")},
			},
		}, nil)

	mockClient.On("DescribeEntitlements", mock.Anything, mock.Anything).
		Return(&appstream.DescribeEntitlementsOutput{
			Entitlements: []appstreamtypes.Entitlement{
				{Name: ptr.String("test-name"), Description: ptr.String("test-description")},
			},
		}, nil)

	lister := &AppStreamEntitlementLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testAppStreamListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*AppStreamEntitlement)
	a.Equal("test-stackname", *r.StackName)
	mockClient.AssertExpectations(t)
}

func Test_Mock_AppStreamEntitlement_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAppStreamClient)

	mockClient.On("DescribeStacks", mock.Anything, mock.Anything).
		Return(&appstream.DescribeStacksOutput{
			Stacks: []appstreamtypes.Stack{},
		}, nil)

	lister := &AppStreamEntitlementLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testAppStreamListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_AppStreamEntitlement_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAppStreamClient)

	r := &AppStreamEntitlement{
		svc:       mockClient,
		Name:      ptr.String("test-name"),
		StackName: ptr.String("test-stackname"),
	}

	mockClient.On("DeleteEntitlement", mock.Anything,
		&appstream.DeleteEntitlementInput{
			Name:      r.Name,
			StackName: r.StackName,
		}).Return(&appstream.DeleteEntitlementOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_AppStreamEntitlement_Properties(t *testing.T) {
	a := assert.New(t)
	r := &AppStreamEntitlement{
		StackName:   ptr.String("test-stackname"),
		Name:        ptr.String("test-name"),
		Description: ptr.String("test-description"),
	}
	props := r.Properties()
	a.Equal("test-stackname", props.Get("StackName"))
	a.Equal("test-name", props.Get("Name"))
	a.Equal("test-description", props.Get("Description"))
}

func Test_Mock_AppStreamEntitlement_String(t *testing.T) {
	a := assert.New(t)
	r := &AppStreamEntitlement{
		Name: ptr.String("test-name"),
	}
	a.Equal("test-name", r.String())
}
