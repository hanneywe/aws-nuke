package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	cloudtrailtypes "github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"

	libsettings "github.com/ekristen/libnuke/pkg/settings"
)

func Test_Mock_CloudTrailDashboard_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCloudTrailClient)
	arn := ptr.String("arn:aws:cloudtrail:us-east-1:123456789012:dashboard/my-dashboard")
	mockClient.On("ListDashboards", mock.Anything, mock.Anything).
		Return(&cloudtrail.ListDashboardsOutput{
			Dashboards: []cloudtrailtypes.DashboardDetail{
				{DashboardArn: arn},
			},
		}, nil)
	mockClient.On("GetDashboard", mock.Anything, &cloudtrail.GetDashboardInput{
		DashboardId: arn,
	}).Return(&cloudtrail.GetDashboardOutput{
		TerminationProtectionEnabled: ptr.Bool(false),
	}, nil)
	lister := &CloudTrailDashboardLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCloudTrailListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CloudTrailDashboard_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCloudTrailClient)
	mockClient.On("ListDashboards", mock.Anything, mock.Anything).
		Return(&cloudtrail.ListDashboardsOutput{Dashboards: []cloudtrailtypes.DashboardDetail{}}, nil)
	lister := &CloudTrailDashboardLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCloudTrailListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CloudTrailDashboard_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCloudTrailClient)
	r := &CloudTrailDashboard{svc: mockClient, DashboardArn: ptr.String("arn:aws:cloudtrail:us-east-1:123456789012:dashboard/my-dashboard")}
	mockClient.On("DeleteDashboard", mock.Anything, &cloudtrail.DeleteDashboardInput{DashboardId: r.DashboardArn}).
		Return(&cloudtrail.DeleteDashboardOutput{}, nil)
	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_CloudTrailDashboard_Remove_WithTerminationProtection(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCloudTrailClient)

	settings := &libsettings.Setting{}
	settings.Set("DisableTerminationProtection", true)

	r := &CloudTrailDashboard{
		svc:                          mockClient,
		settings:                     settings,
		DashboardArn:                 ptr.String("arn:aws:cloudtrail:us-east-1:123456789012:dashboard/my-dashboard"),
		TerminationProtectionEnabled: ptr.Bool(true),
	}

	mockClient.On("UpdateDashboard", mock.Anything, &cloudtrail.UpdateDashboardInput{
		DashboardId:                  r.DashboardArn,
		TerminationProtectionEnabled: ptr.Bool(false),
	}).Return(&cloudtrail.UpdateDashboardOutput{}, nil)

	mockClient.On("DeleteDashboard", mock.Anything, &cloudtrail.DeleteDashboardInput{DashboardId: r.DashboardArn}).
		Return(&cloudtrail.DeleteDashboardOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_CloudTrailDashboard_Remove_ProtectionNotDisabled(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCloudTrailClient)

	settings := &libsettings.Setting{}
	settings.Set("DisableTerminationProtection", false)

	r := &CloudTrailDashboard{
		svc:                          mockClient,
		settings:                     settings,
		DashboardArn:                 ptr.String("arn:aws:cloudtrail:us-east-1:123456789012:dashboard/my-dashboard"),
		TerminationProtectionEnabled: ptr.Bool(true),
	}

	// Should NOT call UpdateDashboard since setting is false
	mockClient.On("DeleteDashboard", mock.Anything, &cloudtrail.DeleteDashboardInput{DashboardId: r.DashboardArn}).
		Return(&cloudtrail.DeleteDashboardOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_CloudTrailDashboard_Properties(t *testing.T) {
	a := assert.New(t)
	dashArn := "arn:aws:cloudtrail:us-east-1:123456789012:dashboard/my-dashboard"
	r := CloudTrailDashboard{DashboardArn: ptr.String(dashArn)}
	a.Equal(dashArn, r.Properties().Get("DashboardArn"))
}

func Test_Mock_CloudTrailDashboard_String(t *testing.T) {
	a := assert.New(t)
	dashArn := ptr.String("arn:aws:cloudtrail:us-east-1:123456789012:dashboard/my-dashboard")
	a.Equal(*dashArn, (&CloudTrailDashboard{DashboardArn: dashArn}).String())
}

func Test_Mock_CloudTrailDashboard_Filter_AWSManaged(t *testing.T) {
	a := assert.New(t)
	r := CloudTrailDashboard{DashboardArn: ptr.String("arn:aws:cloudtrail:us-west-2:054148631998:dashboard/AWSCloudTrail-Overview")}
	err := r.Filter()
	a.Error(err)
	a.Contains(err.Error(), "cannot delete AWS-managed dashboard")
}

func Test_Mock_CloudTrailDashboard_Filter_Custom(t *testing.T) {
	a := assert.New(t)
	r := CloudTrailDashboard{DashboardArn: ptr.String("arn:aws:cloudtrail:us-east-1:123456789012:dashboard/my-dashboard")}
	a.NoError(r.Filter())
}
