package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/codedeploy"
)

func Test_Mock_CodeDeployOnPremisesInstance_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCodeDeployV2Client)

	mockClient.On("ListOnPremisesInstances", mock.Anything, mock.Anything).
		Return(&codedeploy.ListOnPremisesInstancesOutput{
			InstanceNames: []string{"instance-1", "instance-2"},
		}, nil)

	lister := &CodeDeployOnPremisesInstanceLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCodeDeployV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 2)

	r := resources[0].(*CodeDeployOnPremisesInstance)
	a.Equal("instance-1", *r.InstanceName)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CodeDeployOnPremisesInstance_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCodeDeployV2Client)

	mockClient.On("ListOnPremisesInstances", mock.Anything, mock.Anything).
		Return(&codedeploy.ListOnPremisesInstancesOutput{
			InstanceNames: []string{},
		}, nil)

	lister := &CodeDeployOnPremisesInstanceLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCodeDeployV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CodeDeployOnPremisesInstance_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCodeDeployV2Client)

	r := &CodeDeployOnPremisesInstance{
		svc:          mockClient,
		InstanceName: ptr.String("instance-1"),
	}

	mockClient.On("DeregisterOnPremisesInstance", mock.Anything,
		&codedeploy.DeregisterOnPremisesInstanceInput{
			InstanceName: r.InstanceName,
		}).Return(&codedeploy.DeregisterOnPremisesInstanceOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_CodeDeployOnPremisesInstance_Properties(t *testing.T) {
	a := assert.New(t)
	r := &CodeDeployOnPremisesInstance{
		InstanceName: ptr.String("instance-1"),
	}
	props := r.Properties()
	a.Equal("instance-1", props.Get("InstanceName"))
}

func Test_Mock_CodeDeployOnPremisesInstance_String(t *testing.T) {
	a := assert.New(t)
	r := &CodeDeployOnPremisesInstance{
		InstanceName: ptr.String("instance-1"),
	}
	a.Equal("instance-1", r.String())
}
