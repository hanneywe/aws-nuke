package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/devicefarm"
	devicefarmtypes "github.com/aws/aws-sdk-go-v2/service/devicefarm/types"
)

func Test_Mock_DeviceFarmTestGridProject_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockDeviceFarmClient)

	mockClient.
		On("ListTestGridProjects", mock.Anything, mock.Anything).
		Return(&devicefarm.ListTestGridProjectsOutput{
			TestGridProjects: []devicefarmtypes.TestGridProject{
				{
					Arn:  ptr.String("arn:aws:devicefarm:us-west-2:123456789012:testgrid-project:example-1"),
					Name: ptr.String("my-test-grid-project"),
				},
			},
		}, nil)

	lister := &DeviceFarmTestGridProjectLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testDeviceFarmListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	testGridProject := resources[0].(*DeviceFarmTestGridProject)
	assertions.Equal("arn:aws:devicefarm:us-west-2:123456789012:testgrid-project:example-1", *testGridProject.Arn)
	assertions.Equal("my-test-grid-project", *testGridProject.Name)

	mockClient.AssertExpectations(t)
}

func Test_Mock_DeviceFarmTestGridProject_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockDeviceFarmClient)

	mockClient.
		On("ListTestGridProjects", mock.Anything, mock.Anything).
		Return(&devicefarm.ListTestGridProjectsOutput{
			TestGridProjects: []devicefarmtypes.TestGridProject{},
		}, nil)

	lister := &DeviceFarmTestGridProjectLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testDeviceFarmListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_DeviceFarmTestGridProject_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockDeviceFarmClient)

	testGridProject := &DeviceFarmTestGridProject{
		svc:  mockClient,
		Arn:  ptr.String("arn:aws:devicefarm:us-west-2:123456789012:testgrid-project:example-1"),
		Name: ptr.String("my-test-grid-project"),
	}

	mockClient.
		On("DeleteTestGridProject", mock.Anything, &devicefarm.DeleteTestGridProjectInput{
			ProjectArn: testGridProject.Arn,
		}).
		Return(&devicefarm.DeleteTestGridProjectOutput{}, nil)

	err := testGridProject.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_DeviceFarmTestGridProject_Properties(t *testing.T) {
	assertions := assert.New(t)

	testGridProject := DeviceFarmTestGridProject{
		Arn:  ptr.String("arn:aws:devicefarm:us-west-2:123456789012:testgrid-project:example-1"),
		Name: ptr.String("my-test-grid-project"),
	}

	properties := testGridProject.Properties()

	assertions.Equal("arn:aws:devicefarm:us-west-2:123456789012:testgrid-project:example-1", properties.Get("Arn"))
	assertions.Equal("my-test-grid-project", properties.Get("Name"))
}

func Test_Mock_DeviceFarmTestGridProject_String(t *testing.T) {
	assertions := assert.New(t)

	testGridProject := DeviceFarmTestGridProject{
		Name: ptr.String("my-test-grid-project"),
	}

	assertions.Equal("my-test-grid-project", testGridProject.String())
}
