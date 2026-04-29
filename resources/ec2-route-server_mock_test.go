package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func Test_Mock_EC2RouteServer_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockEC2Client)
	mockClient.On("DescribeRouteServers", mock.Anything, mock.Anything).
		Return(&ec2.DescribeRouteServersOutput{
			RouteServers: []ec2types.RouteServer{
				{RouteServerId: ptr.String("rs-12345678")},
			},
		}, nil)
	lister := &EC2RouteServerLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2RouteServer_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockEC2Client)
	mockClient.On("DescribeRouteServers", mock.Anything, mock.Anything).
		Return(&ec2.DescribeRouteServersOutput{RouteServers: []ec2types.RouteServer{}}, nil)
	lister := &EC2RouteServerLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2RouteServer_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockEC2Client)
	r := &EC2RouteServer{svc: mockClient, RouteServerID: ptr.String("rs-12345678")}
	mockClient.On("DeleteRouteServer", mock.Anything, &ec2.DeleteRouteServerInput{
		RouteServerId: r.RouteServerID,
	}).Return(&ec2.DeleteRouteServerOutput{}, nil)
	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2RouteServer_Properties(t *testing.T) {
	a := assert.New(t)
	r := EC2RouteServer{RouteServerID: ptr.String("rs-12345678")}
	a.Equal("rs-12345678", r.Properties().Get("RouteServerId"))
}

func Test_Mock_EC2RouteServer_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("rs-12345678", (&EC2RouteServer{RouteServerID: ptr.String("rs-12345678")}).String())
}

func Test_Mock_EC2RouteServer_Filter_Deleted(t *testing.T) {
	a := assert.New(t)
	r := EC2RouteServer{RouteServerID: ptr.String("rs-12345678"), State: ptr.String("deleted")}
	err := r.Filter()
	a.Error(err)
	a.Contains(err.Error(), "already deleted")
}

func Test_Mock_EC2RouteServer_Filter_Deleting(t *testing.T) {
	a := assert.New(t)
	r := EC2RouteServer{RouteServerID: ptr.String("rs-12345678"), State: ptr.String("deleting")}
	err := r.Filter()
	a.Error(err)
	a.Contains(err.Error(), "already deleting")
}

func Test_Mock_EC2RouteServer_Filter_Available(t *testing.T) {
	a := assert.New(t)
	r := EC2RouteServer{RouteServerID: ptr.String("rs-12345678"), State: ptr.String("available")}
	a.NoError(r.Filter())
}
