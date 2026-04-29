package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/codebuild"
	cbtypes "github.com/aws/aws-sdk-go-v2/service/codebuild/types"
)

func Test_Mock_CodeBuildFleet_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCodeBuildClient)

	mockClient.On("ListFleets", mock.Anything, mock.Anything).
		Return(&codebuild.ListFleetsOutput{
			Fleets: []string{"arn:aws:codebuild:us-east-1:123456789012:fleet/my-fleet"},
		}, nil)

	mockClient.On("BatchGetFleets", mock.Anything, mock.Anything).
		Return(&codebuild.BatchGetFleetsOutput{
			Fleets: []cbtypes.Fleet{
				{
					Name: ptr.String("my-fleet"),
					Arn:  ptr.String("arn:aws:codebuild:us-east-1:123456789012:fleet/my-fleet"),
					Status: &cbtypes.FleetStatus{
						StatusCode: cbtypes.FleetStatusCodeActive,
					},
					Tags: []cbtypes.Tag{
						{Key: ptr.String("env"), Value: ptr.String("test")},
					},
				},
			},
		}, nil)

	lister := &CodeBuildFleetLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCodeBuildListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	fleet := resources[0].(*CodeBuildFleet)
	a.Equal("my-fleet", *fleet.Name)
	a.Equal("arn:aws:codebuild:us-east-1:123456789012:fleet/my-fleet", *fleet.ARN)
	a.Equal(cbtypes.FleetStatusCodeActive, fleet.Status)
	a.Equal("test", fleet.Tags["env"])
	mockClient.AssertExpectations(t)
}

func Test_Mock_CodeBuildFleet_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCodeBuildClient)

	mockClient.On("ListFleets", mock.Anything, mock.Anything).
		Return(&codebuild.ListFleetsOutput{Fleets: []string{}}, nil)

	lister := &CodeBuildFleetLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCodeBuildListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CodeBuildFleet_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCodeBuildClient)

	fleet := &CodeBuildFleet{
		svc:  mockClient,
		Name: ptr.String("my-fleet"),
		ARN:  ptr.String("arn:aws:codebuild:us-east-1:123456789012:fleet/my-fleet"),
	}

	mockClient.On("DeleteFleet", mock.Anything, &codebuild.DeleteFleetInput{
		Arn: fleet.ARN,
	}).Return(&codebuild.DeleteFleetOutput{}, nil)

	a.NoError(fleet.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_CodeBuildFleet_Properties(t *testing.T) {
	a := assert.New(t)

	fleet := CodeBuildFleet{
		Name:   ptr.String("my-fleet"),
		ARN:    ptr.String("arn:aws:codebuild:us-east-1:123456789012:fleet/my-fleet"),
		Status: cbtypes.FleetStatusCodeActive,
		Tags:   map[string]string{"env": "test"},
	}

	props := fleet.Properties()
	a.Equal("my-fleet", props.Get("Name"))
	a.Equal("arn:aws:codebuild:us-east-1:123456789012:fleet/my-fleet", props.Get("ARN"))
	a.Equal("test", props.Get("tag:env"))
}

func Test_Mock_CodeBuildFleet_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("my-fleet", (&CodeBuildFleet{Name: ptr.String("my-fleet")}).String())
}
